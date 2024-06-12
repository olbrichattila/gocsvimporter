package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func NewCsvReader(fileName string, spearator rune) (*readCsv, error) {
	c := &readCsv{}
	c.init(fileName, spearator)
	return c, nil
}

type csvReader interface {
	Header() CSVFields
	Next() bool
	Row() []any
	RowCount() int
	Close()
}

type CSVField struct {
	Name string
	Type string
}

type CSVFields = []CSVField

type readCsv struct {
	file     *os.File
	reader   *csv.Reader
	header   []string
	lengths  []int
	types    []string
	row      []any
	rowCount int
}

func (r *readCsv) Header() CSVFields {
	fields := make(CSVFields, len(r.header))
	for i, fieldName := range r.header {
		fields[i].Name = fieldName
		fields[i].Type = r.constructType(r.types[i], r.lengths[i])
	}

	return fields
}

func (r *readCsv) Next() bool {
	r.row = nil
	record, err := r.reader.Read()
	if err == io.EOF {
		return false
	}

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, v := range record {
		r.row = append(r.row, v)
	}

	return true
}

func (r *readCsv) Row() []any {
	rowsWithNil := make([]any, len(r.row))
	for i, v := range r.row {
		if r.row[i] == "" {
			rowsWithNil[i] = nil
		} else {
			rowsWithNil[i] = v
		}
	}
	return rowsWithNil
}

func (r *readCsv) RowCount() int {
	return r.rowCount
}

func (r *readCsv) Close() {
	r.file.Close()
}

func (r *readCsv) init(f string, c rune) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}

	r.file = file

	reader := csv.NewReader(r.file)
	reader.Comma = c
	r.reader = reader

	err = r.setHeader()
	if err != nil {
		return err
	}

	return nil
}

func (r *readCsv) setHeader() error {
	record, err := r.reader.Read()
	if err != nil {
		return err
	}

	r.header = record

	err = r.fillLengths()
	if err != nil {
		return err
	}

	return nil
}

func (r *readCsv) fillLengths() error {
	headerLen := len(r.header)
	lengths := make([]int, headerLen)
	types := make([]string, headerLen)

	r.rowCount = 0
	for r.Next() {
		r.rowCount++
		if len(r.Row()) > 0 {
			lengths, types = r.constructTypeAndLengths(r.Row(), lengths, types)
		}
	}

	r.lengths = lengths
	r.types = types
	r.file.Seek(0, io.SeekStart)
	_ = r.Next()

	return nil
}

func (r *readCsv) constructTypeAndLengths(row []any, lengths []int, types []string) ([]int, []string) {
	for i, v := range r.Row() {
		st := fmt.Sprintf("%v", v)
		stLn := len(st)
		if lengths[i] < stLn {
			lengths[i] = stLn
		}

		types[i] = r.constructFieldType(st, types[i])
	}

	return lengths, types
}

func (r *readCsv) constructFieldType(fieldContent, lastConstructedFieldType string) string {
	currentRowType := r.getType(fieldContent)
	if currentRowType == "VARCHAR" ||
		(currentRowType == "FLOAT" && lastConstructedFieldType != "VARCHAR") ||
		(currentRowType == "INT" && lastConstructedFieldType != "FLOAT" && lastConstructedFieldType != "VARCHAR") {
		return currentRowType
	}

	if lastConstructedFieldType == "INT" {
		return lastConstructedFieldType
	}

	return currentRowType
}

func (r *readCsv) constructType(fieldType string, length int) string {
	if fieldType == "VARCHAR" {
		return "VARCHAR(" + strconv.Itoa(length) + ")"
	}

	return fieldType
}

func (r *readCsv) isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (r *readCsv) isFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (r *readCsv) isBool(s string) bool {
	return s == "0" || s == "1"
}

func (r *readCsv) getType(s string) string {
	if r.isBool(s) {
		return "TINYINT(1)"
	}

	if r.isInt(s) {
		return "INT"
	}

	if r.isFloat(s) {
		return "FLOAT"
	}

	return "VARCHAR"
}
