package importer

import (
	"fmt"
	"os"
)

const (
	defaultDropTableFormat = "DROP TABLE IF EXISTS %s%s%s"
	driverNameFirebird     = "firebirdsql"
	driverNameSqLite       = "sqlite3"
	driverNameMySQL        = "mysql"
	driverNamePostgres     = "postgres"

	dbConnectionTypeSqLite   = "sqlite"
	dbConnectionTypeMySQL    = "mysql"
	dbConnectionTypePgSQL    = "pgsql"
	dbConnectionTypeFirebird = "firebird"
	dbConnectionTypeMemory   = "memory"
)

func newDbConfig() (dBConfiger, error) {
	conf := &dBconf{}
	err := conf.init()
	return conf.config, err
}

type dBconf struct {
	config dBConfiger
}

func (c *dBconf) init() error {
	dbConnection := os.Getenv(envdbConnection)

	switch dbConnection {
	case dbConnectionTypeSqLite:
		c.config = newSqliteConfig()
		return nil
	case dbConnectionTypeMySQL:
		c.config = newMySQLConfig()
		return nil
	case dbConnectionTypePgSQL:
		c.config = newPgsqlConfig()
		return nil
	case dbConnectionTypeFirebird:
		c.config = newFirebirdConfig()
		return nil
	case dbConnectionTypeMemory:
		c.config = newMemoryDBConfig()
		return nil
	default:
		return fmt.Errorf("invalid DB_CONNECTION %s", dbConnection)
	}
}
