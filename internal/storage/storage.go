// Package storage abstracts database execution
package storage

import (
	"database/sql"

	database "github.com/olbrichattila/gocsvimporter/internal/db"
)

// SQLExecutor executor abstracts database prepare and exec
type SQLExecutor interface {
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...any) (sql.Result, error)
}

// Storager is the database executor
type Storager interface {
	Execute(SQLExecutor, string, ...any) error
}

type store struct {
	dBConf database.DBConfiger
}

// New creates a new storager
func New(dBconf database.DBConfiger) Storager {
	return &store{
		dBConf: dBconf,
	}
}

// Executes execute an SQL expression
func (s store) Execute(exec SQLExecutor, query string, args ...any) error {
	stmt, err := exec.Prepare(query)
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
