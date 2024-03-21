package template

// Vars defines a template for var block in model
var Vars = `
{{if .withCache}}
var (
	{{.cacheKeys}}
){{end}}
`
