package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// CmdBuild represents the source command.
	CmdBuild = &cobra.Command{
		Use:                "build",
		Short:              "build the proto build ",
		Long:               "build the proto. Example: et proto build helloworld.proto",
		DisableFlagParsing: true,
		Run:                run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	var (
		err   error
		proto = strings.TrimSpace(args[0])
	)
	if err = look("protoc-gen-go", "protoc-gen-go-grpc"); err != nil {
		// update the Et plugins
		cmd := exec.Command("et", "upgrade")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if strings.HasSuffix(proto, ".proto") {
		err = Generate(proto, args)
	} else {
		err = walk(proto, args)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" {
			return nil
		}
		return Generate(path, args)
	})
}

// Generate is used to execute the generate command for the specified proto file
func Generate(proto string, args []string) error {
	path, name := filepath.Split(proto)
	input := []string{
		"--proto_path=.",
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		name,
	}
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = path
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}
