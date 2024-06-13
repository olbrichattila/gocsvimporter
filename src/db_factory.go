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
)

func getDbConnector() (dBConfiger, error) {
	dbConnection := os.Getenv(envdbConnection)

	switch dbConnection {
	case "sqlite":
		return newSqliteConfig(), nil
	case "mysql":
		return newMySQLConfig(), nil
	case "pgsql":
		return newPgsqlConfig(), nil
	case "firebird":
		return newFirebirdConfig(), nil
	default:
		return nil, fmt.Errorf("invalid DB_CONNECTION %s", dbConnection)
	}
}
