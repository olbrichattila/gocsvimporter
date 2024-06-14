package importer

import (
	"database/sql"
	"fmt"
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
	getConnection() *sql.DB
}

type dataStore struct {
	databaseConfig dBConfiger
	sQLGenerator   sQLGenerator
	db             *sql.DB
	tx             *sql.Tx
	tableName      string
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

	s.sQLGenerator = newSQLGenerator(config, tableName)
	s.databaseConfig = config
	s.quote = config.getFieldQuote()

	db, err := newDbConnection(config)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *dataStore) getConnection() *sql.DB {
	return s.db
}

func (s *dataStore) dBConfig() dBConfiger {
	return s.databaseConfig
}

func (s *dataStore) close() {
	s.db.Close()
}

func (s *dataStore) create(fieldNames cSVFields) error {
	sql := s.sQLGenerator.ceateTableSQL(fieldNames)
	s.insertSQL = s.sQLGenerator.createInsertSQL()
	err := s.execute(s.sQLGenerator.getDropTableSQL())
	if err != nil {
		return err
	}

	return s.execute(sql)
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

func (s *dataStore) insert(args ...any) error {
	err := s.execute(s.insertSQL, args...)
	if err != nil {
		return s.getSQLError(err, s.insertSQL)
	}

	return nil
}

func (s *dataStore) batchInsert(data [][]any) error {
	insertSQL, pars := s.sQLGenerator.createBatchInsertSQL(data)
	err := s.execute(insertSQL, pars...)
	if err != nil {
		return s.getSQLError(err, insertSQL)
	}

	return nil
}

func (s *dataStore) getSQLError(err error, insertSQL string) error {
	return fmt.Errorf("%s,\n%s", err.Error(), insertSQL)
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
