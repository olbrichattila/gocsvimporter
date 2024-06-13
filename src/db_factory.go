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

func getDbConnector() (dBConfiger, error) {
	dbConnection := os.Getenv(envdbConnection)

	switch dbConnection {
	case dbConnectionTypeSqLite:
		return newSqliteConfig(), nil
	case dbConnectionTypeMySQL:
		return newMySQLConfig(), nil
	case dbConnectionTypePgSQL:
		return newPgsqlConfig(), nil
	case dbConnectionTypeFirebird:
		return newFirebirdConfig(), nil
	case dbConnectionTypeMemory:
		return newMemoryDBConfig(), nil
	default:
		return nil, fmt.Errorf("invalid DB_CONNECTION %s", dbConnection)
	}
}
