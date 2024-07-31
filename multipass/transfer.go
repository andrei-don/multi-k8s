package multipass

import (
	"fmt"
	"os/exec"
)

// TransferReq serves as an input to the Transfer function. It contains the Files field which should contain the source server/file and the destination server/file as expected by the Multipass transfer command.
type TransferReq struct {
	Files string
}

// Transfer function transfers files from the local machine to a multipass VM or viceversa.
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
