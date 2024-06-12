package importer

import "fmt"

type dBConfigerMock struct {
}

func newMockDBConfig() *dBConfigerMock {
	return &dBConfigerMock{}
}

func (c *dBConfigerMock) getConnectionString() string {
	return ":memory:"
}

func (c *dBConfigerMock) getConnectionName() string {
	return driverNameSqLite
}

func (c *dBConfigerMock) getFieldQuote() string {
	return "\""
}

func (c *dBConfigerMock) getBinding() string {
	return "?"
}

func (c *dBConfigerMock) getDropTableString(tableName string) string {
	quote := c.getFieldQuote()
	return fmt.Sprintf(defaultDropTableFormat, quote, tableName, quote)
}
