package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	dbFieldBoolean = "TINYINT(1)"
	dbFieldInt     = "INT"
	dbFieldFloat   = "FLOAT"
	dbFieldText    = "VARCHAR"
)

func newCsvReader() *readCsv {
	return &readCsv{}
}

type csvReader interface {
	init(f string, c rune) error
	header() cSVFields
	next() bool
	row() []any
	rowCount() int
	close()
}

type cSVField struct {
	Name string
	Type string
}

type cSVFields = []cSVField

type readCsv struct {
	file          *os.File
	reader        *csv.Reader
	headers       []string
	lengths       []int
	types         []string
	fields        []any
	totalRowCount int
}

func (r *readCsv) header() cSVFields {
	fields := make(cSVFields, len(r.headers))
	for i, fieldName := range r.headers {
		fields[i].Name = fieldName
		fields[i].Type = r.constructType(r.types[i], r.lengths[i])
	}

	return fields
}

func (r *readCsv) next() bool {
	r.fields = nil
	record, err := r.reader.Read()
	if err == io.EOF {
		return false
	}

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, v := range record {
		r.fields = append(r.fields, v)
	}

	return true
}

func (r *readCsv) row() []any {
	rowsWithNil := make([]any, len(r.fields))
	for i, v := range r.fields {
		if r.fields[i] == "" {
			rowsWithNil[i] = nil
		} else {
			rowsWithNil[i] = v
		}
	}
	return rowsWithNil
}

func (r *readCsv) rowCount() int {
	return r.totalRowCount
}

func (r *readCsv) close() {
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

	r.headers = record

	err = r.fillLengths()
	if err != nil {
		return err
	}

	return nil
}

func (r *readCsv) fillLengths() error {
	headerLen := len(r.headers)
	lengths := make([]int, headerLen)
	types := make([]string, headerLen)

	r.totalRowCount = 0
	for r.next() {
		r.totalRowCount++
		if len(r.row()) > 0 {
			lengths, types = r.constructTypeAndLengths(lengths, types)
		}
	}

	r.lengths = lengths
	r.types = types
	_, err := r.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_ = r.next()

	return nil
}

func (r *readCsv) constructTypeAndLengths(lengths []int, types []string) ([]int, []string) {
	for i, v := range r.row() {
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
	if currentRowType == dbFieldText ||
		(currentRowType == dbFieldFloat && lastConstructedFieldType != dbFieldText) ||
		(currentRowType == dbFieldInt && lastConstructedFieldType != dbFieldFloat && lastConstructedFieldType != dbFieldText) {
		return currentRowType
	}

	if lastConstructedFieldType == dbFieldInt {
		return lastConstructedFieldType
	}

	return currentRowType
}

func (r *readCsv) constructType(fieldType string, length int) string {
	if fieldType == dbFieldText {
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
		return dbFieldBoolean
	}

	if r.isInt(s) {
		return dbFieldInt
	}

	if r.isFloat(s) {
		return dbFieldFloat
	}

	return dbFieldText
}
