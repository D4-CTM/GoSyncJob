package jobs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	database "syncjob/Database"
	logger "syncjob/Logger"
)

type SyncQueue struct {
	Key   string
	Owner database.TableOwner
	Table database.TableMapping
}

var Queue chan SyncQueue = make(chan SyncQueue, 100)

func Init() {
	go QueueSyncJobs(Queue)
}

func QueueSyncJobs(queue chan SyncQueue) {
	var source *database.Credentials = nil
	var target *database.Credentials = nil
	for {
		task := <-queue
		pair := database.SMMap[task.Key]

		switch task.Owner {
		case database.SLAVE:
			source = &pair.Slave
			target = &pair.Master
		case database.MASTER:
			source = &pair.Master
			target = &pair.Slave
		}

		logger.LogDebug("Source cred: %", *source)
		logger.LogDebug("Target cred: %", *target)
		logger.LogInfo("Starting Sync of: %s", task.Table.GetSourceTableName())

		switch task.Owner {
		case database.SLAVE:
			if err := syncSlaveTablesByChunks(source, target, task.Table, 5000); err != nil {
				logger.LogErr("%v", err)
			}
		case database.MASTER:
			if err := syncMasterTablesByChunks(source, target, task.Table, 5000); err != nil {
				logger.LogErr("%v", err)
			}
		}

		source = nil
		target = nil
	}
}

func syncMasterTablesByChunks(
	source *database.Credentials,
	target *database.Credentials,
	table database.TableMapping,
	chunkSize int) error {
	if err := source.Ping(); err != nil {
		return err
	}

	if err := target.Ping(); err != nil {
		return err
	}

	const selectMasterQuery string = `
		SELECT %s
		FROM %s
		WHERE last_update > $1
		%s
		`

	tableName := table.GetSourceTableName()

	sourceDb := source.GetDb()
	targetDb := target.GetDb()

	var slaveColumns []string
	var masterColumns []string

	for _, c := range table.ColumnsMapped {
		slaveColumns = append(slaveColumns, c.SlaveName)
		masterColumns = append(masterColumns, c.MasterName)
	}

	masterRows := strings.Join(masterColumns, ",")
	logger.LogInfo("Searching data from: %", table.LastSync)
	resultCount := sourceDb.QueryRow(fmt.Sprintf("SELECT COUNT(1) AS count FROM %s WHERE last_update > $1", tableName), table.LastSync)
	var rowCount int

	if err := resultCount.Scan(&rowCount); err != nil {
		return err
	}

	if rowCount == 0 {
		logger.LogInfo("Not found any row to insert")
		return nil
	}

	logger.LogInfo("Found %d possible rows", rowCount)
	iterations := int(math.Ceil(float64(rowCount) / float64(chunkSize)))

	targetTN := table.GetTargetTableName()
	for i := 0; i < iterations; i++ {

		selectQuery := fmt.Sprintf(selectMasterQuery,
			masterRows,
			tableName,
			source.CreateOffsetStmt(
				i*chunkSize,
				chunkSize,
			))
		logger.LogInfo("Executing SELECT: %s", selectQuery)
		logger.LogInfo("With param: %v", table.LastSync)

		rows, err := sourceDb.Query(
			selectQuery,
			table.LastSync,
		)
		if err != nil {
			return err
		}
		pkColumn := targetTN + "_id"

		rowCount := 0
		for rows.Next() {
			columns := make([]any, len(masterColumns))
			columnPointers := make([]any, len(masterColumns))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				rows.Close()
				return err
			}

			placeholders := []string{}
			values := []string{}
			chunkData := []any{}
			col := 1
			for _, c := range columns {
				colName := masterColumns[col-1]
				placeholders = append(placeholders, target.Placeholder(col)+" AS "+colName)
				values = append(values, "s."+colName)
				chunkData = append(chunkData, c)
				col += 1
			}

			updateSet := []string{}
			for _, colName := range masterColumns {
				if colName == pkColumn {
					continue
				}
				updateSet = append(updateSet, fmt.Sprintf("t.%s = s.%s", colName, colName))
			}

			mergeStmt := fmt.Sprintf(`
				MERGE INTO %s t
				USING (SELECT %s FROM dual) s
				ON (t.%s = s.%s)
				WHEN MATCHED THEN UPDATE SET %s
				WHEN NOT MATCHED THEN INSERT (%s) VALUES (%s)`,
				targetTN,
				strings.Join(placeholders, ","),
				pkColumn,
				pkColumn,
				strings.Join(updateSet, ","),
				masterRows,
				strings.Join(values, ","),
			)

			result, err := targetDb.Exec(mergeStmt, chunkData...)
			if err != nil {
				logger.LogErr("MERGE error: %v, stmt: %s", err, mergeStmt)
				rows.Close()
				return err
			}
			_, _ = result.RowsAffected()
			rowCount++
			if rowCount%100 == 0 {
				logger.LogInfo("Merged %d rows", rowCount)
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return fmt.Errorf("error iterating rows: %w", err)
		}
		logger.LogInfo("Total merged rows: %d", rowCount)

		rows.Close()
	}

	return nil
}

