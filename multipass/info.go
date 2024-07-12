package multipass

import (
	"os/exec"
	"strings"
)

const (
	Name    = "Name:"
	State   = "State:"
	IPv4    = "IPv4:"
	Release = "Release:"
)

type Instance struct {
	Name    string
	State   string
	IPv4    string
	Release string
}

type InfoReq struct {
	Name string
}

func InstanceOutput(stdout string) *Instance {
	var instance Instance
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, Name) {
			// The combination between TrimSpace and ReplaceAll is used to replace "Name:" with an empty string and to eliminate the spaces,
			// therefore keeping only the value we are interested in
			instance.Name = strings.TrimSpace(strings.ReplaceAll(line, Name, ""))
		}
		if strings.Contains(line, State) {
			instance.State = strings.TrimSpace(strings.ReplaceAll(line, State, ""))
		}
		if strings.Contains(line, IPv4) {
			instance.IPv4 = strings.TrimSpace(strings.ReplaceAll(line, IPv4, ""))
		}
		if strings.Contains(line, Release) {
			instance.Release = strings.TrimSpace(strings.ReplaceAll(line, Release, ""))
		}

	}
	return &instance
}

func InstanceInfo(req *InfoReq) (*Instance, error) {
	infoCmd := "multipass info " + req.Name
	cmdExec := exec.Command("sh", "-c", infoCmd)

	stdout, err := cmdExec.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return InstanceOutput(string(stdout)), nil
}
