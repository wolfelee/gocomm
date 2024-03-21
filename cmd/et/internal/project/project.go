package project

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  "Create a service project using the repository template. Example: et new hello",
	Run:   run,
}

var repoUrl string

func init() {
	CmdNew.Flags().StringVarP(&repoUrl, "-repo-url", "r",
		"https://github.com/wolfelee/goinit.git", "layout repo")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: project name is required.\033[m Example: et new helloworld\n")
		return
	}
	p := &Project{Name: args[0]}
	if err := p.New(ctx, wd, repoUrl); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		return
	}
}
