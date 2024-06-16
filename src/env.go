package importer

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	envFileName     = "./.env"
	envdbConnection = "DB_CONNECTION"
	envdbUserName   = "DB_USERNAME"
	envdbPassword   = "DB_PASSWORD"
	envdbHost       = "DB_HOST"
	envdbPort       = "DB_PORT"
	envdbDatabase   = "DB_DATABASE"
	envdbSSLMode    = "DB_SSLMODE"

	envBatchSize          = "BATCH_SIZE"
	envMaxConnectionCount = "MAX_CONNECTION_COUNT"

	envBatchInsert         = "BATCH_INSERT"
	envMultipleConnections = "MULTIPLE_CONNECTIONS"
	envTransactional       = "TRANSACTIONAL"
)

const ()

type enver interface {
	loadEnv() error
}

func newEnv(fileName string) enver {
	return &env{fileName: fileName}
}

type env struct {
	fileName string
}

func (e *env) loadEnv() error {
	_, err := os.Stat(e.fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return godotenv.Load(e.fileName)
}
