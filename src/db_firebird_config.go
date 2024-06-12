package importer

import (
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/nakagami/firebirdsql"
)

type firebirdConfig struct {
}

func newFirebirdConfig() *firebirdConfig {
	return &firebirdConfig{}
}

func (c *firebirdConfig) getConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s%s",
		os.Getenv(envdbUserName),
		os.Getenv(envdbPassword),
		os.Getenv(envdbHost),
		os.Getenv(envdbPort),
		os.Getenv(envdbDatabase),
	)
}

func (c *firebirdConfig) getConnectionName() string {
	return driverNameFirebird
}

func (c *firebirdConfig) getFieldQuote() string {
	return "\""
}

func (c *firebirdConfig) getBinding() string {
	return "?"
}

func (c *firebirdConfig) getDropTableString(tableName string) string {
	quote := c.getFieldQuote()
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
