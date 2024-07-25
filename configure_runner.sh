#!/bin/bash

set -euo

PAT_TOKEN=$1

mkdir actions-runner && cd actions-runner

curl -o actions-runner-linux-arm64-2.317.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.317.0/actions-runner-linux-arm64-2.317.0.tar.gz

tar xzf ./actions-runner-linux-arm64-2.317.0.tar.gz

REPO_TOKEN=$(curl -L   -X POST   -H "Accept: application/vnd.github+json"   -H "Authorization: Bearer $PAT_TOKEN"   -H "X-GitHub-Api-Version: 2022-11-28"   https://api.github.com/repos/andrei-don/multi-k8s/actions/runners/registration-token | jq -r '.token')

./config.sh --unattended --url https://github.com/andrei-don/multi-k8s --token $REPO_TOKEN

./run.sh