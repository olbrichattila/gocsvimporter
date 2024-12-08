package importer

const testTableName = "test_table"

type mockParser struct {
}

func newMockParser() argParser {
	return &mockParser{}
}

func (*mockParser) parse() (string, rune, string, error) {
	return "./fixtures/testfile.csv", ',', testTableName, nil
}
