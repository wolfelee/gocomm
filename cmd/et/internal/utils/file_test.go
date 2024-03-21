package utils

import (
	"fmt"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/template"
	"testing"
)

func TestGetTemplateDir(t *testing.T) {
	path, err := GetTemplateDir("model")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
}

func TestLoadTemplate(t *testing.T) {
	text, err := LoadTemplate("model", "var.tpl", template.Vars)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(text)
}
