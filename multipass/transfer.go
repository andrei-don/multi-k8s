package multipass

import (
	"fmt"
	"os/exec"
)

type TransferReq struct {
	Files string
}

func Transfer(req *TransferReq) error {
	cmd := fmt.Sprintf("multipass transfer %v", req.Files)
	cmdExec := exec.Command("sh", "-c", cmd)

	stdout, err := cmdExec.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	return nil
}
