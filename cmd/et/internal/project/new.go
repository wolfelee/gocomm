package project

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/wolfelee/gocomm/cmd/et/internal/base"
)

// Project is a project template.
type Project struct {
	Name string
}

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", p.Name)
	}
	fmt.Printf("Creating service %s, layout repo is %s\n", p.Name, layout)
	repo := base.NewRepo(layout)
	if err := repo.CopyTo(ctx, to, p.Name, []string{".git", ".github"}); err != nil {
		return err
	}
	_ = os.Rename(
		path.Join(to, "cmd", "server"),
		path.Join(to, "cmd", p.Name),
	)

	// shells := [][]string{
	// 	// {"go", "mod", "tidy"},
	// 	// {"go", "mod", "download"},
	// 	// {"gofmt", "-w", "."},
	// }

	return nil
}
