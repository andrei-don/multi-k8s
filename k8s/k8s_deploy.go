package k8s

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/andrei-don/multi-k8s/multipass"
)

const (
	BootstrapRepoRaw = "https://raw.githubusercontent.com/andrei-don/multi-k8s-provisioning-scripts/main"
)

var setupCommonScripts = []string{"setup-kernel.sh", "setup-cri.sh", "kube-components.sh"}

var setupControllerScripts = []string{"calico.yaml", "configure-single-controlplane.sh"}

var setupHAControllerScripts = []string{"calico.yaml", "setup-controlplane-first.sh", "setup-secondary-controlplanes.sh", "copy-secondary-controlplane-pki.sh"}

var setupHAProxyScript string = "setup-haproxy.sh"

var setupPostDeploymentScripts = []string{"common-tasks-controlplane.sh", "approve-worker-csr.sh"}

func execCommand(req *multipass.ExecReq) {
	cmd := multipass.Exec(req)
	if cmd != nil {
		log.Fatal(cmd)
	}
}

func transferCommand(req *multipass.TransferReq) {
	cmd := multipass.Transfer(req)
	if cmd != nil {
		log.Fatal(cmd)
	}
}

func AnimateDots(done chan bool, prompt string) {
	go func() {
		dots := ""
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s%s   ", prompt, dots)
				time.Sleep(500 * time.Millisecond)
				dots += "."
				if len(dots) > 3 {
					fmt.Printf("\r%s", prompt)
					dots = ""
				}
			}
		}
	}()
}

