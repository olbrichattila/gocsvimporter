package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type envTestSuite struct {
	suite.Suite
}

func TestEnvRunner(t *testing.T) {
	suite.Run(t, new(envTestSuite))
}

func (t *envTestSuite) TestMissingArgsReturnsError() {
	env := New("./fixtures/env")
	err := env.LoadEnv()
	t.NoError(err)

	t.Equal("conection", os.Getenv(DbConnection))
	t.Equal("usenName", os.Getenv(DbUserName))
	t.Equal("password", os.Getenv(DbPassword))
	t.Equal("host", os.Getenv(DbHost))
	t.Equal("port", os.Getenv(DbPort))
	t.Equal("db", os.Getenv(DbDatabase))
	t.Equal("sslmode", os.Getenv(DbSSLMode))
}
