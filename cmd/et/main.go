package main

import (
	"log"

	"github.com/wolfelee/gocomm/cmd/et/internal/model"

	"github.com/spf13/cobra"
	"github.com/wolfelee/gocomm/cmd/et/internal/project"
	"github.com/wolfelee/gocomm/cmd/et/internal/proto"
	"github.com/wolfelee/gocomm/cmd/et/internal/upgrade"
)

var (
	version string = "v0.1.1"

	rootCmd = &cobra.Command{
		Use:     "et",
		Short:   "et: An elegant toolkit for Go microservices.",
		Long:    `et: An elegant toolkit for Go microservices.`,
		Version: version,
	}
)

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(model.CmdModel)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
