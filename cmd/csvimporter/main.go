// Package main is running the CSV importer supporting large CSV files
package main

import (
	"flag"
	"fmt"
	"os"

	importer "github.com/olbrichattila/gocsvimporter/internal"
	"github.com/olbrichattila/gocsvimporter/internal/arg"
	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

const (
	envFileName = ".env.csvimporter"
)

func main() {
	help := flag.Bool("help", false, "Display help")

	// Parse the flags
	flag.Parse()

	// If help flag is set, display usage and exit
	if *help {
		displayHelp()
		os.Exit(0)
	}

	env := env.New(envFileName)
	err := env.LoadEnv()
	if err != nil {
		fmt.Println(err.Error())
	}

	dBConfig, err := database.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	arg := arg.New()

	err = arg.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	importer.Import(env, dBConfig, arg)
}
