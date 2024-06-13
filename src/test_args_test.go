package importer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type argsTestSuite struct {
	suite.Suite
	originalArgs []string
	pharser      argParser
}

func TestArgsRunner(t *testing.T) {
	suite.Run(t, new(argsTestSuite))

}

func (t *argsTestSuite) SetupTest() {
	t.originalArgs = os.Args
	t.pharser = newArgParser()
}

func (t *argsTestSuite) TearDownTest() {
	os.Args = t.originalArgs
}

func (t *argsTestSuite) TestMissingArgsReturnsError() {
	os.Args = []string{}
	_, _, _, err := t.pharser.pharse()
	t.Error(err)
}

func (t *argsTestSuite) TestArgsReturnedWithDefaultSeparator() {
	os.Args = []string{"1", "./fixtures/testfile.csv", "test_table"}
	fileName, tableName, separator, err := t.pharser.pharse()

	t.NoError(err)
	t.Equal("./fixtures/testfile.csv", fileName)
	t.Equal("test_table", tableName)
	t.Equal(',', separator)
}

func (t *argsTestSuite) TestArgsReturnedWithCustomSeparator() {
	os.Args = []string{"1", "./fixtures/testfile.csv", "test_table", ";"}
	fileName, tableName, separator, err := t.pharser.pharse()

	t.NoError(err)
	t.Equal("./fixtures/testfile.csv", fileName)
	t.Equal("test_table", tableName)
	t.Equal(';', separator)
}

func (t *argsTestSuite) TestFileNotExistsReturnsError() {
	os.Args = []string{"1", "testfile-missing.csv", "test_table", ";"}
	_, _, _, err := t.pharser.pharse()

	t.Error(err)
}
