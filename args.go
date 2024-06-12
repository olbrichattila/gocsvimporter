package main

import (
	"fmt"
	"os"
)

type parseArgs struct {
}

func NewArgs() *parseArgs {
	return &parseArgs{}
}

func (*parseArgs) pharse() (string, string, rune, error) {
	separator := ','
	argLen := len(os.Args)
	if argLen < 3 {
		return "", "", ',', fmt.Errorf("usage: go run main.go <filename> <tablename> \";\"")
	}

	if argLen >= 4 {
		separator = rune(os.Args[3][0])
	}
	fileName := os.Args[1]
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return "", "", ',', fmt.Errorf("file %s does not exist", fileName)
	}

	tableName := os.Args[2]

	return fileName, tableName, separator, nil
}
