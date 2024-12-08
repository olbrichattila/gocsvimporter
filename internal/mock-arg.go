package importer

import "github.com/olbrichattila/gocsvimporter/internal/arg"

const testTableName = "test_table"

type mockParser struct {
}

func newMockArgParser() arg.Parser {
	return &mockParser{}
}

func (*mockParser) Validate() error {
	return nil
}

func (*mockParser) FileName() string {
	return "./fixtures/testfile.csv"
}

func (*mockParser) TableName() string {
	return testTableName
}

func (*mockParser) Separator() rune {
	return ','
}
