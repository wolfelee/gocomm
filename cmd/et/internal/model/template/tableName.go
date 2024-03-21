package template

var TableName = `
func (*{{.upperStartCamelObject}}) TableName() string  {
	return "{{.tableName}}"
}
`
