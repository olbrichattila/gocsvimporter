package importer

import "database/sql"

type SQLExecutor interface {
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...any) (sql.Result, error)
	// Query(string, ...any) (*sql.Rows, error)
	// QueryRow(string, ...any) *sql.Row
	// we may need to add commit and rollback, The begin is only on sql.DB I think
}

type storager interface {
	execute(SQLExecutor, string, ...any) error
}

type store struct {
	dBconf dBConfiger
}

func newStorager(dBconf dBConfiger) storager {
	return &store{
		dBconf: dBconf,
	}
}

func (s store) execute(exec SQLExecutor, query string, args ...any) error {
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
