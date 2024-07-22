# multi-k8s

![Alt text](multi-k8s.png)

multi-k8s is a Go CLI which uses Multipass to deploy k8s clusters on your MacOS. It uses a suite of shell scripts hosted in https://github.com/andrei-don/multi-k8s-provisioning-scripts together with kubeadm to provision single/highly available control-plane clusters.

The CLI is based on the Cobra framework (https://cobra.dev/).

### Pre-requisites

Make sure to:
- have Go installed on your Mac (https://go.dev/doc/install)
- have latest Multipass version installed on your Mac (https://multipass.run/docs/install-multipass)

### Installation

Install it using the commands below:

```
# Use go install to download, build, and install the binary
go install github.com/andrei-don/multi-k8s@latest

# Ensure the Go bin directory is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Verify the installation
multi-k8s --help
```

You can enable autocomplete by adding the autocomplete script to your shell of choice. Example below for zsh:
```
echo "source <(multi-k8s completion zsh)" | tee -a ~/.zshrc
source ~/.zshrc
```