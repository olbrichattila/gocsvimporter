package importer

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/mattn/go-sqlite3"
)

func newSqliteConfig() *sqLiteConfig {
	return &sqLiteConfig{}
}

type sqLiteConfig struct {
	connection
}

func (c *sqLiteConfig) getConnectionString() string {
	return os.Getenv(envdbDatabase)
}

func (c *sqLiteConfig) getConnectionName() string {
	return driverNameSqLite
}

func (c *sqLiteConfig) getFieldQuote() string {
	return "\""
}

func (c *sqLiteConfig) getBinding() string {
	return "?"
}

func (c *sqLiteConfig) getDropTableString(tableName string) string {
	quote := c.getFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *sqLiteConfig) getNewConnection() (*sql.DB, error) {
	return c.connect(c.getConnectionName(), c.getConnectionString())
}

func (c *sqLiteConfig) haveBatchInsert() bool {
	c.isOnByEnv(envBatchInsert, true)
	return true
}

func (c *sqLiteConfig) haveMultipleThreads() bool {
	// this will cause a file lock not supported here
	return false
}

func (c *sqLiteConfig) needTransactions() bool {
	return c.isOnByEnv(envTransactional, true)
}
