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
	err := t.parser.Validate()
	t.Error(err)
}

func (t *argsTestSuite) TestArgsReturnedWithDefaultSeparator() {
	os.Args = []string{"1", "./fixtures/testfile.csv", testTableName}
	separator := t.parser.Separator()

	t.Equal(',', separator)
}

func (t *argsTestSuite) TestArgsReturnedWithCustomSeparator() {
	os.Args = []string{"1", "./fixtures/testfile.csv", testTableName, ";"}
	separator := t.parser.Separator()

	t.Equal(';', separator)
}

func (t *argsTestSuite) TestFileNotExistsReturnsError() {
	os.Args = []string{"1", "testfile-missing.csv", testTableName, ";"}
	err := t.parser.Validate()

	t.Error(err)
}

func (t *argsTestSuite) TestValidationErrorIfSeparatorIsTooLong() {
	os.Args = []string{"1", "./fixtures/testfile.csv", testTableName, "long"}
	err := t.parser.Validate()

	t.Error(err)
}

func (t *argsTestSuite) TestFlags() {
	os.Args = []string{"1", "./fixtures/testfile.csv", testTableName, ";", "-flag1=hasValue1", "-flagNoValue1=", "-flagNoValue2", "-flag2=hasValue2"}

	_, err := t.parser.Flag("missing")
	t.Error(err)

	result, err := t.parser.Flag("flag1")
	t.Nil(err)
	t.Equal("hasValue1", result)

	result, err = t.parser.Flag("flagNoValue1")
	t.Nil(err)
	t.Equal("", result)

	result, err = t.parser.Flag("flagNoValue2")
	t.Nil(err)
	t.Equal("", result)

	result, err = t.parser.Flag("flag2")
	t.Nil(err)
	t.Equal("hasValue2", result)
}
