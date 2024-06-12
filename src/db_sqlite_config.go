package importer

import (
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/mattn/go-sqlite3"
)

type sqLiteConfig struct {
}

func newSqliteConfig() *sqLiteConfig {
	return &sqLiteConfig{}
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
