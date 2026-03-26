package jobs

import (
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
			break
		case database.MASTER:
			source = &pair.Master
			target = &pair.Slave
			break
		}

		if err := syncMasterTablesByChunks(source, target, task.Table, 5000); err != nil {
			logger.LogErr("%v", err)
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

	const selectMasterQuery string = 
		`
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

	for _, c  := range table.ColumnsMapped {
		slaveColumns = append(slaveColumns, c.SlaveName)
		masterColumns = append(masterColumns, c.MasterName)
	}

	masterRows := strings.Join(masterColumns, ",")
	resultCount := sourceDb.QueryRow(fmt.Sprintf("SELECT COUNT(1) AS count FROM %s WHERE last_update > $1", tableName), table.LastSync)
	var rowCount int

	if err := resultCount.Scan(&rowCount); err != nil {
		return err
	}

	if rowCount == 0 {
		return nil
	}

	iterations := int(math.Ceil(float64(rowCount)/float64(chunkSize)))	

	targetTN := table.GetTargetTableName()
	for i := 0; i < iterations; i++ {

		rows, err := sourceDb.Query(
			fmt.Sprintf(selectMasterQuery, 
				masterRows, 
				tableName, 
				source.CreateOffsetStmt(
					chunkSize,
					i * chunkSize,
				)),
			table.LastSync,
		)
		if err != nil {
			return err
		}
		insertStmt := fmt.Sprintf("INSERT INTO %s(%s) VALUES", targetTN, slaveColumns)

		valuesStmt := []string{}
		chunkData := []any{}
		col := 1
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

			inserts := []string{}
			values := "("
			for _, c := range columns {
				inserts = append(inserts, source.Placeholder(col))
				chunkData = append(chunkData, c)
				col += 1
			}
			values += strings.Join(inserts, ",")
			values += ")"

			valuesStmt = append(valuesStmt, values)
		}
		insertStmt += strings.Join(valuesStmt, ",")
		result, err := targetDb.Exec(insertStmt, chunkData...)
		if (err != nil) {
			rows.Close()
			return err
		}
		affectedRows, _ := result.RowsAffected()
		logger.LogInfo("Inserted %d rows!", affectedRows)

		rows.Close()
	}

	return nil
}
