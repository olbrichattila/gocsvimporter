package importer

type mockParser struct {
}

func newMockParser() argParser {
	return &mockParser{}
}

func (*mockParser) parse() (string, rune, string, error) {
	return "./fixtures/testfile.csv", ',', "test_table", nil
}
