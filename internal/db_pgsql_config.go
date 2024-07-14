package importer

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/lib/pq"
)

type pgsqlConfig struct {
	connection
}

func newPgsqlConfig() *pgsqlConfig {
	return &pgsqlConfig{}
}

func (c *pgsqlConfig) getConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv(envdbUserName),
		os.Getenv(envdbPassword),
		os.Getenv(envdbHost),
		os.Getenv(envdbPort),
		os.Getenv(envdbDatabase),
		os.Getenv(envdbSSLMode),
	)
}

func (c *pgsqlConfig) getConnectionName() string {
	return driverNamePostgres
}

func (c *pgsqlConfig) getFieldQuote() string {
	return "\""
}

func (c *pgsqlConfig) getBinding() string {
	return "$"
}

func (c *pgsqlConfig) getDropTableString(tableName string) string {
	quote := c.getFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}

func (c *pgsqlConfig) getNewConnection() (*sql.DB, error) {
	return c.connect(c.getConnectionName(), c.getConnectionString())
}

func (c *pgsqlConfig) haveBatchInsert() bool {
	return c.isOnByEnv(envBatchInsert, true)
}

func (c *pgsqlConfig) haveMultipleThreads() bool {
	return c.isOnByEnv(envMultipleConnections, true)
}

func (c *pgsqlConfig) needTransactions() bool {
	return c.isOnByEnv(envTransactional, true)
}
