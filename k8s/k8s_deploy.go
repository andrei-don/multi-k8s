package k8s

import (
	"fmt"
	"log"
	"strings"

	"github.com/andrei-don/multi-k8s/multipass"
)

const (
	BootstrapRepoRaw = "https://raw.githubusercontent.com/andrei-don/multi-k8s-provisioning-scripts/main"
)

var setupCommonScripts = []string{"setup-kernel.sh", "setup-cri.sh", "kube-components.sh"}

var setupControllerScripts = []string{"calico.yaml", "configure-single-controlplane.sh"}

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

// DownloadBootstrapScripts downloads the scripts located in the multi-k8s-provisioning-scripts repo and runs them on all nodes. It installs kubelet, kubeadm and containerd.
func DownloadAndRunBootstrapScripts(instances []*multipass.Instance) {
	var downloadCommands []string
	var runCommands []string
	for _, script := range setupCommonScripts {
		downloadCommand := fmt.Sprintf("\"wget -O /tmp/%v %v/%v\"", script, BootstrapRepoRaw, script)
		runCommand := fmt.Sprintf("\"chmod +x /tmp/%v && /tmp/%v\"", script, script)
		downloadCommands = append(downloadCommands, downloadCommand)
		runCommands = append(runCommands, runCommand)
	}

	for _, instance := range instances {

		for _, command := range downloadCommands {
			downloadBootstrapScript := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: command})
			if downloadBootstrapScript != nil {
				log.Fatal(downloadBootstrapScript)
			}
		}
		fmt.Printf("Downloaded bootstrap scripts for node %v\n", instance.Name)

		for _, command := range runCommands {
			runBootstrapScript := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: command})
			if runBootstrapScript != nil {
				log.Fatal(runBootstrapScript)
			}
		}
		fmt.Printf("Ran bootstrap scripts for node %v\n", instance.Name)
	}
}

// FilterNodes takes the list of all instance structs as inputs and returns a list of instance structs corresponding to controller or worker nodes only.
func FilterNodes(instances []*multipass.Instance, nodeType string) []*multipass.Instance {
	var controllers []*multipass.Instance
	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, nodeType) {
			controllers = append(controllers, instance)
		}
	}
	return controllers
}

// ConfigureControlPlane takes the list of controller instances structs from FilterNodes. It downloads the controlplane configuration script and calico manifest from the multi-k8s-provisioning-scripts repo.
// It generates the join-command and transfers it to the local machine.
func ConfigureControlPlane(instances []*multipass.Instance) {
	var downloadCommands []string
	controllerConfigScript := setupControllerScripts[1]
	for _, script := range setupControllerScripts {
		downloadCommand := fmt.Sprintf("\"wget -O /tmp/%v %v/k8s/%v\"", script, BootstrapRepoRaw, script)
		downloadCommands = append(downloadCommands, downloadCommand)
	}

	for _, instance := range instances {
		for _, command := range downloadCommands {
			downloadControllerConfigScript := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: command})
			if downloadControllerConfigScript != nil {
				log.Fatal(downloadControllerConfigScript)
			}
		}
		command := fmt.Sprintf("\"chmod +x /tmp/%v && /tmp/%v\"", controllerConfigScript, controllerConfigScript)
		runControllerConfigScript := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: command})
		if runControllerConfigScript != nil {
			log.Fatal(runControllerConfigScript)
		}
		fmt.Printf("Ran configuration script for controller node %v\n", instance.Name)

		transferFiles := fmt.Sprintf("%v:/tmp/join-command.sh /tmp/join-command.sh", instance.Name)
		transferCommand := multipass.Transfer(&multipass.TransferReq{Files: transferFiles})
		if transferCommand != nil {
			log.Fatal(transferCommand)
		}
		fmt.Printf("Copied join script from controller node %v to your local machine\n", instance.Name)
	}
}

// ConfigureWorkerNodes takes the list of worker instances structs from FilterNodes. It transfers the join-command script from the local machine to the worker node and runs it.
func ConfigureWorkerNodes(instances []*multipass.Instance) {
	for _, instance := range instances {
		transferFiles := fmt.Sprintf("/tmp/join-command.sh %v:/tmp/join-command.sh", instance.Name)
		transferCommand := multipass.Transfer(&multipass.TransferReq{Files: transferFiles})
		if transferCommand != nil {
			log.Fatal(transferCommand)
		}
		fmt.Printf("Copied join script from your local machine to worker node %v\n", instance.Name)
		commandJoin := "\"chmod +x /tmp/join-command.sh && sudo /tmp/join-command.sh\""
		runWorkerJoin := multipass.Exec(&multipass.ExecReq{Name: instance.Name, Script: commandJoin})
		if runWorkerJoin != nil {
			log.Fatal(runWorkerJoin)
		}
		fmt.Printf("Joined worker node %v to cluster", instance.Name)
	}
}
