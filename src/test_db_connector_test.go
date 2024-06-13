package importer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type dbConnectorTestSuite struct {
	suite.Suite
}

func TestDbConnectorRunner(t *testing.T) {
	suite.Run(t, new(dbConnectorTestSuite))
}

func (t *dbConnectorTestSuite) TestConnectToMemoryDatabase() {
	dBConfig := newMemoryDBConfig()

	conn, err := newDbConnection(dBConfig)

	t.NoError(err)
	t.NotNil(conn)
}
