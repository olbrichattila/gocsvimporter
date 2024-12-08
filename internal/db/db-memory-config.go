package database

import (
	"database/sql"
	"fmt"

	"github.com/olbrichattila/gocsvimporter/internal/env"
)

type memoryConfig struct {
	connection
}

func newMemoryDBConfig() *memoryConfig {
	return &memoryConfig{}
}

func (c *memoryConfig) GetConnectionString() string {
	return ":memory:"
}

func (c *memoryConfig) GetConnectionName() string {
	return driverNameSqLite
}

func (c *memoryConfig) GetFieldQuote() string {
	return "\""
}

func (c *memoryConfig) GetBinding() string {
	return "?"
}

func (c *memoryConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *memoryConfig) GetNewConnection() (*sql.DB, error) {
	return c.connect(c.GetConnectionName(), c.GetConnectionString())
}

func (c *memoryConfig) HaveBatchInsert() bool {
	return c.isOnByEnv(env.BatchInsert, true)
}

func (c *memoryConfig) HaveMultipleThreads() bool {
	// not supported
	return false
}

func (c *memoryConfig) NeedTransactions() bool {
	return c.isOnByEnv(env.Transactional, true)
}
