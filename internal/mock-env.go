package importer

import (
	"os"

	database "github.com/olbrichattila/gocsvimporter/internal/db"
	"github.com/olbrichattila/gocsvimporter/internal/env"
)

type envMock struct {
	called int
}

func newEnvMock() env.Enver {
	return &envMock{called: 0}
}

func (e *envMock) LoadEnv() error {
	e.called++
	os.Setenv(env.DbConnection, database.DbConnectionTypeSqLite)
	os.Setenv(env.DbDatabase, "./test_database.sqlite")
	return nil
}
