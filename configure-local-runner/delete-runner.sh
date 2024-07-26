#!/bin/bash

set -euo pipefail

PAT_TOKEN=$1

RUNNER_ID=$(curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${PAT_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/andrei-don/multi-k8s/actions/runners | jq '.runners[0].id')

curl -L \
  -X DELETE \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${PAT_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/andrei-don/multi-k8s/actions/runners/${RUNNER_ID}

multipass delete -p self-hosted-runner