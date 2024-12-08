package database

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

type mySQLConfig struct {
	connection
}

func newMySQLConfig() *mySQLConfig {
	return &mySQLConfig{}
}

func (c *mySQLConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv(env.DbUserName),
		os.Getenv(env.DbPassword),
		os.Getenv(env.DbHost),
		os.Getenv(env.DbPort),
		os.Getenv(env.DbDatabase),
	)
}

func (c *mySQLConfig) GetConnectionName() string {
	return driverNameMySQL
}

func (c *mySQLConfig) GetFieldQuote() string {
	return "`"
}

func (c *mySQLConfig) GetBinding() string {
	return "?"
}

func (c *mySQLConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *mySQLConfig) GetNewConnection() (*sql.DB, error) {
	return c.connect(c.GetConnectionName(), c.GetConnectionString())
}

func (c *mySQLConfig) HaveBatchInsert() bool {
	return c.isOnByEnv(env.BatchInsert, true)
}

func (c *mySQLConfig) HaveMultipleThreads() bool {
	return c.isOnByEnv(env.MultipleConnections, true)
}

func (c *mySQLConfig) NeedTransactions() bool {
	return c.isOnByEnv(env.Transactional, true)
}
