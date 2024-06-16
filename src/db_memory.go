package importer

import (
	"database/sql"
	"fmt"
)

type memoryConfig struct {
	connection
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
	return c.connect(c.getConnectionName(), c.getConnectionString())
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
