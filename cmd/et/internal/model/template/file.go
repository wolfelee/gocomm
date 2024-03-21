package template

var (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"fmt"
	"github.com/wolfelee/gocomm/pkg/cache"
	{{if .time}}"time"{{end}}
	"xorm.io/xorm"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"github.com/wolfelee/gocomm/pkg/cache"
	{{if .time}}"time"{{end}}
	"xorm.io/xorm"
)
`
)
