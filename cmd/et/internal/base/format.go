package base

import (
	"fmt"
	"os"
	"os/exec"
)

func FormatCode(dir string, shells ...[]string) error {
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	for _, shell := range shells {
		if len(shell) < 1 {
			return fmt.Errorf("this shell is filed:%v", shell)
		}
		cmd := exec.Command(shell[0], shell[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
