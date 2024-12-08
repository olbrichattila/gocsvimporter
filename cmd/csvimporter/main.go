// Package main is running the CSV importer supporting large CSV files
package main

import (
	"fmt"

	importer "github.com/olbrichattila/gocsvimporter/internal"
	"github.com/olbrichattila/gocsvimporter/internal/arg"
	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

const (
	envFileName = ".env.csvimporter"
)

func main() {
	arg := arg.New()
	_, err := arg.Flag("help")
	if err == nil {
		displayHelp()
		return
	}

	err = arg.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	env := env.New(envFileName)
	err = env.LoadEnv()
	if err != nil {
		fmt.Println(err.Error())
	}

	dBConfig, err := database.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	importer.Import(env, dBConfig, arg)
}
