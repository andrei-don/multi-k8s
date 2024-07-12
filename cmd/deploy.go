/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/andrei-don/multi-k8s/multipass"
	"github.com/spf13/cobra"
)

// Flags
var controlNodes int
var workerNodes int

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Command for deploying a cluster. Run multi-k8s deploy -h for further details",
	Long: `Deploy a k8s cluster with specified number of control nodes and worker nodes.
	
	You can create a cluster with a single control-plane node:
	multi-k8s deploy --control-nodes 1
	
	Or a cluster with highly-available 3 control-plane nodes:
	multi-k8s deploy --control-nodes 3`,

	Run: func(cmd *cobra.Command, args []string) {
		if controlNodes != 1 && controlNodes != 3 {
			log.Fatal("Control nodes must be either 1 or 3.")
		}
		if workerNodes < 1 || workerNodes > 3 {
			log.Fatal("Worker nodes must be between 1 and 3.")
		}
		if workerNodes+controlNodes > 4 {
			log.Fatal("Cannot have more than 4 local nodes.")
		}
		deployCluster(controlNodes, workerNodes)
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

func deployCluster(controlNodes int, workerNodes int) []byte {
	fmt.Printf("Deploying Kubernetes cluster with %d control node(s) and %d worker node(s)...", controlNodes, workerNodes)
	// Create LaunchReqs with default values (including default Image)
	launchReq := multipass.NewLaunchReqs("50G", "2G", "2", "controller-node-2")
	instance, err := multipass.Launch(launchReq)
	if err != nil {
		log.Fatal(err)
	}
	return instance
}
