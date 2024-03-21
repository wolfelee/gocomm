package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
)

func genTableName(table Table) (string, error) {
	text, err := utils.LoadTemplate(category, tableNameFile, template.TableName)
	if err != nil {
		return "", err
	}

	output, err := utils.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject": table.Name.ToCamel(),
			"tableName":             table.Name.Source(),
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil

}
