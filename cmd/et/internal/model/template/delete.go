package template

var Delete = `
func (m *{{.upperStartCamelObject}}Model) Delete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{.keys}}
{{end}}return m.Exec(func(session *xorm.Session) error {
		data := new({{.upperStartCamelObject}})
		_, err := session.ID({{.lowerStartCamelPrimaryKey}}).Delete(data)
		return err
	}{{if .withCache}}, {{.keyValues}}{{end}})
}`
