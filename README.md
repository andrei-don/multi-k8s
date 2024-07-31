# multi-k8s

![Alt text](multi-k8s.png)

multi-k8s is a Go CLI which uses Multipass to deploy k8s clusters on your MacOS. It uses a suite of shell scripts hosted in https://github.com/andrei-don/multi-k8s-provisioning-scripts together with kubeadm to provision single/highly available control-plane clusters.

The CLI is based on the Cobra framework (https://cobra.dev/).

### Pre-requisites

Make sure to:
- have latest Multipass version installed on your Mac (https://multipass.run/docs/install-multipass)

### Installation

#### Using brew

```
# Add the tap
brew tap andrei-don/tap

# Install the formula
brew install andrei-don/tap/multi-k8s

# Verify the installation
multi-k8s --help
```

#### Using go install

```
# Use go install to download, build, and install the binary
go install github.com/andrei-don/multi-k8s@latest

# Ensure the Go bin directory is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Verify the installation
multi-k8s --help
```

#### Using the release binaries

Go to the Releases section of this repo and download the latest binary for your OS (currently just ARM64 for MacOS and Linux).

```
# Make the binary executable
chmod +x multi-k8s

# Add it to a directory within your PATH.
mv multi-k8s /usr/local/bin

# Verify the installation
multi-k8s --help
```

### How to use

You can enable autocomplete by adding the autocomplete script to your shell of choice. Example below for zsh:
```
echo "source <(multi-k8s completion zsh)" | tee -a ~/.zshrc
source ~/.zshrc
```

During the bootstrap process, you will be asked if you would like the CLI to create/replace your kubeconfig file. If you agree to it, you will have cluster admin access from your local machine and can interact with the cluster straightaway. You will need to have kubectl installed.

If you do not want the CLI to change your kubeconfig file in case you already have one, you can shell into the controller node and run your kubectl commands from there:

```
$ multipass shell controller-node-1
```
