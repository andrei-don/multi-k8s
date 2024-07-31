package multipass

import (
	"fmt"
	"os/exec"
)

// DeleteReq has the fields needed as input for the Delete function.
type DeleteReq struct {
	Name string
}

// Delete function deletes the Multipass VM.
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
