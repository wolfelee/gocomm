package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/jxorm"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
	"strings"
)

func genFields(fields []*jxorm.Field, goTable jxorm.GoTable) (string, error) {
	var list []string

	for _, field := range fields {
		var isPrimaryKey bool
		var primaryKeyAutoIncr bool
		if field.Name.Source() == goTable.PrimaryKey.Field.Name.Source() {
			isPrimaryKey = true
			primaryKeyAutoIncr = goTable.PrimaryKey.AutoIncrement
		}
		result, err := genField(field, isPrimaryKey, primaryKeyAutoIncr)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}

func genField(field *jxorm.Field, isPrimaryKey, primaryKeyAutoIncr bool) (string, error) {
	var (
		isUpdateTime bool
		isInsertTime bool
	)
	if field.MakeType == "updateTime" {
		isUpdateTime = true
	} else if field.MakeType == "insertTime" {
		isInsertTime = true
	}

	tag, err := genTag(field.Name.Source(), isPrimaryKey, primaryKeyAutoIncr, isInsertTime, isUpdateTime)
	if err != nil {
		return "", err
	}

	text, err := utils.LoadTemplate(category, fieldTemplateFile, template.Field)
	if err != nil {
		return "", err
	}

	output, err := utils.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"name":       field.Name.ToCamel(),
			"type":       field.DataType,
			"tag":        tag,
			"hasComment": field.Comment != "",
			"comment":    field.Comment,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
