package template

var FindOne = `
func (m *{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) ({{.upperStartCamelObject}}, bool, error) {
	{{if .withCache}}{{.cacheKey}}
{{end}}var resp {{.upperStartCamelObject}}
	has, err := m.{{if .withCache}}QueryRow({{.cacheKeyVariable}},{{else}}QueryRowNotCache({{end}}&resp, func(session *xorm.Session, v interface{}) (bool, error) {
		return session.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Get(&resp)
	})
	if err != nil {
		return resp, false, err
	}
	return resp, has, nil
}`
