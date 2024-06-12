package main

import (
	"fmt"
	"math"
	"strings"
)

const (
	batchSize = 100
)

type importing struct {
	app        *application
	progress   int
	rowNr      int
	rowCount   int
	batch      [][]any
	batchIndex int
}

func NewImporter(app *application) *importing {
	return &importing{app: app}
}

func (i *importing) importCsv() error {
	// TODO: Refactor this block, use early returns, merge ifs and extract logics.
	defer func() {
		i.app.store.Close()
		i.app.csv.Close()
	}()

	headers := i.app.csv.Header()
	fmt.Printf("Found %d fields\n(%s)\n", len(headers), i.headerList(headers))
	err := i.app.store.Create(i.app.csv.Header())
	if err != nil {
		return err
	}

	err = i.app.store.StartTransaction()
	if err != nil {
		return err
	}
	defer i.app.store.CommitTransaction()

	i.resetProgress()
	connectionName := i.app.store.dBConfig.GetConnectionName()

	for i.app.csv.Next() {
		i.showProgress()
		// Firebird does not seem to support batch insert?
		if connectionName == "firebirdsql" {
			err = i.insertSQL()
			if err != nil {
				return err
			}
			continue
		}

		i.batchInsertSQL()
	}

	if len(i.batch) > 0 && connectionName != "firebirdsql" {
		err = i.app.store.BatchInsert(i.batch)
		if err != nil {
			i.app.store.RollbackTransaction()
			return err
		}
	}

	return nil
}

func (i *importing) headerList(headers CSVFields) string {
	fields := make([]string, len(headers))
	for i, f := range headers {
		fields[i] = f.Name
	}

	return strings.Join(fields, ", ")
}

func (i *importing) resetProgress() {
	i.rowNr = 0
	i.progress = 0
	i.rowCount = i.app.csv.RowCount()
	i.batchIndex = 0
}

func (i *importing) showProgress() {
	i.rowNr++
	percent := int(math.Ceil(float64(i.rowNr) / float64(i.rowCount) * 100))
	if percent != i.progress {
		i.progress = percent
		fmt.Printf("\rImporting: %d%%", i.progress)
	}
}

func (i *importing) insertSQL() error {
	err := i.app.store.Insert(i.app.csv.Row()...)
	if err != nil {
		i.app.store.RollbackTransaction()
		return err
	}
	return nil
}

func (i *importing) batchInsertSQL() error {
	i.batch = append(i.batch, i.app.csv.Row())
	if i.batchIndex == batchSize {
		i.batchIndex = 0
		err := i.app.store.BatchInsert(i.batch)
		if err != nil {
			i.app.store.RollbackTransaction()
			return err
		}
		i.batch = nil
		return nil
	}
	i.batchIndex++

	return nil
}
