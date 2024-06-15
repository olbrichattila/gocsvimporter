package importer

import "os"

type envMock struct {
	called int
}

func newEnvMock() enver {
	return &envMock{called: 0}
}

func (e *envMock) loadEnv() error {
	e.called++
	os.Setenv(envdbConnection, dbConnectionTypeSqLite)
	os.Setenv(envdbDatabase, "./test_database.sqlite")
	return nil
}
