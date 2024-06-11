package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type sqLiteConfig struct {
}

func newSqliteConfig() *sqLiteConfig {
	return &sqLiteConfig{}
}

func (c *sqLiteConfig) GetConnectionString() string {
	return os.Getenv("DB_DATABASE")
}

func (c *sqLiteConfig) GetConnectionName() string {
	return "sqlite3"
}

func (c *sqLiteConfig) GetFieldQuote() string {
	return "\""
}

func (c *sqLiteConfig) GetBinding() string {
	return "?"
}

func (c *sqLiteConfig) GetDropTableString(tableName string) string {
	quote := c.GetFieldQuote()
	return fmt.Sprintf("DROP TABLE IF EXISTS %s%s%s", quote, tableName, quote)
}
