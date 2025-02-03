// Package main runs the CSV importer, supporting large CSV files
package main

import (
	"fmt"
	"log"

	importer "github.com/olbrichattila/gocsvimporter/internal"
	"github.com/olbrichattila/gocsvimporter/internal/arg"
	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

const envFileName = ".env.csvimporter"

func main() {
	// Parse command-line arguments
	args := parseArgs()

	// Load environment configuration
	environment := loadEnvConfig()

	// Set up database connection
	dbConfig := setupDatabase()

	// Perform CSV import
	importCSV(environment, dbConfig, args)
}

func parseArgs() arg.Parser {
	args := arg.New()
	if err := displayHelpIfNeeded(args); err != nil {
		log.Fatal(err)
	}
	if err := args.Validate(); err != nil {
		log.Fatal(err)
	}
	return args
}

func displayHelpIfNeeded(args arg.Parser) error {
	_, err := args.Flag("help")
	if err == nil {
		displayHelp()
		return fmt.Errorf("help requested")
	}
	return nil
}

func loadEnvConfig() env.Enver {
	envConfig := env.New(envFileName)
	if err := envConfig.LoadEnv(); err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}
	return envConfig
}

func setupDatabase() database.DBConfiger {
	dbConfig, err := database.New()
	if err != nil {
		log.Fatal("Error setting up database: ", err)
	}
	return dbConfig
}

func importCSV(environment env.Enver, dbConfig database.DBConfiger, args arg.Parser) {
	importer.Import(environment, dbConfig, args)
}
