package k8s

import (
	"fmt"
	"log"

	"github.com/andrei-don/multi-k8s/multipass"
)

const (
	BootstrapRepoRaw = "https://raw.githubusercontent.com/andrei-don/multi-k8s-provisioning-scripts/main"
)

var setupScripts = []string{"setup-kernel.sh", "setup-cri.sh", "kube-components.sh"}

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

// DownloadBootstrapScripts downloads the scripts located in the multi-k8s-provisioning-scripts repo
func DownloadBootstrapScripts(instances []*multipass.Instance) {
	var downloadCommands []string
	for _, script := range setupScripts {
		command := fmt.Sprintf("\"wget -O /tmp/%v %v/%v\"", script, BootstrapRepoRaw, script)
		downloadCommands = append(downloadCommands, command)
	}

	for _, instance := range instances {
		for _, command := range downloadCommands {
			downloadBootstrapScript := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: command})
			if downloadBootstrapScript != nil {
				log.Fatal(downloadBootstrapScript)
			}
		}
		fmt.Printf("Downloaded bootstrap scripts for node %v\n", instance.Name)
	}
}
