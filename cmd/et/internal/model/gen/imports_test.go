package gen

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"github.com/wolfelee/gocomm/cmd/et/internal/utils"
	"testing"
)

func TestLoadTimplate(t *testing.T) {
	text, err := utils.LoadTemplate(category, importsTemplateFile, template.Imports)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(text)
}
