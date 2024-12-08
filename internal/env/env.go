// Package env manages .env file
package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Definition of different environment variable names
const (
	DbConnection = "DB_CONNECTION"
	DbUserName   = "DB_USERNAME"
	DbPassword   = "DB_PASSWORD"
	DbHost       = "DB_HOST"
	DbPort       = "DB_PORT"
	DbDatabase   = "DB_DATABASE"
	DbSSLMode    = "DB_SSLMODE"

	BatchSize          = "BATCH_SIZE"
	MaxConnectionCount = "MAX_CONNECTION_COUNT"

	BatchInsert         = "BATCH_INSERT"
	MultipleConnections = "MULTIPLE_CONNECTIONS"
	Transactional       = "TRANSACTIONAL"
)

const ()

// Enver is the interface for managing environment variables (linux and .env)
type Enver interface {
	LoadEnv() error
}

func New(fileName string) Enver {
	return &env{fileName: fileName}
}

type env struct {
	fileName string
}

// LoadEnv loads environment variables from .env if exists
func (e *env) LoadEnv() error {
	fmt.Println("load env ", e.fileName)
	_, err := os.Stat(e.fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return godotenv.Load(e.fileName)
}
