package multipass

import (
	"fmt"
	"os/exec"
)

type DeleteReq struct {
	Name string
}

func Delete(req *DeleteReq) error {
	cmd := fmt.Sprintf("multipass delete -p %v", req.Name)
	cmdExec := exec.Command("sh", "-c", cmd)

	stdout, err := cmdExec.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	return nil
}
