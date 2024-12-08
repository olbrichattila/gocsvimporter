package importer

import (
	"database/sql"

	"github.com/olbrichattila/gocsvimporter/internal/storage"
)

type threadConnection struct {
	db *sql.DB
	tx *sql.Tx
}

type threadConnections []*threadConnection

func (t *threadConnection) getExecutor() storage.SQLExecutor {
	if t.tx != nil {
		return t.tx
	}

	return t.db
}
