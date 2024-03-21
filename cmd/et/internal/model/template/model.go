package template

// Model defines a template for model
var Model = `package {{.pkg}}
{{.imports}}
{{.vars}}
{{.types}}
{{.tableName}}
{{.new}}
{{.insert}}
{{.find}}
{{.update}}
{{.delete}}
{{.extraMethod}}
`
