package multipass

import (
	"fmt"
	"os/exec"
)

// ExecReq has the fields needed as input for the Exec function.
type ExecReq struct {
	Name   string
	Script string
}

// Exec function executes commands on the Multipass VM.
func Exec(req *ExecReq) error {
	cmd := fmt.Sprintf("multipass exec %v -- sh -c %v", req.Name, req.Script)
	cmdExec := exec.Command("sh", "-c", cmd)

	stdout, err := cmdExec.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	return nil
}
