package template

var Insert = `
func (m *{{.upperStartCamelObject}}Model) Insert(data {{.upperStartCamelObject}}) error {
	return m.Exec(func(session *xorm.Session) error {
		_, err := session.Insert(data)
		return err
	})
}`
