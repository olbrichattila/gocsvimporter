// Package importer is a CSV importer supporting large CSV files
package importer

import (
	"fmt"
)

// Import process the file
func Import() {
	err := newEnv(envFileName).loadEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	dBconfig, err := newDbConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	csvFileName, separator, tableName, err := newArgParser().parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Analising CSV...")

	app, err := newApplication(
		newImporter(
			dBconfig,
			newCsvReader(csvFileName, separator),
			newSQLGenerator(
				dBconfig,
				tableName,
			),
			newStorager(dBconfig),
		),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	phraseTime, importTime, totalTime, err := app.importer.importCsv()
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: tests stats, time may not be correctly displayed
	app.displayTimeStat(phraseTime, importTime, totalTime)
}
