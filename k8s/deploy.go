package k8s

import (
	"fmt"
	"log"

	"github.com/andrei-don/multi-k8s/multipass"
)

func DeployClusterVMs(controlNodes int, workerNodes int) {
	fmt.Printf("Deploying Kubernetes cluster with %d control node(s) and %d worker node(s)...\n", controlNodes, workerNodes)
	for i := 1; i <= controlNodes; i++ {
		nodeName := fmt.Sprintf("controller-node-%d", i)
		fmt.Printf("Deploying node %v\n", nodeName)
		launchReq := multipass.NewLaunchReqs("50G", "2G", "2", nodeName)
		instance, err := multipass.Launch(launchReq)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The IP address of %v is %v\n", nodeName, instance.IPv4)
	}
	for i := 1; i <= workerNodes; i++ {
		nodeName := fmt.Sprintf("worker-node-%d", i)
		fmt.Printf("Deploying node %v\n", nodeName)
		launchReq := multipass.NewLaunchReqs("50G", "2G", "2", nodeName)
		instance, err := multipass.Launch(launchReq)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The IP address of %v is %v\n", nodeName, instance.IPv4)
	}
}
