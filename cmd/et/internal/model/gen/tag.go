package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
)

func genTag(in string, isPrimaryKey, primaryKeyAutoIncr, isInsertTime, isUpdateTime bool) (string, error) {
	if in == "" {
		return in, nil
	}

	text, err := utils.LoadTemplate(category, tagTemplateFile, template.Tag)
	if err != nil {
		return "", err
	}

	output, err := utils.With("tag").Parse(text).Execute(map[string]interface{}{
		"field":              in,
		"isPrimaryKey":       isPrimaryKey,
		"primaryKeyAutoIncr": primaryKeyAutoIncr,
		"isInsertTime":       isInsertTime,
		"isUpdateTime":       isUpdateTime,
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
