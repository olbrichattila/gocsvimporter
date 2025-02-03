package database

import (
	"database/sql"
	"os"
	"strings"
)

type connection struct {
}

func (c *connection) connect(connectionName, connectionString string) (*sql.DB, error) {
	db, err := sql.Open(connectionName, connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *connection) isOnByEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		lowerCaseConfigValue := strings.ToLower(value)

		if lowerCaseConfigValue == "on" {
			return true
		}

		if lowerCaseConfigValue == "off" {
			return false
		}
	}

	return defaultValue
}
