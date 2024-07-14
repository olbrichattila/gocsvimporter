package importer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type sQLGenerator interface {
	cerateTableSQL(cSVFields) string
	getDropTableSQL() string
	createInsertSQL() string
	createBatchInsertSQL([][]any, bool) (string, []any)
}

type sQLgen struct {
	databaseConfig       dBConfiger
	bindingChar          string
	tableName            string
	quote                string
	cachedBatchInsertSQL string
	batchInsertSQLCached bool
	normalizedFieldNames []string
	fieldCount           int
}

func newSQLGenerator(dBConfig dBConfiger, tableName string) sQLGenerator {
	return &sQLgen{
		databaseConfig: dBConfig,
		tableName:      tableName,
		quote:          dBConfig.getFieldQuote(),
		bindingChar:    dBConfig.getBinding(),
	}
}

func (g *sQLgen) cerateTableSQL(fieldNames cSVFields) string {
	var crDecl []string
	for _, n := range fieldNames {
		g.fieldCount++
		fn := g.normalizeFieldName(n.Name)
		g.normalizedFieldNames = append(g.normalizedFieldNames, fn)
		crDecl = append(crDecl, fmt.Sprintf("%s%s%s %s", g.quote, fn, g.quote, n.Type))
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

func (g *sQLgen) normalizeFieldName(str string) string {
	p := strings.Split(str, " ")
	var np []string
	for _, pc := range p {
		reg := regexp.MustCompile("[^a-zA-Z0-9]+")
		result := reg.ReplaceAllString(pc, "")
		if len(result) > 0 && unicode.IsDigit(rune(result[0])) {
			result = "a" + result
		}

		if result != "" {
			np = append(np, strings.ToLower(result))
		}
	}

	return strings.Join(np, "_")
}

func (g *sQLgen) fieldNamesAsString() string {
	quotedFieldNames := make([]string, g.fieldCount)
	for i, f := range g.normalizedFieldNames {
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
	return g.databaseConfig.getDropTableString(g.tableName)
}