func syncSlaveTablesByChunks(
	source *database.Credentials,
	target *database.Credentials,
	table database.TableMapping,
	chunkSize int) error {
	if err := source.Ping(); err != nil {
		return err
	}
	if err := target.Ping(); err != nil {
		return err
	}

	sourceDb := source.GetDb()
	targetDb := target.GetDb()

	tableName := table.GetSourceTableName()
	logTableName := strings.ToUpper(tableName) + "_LOG"
	pkColumn := tableName + "_id"

	for {
		rows, err := sourceDb.Query(fmt.Sprintf(`
			SELECT audit_id, operation, old_value, new_value
			FROM %s
			ORDER BY audit_id
			FETCH FIRST %d ROWS ONLY
		`, logTableName, chunkSize))
		if err != nil {
			return err
		}

		type logEntry struct {
			auditId   int64
			operation string
			oldValue  string
			newValue  string
		}
		entries := []logEntry{}

		for rows.Next() {
			var auditId int64
			var operation, oldValue, newValue string
			if err := rows.Scan(&auditId, &operation, &oldValue, &newValue); err != nil {
				rows.Close()
				return err
			}

			entries = append(entries, logEntry{
				auditId:   auditId,
				operation: operation,
				oldValue:  oldValue,
				newValue:  newValue,
			})
		}
		rows.Close()

		if len(entries) == 0 {
			break
		}

		for _, entry := range entries {
			if entry.operation == "DELETE" {
				// For DELETE: use oldValue to get PK and execute DELETE
				jsonData := entry.oldValue
				lowercaseJson, err := transformJsonToLowercase(jsonData)
				if err != nil {
					logger.LogErr("Failed to transform JSON for DELETE %s (audit_id=%d): %v",
						tableName, entry.auditId, err)
					continue
				}
				if err := applyDeleteToMaster(targetDb, table, lowercaseJson, pkColumn); err != nil {
					logger.LogErr("DELETE failed for %s (audit_id=%d): %v",
						tableName, entry.auditId, err)
					continue
				}
			} else {
				// For INSERT/UPDATE: use newValue
				jsonData := entry.newValue

				// Transform JSON keys to lowercase (Oracle uses uppercase)
				lowercaseJson, err := transformJsonToLowercase(jsonData)
				if err != nil {
					logger.LogErr("Failed to transform JSON for %s (audit_id=%d): %v",
						tableName, entry.auditId, err)
					continue
				}

				if err := applyMergeJsonToMaster(targetDb, table, lowercaseJson, pkColumn); err != nil {
					logger.LogErr("MERGE failed for %s (audit_id=%d): %v",
						tableName, entry.auditId, err)
					continue
				}
			}
		}

		logger.LogInfo("Processed %d entries for %s", len(entries), tableName)

		_, err = sourceDb.Exec(fmt.Sprintf("TRUNCATE TABLE %s", logTableName))
		if err != nil {
			logger.LogErr("Failed to truncate %s: %v", logTableName, err)
		}

		if len(entries) < chunkSize {
			break
		}
	}

	logger.LogInfo("Sync-OUT completed for %s", tableName)
	return nil
}

func transformJsonToLowercase(jsonStr string) (string, error) {
	var result map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return "", err
	}

	lowercaseResult := make(map[string]any)
	for key, value := range result {
		lowercaseResult[strings.ToLower(key)] = value
	}

	output, err := json.Marshal(lowercaseResult)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func applyMergeJsonToMaster(db *sql.DB, table database.TableMapping, jsonData string, pkColumn string) error {
	tableName := table.GetSourceTableName()

	upsertStmt := fmt.Sprintf(`
		INSERT INTO %s
		SELECT * FROM json_populate_record(null::%s, $1::json)
		ON CONFLICT (%s) DO UPDATE SET
			%s`,
		tableName,
		tableName,
		pkColumn,
		buildOnConflictSet(table, pkColumn),
	)

	logger.LogDebug("Upsert SQL for %s: %s", tableName, upsertStmt)
	logger.LogDebug("JSON data: %s", jsonData)

	_, err := db.Exec(upsertStmt, jsonData)
	return err
}

func applyDeleteToMaster(db *sql.DB, table database.TableMapping, jsonData string, pkColumn string) error {
	tableName := table.GetSourceTableName()

	// Extract PK value from JSON using json_populate_record
	deleteStmt := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = (SELECT %s FROM json_populate_record(null::%s, $1::json))`,
		tableName,
		pkColumn,
		pkColumn,
		tableName,
	)

	logger.LogDebug("Delete SQL for %s: %s", tableName, deleteStmt)
	logger.LogDebug("JSON data for delete: %s", jsonData)

	_, err := db.Exec(deleteStmt, jsonData)
	return err
}

func buildOnConflictSet(table database.TableMapping, pkColumn string) string {
	updateSet := []string{}
	for _, colMap := range table.ColumnsMapped {
		if colMap.MasterName == pkColumn {
			continue // Skip the primary key in UPDATE
		}
		updateSet = append(updateSet, fmt.Sprintf("%s = EXCLUDED.%s", colMap.MasterName, colMap.MasterName))
	}
	return strings.Join(updateSet, ",")
}

func buildColumnList(table database.TableMapping) string {
	columns := []string{}
	for _, colMap := range table.ColumnsMapped {
		columns = append(columns, colMap.MasterName)
	}
	return strings.Join(columns, ",")
}

func buildUpdateSet(table database.TableMapping) string {
	updateSet := []string{}
	for _, colMap := range table.ColumnsMapped {
		updateSet = append(updateSet, fmt.Sprintf("t.%s = s.%s", colMap.MasterName, colMap.MasterName))
	}
	return strings.Join(updateSet, ",")
}
