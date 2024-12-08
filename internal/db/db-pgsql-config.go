package database

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/lib/pq"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

type pgsqlConfig struct {
	connection
}

func newPgsqlConfig() *pgsqlConfig {
	return &pgsqlConfig{}
}

func (c *pgsqlConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv(env.DbUserName),
		os.Getenv(env.DbPassword),
		os.Getenv(env.DbHost),
		os.Getenv(env.DbPort),
		os.Getenv(env.DbDatabase),
		os.Getenv(env.DbSSLMode),
	)
}

func (c *pgsqlConfig) GetConnectionName() string {
	return driverNamePostgres
}

func (c *pgsqlConfig) GetFieldQuote() string {
	return "\""
}

func (c *pgsqlConfig) GetBinding() string {
	return "$"
}

func (c *pgsqlConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *pgsqlConfig) GetNewConnection() (*sql.DB, error) {
	return c.connect(c.GetConnectionName(), c.GetConnectionString())
}

func (c *pgsqlConfig) HaveBatchInsert() bool {
	return c.isOnByEnv(env.BatchInsert, true)
}

func (c *pgsqlConfig) HaveMultipleThreads() bool {
	return c.isOnByEnv(env.MultipleConnections, true)
}

func (c *pgsqlConfig) NeedTransactions() bool {
	return c.isOnByEnv(env.Transactional, true)
}
