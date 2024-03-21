package proto

import (
	"github.com/wolfelee/gocomm/cmd/et/internal/proto/add"
	"github.com/wolfelee/gocomm/cmd/et/internal/proto/build"
	"github.com/wolfelee/gocomm/cmd/et/internal/proto/server"

	"github.com/spf13/cobra"
)

// CmdProto represents the proto command.
var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files.",
	Run:   run,
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(build.CmdBuild)
	CmdProto.AddCommand(server.CmdServer)
}

func run(cmd *cobra.Command, args []string) {

}
