// Package main is a CSV importer supporting large CSV files
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Analising CSV...")
	startTime := time.Now()

	app, err := newApplication(
		newArgParser(),
		newEnv(),
		newImporter(
			newDataStore(),
			newCsvReader(),
		),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	analysisTime := time.Now()
	err = app.importer.importCsv()
	if err != nil {
		fmt.Println(err)
		return
	}

	app.displayTimeStat(startTime, analysisTime)
}
