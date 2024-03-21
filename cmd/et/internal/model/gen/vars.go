package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/jxorm"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
	"strings"
)

func genVars(table Table, withCache bool) (string, error) {
	keys := make([]string, 0)
	keys = append(keys, table.PrimaryCacheKey.VarExpression)

	camel := table.Name.ToCamel()
	text, err := utils.LoadTemplate(category, varTemplateFile, template.Vars)
	if err != nil {
		return "", err
	}

	output, err := utils.With("var").Parse(text).
		GoFmt(true).Execute(map[string]interface{}{
		"lowerStartCamelObject": jxorm.From(camel).Untitle(),
		"upperStartCamelObject": camel,
		"cacheKeys":             strings.Join(keys, "\n"),
		"autoIncrement":         table.PrimaryKey.AutoIncrement,
		"originalPrimaryKey":    wrapWithRawString(table.PrimaryKey.Name.Source()),
		"withCache":             withCache,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
