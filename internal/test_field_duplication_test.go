package importer

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/stretchr/testify/suite"
)

type duplicateFieldStruct struct {
	cid       int
	name      string
	ctype     string
	notnull   int
	dfltValue sql.NullString
	pk        int
}

type duplicateFieldsStucts []duplicateFieldStruct

type duplicateRecord struct {
	fieldvarchar string
	fieldint     int
	fieldfloat   float64
	fieldbool    bool
	fieldbool2   bool
}

type duplicateRecords = []duplicateRecord

type duplicationTestSuite struct {
	suite.Suite
	importer importer
	dBConfig database.DBConfiger
}

func TestDuplicationSuiteRunner(t *testing.T) {
	suite.Run(t, new(duplicationTestSuite))
}

func (t *duplicationTestSuite) SetupTest() {
	err := newEnvMock().LoadEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	duplicateMockArgParser := newDuplicateMockArgParser()

	dBConfig, err := database.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	t.dBConfig = dBConfig
	// TODO review when getters are done
	_, _, _, err = duplicateMockArgParser.Parse()
	t.NoError(err)

	t.importer = newImporter(
		dBConfig,
		newCsvReader(duplicateMockArgParser),
		newSQLGenerator(
			dBConfig,
			duplicateMockArgParser,
		),
		newStorager(dBConfig),
	)
}

func (t *duplicationTestSuite) TestImportsCorrectly() {
	_, _, _, err := t.importer.importCsv()
	t.NoError(err)

	db, err := t.reConnect()
	t.NoError(err)
	defer db.Close()

	// Act
	fieldNames, err := t.fieldNames(db)
	t.NoError(err)

	// Assert correct amount of fields created
	t.Len(fieldNames, 5)

	// Assert fields and their type correctly determinded
	t.Equal("field1", fieldNames[0].name)
	t.Equal("VARCHAR(15)", fieldNames[0].ctype)

	t.Equal("field", fieldNames[1].name)
	t.Equal("INT", fieldNames[1].ctype)

	t.Equal("field2", fieldNames[2].name)
	t.Equal("FLOAT", fieldNames[2].ctype)

	t.Equal("field_1", fieldNames[3].name)
	t.Equal("TINYINT(1)", fieldNames[3].ctype)

	t.Equal("field_2", fieldNames[4].name)
	t.Equal("TINYINT(1)", fieldNames[4].ctype)

	// Act
	rows, err := t.fetchAll(db)
	t.NoError(err)
	t.Len(rows, 4)

	// Assert the impored data is identical with CSV
	t.Equal("1", rows[0].fieldvarchar)
	t.Equal(1, rows[0].fieldint)
	t.Equal(1.0, rows[0].fieldfloat)
	t.Equal(true, rows[0].fieldbool)
	t.Equal(true, rows[0].fieldbool2)

	t.Equal("15", rows[1].fieldvarchar)
	t.Equal(15, rows[1].fieldint)
	t.Equal(15.0, rows[1].fieldfloat)
	t.Equal(false, rows[1].fieldbool)
	t.Equal(true, rows[1].fieldbool2)

	t.Equal("16", rows[2].fieldvarchar)
	t.Equal(17, rows[2].fieldint)
	t.Equal(4.9, rows[2].fieldfloat)
	t.Equal(true, rows[2].fieldbool)
	t.Equal(true, rows[2].fieldbool2)

	t.Equal("Hello, John Doe", rows[3].fieldvarchar)
	t.Equal(19, rows[3].fieldint)
	t.Equal(8.9, rows[3].fieldfloat)
	t.Equal(false, rows[3].fieldbool)
	t.Equal(true, rows[3].fieldbool2)

}

func (t *duplicationTestSuite) reConnect() (*sql.DB, error) {
	return t.dBConfig.GetNewConnection()
}

func (t *duplicationTestSuite) fieldNames(database *sql.DB) (duplicateFieldsStucts, error) {
	_, _, tableName, _ := newDuplicateMockArgParser().Parse()
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var fStructs duplicateFieldsStucts

	for rows.Next() {
		var fStruct duplicateFieldStruct
		err := rows.Scan(&fStruct.cid, &fStruct.name, &fStruct.ctype, &fStruct.notnull, &fStruct.dfltValue, &fStruct.pk)
		if err != nil {
			return nil, err
		}
		fStructs = append(fStructs, fStruct)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return fStructs, nil
}

func (t *duplicationTestSuite) fetchAll(database *sql.DB) (duplicateRecords, error) {
	_, _, tableName, _ := newDuplicateMockArgParser().Parse()
	query := fmt.Sprintf("SELECT field1,field,field2,field_1,field_2 FROM %s", tableName)
	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var records duplicateRecords

	for rows.Next() {
		var record duplicateRecord
		err := rows.Scan(&record.fieldvarchar, &record.fieldint, &record.fieldfloat, &record.fieldbool, &record.fieldbool2)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return records, nil
}
