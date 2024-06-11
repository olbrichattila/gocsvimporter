package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type mySqlConfig struct {
}

func newMySqlConfig() *mySqlConfig {
	return &mySqlConfig{}
}

func (c *mySqlConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
}

func (c *mySqlConfig) GetConnectionName() string {
	return "mysql"
}

func (c *mySqlConfig) GetFieldQuote() string {
	return "`"
}

func (c *mySqlConfig) GetBinding() string {
	return "?"
}

func (c *mySqlConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf("DROP TABLE IF EXISTS %s%s%s", quote, tableName, quote)
}
