package multipass

import (
	"fmt"
	"os/exec"
)

func List() (string, error) {
	cmd := "multipass list"
	cmdExec := exec.Command("sh", "-c", cmd)

	stdout, err := cmdExec.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return "", err
	}
	return string(stdout), nil
}
