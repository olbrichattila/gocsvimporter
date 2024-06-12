package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	batchSize = 100
)

type application struct {
	csv   *readCsv
	store *dataStore
}

func main() {
	// TODO: Refactoring, simpify, extract
	err := loadEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName, tableName, separator, err := pharseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Analising CSV...")
	startTime := time.Now()

	app, err := NewApplication(fileName, "./database/database.sqlite", tableName, separator)
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

	finisedTime := time.Now()
	fullAnalysisTime := durasionAsString(analysisTime.Sub(startTime).Seconds())
	fullDurationTime := durasionAsString(finisedTime.Sub(analysisTime).Seconds())
	totalTime := durasionAsString(finisedTime.Sub(startTime).Seconds())

	fmt.Printf("\nDone\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n", fullAnalysisTime, fullDurationTime, totalTime)
}

func durasionAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", elapsed/60, int64(elapsed)%60)
}

func NewApplication(csvFile, dbFile, tableName string, csvSeparator rune) (*application, error) {
	app := &application{}
	var err error
	app.csv, err = NewCsvReader(csvFile, csvSeparator)
	if err != nil {
		return nil, err
	}

	app.store, err = NewDataStore(dbFile, tableName)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func pharseArgs() (string, string, rune, error) {
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func loadEnv() error {
	if fileExists("./.env") {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	return nil
}
