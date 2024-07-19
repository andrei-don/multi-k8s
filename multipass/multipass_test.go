package multipass

import (
	"os/exec"
	"testing"
)

func TestLaunch(t *testing.T) {
	launchReq := NewLaunchReq("10G", "1G", "1", "test")
	instance, err := Launch(launchReq)

	defer func() {
		if instance != nil {
			cmd := exec.Command("multipass", "delete", "--purge", instance.Name)
			if err := cmd.Run(); err != nil {
				t.Errorf("Failed to delete instance %s: %v", instance.Name, err)
			}
		}
	}()

	if err != nil {
		t.Errorf("Multipass Launch errored with %v", err)
	}

	if instance.Name != "test" {
		t.Error("The instance launched with the wrong name.")
	}
}
