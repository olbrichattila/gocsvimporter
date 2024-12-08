package importer

import "github.com/olbrichattila/gocsvimporter/internal/arg"

type duplicateMockParser struct {
}

func newDuplicateMockArgParser() arg.Parser {
	return &duplicateMockParser{}
}

func (*duplicateMockParser) Validate() error {
	return nil
}

func (*duplicateMockParser) FileName() string {
	return "./fixtures/duplicate_testfile.csv"
}

func (*duplicateMockParser) TableName() string {
	return testTableName
}

func (*duplicateMockParser) Separator() rune {
	return ','
}
