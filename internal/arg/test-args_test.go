package arg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

const testTableName = "test_table"

type argsTestSuite struct {
	suite.Suite
	originalArgs []string
	parser       Parser
}

func TestArgsRunner(t *testing.T) {
	suite.Run(t, new(argsTestSuite))

}

func (t *argsTestSuite) SetupTest() {
	t.originalArgs = os.Args
	t.parser = New()
}

func (t *argsTestSuite) TearDownTest() {
	os.Args = t.originalArgs
}

func (t *argsTestSuite) TestMissingArgsReturnsError() {
	os.Args = []string{}
	_, _, _, err := t.parser.Parse()
	t.Error(err)
}

func (t *argsTestSuite) TestArgsReturnedWithDefaultSeparator() {
	os.Args = []string{"1", "./fixtures/testfile.csv", testTableName}
	fileName, separator, tableName, err := t.parser.Parse()

	t.NoError(err)
	t.Equal("./fixtures/testfile.csv", fileName)
	t.Equal(testTableName, tableName)
	t.Equal(',', separator)
}

func (t *argsTestSuite) TestArgsReturnedWithCustomSeparator() {
	os.Args = []string{"1", "./fixtures/testfile.csv", testTableName, ";"}
	fileName, separator, tableName, err := t.parser.Parse()

	t.NoError(err)
	t.Equal("./fixtures/testfile.csv", fileName)
	t.Equal(testTableName, tableName)
	t.Equal(';', separator)
}

func (t *argsTestSuite) TestFileNotExistsReturnsError() {
	os.Args = []string{"1", "testfile-missing.csv", testTableName, ";"}
	_, _, _, err := t.parser.Parse()

	t.Error(err)
}
