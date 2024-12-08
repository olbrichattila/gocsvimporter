package importer

import "github.com/olbrichattila/gocsvimporter/internal/arg"

const testTableName = "test_table"

type mockParser struct {
}

func newMockArgParser() arg.Parser {
	return &mockParser{}
}

func (*mockParser) Parse() (string, rune, string, error) {
	return "./fixtures/testfile.csv", ',', testTableName, nil
}
