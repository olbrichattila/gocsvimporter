package importer

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type dataStorer interface {
	init(string) error
	close()
	create(cSVFields) error
	startTransaction() error
	commitTransaction() error
	batchInsert([][]any) error
	rollbackTransaction() error
	insert(...any) error
	dBConfig() dBConfiger
}

type dataStore struct {
	databaseConfig dBConfiger
	db             *sql.DB
	tx             *sql.Tx
	tableName      string
	fieldNames     []string
	insertSQL      string
	quote          string
}

func newDataStore() *dataStore {
	return &dataStore{}
}

func (s *dataStore) init(tableName string) error {
	s.tableName = tableName
	config, err := getDbConnector()
	if err != nil {
		return err
	}
	s.databaseConfig = config
	s.quote = config.getFieldQuote()

	db, err := newDbConnection(config)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *dataStore) dBConfig() dBConfiger {
	return s.databaseConfig
}

func (s *dataStore) close() {
	s.db.Close()
}

func (s *dataStore) create(fieldNames cSVFields) error {
	sql := s.createSQL(fieldNames)
	s.createInsertSQL()
	err := s.execute(s.databaseConfig.getDropTableString(s.tableName))
	if err != nil {
		return err
	}

	return s.execute(sql)
}

func (s *dataStore) fieldNamesAsString() string {
	quotedFieldNames := make([]string, len(s.fieldNames))
	for i, f := range s.fieldNames {
		quotedFieldNames[i] = fmt.Sprintf("%s%s%s", s.quote, f, s.quote)

	}

	return strings.Join(quotedFieldNames, ",")
}

func (s *dataStore) createSQL(fieldNames cSVFields) string {
	var crDecl []string
	for _, n := range fieldNames {
		fn := s.normalizeFieldName(n.Name)
		s.fieldNames = append(s.fieldNames, fn)
		crDecl = append(crDecl, fmt.Sprintf("%s%s%s %s", s.quote, fn, s.quote, n.Type))

	}

	body := strings.Join(crDecl, ",\n")

	return fmt.Sprintf("CREATE TABLE %s%s%s (\n%s\n)", s.quote, s.tableName, s.quote, body)
}

func (s *dataStore) normalizeFieldName(str string) string {
	p := strings.Split(str, " ")
	var np []string
	for _, pc := range p {
		reg := regexp.MustCompile("[^a-zA-Z0-9]+")
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

func (s *dataStore) execute(query string, args ...any) error {
	var err error
	var stmt *sql.Stmt

	if s.tx != nil {
		stmt, err = s.tx.Prepare(query)
	} else {
		stmt, err = s.db.Prepare(query)
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *dataStore) createInsertSQL() {
	binding := s.databaseConfig.getBinding()
	bindings := make([]string, len(s.fieldNames))
	for i := range bindings {
		if binding == "?" {
			bindings[i] = binding
		} else {
			bindings[i] = binding + strconv.Itoa(i+1)
		}
	}

	bindingStr := strings.Join(bindings, ",")

	s.insertSQL = fmt.Sprintf("INSERT INTO %s%s%s (%s) VALUES (%s)", s.quote, s.tableName, s.quote, s.fieldNamesAsString(), bindingStr)

}

func (s *dataStore) insert(args ...any) error {
	err := s.execute(s.insertSQL, args...)
	if err != nil {
		return s.getSqlError(err, s.insertSQL)
	}

	return nil
}

func (s *dataStore) batchInsert(data [][]any) error {
	bindinStr := s.getBatchBindings(len(data), len(s.fieldNames))

	insertSQL := fmt.Sprintf("INSERT INTO %s%s%s (%s) VALUES %s", s.quote, s.tableName, s.quote, s.fieldNamesAsString(), bindinStr)
	var pars []any
	for _, val := range data {
		pars = append(pars, val...)
	}

	err := s.execute(insertSQL, pars...)
	if err != nil {
		return s.getSqlError(err, insertSQL)
	}

	return nil
}

func (s *dataStore) getSqlError(err error, insertSQL string) error {
	return fmt.Errorf("%s,\n%s", err.Error(), insertSQL)
}

func (s *dataStore) getBatchBindings(dataLen, fieldsLen int) string {
	bindings := make([]string, dataLen)
	binding := make([]string, fieldsLen)
	bindingChar := s.databaseConfig.getBinding()

	bindingPos := 0
	for i := range bindings {
		for x := range binding {
			if bindingChar == "?" {
				binding[x] = bindingChar
			} else {
				bindingPos++
				binding[x] = bindingChar + strconv.Itoa(bindingPos)
			}
		}
		bindings[i] = fmt.Sprintf("(%s)", strings.Join(binding, ","))
	}

	return strings.Join(bindings, ",")
}

func (s *dataStore) startTransaction() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	s.tx = tx
	return nil
}

func (s *dataStore) commitTransaction() error {
	if s.tx != nil {
		err := s.tx.Commit()
		if err != nil {
			return err
		}

		s.tx = nil
	}

	return nil
}

func (s *dataStore) rollbackTransaction() error {
	if s.tx != nil {
		err := s.tx.Rollback()
		if err != nil {
			return err
		}
		s.tx = nil
	}

	return nil
}
