package multipass

import (
	"log"
	"os/exec"
	"testing"
)

func launchInstance(name string) {
	cmd := exec.Command("multipass", "launch", "--name", name)
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to launch instance")
	}
}

func deleteInstance(name string) {
	cmd := exec.Command("multipass", "delete", "--purge", name)
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to delete instance")
	}
}

func TestLaunch(t *testing.T) {
	launchReq := NewLaunchReq("10G", "1G", "1", "testLaunch")
	instance, err := Launch(launchReq)

	defer func() {
		if instance != nil {
			deleteInstance("testLaunch")
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
	launchInstance("testDelete")
	err := Delete(&DeleteReq{Name: "testDelete"})
	if err != nil {
		t.Error("Failed to delete instance!")
	}
}

func TestExec(t *testing.T) {
	launchInstance("testExec")
	defer func() {
		deleteInstance("testExec")
	}()
	err := Exec(&ExecReq{Name: "testExec", Script: "\"echo 'testExec' | sudo tee -a -i /tmp/test\""})
	if err != nil {
		t.Error("Failed to run exec command!")
	}
	cmd := exec.Command("multipass", "exec", "testExec", "--", "sh", "-c", "cat /tmp/test")
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Failed to launch instance! got %v", err)
	}
	if string(output) != "testExec\n" {
		t.Errorf("The test failed! Expected %v, got %v instead.", "testExec", string(output))
	}
}

func TestTransfer(t *testing.T) {
	launchInstance("testTransfer")
	defer func() {
		deleteInstance("testTransfer")
	}()

	cmd := exec.Command("touch", "/tmp/testTransfer")
	if err := cmd.Run(); err != nil {
		t.Error("Failed to create local file!")
	}

	err := Transfer(&TransferReq{Files: "/tmp/testTransfer testTransfer:/tmp/testTransfer"})
	if err != nil {
		t.Error("Failed to transfer file.")
	}

	cmdCheck := exec.Command("multipass", "exec", "testTransfer", "--", "sh", "-c", "test -f /tmp/testTransfer && echo \"success\" || echo \"failure\"")
	output, err := cmdCheck.Output()
	if err != nil {
		t.Error("Failed to execute command on instance.")
	}
	defer func() {
		cmdCleanup := exec.Command("rm", "/tmp/testTransfer")
		if err := cmdCleanup.Run(); err != nil {
			t.Errorf("Failed to remove local test file: %v", err)
		}
	}()

	if string(output) != "success\n" {
		t.Error("Could not transfer file.")
	}
}
