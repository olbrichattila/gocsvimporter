package main

import (
	"fmt"
	"os"
)

func GetDbConnector() (DBConfiger, error) {
	dbConnection := os.Getenv("DB_CONNECTION")

	switch dbConnection {
	case "sqlite":
		return newSqliteConfig(), nil
	case "mysql":
		return newMySqlConfig(), nil
	case "pgsql":
		return newPgsqlConfig(), nil
	case "firebird":
		return newFirebirdConfig(), nil
	default:
		return nil, fmt.Errorf("invalid DB_CONNECTION %s", dbConnection)
	}
}
