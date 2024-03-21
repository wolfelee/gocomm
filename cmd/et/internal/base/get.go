package base

import (
	"fmt"
	"os"
	"os/exec"
)

// GoGet go get path.
func GoGet(path ...string) error {
	isGte, err := VersionGte116()
	if err != nil {
		return err
	}
	for _, p := range path {
		var cmd *exec.Cmd
		if isGte {
			fmt.Printf("go install %s\n", p+"@latest")
			cmd = exec.Command("go", "install", p+"@latest")
		} else {
			fmt.Printf("go get -u %s\n", p)
			cmd = exec.Command("go", "get", "-u", p)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func VersionGte116() (bool, error) {
	cacheOut, err := exec.Command("go", "version").Output()
	if err != nil {
		return false, err
	}
	return string(cacheOut) >= "go version go1.16", nil
}
