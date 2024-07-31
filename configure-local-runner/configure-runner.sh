#!/bin/bash

set -euo pipefail

PAT_TOKEN=$1

if [ -d ~/actions-runner ]; then 
   cd ~/actions-runner
   REPO_TOKEN=$(curl -L   -X POST   -H "Accept: application/vnd.github+json"   -H "Authorization: Bearer ${PAT_TOKEN}"   -H "X-GitHub-Api-Version: 2022-11-28"   https://api.github.com/repos/andrei-don/multi-k8s/actions/runners/registration-token | jq -r '.token')
   ./config.sh --unattended --url https://github.com/andrei-don/multi-k8s --token $REPO_TOKEN
   ./svc.sh install
   ./svc.sh start
else
    mkdir -p ~/actions-runner
    cd ~/actions-runner

    # Download the latest runner package
    RUNNER_VERSION=$(curl -s https://api.github.com/repos/actions/runner/releases/latest | jq -r '.tag_name')
    RUNNER_VERSION_WITHOUT_V=$(echo $RUNNER_VERSION | sed 's/^v//')
    curl -o actions-runner-osx-x64-${RUNNER_VERSION_WITHOUT_V}.tar.gz -L https://github.com/actions/runner/releases/download/${RUNNER_VERSION}/actions-runner-osx-x64-${RUNNER_VERSION_WITHOUT_V}.tar.gz

    # Extract the runner package
    tar xzf ./actions-runner-osx-x64-${RUNNER_VERSION_WITHOUT_V}.tar.gz

    # Get the registration token
    REPO_TOKEN=$(curl -L   -X POST   -H "Accept: application/vnd.github+json"   -H "Authorization: Bearer ${PAT_TOKEN}"   -H "X-GitHub-Api-Version: 2022-11-28"   https://api.github.com/repos/andrei-don/multi-k8s/actions/runners/registration-token | jq -r '.token')

    # Configure the runner
    ./config.sh --unattended --url https://github.com/andrei-don/multi-k8s --token $REPO_TOKEN

    # Enabling the service to run on server

    ./svc.sh install
    ./svc.sh start
fi