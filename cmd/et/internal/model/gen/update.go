package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/jxorm"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
	"strings"
)

func genUpdate(table Table, withCache bool) (string, error) {
	keySet := jxorm.NewSet()
	keyVariableSet := jxorm.NewSet()
	keySet.AddStr(table.PrimaryCacheKey.DataKeyExpression)
	keyVariableSet.AddStr(table.PrimaryCacheKey.KeyLeft)

	camelTableName := table.Name.ToCamel()
	text, err := utils.LoadTemplate(category, updateTemplateFile, template.Update)
	if err != nil {
		return "", err
	}

	output, err := utils.With("update").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": camelTableName,
			"keys":                  strings.Join(keySet.KeysStr(), "\n"),
			"primaryCacheKey":       table.PrimaryCacheKey.DataKeyExpression,
			"primaryKeyVariable":    table.PrimaryCacheKey.KeyLeft,
			"PrimaryKeyToCamel":     table.PrimaryKey.Name.ToCamel(),
			"keyValues":             strings.Join(keyVariableSet.KeysStr(), ", "),
		})
	if err != nil {
		return "", nil
	}

	return output.String(), nil
}
