package importer

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/go-sql-driver/mysql"
)

type mySQLConfig struct {
	connection
}

func newMySQLConfig() *mySQLConfig {
	return &mySQLConfig{}
}

func (c *mySQLConfig) getConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv(envdbUserName),
		os.Getenv(envdbPassword),
		os.Getenv(envdbHost),
		os.Getenv(envdbPort),
		os.Getenv(envdbDatabase),
	)
}

func (c *mySQLConfig) getConnectionName() string {
	return driverNameMySQL
}

func (c *mySQLConfig) getFieldQuote() string {
	return "`"
}

func (c *mySQLConfig) getBinding() string {
	return "?"
}

func (c *mySQLConfig) getDropTableString(tableName string) string {
	quote := c.getFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *mySQLConfig) getNewConnection() (*sql.DB, error) {
	return c.connect(c.getConnectionName(), c.getConnectionString())
}

func (c *mySQLConfig) haveBatchInsert() bool {
	return true
}

func (c *mySQLConfig) haveMultipleThreads() bool {
	return true
}

func (c *mySQLConfig) needTransactions() bool {
	return true
}
