package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wolfelee/gocomm/cmd/et/internal/base"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the Et tools",
	Long:  "Upgrade the Et tools. Example: Et upgrade",
	Run:   Run,
}

// Run upgrade the Et tools.
func Run(cmd *cobra.Command, args []string) {
	err := base.GoGet(
		"github.com/wolfelee/gocomm/cmd/et",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
	)
	if err != nil {
		fmt.Println(err)
	}
}
