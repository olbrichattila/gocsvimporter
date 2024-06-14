// Package importer is a CSV importer supporting large CSV files
package importer

import (
	"fmt"
	"time"
)

// Import process the file
func Import() {
	fmt.Println("Analising CSV...")
	startTime := time.Now()

	app, err := newApplication(
		newArgParser(),
		newEnv(envFileName),
		newDbConnector(),
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

	// TODO: tests stats, time may not be correctly dispayed
	app.displayTimeStat(startTime, analysisTime)
}
