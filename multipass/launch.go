package multipass

import (
	"fmt"
	"os/exec"
)

type LaunchReq struct {
	Disk   string
	Memory string
	Cpus   string
	Name   string
	Image  string
}

func NewLaunchReq(disk, memory, cpus, name string) *LaunchReq {
	return &LaunchReq{
		Disk:   disk,
		Memory: memory,
		Cpus:   cpus,
		Name:   name,
		Image:  "jammy",
	}
}

func Launch(launchReq *LaunchReq) (*Instance, error) {
	var args = ""
	if launchReq.Image != "" {
		args = args + fmt.Sprintf(" %v", launchReq.Image)
	}
	if launchReq.Disk != "" {
		args = args + fmt.Sprintf(" --disk %v", launchReq.Disk)
	}
	if launchReq.Memory != "" {
		args = args + fmt.Sprintf(" --memory %v", launchReq.Memory)
	}
	if launchReq.Cpus != "" {
		args = args + fmt.Sprintf(" --cpus %v", launchReq.Cpus)
	}
	if launchReq.Name != "" {
		args = args + fmt.Sprintf(" --name %v", launchReq.Name)
	}

	cmd := fmt.Sprintf("multipass launch" + args)

	cmdExec := exec.Command("sh", "-c", cmd)
	stdout, err := cmdExec.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return nil, err
	}

	instance, err := InstanceInfo(&InfoReq{Name: launchReq.Name})
	if err != nil {
		return nil, err
	}

	return instance, nil
}
