package database

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/godror/godror"
	pq "github.com/lib/pq"
)

type DbType int

const (
	ORACLE DbType = iota
	POSTGRES
)

type Credentials struct {
	Database string
	Server   string
	Port     int

	Password string
	User     string

	DbType DbType
	db     *sql.DB
}

func (c *Credentials) createOracleDb() error {
	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`, c.User, c.Password, c.Server, c.Port, c.Database))
	if err != nil {
		return fmt.Errorf("Unable to stablish connection!\n%v", err)
	}
	c.db = db
	return nil
}

func (c *Credentials) createPostgreDb() error {
	cfg := pq.Config{
		Host: c.Server,
		Port: uint16(c.Port),
		Database: c.Database,
		Password: c.Password,
		User: c.User,
	}
	
	con, err := pq.NewConnectorConfig(cfg)
	if err != nil {
		return err
	}

	c.db = sql.OpenDB(con)
	return nil
}

func (c *Credentials) CreateOffsetStmt(skip int, take int) string {
	switch c.DbType {
	case POSTGRES:
		return fmt.Sprintf(
			"LIMIT %d OFFSET %d",
			take,
			skip,
		)
	case ORACLE:
		return fmt.Sprintf(
			"OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			take,
			skip,
		)
	}

	return ""
}

func (c *Credentials) GetDb() *sql.DB {
	return c.db
}

func (c *Credentials) connect() error {
	if c.db != nil {
		return nil
	}

	if (c.DbType == ORACLE) {
		return c.createOracleDb()
	} else if (c.DbType == POSTGRES) {
		return c.createPostgreDb()
	}

	return fmt.Errorf("Database type not supported")
}

func (c *Credentials) Ping() error {
	if c.db == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}

	return c.db.Ping()
}

func (c *Credentials) Placeholder(paramIdx int) string {
	switch (c.DbType) {
		case ORACLE: return fmt.Sprintf(":%d", paramIdx)
		case POSTGRES: return fmt.Sprintf("$%d", paramIdx)
	}
	return "?"
}

func CreateCredFromGin(c *gin.Context) (*Credentials, error) {	
	cred := Credentials{}
	if err := c.ShouldBind(&cred); err != nil {
		return nil, fmt.Errorf("Unable to bind 'Credential': %v", err)
	}

	if err := cred.Ping(); err != nil {
		return nil, err
	}

	return &cred, nil
}

func (c *Credentials) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

