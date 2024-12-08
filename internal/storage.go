package importer

import (
	"database/sql"

	database "github.com/olbrichattila/gocsvimporter/internal/db"
)

type sQLExecutor interface {
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...any) (sql.Result, error)
}

type storager interface {
	execute(sQLExecutor, string, ...any) error
}

type store struct {
	dBConf database.DBConfiger
}

func newStorager(dBconf database.DBConfiger) storager {
	return &store{
		dBConf: dBconf,
	}
}

func (s store) execute(exec sQLExecutor, query string, args ...any) error {
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
