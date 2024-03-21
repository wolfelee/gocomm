package template

var New = `
func New{{.upperStartCamelObject}}Model(dbName{{if .withCache}}, cacheName{{end}} string, ops cache.Options) *{{.upperStartCamelObject}}Model {
	return &{{.upperStartCamelObject}}Model{
		cache.NewDataBase(dbName, {{if .withCache}}cacheName{{else}}""{{end}}, ops), "{{.table}}",
	}
}`
