package importer

import (
	"fmt"
	"os"
)

type argParser interface {
	parse() (string, rune, string, error)
}

type parseArgs struct {
}

func newArgParser() argParser {
	return &parseArgs{}
}

func (*parseArgs) parse() (string, rune, string, error) {
	separator := ','
	argLen := len(os.Args)
	if argLen < 3 {
		return "", ',', "", fmt.Errorf("usage: csvimporter <filename> <tablename> \";\"\nFor help:\n csvimporter --help")
	}

	if argLen >= 4 {
		separator = rune(os.Args[3][0])
	}
	fileName := os.Args[1]
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", ',', "", fmt.Errorf("file %s does not exist", fileName)
	}

	tableName := os.Args[2]

	return fileName, separator, tableName, nil
}
