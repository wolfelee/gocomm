package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
)

func genNew(table Table, withCache bool) (string, error) {
	text, err := utils.LoadTemplate(category, modelNewTemplateFile, template.New)
	if err != nil {
		return "", err
	}

	output, err := utils.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"table":                 table.Name.Source(),
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
