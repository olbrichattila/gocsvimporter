package importer

type dBConfiger interface {
	getConnectionName() string
	getConnectionString() string
	getFieldQuote() string
	getBinding() string
	getDropTableString(string) string
}
