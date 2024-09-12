package importer

type duplicateMockParser struct {
}

func newDuplicateMockParser() argParser {
	return &duplicateMockParser{}
}

func (*duplicateMockParser) parse() (string, rune, string, error) {
	return "./fixtures/duplicate_testfile.csv", ',', "test_table", nil
}
