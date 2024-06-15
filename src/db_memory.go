package importer

import (
	"database/sql"
	"fmt"
)

type memoryConfig struct {
}

func newMemoryDBConfig() *memoryConfig {
	return &memoryConfig{}
}

func (c *memoryConfig) getConnectionString() string {
	return ":memory:"
}

func (c *memoryConfig) getConnectionName() string {
	return driverNameSqLite
}

func (c *memoryConfig) getFieldQuote() string {
	return "\""
}

func (c *memoryConfig) getBinding() string {
	return "?"
}

func (c *memoryConfig) getDropTableString(tableName string) string {
	quote := c.getFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *memoryConfig) getNewConnection() (*sql.DB, error) {
	db, err := sql.Open(c.getConnectionName(), c.getConnectionString())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *memoryConfig) haveBatchInsert() bool {
	return true
}

func (c *memoryConfig) haveMultipleThreads() bool {
	return false
}

func (c *memoryConfig) needTransactions() bool {
	return true
}
