package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type dbConnection struct {
	db     *sql.DB
	config DBConfiger
}

func NewDbConnection(config DBConfiger) (*sql.DB, error) {
	conn := &dbConnection{
		config: config,
	}
	err := conn.Connect()
	if err != nil {
		return nil, err
	}

	return conn.db, nil
}

func (d *dbConnection) Connect() error {
	db, err := sql.Open(d.config.GetConnectionName(), d.config.GetConnectionString())
	if err != nil {
		return err
	}
	d.db = db
	return nil
}
