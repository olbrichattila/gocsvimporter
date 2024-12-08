package database

import (
	"fmt"
	"os"

	"github.com/olbrichattila/gocsvimporter/internal/env"
)

// Database connection types
const (
	defaultDropTableFormat = "DROP TABLE IF EXISTS %s%s%s"
	driverNameFirebird     = "firebirdsql"
	driverNameSqLite       = "sqlite3"
	driverNameMySQL        = "mysql"
	driverNamePostgres     = "postgres"

	DbConnectionTypeSqLite   = "sqlite"
	DbConnectionTypeMySQL    = "mysql"
	DbConnectionTypePgSQL    = "pgsql"
	DbConnectionTypeFirebird = "firebird"
	DbConnectionTypeMemory   = "memory"
)

// New creates a new database configuration manager
func New() (DBConfiger, error) {
	conf := &dBConf{}
	err := conf.init()
	return conf.config, err
}

type dBConf struct {
	config DBConfiger
}

func (c *dBConf) init() error {
	dbConnection := os.Getenv(env.DbConnection)

	switch dbConnection {
	case DbConnectionTypeSqLite:
		c.config = newSqliteConfig()
		return nil
	case DbConnectionTypeMySQL:
		c.config = newMySQLConfig()
		return nil
	case DbConnectionTypePgSQL:
		c.config = newPgsqlConfig()
		return nil
	case DbConnectionTypeFirebird:
		c.config = newFirebirdConfig()
		return nil
	case DbConnectionTypeMemory:
		c.config = newMemoryDBConfig()
		return nil
	default:
		return fmt.Errorf("invalid DB_CONNECTION `%s`", dbConnection)
	}
}
