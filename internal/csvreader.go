package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/olbrichattila/gocsvimporter/internal/arg"
)

const (
	dbFieldText    = 4
	dbFieldFloat   = 3
	dbFieldInt     = 2
	dbFieldBoolean = 1
)

var dbFieldMap = map[int]string{
	1: "TINYINT(1)",
	2: "INT",
	3: "FLOAT",
	4: "VARCHAR",
}

type fieldDef struct {
	Name   string `json:"name"`
	Type   string `type:"name"`
	Length int    `length:"name2"`
}

type fieldDefs = []fieldDef

func newCsvReader(parser arg.Parser) csvReader {
	return &readCsv{
		fileName:  parser.FileName(),
		separator: parser.Separator(),
	}
}

type csvReader interface {
	init() error
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
	fileName      string
	separator     rune
	file          *os.File
	reader        *csv.Reader
	headers       []string
	headerLen     int
	lengths       []int
	types         []int
	fields        []any
	totalRowCount int
}

func (r *readCsv) header() cSVFields {
	fields := make(cSVFields, r.headerLen)
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

func (r *readCsv) init() error {
	file, err := os.Open(r.fileName)
	if err != nil {
		return err
	}

	r.file = file

	reader := csv.NewReader(r.file)
	reader.Comma = r.separator
	r.reader = reader

	// err = r.loadHeaders()
	// if err != nil {
	// 	return err
	// }

	err = r.setHeader()
	if err != nil {
		return err
	}

	// TODO save only if flag set
	err = r.saveHeaders()
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

	r.headers = r.deDuplicateHeader(record)
	r.headerLen = len(record)

	err = r.fillLengths()
	if err != nil {

		return err
	}

	return nil
}

func (r *readCsv) fillLengths() error {
	lengths := make([]int, r.headerLen)
	types := make([]int, r.headerLen)

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

func (r *readCsv) constructTypeAndLengths(lengths []int, types []int) ([]int, []int) {
	for i, v := range r.row() {
		if v != nil {
			st := v.(string)
			stLn := len(st)
			if lengths[i] < stLn {
				lengths[i] = stLn
			}

			types[i] = r.constructFieldType(st, types[i])
		}
	}

	return lengths, types
}

func (r *readCsv) constructFieldType(fieldContent string, lastConstructedFieldType int) int {
	currentType := r.getType(fieldContent)
	if currentType > lastConstructedFieldType {
		return currentType
	}

	return lastConstructedFieldType
}

func (r *readCsv) constructType(fieldType int, length int) string {
	if fieldType == dbFieldText {
		return "VARCHAR(" + strconv.Itoa(length) + ")"
	}

	return dbFieldMap[fieldType]
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

func (r *readCsv) getType(s string) int {
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

func (r *readCsv) deDuplicateHeader(header []string) []string {
	deDupedHeader := make([]string, 0)

	for _, fieldName := range header {
		fn := strings.TrimSpace(fieldName)
		if fn == "" {
			fn = "unknown"
		}
		normalizedFieldName := r.normalizeFieldName(fieldName)
		uniqueFieldName := r.uniqueFieldName(normalizedFieldName, deDupedHeader)
		deDupedHeader = append(deDupedHeader, uniqueFieldName)
	}

	return deDupedHeader
}

func (r *readCsv) uniqueFieldName(fieldName string, header []string) string {
	if !r.fieldExists(fieldName, header) {
		return fieldName
	}
	index := 1
	for {
		newFieldName := fieldName + "_" + strconv.Itoa(index)
		if !r.fieldExists(newFieldName, header) {
			return newFieldName
		}
		index++
	}
}

func (r *readCsv) fieldExists(fieldName string, header []string) bool {
	for _, fn := range header {
		if fn == fieldName {
			return true
		}
	}

	return false
}

func (*readCsv) normalizeFieldName(str string) string {
	p := strings.Split(str, " ")
	var np []string
	for _, pc := range p {
		reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
		result := reg.ReplaceAllString(pc, "")
		if len(result) > 0 && unicode.IsDigit(rune(result[0])) {
			result = "a" + result
		}

		if result != "" {
			np = append(np, strings.ToLower(result))
		}
	}

	return strings.Join(np, "_")
}

func (r *readCsv) saveHeaders() error {
	aggregatedFieldDefs := make(fieldDefs, 0)
	for i, typeID := range r.types {
		aggregatedFieldDefs = append(
			aggregatedFieldDefs,
			fieldDef{
				Type:   dbFieldMap[typeID],
				Name:   r.headers[i],
				Length: r.lengths[i],
			},
		)

	}

	jsonData, _ := json.MarshalIndent(aggregatedFieldDefs, "", " ")
	file, err := os.Create("header-info.json")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the string to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (r *readCsv) loadHeaders() error {
	data, err := os.ReadFile("header-info.json")
	if err != nil {
		return err
	}

	var aggregatedFieldDefs fieldDefs
	err = json.Unmarshal(data, &aggregatedFieldDefs)
	if err != nil {
		return err
	}

	r.headerLen = len(aggregatedFieldDefs)

	r.types = make([]int, r.headerLen)
	r.lengths = make([]int, r.headerLen)
	r.headers = make([]string, r.headerLen)

	for i, fDef := range aggregatedFieldDefs {
		r.lengths[i] = fDef.Length
		r.headers[i] = fDef.Name
		fieldType, err := r.getFieldTypeByName(fDef.Type)
		if err != nil {
			return err
		}

		r.types[i] = fieldType
	}

	r.totalRowCount = -1
	return nil
}

func (*readCsv) getFieldTypeByName(name string) (int, error) {
	for key, value := range dbFieldMap {
		if value == name {
			return key, nil
		}
	}

	return 0, fmt.Errorf("cannot find field " + name)
}
