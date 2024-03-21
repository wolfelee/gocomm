package template

var Types = `
type (
	{{.upperStartCamelObject}}Model struct {
		cache.DataBase
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)`
