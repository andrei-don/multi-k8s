/*
Copyright © 2024 Alex Stan
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/andrei-don/multi-k8s/k8s"
	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

// Flags
var controlNodes int
var workerNodes int

// readInputFunc is a function type for reading input
type readInputFunc func() string

// Deployer encapsulates the deploy logic and its dependencies
type Deployer struct {
	readInput readInputFunc
}

// NewDeployer creates a new Deployer with the given input function
func NewDeployer(inputFunc readInputFunc) *Deployer {
	return &Deployer{readInput: inputFunc}
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	input = strings.TrimSpace(input)
	return input
}

func CheckMultipass() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("which", "multipass")
	default:
		fmt.Println("Unsupported operating system")
		os.Exit(1)
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("Multipass not found, download it from https://multipass.run/docs/install-multipass#install")
		os.Exit(1)
	}
}

// RunDeploy performs the deployment logic, it is a method on the Deployer type. This was needed because we needed to replace the readInput function in the E2E tests to mock user input.
func (d *Deployer) RunDeploy(cmd *cobra.Command, args []string) {
	CheckMultipass()
	if controlNodes != 1 && controlNodes != 3 {
		log.Fatal("Control nodes must be either 1 or 3.")
	}
	if workerNodes < 1 || workerNodes > 3 {
		log.Fatal("Worker nodes must be between 1 and 3.")
	}
	if workerNodes+controlNodes > 4 {
		log.Fatal("Cannot have more than 4 local nodes.")
	}
	multipassList, err := multipass.List()
	if err != nil {
		fmt.Println("Error listing multipass nodes:", err)
		return
	}
	if k8s.FilterNodesListCmd(multipassList) != "" {
		fmt.Println("There is a cluster currently running! Delete the nodes and deploy a new cluster?(y/n)")
		input := d.readInput()
		if input == "y" {
			k8s.DeleteClusterVMs(k8s.GetCurrentNodes(k8s.FilterNodesListCmd(multipassList)))
		} else {
			fmt.Println("Did not delete current cluster. Delete it if you want to deploy a new one.")
			os.Exit(0)
		}
	}
	deployedInstances := k8s.DeployClusterVMs(controlNodes, workerNodes)
	k8s.DownloadAndRunBootstrapScripts(deployedInstances)
	controllerInstances := k8s.FilterNodes(deployedInstances, "controller")
	workerInstances := k8s.FilterNodes(deployedInstances, "worker")
	if controlNodes == 1 {
		k8s.CreateHostnamesFile(deployedInstances)
		k8s.ConfigureControlPlane(controllerInstances)
		fmt.Println("The script is about to replace the contents of your ~/.kube/config file. If you have other entries from other clusters that you still want to connect to, please do not proceed. Do you want to proceed (y/n)?")
		input := d.readInput()
		if input == "y" {
			k8s.CreateLocalAdmin(0)
		} else {
			fmt.Println("Did not add a kubeconfig file. Please shell into the controller node to get access to your cluster.")
		}
	} else {
		haproxy := k8s.DeployHAProxy(controllerInstances)
		k8s.CreateHostnamesFile(append(deployedInstances, haproxy))
		k8s.ConfigureControlPlaneHA(controllerInstances)
		fmt.Println("The script is about to replace the contents of your ~/.kube/config file. If you have other entries from other clusters that you still want to connect to, please do not proceed. Do you want to proceed (y/n)?")
		input := d.readInput()
		if input == "y" {
			k8s.CreateLocalAdmin(1)
		} else {
			fmt.Println("Did not add a kubeconfig file. Please shell into the controller node to get access to your cluster.")
		}
	}
	k8s.ConfigureWorkerNodes(workerInstances)
	k8s.ConfigurePostDeploy(controllerInstances)
	k8s.PostDeployCleanup()
}

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Command for deploying a cluster",
	Long: `
	Deploy a k8s cluster with specified number of control nodes and worker nodes.
	You can deploy maximum 4 nodes.`,
	Example: `
	You can create a cluster with a single control-plane node:
	multi-k8s deploy --control-nodes 1
	
	Or a cluster with highly-available 3 control-plane nodes:
	multi-k8s deploy --control-nodes 3`,
	Run: func(cmd *cobra.Command, args []string) {
		deployer := NewDeployer(readInput)
		deployer.RunDeploy(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().IntVarP(&controlNodes, "control-nodes", "c", 1, "Number of control nodes (1 or 3)")
	deployCmd.Flags().IntVarP(&workerNodes, "worker-nodes", "w", 1, "Number of worker nodes (1-3)")
}
