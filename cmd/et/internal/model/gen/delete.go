package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/jxorm"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
	"strings"
)

func genDelete(table Table, withCache bool) (string, error) {
	keySet := jxorm.NewSet()
	keyVariableSet := jxorm.NewSet()
	keySet.AddStr(table.PrimaryCacheKey.KeyExpression)
	keyVariableSet.AddStr(table.PrimaryCacheKey.KeyLeft)

	camel := table.Name.ToCamel()
	text, err := utils.LoadTemplate(category, deleteTemplateFile, template.Delete)
	if err != nil {
		return "", err
	}

	output, err := utils.With("delete").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject":     camel,
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey,
			"lowerStartCamelPrimaryKey": jxorm.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
			"dataType":                  table.PrimaryKey.DataType,
			"keys":                      strings.Join(keySet.KeysStr(), "\n"),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source()),
			"keyValues":                 strings.Join(keyVariableSet.KeysStr(), ", "),
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
