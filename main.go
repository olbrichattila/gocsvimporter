package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Analising CSV...")
	startTime := time.Now()

	fileName, tableName, separator, err := NewArgs().pharse()
	if err != nil {
		fmt.Println(err)
		return
	}

	app, err := NewApplication(fileName, tableName, separator)
	if err != nil {
		fmt.Println(err)
		return
	}

	analysisTime := time.Now()
	importer := NewImporter(app)
	err = importer.importCsv()
	if err != nil {
		fmt.Println(err)
		return
	}

	app.displayTimeStat(startTime, analysisTime)
}
