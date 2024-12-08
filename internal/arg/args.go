// Package Arg parser parses command line arguments
package arg

import (
	"fmt"
	"os"
)

// Parser is the interface to abstract argument parser
type Parser interface {
	Parse() (string, rune, string, error)
}

type parseArgs struct {
}

// New creates a new parser
func New() Parser {
	return &parseArgs{}
}

func (*parseArgs) Parse() (string, rune, string, error) {
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
