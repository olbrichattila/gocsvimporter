package database

import (
	"database/sql"
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/nakagami/firebirdsql"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

type firebirdConfig struct {
	connection
}

func newFirebirdConfig() *firebirdConfig {
	return &firebirdConfig{}
}

func (c *firebirdConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s%s",
		os.Getenv(env.DbUserName),
		os.Getenv(env.DbPassword),
		os.Getenv(env.DbHost),
		os.Getenv(env.DbPort),
		os.Getenv(env.DbDatabase),
	)
}

func (c *firebirdConfig) GetConnectionName() string {
	return driverNameFirebird
}

func (c *firebirdConfig) GetFieldQuote() string {
	return "\""
}

func (c *firebirdConfig) GetBinding() string {
	return "?"
}

func (c *firebirdConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf(
		`EXECUTE BLOCK AS
			BEGIN
			IF (EXISTS (SELECT 1 FROM rdb$relations WHERE rdb$relation_name = '%s')) THEN
			BEGIN
				EXECUTE STATEMENT 'DROP TABLE %s%s%s';
			END
		END`,
		tableName,
		quote,
		tableName,
		quote,
	)
}

func (c *firebirdConfig) GetNewConnection() (*sql.DB, error) {
	return c.connect(c.GetConnectionName(), c.GetConnectionString())
}

func (c *firebirdConfig) HaveBatchInsert() bool {
	// Firebird does not support bach insert, setting it to true will break the import
	return false
}

func (c *firebirdConfig) HaveMultipleThreads() bool {
	return c.isOnByEnv(env.MultipleConnections, true)
}

func (c *firebirdConfig) NeedTransactions() bool {
	return c.isOnByEnv(env.Transactional, true)
}
