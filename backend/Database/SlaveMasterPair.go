package database

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type TableOwner int

const (
	MASTER TableOwner = iota
	SLAVE	
)

type ColumnMapping struct {
	SlaveName string
	MasterName string
}

type TableMapping struct {
	Owner TableOwner
	MasterTableName string
	SlaveTableName string
	LastSync time.Time
	ColumnsMapped []ColumnMapping
}

func (tm *TableMapping) GetSourceTableName() string {
	switch tm.Owner {
	case SLAVE:
		return tm.SlaveTableName
	case MASTER:
		return tm.MasterTableName
	}

	return ""
}

func (tm *TableMapping) GetTargetTableName() string {
	switch tm.Owner {
	case SLAVE:
		return tm.MasterTableName
	case MASTER:
		return tm.SlaveTableName
	}

	return ""
}

type Mapping struct {
	Tables []TableMapping
}

func (m *Mapping) GetTableNamesFromOwner(owner TableOwner) []string {
	names := []string{}
	for _, tm := range m.Tables {
		if tm.Owner != owner {
			continue
		}
		names = append(names, tm.GetSourceTableName())
	}

	return names
}

func (m *Mapping) ContainsTable(table string, owner TableOwner) *TableMapping{
	for _, tn := range m.Tables {
		if tn.Owner != owner {
			continue
		}

		if tn.GetSourceTableName() == table {
			return &tn
		}
	}

	return nil
}

type SlaveMasterPair struct {
	Name string
	Slave Credentials
	Master Credentials
	Mappings Mapping
}

func CreatePairFromGin(c *gin.Context) (*SlaveMasterPair, error) {
	pair := SlaveMasterPair{}
	if err := c.ShouldBind(&pair); err != nil {
		return nil, fmt.Errorf("Unable to bind 'SlaveMasterPair': %v", err)
	}

	if err := pair.Ping(); err != nil {
		return nil, err
	}

	return &pair, nil
}

func (p *SlaveMasterPair) Ping() error {
	if err := p.Master.Ping(); err != nil {
		return err;
	}

	if err := p.Slave.Ping(); err != nil {
		return err;
	}

	return nil
}

func (p *SlaveMasterPair) Close() error {
	var errStr string = ""
	if err := p.Master.Close(); err != nil {
		errStr += fmt.Sprintf("Failed to close master: %v\n", err)
	}

	if err := p.Slave.Close(); err != nil {
		errStr += fmt.Sprintf("Failed to close slave: %v\n", err)
	}

	if errStr == "" {
		return nil
	}

	return fmt.Errorf(errStr)
}
