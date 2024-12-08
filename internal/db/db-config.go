// Package database encapsulates database logic
package database

import "database/sql"

// DBConfiger contains database management functions
type DBConfiger interface {
	GetConnectionName() string
	GetConnectionString() string
	GetFieldQuote() string
	GetBinding() string
	GetDropTableString(string) string
	GetNewConnection() (*sql.DB, error)
	HaveBatchInsert() bool
	HaveMultipleThreads() bool
	NeedTransactions() bool
}
