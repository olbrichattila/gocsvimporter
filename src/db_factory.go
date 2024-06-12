package importer

import (
	"fmt"
	"os"
)

const (
	envdbConnection        = "DB_CONNECTION"
	envdbUserName          = "DB_USERNAME"
	envdbPassword          = "DB_PASSWORD"
	envdbHost              = "DB_HOST"
	envdbPort              = "DB_PORT"
	envdbDatabase          = "DB_DATABASE"
	envdbSSLMode           = "DB_SSLMODE"
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
