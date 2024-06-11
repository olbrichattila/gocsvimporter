package main

import (
	"fmt"
	"os"

	_ "github.com/nakagami/firebirdsql"
)

type firebirdConfig struct {
}

func newFirebirdConfig() *firebirdConfig {
	return &firebirdConfig{}
}

func (c *firebirdConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
}

func (c *firebirdConfig) GetConnectionName() string {
	return "firebirdsql"
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
