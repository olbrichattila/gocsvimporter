package main

import (
	"fmt"
	"time"
)

type application struct {
	csv   csvReader
	store *dataStore
}

func NewApplication(csvFile, tableName string, csvSeparator rune) (*application, error) {
	env := NewEnv()
	err := env.LoadEnv()
	if err != nil {
		return nil, err
	}

	app := &application{}
	app.csv, err = NewCsvReader(csvFile, csvSeparator)
	if err != nil {
		return nil, err
	}

	app.store, err = NewDataStore(tableName)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *application) displayTimeStat(startTime, analysisTime time.Time) {
	finisedTime := time.Now()
	fullAnalysisTime := a.durasionAsString(analysisTime.Sub(startTime).Seconds())
	fullDurationTime := a.durasionAsString(finisedTime.Sub(analysisTime).Seconds())
	totalTime := a.durasionAsString(finisedTime.Sub(startTime).Seconds())

	fmt.Printf("\nDone\nFull Analysis time: %s\nFull duration time: %s\nTotal: %s\n", fullAnalysisTime, fullDurationTime, totalTime)
}

func (application) durasionAsString(elapsed float64) string {
	return fmt.Sprintf("%.0f minutes %d seconds", elapsed/60, int64(elapsed)%60)
}
