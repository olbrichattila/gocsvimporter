package importer

// TODO: this will be deprecated
import (
	"database/sql"

	// This is required to work, vscode auto removes it
	_ "github.com/mattn/go-sqlite3"
)

type dbConnection struct {
	db     *sql.DB
	config dBConfiger
}

func newDbConnection(config dBConfiger) (*sql.DB, error) {
	conn := &dbConnection{
		config: config,
	}
	err := conn.connect()
	if err != nil {
		return nil, err
	}

	return conn.db, nil
}

func (d *dbConnection) connect() error {
	db, err := sql.Open(d.config.getConnectionName(), d.config.getConnectionString())
	if err != nil {
		return err
	}
	d.db = db
	return nil
}
