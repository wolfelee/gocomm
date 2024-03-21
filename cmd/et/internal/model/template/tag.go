package template

var Tag = "`xorm:\"{{.field}}{{if .isPrimaryKey}} pk{{end}}{{if .primaryKeyAutoIncr}} autoincr{{end}}{{if .isInsertTime}} created{{end}}{{if .isUpdateTime}} updated{{end}}\" json:\"{{.field}}\"`"
