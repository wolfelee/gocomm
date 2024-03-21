package job

import (
	"github.com/wolfelee/gocomm/pkg/flag"
)

func init() {
	flag.Register(
		&flag.StringFlag{
			Name:    "job",
			Usage:   "--job",
			Default: "",
		},
	)
}

// Runner ...
type Runner interface {
	Run()
}
