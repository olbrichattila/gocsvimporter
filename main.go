package main

import (
	"fmt"
	"math"
	"os"
	"strings"
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
	err := loadEnv()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fileName, tableName, separator, err := pharseArgs()
	if err != nil {
		fmt.Println(err.Error())
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
	importer(app)

	finisedTime := time.Now()

	fullAnalysisTime := durasionAsString(analysisTime.Sub(startTime).Seconds())
	fullDurationTime := durasionAsString(finisedTime.Sub(analysisTime).Seconds())
	totalTime := durasionAsString(finisedTime.Sub(startTime).Seconds())

	fmt.Printf("\nDone\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n", fullAnalysisTime, fullDurationTime, totalTime)
}

func durasionAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", elapsed/60, int64(elapsed)%60)
}

func importer(app *application) {
	defer func() {
		app.store.Close()
		app.csv.Close()
	}()

	headers := app.csv.Header()
	fmt.Printf("Found %d fields\n(%s)\n", len(headers), headerList(headers))
	err := app.store.Create(app.csv.Header())
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.store.StartTransaction()
	if err != nil {
		fmt.Println(err)
	}
	defer app.store.CommitTransaction()

	i := 0
	progress := 0
	rowCount := app.csv.RowCount()
	batchIndex := 0
	var batch [][]any
	connectionName := app.store.dBConfig.GetConnectionName()
	// TODO: Refactor this block, use early returns, merge ifs and extract logics.
	for app.csv.Next() {
		i++
		percent := int(math.Ceil(float64(i) / float64(rowCount) * 100))
		if percent != progress {
			progress = percent
			fmt.Printf("\rImporting: %d%%", progress)
		}

		// Firebird does not seem to support batch insert?
		if connectionName == "firebirdsql" {
			err = app.store.Insert(app.csv.Row()...)
			if err != nil {
				fmt.Println(err)
				app.store.RollbackTransaction()
				break
			}
		} else {
			batch = append(batch, app.csv.Row())
			if batchIndex == batchSize {
				batchIndex = 0
				err = app.store.BatchInsert(batch)
				if err != nil {
					fmt.Println(err)
					app.store.RollbackTransaction()
					break
				}
				batch = nil
			} else {
				batchIndex++
			}
		}
	}

	if len(batch) > 0 && connectionName != "firebirdsql" {
		err = app.store.BatchInsert(batch)
		if err != nil {
			fmt.Println(err)
			app.store.RollbackTransaction()
		}
	}
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

func headerList(headers CSVFields) string {
	fields := make([]string, len(headers))
	for i, f := range headers {
		fields[i] = f.Name
	}

	return strings.Join(fields, ", ")
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
