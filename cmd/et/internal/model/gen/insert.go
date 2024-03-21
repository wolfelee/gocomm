package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
)

func genInsert(table Table, withCache bool) (string, error) {

	camel := table.Name.ToCamel()
	text, err := utils.LoadTemplate(category, insertTemplateFile, template.Insert)
	if err != nil {
		return "", err
	}

	output, err := utils.With("insert").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": camel,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
