package importer

import "database/sql"

type connection struct {
}

func (c *connection) connect(connectionName, connectionString string) (*sql.DB, error) {
	db, err := sql.Open(connectionName, connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
