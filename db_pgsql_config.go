package main

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type pgsqlConfig struct {
}

func newPgsqlConfig() *pgsqlConfig {
	return &pgsqlConfig{}
}

func (c *pgsqlConfig) GetConnectionString() string {

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_SSLMODE"),
	)
}

func (c *pgsqlConfig) GetConnectionName() string {
	return "postgres"
}

func (c *pgsqlConfig) GetFieldQuote() string {
	return "\""
}

func (c *pgsqlConfig) GetBinding() string {
	return "$"
}

func (c *pgsqlConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf("DROP TABLE IF EXISTS %s%s%s", quote, tableName, quote)
}
