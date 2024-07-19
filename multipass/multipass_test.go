package multipass

import (
	"os/exec"
	"testing"
)

func TestLaunch(t *testing.T) {
	launchReq := NewLaunchReq("10G", "1G", "1", "testLaunch")
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
		t.Errorf("multipass.Launch errored with %v", err)
	}

	if instance.Name != "testLaunch" {
		t.Error("The instance launched with the wrong name.")
	}
}

func TestDelete(t *testing.T) {
	cmd := exec.Command("multipass", "launch", "--name", "testDelete")
	if err := cmd.Run(); err != nil {
		t.Error("Failed to launch instance")
	}

	err := Delete(&DeleteReq{Name: "testDelete"})
	if err != nil {
		t.Error("Failed to delete instance!")
	}
}
