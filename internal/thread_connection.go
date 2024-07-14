package importer

import "database/sql"

type threadConnection struct {
	db *sql.DB
	tx *sql.Tx
}

type threadConnections []*threadConnection

func (t *threadConnection) getExecutor() sQLExecutor {
	if t.tx != nil {
		return t.tx
	}

	return t.db
}
