package multipass

import (
	"fmt"
	"os/exec"
)

type ExecReq struct {
	Name   string
	Script string
}

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
