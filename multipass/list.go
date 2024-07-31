package multipass

import (
	"fmt"
	"os/exec"
)

// List function lists the multipass instances.
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
