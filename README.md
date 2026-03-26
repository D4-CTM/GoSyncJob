# SyncJob

A Go-based service for synchronizing data between database pairs (Master/Slave). Supports Oracle and PostgreSQL databases.

## How It Works

1. **Database Pairs**: The service manages pairs of databases (Slave ↔ Master). Each pair contains connection credentials for both databases and table mappings defining what to sync.

2. **Table Mappings**: Each pair defines which tables to synchronize, including:
   - Source and target table names
   - Column mappings between databases
   - Sync direction (Master→Slave or Slave→Master)

3. **Sync Queue**: When a sync is triggered, tasks are queued and processed asynchronously. The sync works differently based on direction:
   - **Master → Slave**: Queries for rows with `last_update` newer than last sync, then MERGEs into slave
   - **Slave → Master**: Reads from an audit/log table (`TABLE_NAME_LOG`), applies changes via MERGE/DELETE, then truncates the log

4. **Persistence**: Pairs and mappings are stored in `mappings.json` (loaded from `$CREDS_SUBDIR`)

## Requirements

- Go 1.26+
- Oracle client library (for Oracle DB support)
- Environment variable: `CREDS_SUBDIR` (directory for storing mappings.json)

## Running

```bash
cd backend
go run main.go
```

Server runs on `http://localhost:5461`

## Endpoints

### Pair Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/pairs` | List all pair names |
| GET | `/api/pairs/:key` | Get pair details by name |
| POST | `/api/pairs` | Create new pair |
| PUT | `/api/pairs/:key` | Update existing pair |
| DELETE | `/api/pairs/:key` | Delete pair |

**Create/Update Pair JSON structure:**
```json
{
  "Name": "my-pair",
  "Slave": {
    "Database": "dbname",
    "Server": "localhost",
    "Port": 5432,
    "User": "user",
    "Password": "pass",
    "DbType": 1
  },
  "Master": {
    "Database": "ORCL",
    "Server": "localhost",
    "Port": 1521,
    "User": "system",
    "Password": "pass",
    "DbType": 0
  },
  "Mappings": {
    "Tables": [
      {
        "Owner": 0,
        "MasterTableName": "USERS",
        "SlaveTableName": "users",
        "ColumnsMapped": [
          {"SlaveName": "id", "MasterName": "USER_ID"},
          {"SlaveName": "name", "MasterName": "NAME"},
          {"SlaveName": "updated_at", "MasterName": "LAST_UPDATE"}
        ]
      }
    ]
  }
}
```

`DbType`: `0` = Oracle, `1` = PostgreSQL  
`Owner`: `0` = Master, `1` = Slave

### Credentials

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/credentials/ping` | Test database connection |

### Sync

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/pairs/:key/sync` | Trigger sync for a pair |

**Sync Request JSON:**
```json
{
  "Owner": 0,
  "Table": "users",
  "All": false
}
```

- `Owner`: `0` = sync from Master to Slave, `1` = sync from Slave to Master
- `Table`: specific table name (required if `All: false`)
- `All`: `true` to sync all tables for the owner

## Example Usage

1. Create a pair:
```bash
curl -X POST http://localhost:5461/api/pairs \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "prod-sync",
    "Slave": {"Server": "localhost", "Port": 5432, "Database": "mydb", "User": "user", "Password": "pass", "DbType": 1},
    "Master": {"Server": "localhost", "Port": 1521, "Database": "ORCL", "User": "system", "Password": "pass", "DbType": 0},
    "Mappings": {"Tables": [{"Owner": 0, "MasterTableName": "USERS", "SlaveTableName": "users", "ColumnsMapped": [{"SlaveName": "id", "MasterName": "ID"}, {"SlaveName": "name", "MasterName": "NAME"}]}]}
  }'
```

2. Trigger sync:
```bash
curl -X POST http://localhost:5461/api/pairs/prod-sync/sync \
  -H "Content-Type: application/json" \
  -d '{"Owner": 0, "Table": "users", "All": false}'
```
