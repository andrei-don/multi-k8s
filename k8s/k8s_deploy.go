package k8s

import (
	"fmt"
	"log"

	"github.com/andrei-don/multi-k8s/multipass"
)

// DeployClusterVMs deploys the VMs needed for the controller/worker nodes. It takes the input from the 'multi-k8s deploy' flags.
func DeployClusterVMs(controlNodes int, workerNodes int) []*multipass.Instance {
	fmt.Printf("Deploying Kubernetes cluster with %d control node(s) and %d worker node(s)...\n", controlNodes, workerNodes)

	var instances []*multipass.Instance

	for i := 1; i <= controlNodes; i++ {
		nodeName := fmt.Sprintf("controller-node-%d", i)
		fmt.Printf("Deploying node %v\n", nodeName)
		launchReq := multipass.NewLaunchReq("50G", "2G", "2", nodeName)
		instance, err := multipass.Launch(launchReq)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The IP address of %v is %v\n", nodeName, instance.IPv4)
		instances = append(instances, instance)
	}
	for i := 1; i <= workerNodes; i++ {
		nodeName := fmt.Sprintf("worker-node-%d", i)
		fmt.Printf("Deploying node %v\n", nodeName)
		launchReq := multipass.NewLaunchReq("50G", "2G", "2", nodeName)
		instance, err := multipass.Launch(launchReq)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The IP address of %v is %v\n", nodeName, instance.IPv4)
		instances = append(instances, instance)
	}

	return instances
}

// CreateHostnameFile takes the list of multipass.Instance structs from DeployClusterVMs and creates the hostnames file on each instance.
func CreateHostnamesFile(instances []*multipass.Instance) {
	var hostnameEntries string
	for _, instance := range instances {
		hostnameEntry := fmt.Sprintf("%v %v\n", instance.IPv4, instance.Name)
		hostnameEntries = hostnameEntries + hostnameEntry
	}
	//We use escape characters for the double quotes because we would like the shell command to be enclosed in double quotes
	createHostnamesFileCmd := fmt.Sprintf("\"echo '%s' | sudo tee -a /etc/hosts\"", hostnameEntries)
	for _, instance := range instances {
		writeHostnamesFileCmd := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: createHostnamesFileCmd})
		if writeHostnamesFileCmd != nil {
			log.Fatal(writeHostnamesFileCmd)
		}
		fmt.Printf("Added hostnames for node %v\n", instance.Name)
	}
}
