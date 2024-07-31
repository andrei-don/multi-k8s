package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func runCommand(t *testing.T, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		t.Logf("Command failed with output: %s", out.String())
		return "", err
	}
	return out.String(), nil
}

func mockReadInput(input string) func() string {
	return func() string {
		return input
	}
}

func TestDeployCmdSingleControler(t *testing.T) {

	_, err := runCommand(t, "which", "multipass")
	if err != nil {
		t.Fatal("Multipass is not installed. Please install it and try again.")
	}

	_, err = runCommand(t, "which", "kubectl")
	if err != nil {
		t.Fatal("kubectl is not installed. Please install it and try again.")
	}

	// create a new root command
	rootCmd := &cobra.Command{Use: "multi-k8s"}
	rootCmd.AddCommand(deployCmd)

	output := &bytes.Buffer{}
	deployCmd.SetOut(output)
	deployCmd.SetErr(output)

	// Answering with 'y' to the interactive prompts because we would like the CLI to replace our local kubeconfig file so that we can test our kubectl commands.
	mockInputFunc := mockReadInput("y")
	deployer := NewDeployer(mockInputFunc)
	deployCmd.Run = func(cmd *cobra.Command, args []string) {
		deployer.RunDeploy(cmd, args)
	}

	rootCmd.SetArgs([]string{"deploy", "--control-nodes", "1", "--worker-nodes", "1"})
	err = rootCmd.Execute()

	assert.NoError(t, err, "Deployment command failed")
	t.Log(output.String())

	defer cleanupVMs("controller-node-1", "worker-node-1")
	time.Sleep(1 * time.Minute)
	verifyClusterDeployment(t)

}

func TestDeployCmdMultiControllers(t *testing.T) {
	_, err := runCommand(t, "which", "multipass")
	if err != nil {
		t.Fatal("Multipass is not installed. Please install it and try again.")
	}

	_, err = runCommand(t, "which", "kubectl")
	if err != nil {
		t.Fatal("kubectl is not installed. Please install it and try again.")
	}
	// create a new root command
	rootCmd := &cobra.Command{Use: "multi-k8s"}
	rootCmd.AddCommand(deployCmd)

	output := &bytes.Buffer{}
	deployCmd.SetOut(output)
	deployCmd.SetErr(output)

	// Answering with 'y' to the interactive prompts because we would like the CLI to replace our local kubeconfig file so that we can test our kubectl commands.
	mockInputFunc := mockReadInput("y")
	deployer := NewDeployer(mockInputFunc)
	deployCmd.Run = func(cmd *cobra.Command, args []string) {
		deployer.RunDeploy(cmd, args)
	}

	rootCmd.SetArgs([]string{"deploy", "--control-nodes", "3", "--worker-nodes", "1"})
	err = rootCmd.Execute()

	assert.NoError(t, err, "Deployment command failed")
	t.Log(output.String())

	defer cleanupVMs("controller-node-1", "controller-node-2", "controller-node-3", "worker-node-1", "haproxy")
	time.Sleep(1 * time.Minute)
	verifyClusterDeployment(t)
}

// verifyClusterDeployment checks if the cluster nodes are in a 'Ready' state and deploys an nginx pod, checking if it is running succesfully.
func verifyClusterDeployment(t *testing.T) {
	nodesOutput, err := runCommand(t, "kubectl", "get", "nodes")
	assert.NoError(t, err, "Failed to get Kubernetes nodes")
	assert.NotContains(t, nodesOutput, "NotReady")
	t.Logf("Nodes output: %s", nodesOutput)
	podsOutput, err := runCommand(t, "kubectl", "run", "test", "--image=nginx")
	assert.NoError(t, err, "Failed to deploy nginx pod")
	assert.Contains(t, podsOutput, "pod/test created")
	time.Sleep(60 * time.Second)
	podStatusOutput, _ := runCommand(t, "kubectl", "get", "pods", "test")
	t.Logf("Pod status output: %s", podStatusOutput)
	assert.Contains(t, podStatusOutput, "Running")
}

func cleanupVMs(nodes ...string) {
	for _, node := range nodes {
		_, err := runCommand(nil, "multipass", "delete", "-p", node)
		if err != nil {
			fmt.Printf("Failed to delete node %s during cleanup.\n", node)
		}
	}
}
