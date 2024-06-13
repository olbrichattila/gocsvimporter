package importer

type mockParser struct {
}

func newMockParser() *mockParser {
	return &mockParser{}
}

func (*mockParser) pharse() (string, string, rune, error) {
	return "./fixtures/testfile.csv", "test_table", ',', nil
}