// DeployClusterVMs deploys the VMs needed for the controller/worker nodes. It takes the input from the 'multi-k8s deploy' flags.
func DeployClusterVMs(controlNodes int, workerNodes int) []*multipass.Instance {
	fmt.Printf("\nDeploying Kubernetes cluster with %d control node(s) and %d worker node(s)...\n", controlNodes, workerNodes)

	var instances []*multipass.Instance

	for i := 1; i <= controlNodes; i++ {
		nodeName := fmt.Sprintf("controller-node-%d", i)
		done := make(chan bool)
		AnimateDots(done, fmt.Sprintf("Deploying node %v", nodeName))
		launchReq := multipass.NewLaunchReq("50G", "2G", "2", nodeName)
		instance, err := multipass.Launch(launchReq)
		if err != nil {
			log.Fatal(err)
		}
		close(done)
		fmt.Printf("\nThe IP address of %v is %v\n", nodeName, instance.IPv4)
		instances = append(instances, instance)
	}
	for i := 1; i <= workerNodes; i++ {
		nodeName := fmt.Sprintf("worker-node-%d", i)
		done := make(chan bool)
		AnimateDots(done, fmt.Sprintf("Deploying node %v", nodeName))
		launchReq := multipass.NewLaunchReq("50G", "2G", "2", nodeName)
		instance, err := multipass.Launch(launchReq)
		if err != nil {
			log.Fatal(err)
		}
		close(done)
		fmt.Printf("\nThe IP address of %v is %v\n", nodeName, instance.IPv4)
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
	fmt.Printf("\n")
	//We use escape characters for the double quotes because we would like the shell command to be enclosed in double quotes
	createHostnamesFileCmd := fmt.Sprintf("\"echo '%s' | sudo tee -a /etc/hosts\"", hostnameEntries)
	for _, instance := range instances {
		execCommand(&multipass.ExecReq{Name: instance.Name, Script: createHostnamesFileCmd})
		fmt.Printf("Added hostnames for %v\n", instance.Name)
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
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
		fmt.Printf("\nDownloaded bootstrap scripts for %v\n", instance.Name)

		done := make(chan bool)
		AnimateDots(done, fmt.Sprintf("Running bootstrap scripts for %v", instance.Name))
		for _, command := range runCommands {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
		close(done)
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
	fmt.Printf("\n")
	for _, instance := range instances {
		done := make(chan bool)
		AnimateDots(done, fmt.Sprintf("Running configuration script for %v", instance.Name))
		for _, command := range downloadCommands {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
		command := fmt.Sprintf("\"chmod +x /tmp/%v && /tmp/%v\"", controllerConfigScript, controllerConfigScript)
		execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		close(done)

		transferFiles := fmt.Sprintf("%v:/tmp/join-command-worker.sh /tmp/join-command-worker.sh", instance.Name)
		transferCommand(&multipass.TransferReq{Files: transferFiles})
		fmt.Printf("\nCopied join script from %v to your local machine\n", instance.Name)
	}
}

func DeployHAProxy(instances []*multipass.Instance) *multipass.Instance {
	var ipList string
	done := make(chan bool)
	AnimateDots(done, "Running configuration script for haproxy")
	launchReqHAProxy := multipass.NewLaunchReq("25G", "1G", "1", "haproxy")
	haproxy, err := multipass.Launch(launchReqHAProxy)
	if err != nil {
		log.Fatal(err)
	}
	haproxyIp := fmt.Sprintf("%v\n", haproxy.IPv4)
	for _, instance := range instances {
		ip := fmt.Sprintf("%v\n", instance.IPv4)
		ipList = ipList + ip
	}
	ipList = ipList + haproxyIp
	createIPListCmd := fmt.Sprintf("\"echo '%s' | sudo tee -a -i /tmp/ip_list\"", ipList)
	execCommand(&multipass.ExecReq{Name: haproxy.Name, Script: createIPListCmd})

	downloadCmd := fmt.Sprintf("\"wget -O /tmp/%v %v/%v\"", setupHAProxyScript, BootstrapRepoRaw, setupHAProxyScript)
	execCommand(&multipass.ExecReq{Name: haproxy.Name, Script: downloadCmd})

	configureHAProxyCmd := fmt.Sprintf("\"chmod +x /tmp/%v && /tmp/%v\"", setupHAProxyScript, setupHAProxyScript)
	execCommand(&multipass.ExecReq{Name: haproxy.Name, Script: configureHAProxyCmd})

	close(done)
	fmt.Println("\nDeployed HAproxy!")
	return haproxy
}

func ConfigureControlPlaneHA(instances []*multipass.Instance) {
	var downloadCommands []string

	for _, script := range setupHAControllerScripts {
		downloadCommand := fmt.Sprintf("\"wget -O /tmp/%v %v/k8s/%v\"", script, BootstrapRepoRaw, script)
		downloadCommands = append(downloadCommands, downloadCommand)
	}

	//Configuring first control node
	for _, instance := range instances[:1] {
		done := make(chan bool)
		AnimateDots(done, fmt.Sprintf("Running configuration script for %v", instance.Name))
		for _, command := range downloadCommands[:2] {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
		command := fmt.Sprintf("\"chmod +x /tmp/%v && /tmp/%v\"", setupHAControllerScripts[1], setupHAControllerScripts[1])
		execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		close(done)

		//We need to change permissions so that we can transfer the certificates
		var commandChmod []string
		commandChmod = append(commandChmod, "\"sudo chmod 644 /home/ubuntu/pki/*.key\"")
		commandChmod = append(commandChmod, "\"sudo chmod 644 /home/ubuntu/pki/*.pub\"")
		commandChmod = append(commandChmod, "\"sudo chmod 644 /home/ubuntu/pki/etcd/*.key\"")
		for _, command := range commandChmod {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}

		var transferFiles []string
		transferFiles = append(transferFiles, fmt.Sprintf("%v:/tmp/join-command-controller.sh /tmp/join-command-controller.sh", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("%v:/tmp/join-command-worker.sh /tmp/join-command-worker.sh", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("--recursive %v:/home/ubuntu/pki /tmp", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("%v:/home/ubuntu/admin.conf /tmp", instance.Name))

		for _, transferFile := range transferFiles {
			transferCommand(&multipass.TransferReq{Files: transferFile})
		}
		fmt.Printf("\n%v finished provisioning!\n", instance.Name)

		//Reverting permissions
		var commandChmodRevert []string
		commandChmodRevert = append(commandChmodRevert, "\"sudo chmod 600 /home/ubuntu/pki/*.key\"")
		commandChmodRevert = append(commandChmodRevert, "\"sudo chmod 600 /home/ubuntu/pki/*.pub\"")
		commandChmodRevert = append(commandChmodRevert, "\"sudo chmod 600 /home/ubuntu/pki/etcd/*.key\"")
		for _, command := range commandChmodRevert {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
	}

	//Configuring secondary control nodes
	for _, instance := range instances[1:] {
		done := make(chan bool)
		AnimateDots(done, fmt.Sprintf("Running configuration script for %v", instance.Name))
		for _, command := range downloadCommands[2:] {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
		commandEtcd := "\"sudo mkdir /home/ubuntu/etcd && sudo chown -R ubuntu:ubuntu /home/ubuntu/etcd\""
		execCommand(&multipass.ExecReq{Name: instance.Name, Script: commandEtcd})

		var transferFiles []string
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/admin.conf %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/ca.crt %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/ca.key %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/sa.pub %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/sa.key %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/front-proxy-ca.crt %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/front-proxy-ca.key %v:/home/ubuntu", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/etcd/ca.crt %v:/home/ubuntu/etcd", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/pki/etcd/ca.key %v:/home/ubuntu/etcd", instance.Name))
		transferFiles = append(transferFiles, fmt.Sprintf("/tmp/join-command-controller.sh %v:/tmp/", instance.Name))

		for _, transferFile := range transferFiles {
			transferCommand(&multipass.TransferReq{Files: transferFile})
		}

		//Reverting permissions
		var commandChmodRevert []string
		commandChmodRevert = append(commandChmodRevert, "\"sudo chmod 600 /home/ubuntu/*.key\"")
		commandChmodRevert = append(commandChmodRevert, "\"sudo chmod 600 /home/ubuntu/*.pub\"")
		commandChmodRevert = append(commandChmodRevert, "\"sudo chmod 600 /home/ubuntu/etcd/*.key\"")
		for _, command := range commandChmodRevert {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}

		var execCommands []string
		execCommands = append(execCommands, "\"chmod +x /tmp/copy-secondary-controlplane-pki.sh && sudo /tmp/copy-secondary-controlplane-pki.sh\"")
		execCommands = append(execCommands, "\"sudo /tmp/join-command-controller.sh\"")
		execCommands = append(execCommands, "\"chmod +x /tmp/setup-secondary-controlplanes.sh && /tmp/setup-secondary-controlplanes.sh\"")

		for _, command := range execCommands {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}

		close(done)
		fmt.Printf("%v finished provisioning!\n", instance.Name)
	}
}

// ConfigureWorkerNodes takes the list of worker instances structs from FilterNodes. It transfers the join-command script from the local machine to the worker node and runs it.
func ConfigureWorkerNodes(instances []*multipass.Instance) {
	for _, instance := range instances {
		transferFiles := fmt.Sprintf("/tmp/join-command-worker.sh %v:/tmp/join-command-worker.sh", instance.Name)
		transferCommand(&multipass.TransferReq{Files: transferFiles})
		fmt.Printf("\nCopied join script from your local machine to %v\n", instance.Name)

		commandJoin := "\"chmod +x /tmp/join-command-worker.sh && sudo /tmp/join-command-worker.sh\""
		execCommand(&multipass.ExecReq{Name: instance.Name, Script: commandJoin})
		fmt.Printf("Joined %v to cluster\n", instance.Name)
	}
}

// ConfigurePostDeploy adds autocomplete and aliases for kubectl
func ConfigurePostDeploy(instances []*multipass.Instance) {
	var downloadCommands []string
	var runCommands []string
	for _, script := range setupPostDeploymentScripts {
		downloadCommand := fmt.Sprintf("\"wget -O /tmp/%v %v/%v\"", script, BootstrapRepoRaw, script)
		runCommand := fmt.Sprintf("\"chmod +x /tmp/%v && /tmp/%v\"", script, script)
		downloadCommands = append(downloadCommands, downloadCommand)
		runCommands = append(runCommands, runCommand)
	}

	for _, instance := range instances {
		for _, command := range downloadCommands {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}

		for _, command := range runCommands {
			execCommand(&multipass.ExecReq{Name: instance.Name, Script: command})
		}
		fmt.Printf("Added kubectl autocomplete and kubectl related aliases to %v\n", instance.Name)
	}
}

func PostDeployCleanup() {
	cmdRemovePKI := exec.Command("rm", "-rf", "/tmp/pki")
	if err := cmdRemovePKI.Run(); err != nil {
		fmt.Printf("Command failed with %v!", err)
	}
	cmdRemoveAdminConf := exec.Command("rm", "-rf", "/tmp/admin.conf")
	if err := cmdRemoveAdminConf.Run(); err != nil {
		fmt.Printf("Command failed with %v!", err)
	}
	cmdRemoveJoinWorker := exec.Command("rm", "-rf", "/tmp/join-command-worker.sh")
	if err := cmdRemoveJoinWorker.Run(); err != nil {
		fmt.Printf("Command failed with %v!", err)
	}
	cmdRemoveJoinController := exec.Command("rm", "-rf", "/tmp/join-command-controller.sh")
	if err := cmdRemoveJoinController.Run(); err != nil {
		fmt.Printf("Command failed with %v!", err)
	}
	fmt.Println("Cleaned up local files!")
}
