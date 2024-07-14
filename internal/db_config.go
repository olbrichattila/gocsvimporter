package importer

import "database/sql"

type dBConfiger interface {
	getConnectionName() string
	getConnectionString() string
	getFieldQuote() string
	getBinding() string
	getDropTableString(string) string
	getNewConnection() (*sql.DB, error)
	haveBatchInsert() bool
	haveMultipleThreads() bool
	needTransactions() bool
}
