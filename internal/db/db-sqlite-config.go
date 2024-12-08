package database

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

func newSqliteConfig() *sqLiteConfig {
	return &sqLiteConfig{}
}

type sqLiteConfig struct {
	connection
}

func (c *sqLiteConfig) GetConnectionString() string {
	return os.Getenv(env.DbDatabase)
}

func (c *sqLiteConfig) GetConnectionName() string {
	return driverNameSqLite
}

func (c *sqLiteConfig) GetFieldQuote() string {
	return "\""
}

func (c *sqLiteConfig) GetBinding() string {
	return "?"
}

func (c *sqLiteConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *sqLiteConfig) GetNewConnection() (*sql.DB, error) {
	return c.connect(c.GetConnectionName(), c.GetConnectionString())
}

func (c *sqLiteConfig) HaveBatchInsert() bool {
	c.isOnByEnv(env.BatchInsert, true)
	return true
}

func (c *sqLiteConfig) HaveMultipleThreads() bool {
	// this will cause a file lock not supported here
	return false
}

func (c *sqLiteConfig) NeedTransactions() bool {
	return c.isOnByEnv(env.Transactional, true)
}
