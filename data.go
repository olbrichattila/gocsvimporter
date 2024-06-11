package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type dataStore struct {
	dBConfig   DBConfiger
	db         *sql.DB
	tx         *sql.Tx
	tableName  string
	fieldNames []string
	insertSql  string
	quote      string
}

func NewDataStore(dbName, tableName string) (*dataStore, error) {
	d := &dataStore{
		tableName: tableName,
	}
	d.init(dbName)

	return d, nil
}

func (s *dataStore) init(dbName string) error {
	config, err := GetDbConnector()
	if err != nil {
		return err
	}
	s.dBConfig = config
	s.quote = config.GetFieldQuote()

	db, err := NewDbConnection(config)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *dataStore) Close() {
	s.db.Close()
}

func (s *dataStore) Create(fieldNames CSVFields) error {
	sql := s.createSql(fieldNames)
	// fmt.Println(sql)
	s.createInsertSql()

	s.execute(s.dBConfig.GetDropTableString(s.tableName))

	return s.execute(sql)
}

func (s *dataStore) FieldNames() []string {
	return s.fieldNames
}

func (s *dataStore) FieldNamesAsString() string {
	quotedFieldNames := make([]string, len(s.fieldNames))
	for i, f := range s.fieldNames {
		quotedFieldNames[i] = fmt.Sprintf("%s%s%s", s.quote, f, s.quote)

	}

	return strings.Join(quotedFieldNames, ",")
}

func (s *dataStore) createSql(fieldNames CSVFields) string {
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

func (s *dataStore) createInsertSql() {

	binding := s.dBConfig.GetBinding()
	bindings := make([]string, len(s.fieldNames))
	for i := range bindings {
		if binding == "?" {
			bindings[i] = binding
		} else {
			bindings[i] = binding + strconv.Itoa(i+1)
		}
	}

	bindingStr := strings.Join(bindings, ",")

	s.insertSql = fmt.Sprintf("INSERT INTO %s%s%s (%s) VALUES (%s)", s.quote, s.tableName, s.quote, s.FieldNamesAsString(), bindingStr)

}

func (s *dataStore) Insert(args ...any) error {
	err := s.execute(s.insertSql, args...)
	if err != nil {
		return fmt.Errorf("%s,\n%s", err.Error(), s.insertSql)
	}

	return nil
}

func (s *dataStore) BatchInsert(data [][]any) error {
	bindings := make([]string, len(data))
	binding := make([]string, len(s.fieldNames))
	bindingChar := s.dBConfig.GetBinding()

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

	insertSql := fmt.Sprintf("INSERT INTO %s%s%s (%s) VALUES %s", s.quote, s.tableName, s.quote, s.FieldNamesAsString(), strings.Join(bindings, ","))
	var pars []any
	for _, val := range data {
		pars = append(pars, val...)
	}

	err := s.execute(insertSql, pars...)
	if err != nil {
		return fmt.Errorf("%s,\n%s", err.Error(), insertSql)
	}

	return nil
}

func (s *dataStore) StartTransaction() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	s.tx = tx
	return nil
}

func (s *dataStore) CommitTransaction() error {
	if s.tx != nil {
		err := s.tx.Commit()
		if err != nil {
			return err
		}

		s.tx = nil
	}

	return nil
}

func (s *dataStore) RollbackTransaction() error {
	if s.tx != nil {
		err := s.tx.Rollback()
		if err != nil {
			return err
		}
		s.tx = nil
	}

	return nil
}
