package importer

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
	env := newEnv("./fixtures/env")
	err := env.loadEnv()
	t.NoError(err)

	t.Equal("conection", os.Getenv(envdbConnection))
	t.Equal("usenName", os.Getenv(envdbUserName))
	t.Equal("password", os.Getenv(envdbPassword))
	t.Equal("host", os.Getenv(envdbHost))
	t.Equal("port", os.Getenv(envdbPort))
	t.Equal("db", os.Getenv(envdbDatabase))
	t.Equal("sslmode", os.Getenv(envdbSSLMode))
}
