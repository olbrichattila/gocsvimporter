package main

import (
	"fmt"
	"math"
	"strings"
)

const (
	batchSize = 100
)

type importer interface {
	importCsv() error
	getCsvReader() csvReader
	getStorer() dataStorer
}

type importing struct {
	storer     dataStorer
	csv        csvReader
	progress   int
	rowNr      int
	rowCount   int
	batch      [][]any
	batchIndex int
}

func newImporter(storer dataStorer, csv csvReader) *importing {
	return &importing{
		storer: storer,
		csv:    csv,
	}
}

func (i *importing) getCsvReader() csvReader {
	return i.csv
}

func (i *importing) getStorer() dataStorer {
	return i.storer
}

func (i *importing) importCsv() error {
	// TODO: Refactor this block, use early returns, merge ifs and extract logics.
	defer func() {
		i.storer.close()
		i.csv.close()
	}()

	headers := i.csv.header()
	fmt.Printf("Found %d fields\n(%s)\n", len(headers), i.headerList(headers))
	err := i.storer.create(i.csv.header())
	if err != nil {
		return err
	}

	err = i.storer.startTransaction()
	if err != nil {
		return err
	}
	defer func() {
		err = i.storer.commitTransaction()
	}()

	i.resetProgress()
	connectionName := i.storer.dBConfig().getConnectionName()

	for i.csv.next() {
		i.showProgress()
		// Firebird does not seem to support batch insert?
		if connectionName == driverNameFirebird {
			err = i.insertSQL()
			if err != nil {
				return err
			}
			continue
		}

		err := i.batchInsertSQL()
		if err != nil {
			return err
		}
	}

	if len(i.batch) > 0 && connectionName != driverNameFirebird {
		err = i.storer.batchInsert(i.batch)
		if err != nil {
			rollbackErr := i.storer.rollbackTransaction()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	return nil
}

func (i *importing) headerList(headers cSVFields) string {
	fields := make([]string, len(headers))
	for i, f := range headers {
		fields[i] = f.Name
	}

	return strings.Join(fields, ", ")
}

func (i *importing) resetProgress() {
	i.rowNr = 0
	i.progress = 0
	i.rowCount = i.csv.rowCount()
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
	err := i.storer.insert(i.csv.row()...)
	if err != nil {
		rollbackErr := i.storer.rollbackTransaction()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	return nil
}

func (i *importing) batchInsertSQL() error {
	i.batch = append(i.batch, i.csv.row())
	if i.batchIndex == batchSize {
		i.batchIndex = 0
		err := i.storer.batchInsert(i.batch)
		if err != nil {
			rollbackErr := i.storer.rollbackTransaction()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
		i.batch = nil
		return nil
	}
	i.batchIndex++

	return nil
}
