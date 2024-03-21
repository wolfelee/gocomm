package model

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wolfelee/gocomm/cmd/et/internal/model/datasource"
)

// CmdModel represents the mysql model.
var CmdModel = &cobra.Command{
	Use:   "model",
	Short: "generate model code",
	Long:  "generate model code. Example: et model datasource -url={datasource} -table={patterns}  -dir={dir} -cache=true",
	Run:   run,
}

func init() {
	CmdModel.AddCommand(datasource.CmdDatasource)
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("error: please input model mode")
}
