package importer

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type fieldStruct struct {
	cid       int
	name      string
	ctype     string
	notnull   int
	dfltValue sql.NullString
	pk        int
}

type fieldsStucts []fieldStruct

type record struct {
	fieldvarchar string
	fieldint     int
	fieldfloat   float64
	fieldbool    bool
}

type records = []record

type integrationTestSuite struct {
	suite.Suite
	app      *application
	dBconfig dBConfiger
}

func TestIntegrationRunner(t *testing.T) {
	suite.Run(t, new(integrationTestSuite))
}

func (t *integrationTestSuite) SetupTest() {

	err := newEnvMock().loadEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	dBconfig, err := newDbConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	t.dBconfig = dBconfig
	csvFileName, separator, tableName, err := newMockParser().parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	app, err := newApplication(
		newImporter(
			dBconfig,
			newCsvReader(csvFileName, separator),
			newSQLGenerator(
				dBconfig,
				tableName,
			),
			newStorager(dBconfig),
		),
	)

	t.NoError(err)
	t.app = app
}

func (t *integrationTestSuite) TestImportsCorrectly() {
	_, _, _, err := t.app.importer.importCsv()
	t.NoError(err)

	db, err := t.reConnect()
	t.NoError(err)
	defer db.Close()

	// Act
	fieldNames, err := t.fieldNames(db)
	t.NoError(err)

	// Assert correct amount of fields created
	t.Len(fieldNames, 4)

	// Assert fields and their type correctly determinded
	t.Equal("fieldvarchar", fieldNames[0].name)
	t.Equal("VARCHAR(15)", fieldNames[0].ctype)

	t.Equal("fieldint", fieldNames[1].name)
	t.Equal("INT", fieldNames[1].ctype)

	t.Equal("fieldfloat", fieldNames[2].name)
	t.Equal("FLOAT", fieldNames[2].ctype)

	t.Equal("fieldbool", fieldNames[3].name)
	t.Equal("TINYINT(1)", fieldNames[3].ctype)

	// Act
	rows, err := t.fetcAll(db)
	t.NoError(err)
	t.Len(rows, 4)

	// Assert the impored data is identical with CSV
	t.Equal("1", rows[0].fieldvarchar)
	t.Equal(1, rows[0].fieldint)
	t.Equal(1.0, rows[0].fieldfloat)
	t.Equal(true, rows[0].fieldbool)

	t.Equal("15", rows[1].fieldvarchar)
	t.Equal(15, rows[1].fieldint)
	t.Equal(15.0, rows[1].fieldfloat)
	t.Equal(false, rows[1].fieldbool)

	t.Equal("16", rows[2].fieldvarchar)
	t.Equal(17, rows[2].fieldint)
	t.Equal(4.9, rows[2].fieldfloat)
	t.Equal(true, rows[2].fieldbool)

	t.Equal("Hello, John Doe", rows[3].fieldvarchar)
	t.Equal(19, rows[3].fieldint)
	t.Equal(8.9, rows[3].fieldfloat)
	t.Equal(false, rows[3].fieldbool)
}

func (t *integrationTestSuite) reConnect() (*sql.DB, error) {
	return t.dBconfig.getNewConnection()
}

func (t *integrationTestSuite) fieldNames(database *sql.DB) (fieldsStucts, error) {
	_, _, tableName, _ := newMockParser().parse()
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var fStructs fieldsStucts

	for rows.Next() {
		var fStruct fieldStruct
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

func (t *integrationTestSuite) fetcAll(database *sql.DB) (records, error) {
	_, _, tableName, _ := newMockParser().parse()
	query := fmt.Sprintf("SELECT fieldvarchar, fieldint, fieldfloat, fieldbool FROM 	%s", tableName)
	rows, err := database.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var records records

	for rows.Next() {
		var record record
		err := rows.Scan(&record.fieldvarchar, &record.fieldint, &record.fieldfloat, &record.fieldbool)
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
