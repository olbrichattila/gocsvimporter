package importer

import "github.com/olbrichattila/gocsvimporter/internal/arg"

type duplicateMockParser struct {
}

func newDuplicateMockArgParser() arg.Parser {
	return &duplicateMockParser{}
}

func (*duplicateMockParser) Parse() (string, rune, string, error) {
	return "./fixtures/duplicate_testfile.csv", ',', testTableName, nil
}
