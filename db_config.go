package main

type DBConfiger interface {
	GetConnectionName() string
	GetConnectionString() string
	GetFieldQuote() string
	GetBinding() string
	GetDropTableString(string) string
}
