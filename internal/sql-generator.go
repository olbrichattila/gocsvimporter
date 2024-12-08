package importer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/olbrichattila/gocsvimporter/internal/arg"
	database "github.com/olbrichattila/gocsvimporter/internal/db"
)

type sQLGenerator interface {
	cerateTableSQL(cSVFields) string
	getDropTableSQL() string
	createInsertSQL() string
	createBatchInsertSQL([][]any, bool) (string, []any)
}

type sQLgen struct {
	databaseConfig       database.DBConfiger
	bindingChar          string
	tableName            string
	quote                string
	cachedBatchInsertSQL string
	batchInsertSQLCached bool
	fieldNames           []string
	fieldCount           int
}

func newSQLGenerator(dBConfig database.DBConfiger, parser arg.Parser) sQLGenerator {
	return &sQLgen{
		databaseConfig: dBConfig,
		tableName:      parser.TableName(),
		quote:          dBConfig.GetFieldQuote(),
		bindingChar:    dBConfig.GetBinding(),
	}
}

func (g *sQLgen) cerateTableSQL(fieldNames cSVFields) string {
	var crDecl []string
	for _, n := range fieldNames {
		g.fieldCount++

		g.fieldNames = append(g.fieldNames, n.Name)
		crDecl = append(crDecl, fmt.Sprintf("%s%s%s %s", g.quote, n.Name, g.quote, n.Type))
	}

	body := strings.Join(crDecl, ",\n")

	return fmt.Sprintf("CREATE TABLE %s%s%s (\n%s\n)", g.quote, g.tableName, g.quote, body)
}

func (g *sQLgen) createInsertSQL() string {
	bindings := make([]string, g.fieldCount)
	for i := range bindings {
		if g.bindingChar == "?" {
			bindings[i] = g.bindingChar
		} else {
			bindings[i] = g.bindingChar + strconv.Itoa(i+1)
		}
	}

	bindingStr := strings.Join(bindings, ",")

	return fmt.Sprintf("INSERT INTO %s%s%s (%s) VALUES (%s)", g.quote, g.tableName, g.quote, g.fieldNamesAsString(), bindingStr)
}

func (g *sQLgen) createBatchInsertSQL(data [][]any, isFullBatch bool) (string, []any) {
	var pars []any
	for _, val := range data {
		pars = append(pars, val...)
	}
	if isFullBatch && g.batchInsertSQLCached {
		return g.cachedBatchInsertSQL, pars
	}

	bindingStr := g.getBatchBindings(len(data), g.fieldCount)
	insertSQL := fmt.Sprintf("INSERT INTO %s%s%s (%s) VALUES %s", g.quote, g.tableName, g.quote, g.fieldNamesAsString(), bindingStr)

	if isFullBatch {
		g.cachedBatchInsertSQL = insertSQL
		g.batchInsertSQLCached = true
	}

	return insertSQL, pars
}

func (g *sQLgen) fieldNamesAsString() string {
	quotedFieldNames := make([]string, g.fieldCount)
	for i, f := range g.fieldNames {
		quotedFieldNames[i] = fmt.Sprintf("%s%s%s", g.quote, f, g.quote)

	}

	return strings.Join(quotedFieldNames, ",")
}

func (g *sQLgen) getBatchBindings(dataLen, fieldsLen int) string {
	bindings := make([]string, dataLen)
	binding := make([]string, fieldsLen)

	bindingPos := 0
	for i := range bindings {
		for x := range binding {
			if g.bindingChar == "?" {
				binding[x] = g.bindingChar
			} else {
				bindingPos++
				binding[x] = g.bindingChar + strconv.Itoa(bindingPos)
			}
		}
		bindings[i] = fmt.Sprintf("(%s)", strings.Join(binding, ","))
	}

	return strings.Join(bindings, ",")
}

func (g *sQLgen) getDropTableSQL() string {
	return g.databaseConfig.GetDropTableString(g.tableName)
}
