// Package arg is an argument parser working with command line arguments
package arg

import (
	"fmt"
	"os"
	"strings"
)

// Parser is the interface to abstract argument parser
type Parser interface {
	Validate() error
	FileName() string
	TableName() string
	Separator() rune
	Flag(name string) (string, error)
}

type parseArgs struct {
}

// New creates a new parser
func New() Parser {
	return &parseArgs{}
}

// Validate validates if all parameters are provided, and file exists, separator is 1 char long
func (t *parseArgs) Validate() error {
	args := t.getParams()
	argLen := len(args)
	if argLen < 2 {
		return fmt.Errorf("usage: csvimporter <filename> <tablename> \";\"\nFor help:\n csvimporter --help")
	}

	if argLen >= 3 && len(args[2]) != 1 {
		return fmt.Errorf("separator must be 1 character long")
	}

	fileName := args[0]
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", fileName)
	}

	return nil
}

func (t *parseArgs) FileName() string {
	args := t.getParams()
	return args[0]
}

func (t *parseArgs) TableName() string {
	args := t.getParams()
	return args[1]
}

func (t *parseArgs) Separator() rune {
	args := t.getParams()
	if len(args) < 3 {
		return ','
	}

	return rune(args[2][0])
}

func (t *parseArgs) Flag(name string) (string, error) {
	flags := t.getFlags()
	for key, value := range flags {
		if key == name {
			return value, nil
		}
	}

	return "", fmt.Errorf("flag " + name + " not provided.")
}

func (*parseArgs) getParams() []string {
	result := make([]string, 0)

	for i, arg := range os.Args {
		if i > 0 {
			if !strings.HasPrefix(arg, "-") {
				result = append(result, arg)
			}
		}
	}

	return result
}

func (*parseArgs) getFlags() map[string]string {
	result := make(map[string]string, 0)

	for i, arg := range os.Args {
		if i > 0 {
			if strings.HasPrefix(arg, "-") {
				value := strings.TrimPrefix(arg, "-")
				parts := strings.Split(value, "=")
				key := strings.TrimSpace(parts[0])
				param := ""
				if len(parts) > 1 {
					param = parts[1]
				}
				result[key] = param
			}
		}
	}

	return result
}
