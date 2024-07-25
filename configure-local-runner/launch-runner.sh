#!/bin/bash

set -euo pipefail

PAT_TOKEN=$1

multipass launch --name self-hosted-runner jammy

multipass exec self-hosted-runner -- wget -O /tmp/configure-runner.sh https://raw.githubusercontent.com/andrei-don/multi-k8s/feature/add-ci/configure-local-runner/configure-runner.sh

multipass exec self-hosted-runner -- chmod +x /tmp/configure-runner.sh

multipass exec self-hosted-runner -- /tmp/configure-runner.sh ${PAT_TOKEN}



