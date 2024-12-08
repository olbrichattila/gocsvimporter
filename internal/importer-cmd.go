// Package importer is a CSV importer supporting large CSV files
package importer

import (
	"fmt"
	"math"

	"github.com/olbrichattila/gocsvimporter/internal/arg"
	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/olbrichattila/gocsvimporter/internal/env"
	"github.com/olbrichattila/gocsvimporter/internal/storage"
)

// Import process the file
func Import(env env.Enver, dbConfig database.DBConfiger, argParser arg.Parser) {
	fmt.Println("Analyzing CSV...")

	importer := newImporter(
		dbConfig,
		newCsvReader(argParser),
		newSQLGenerator(
			dbConfig,
			argParser,
		),
		storage.New(dbConfig),
	)

	phraseTime, importTime, totalTime, err := importer.importCsv()
	if err != nil {
		fmt.Println(err)
		return
	}

	displayTimeStat(phraseTime, importTime, totalTime)
}

func displayTimeStat(analysisTime, importTime, totalTime float64) {
	fmt.Printf(
		"\n\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n",
		durationAsString(analysisTime),
		durationAsString(importTime),
		durationAsString(totalTime),
	)
}

func durationAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", math.Floor(elapsed/60), int64(elapsed)%60)
}
