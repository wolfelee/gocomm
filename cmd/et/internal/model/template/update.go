package template

var Update = `
func (m *{{.upperStartCamelObject}}Model) Update(data {{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
{{end}}return m.Exec(func(session *xorm.Session) error {
		_, err := session.ID(data.{{.PrimaryKeyToCamel}}).Update(&data)
		return err
	}{{if .withCache}}, {{.keyValues}}{{end}})
}
`
