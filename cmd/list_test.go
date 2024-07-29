package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"

	"testing"
)

// Helper function to capture output
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestListCommand_NoClusters(t *testing.T) {
	// Mock functions
	multipassListFunction := func() (string, error) {
		return "", nil
	}
	k8sFilterNodesListCmdFunction := func(s string) string {
		return ""
	}

	// Capture the output
	output := captureOutput(func() {
		listCommand(multipassListFunction, k8sFilterNodesListCmdFunction)
	})

	expected := "You have no multi-k8s clusters!\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestListCommand_WithClusters(t *testing.T) {
	// Mock functions
	multipassListFunction := func() (string, error) {
		return "node1 node2", nil
	}
	k8sFilterNodesListCmdFunction := func(s string) string {
		return "filtered nodes"
	}

	// Capture the output
	output := captureOutput(func() {
		listCommand(multipassListFunction, k8sFilterNodesListCmdFunction)
	})

	expected := "filtered nodes\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}

func TestListCommand_Error(t *testing.T) {
	// Mock functions
	multipassListFunction := func() (string, error) {
		return "", errors.New("mock error")
	}
	k8sFilterNodesListCmdFunction := func(s string) string {
		return ""
	}

	// Capture the output
	output := captureOutput(func() {
		listCommand(multipassListFunction, k8sFilterNodesListCmdFunction)
	})

	expected := "Error listing multipass nodes: mock error\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}
