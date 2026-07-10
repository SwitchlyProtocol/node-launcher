#!/usr/bin/env bash

set -e

source ./scripts/core.sh

get_node_info

if ! node_exists; then
  die "No existing SwitchlyNode found, make sure this is the correct name"
fi

source ./scripts/install.sh

echo
echo "=> Waiting for SwitchlyNode daemon to be ready"
kubectl rollout status -w deployment/switchlynode -n "$NAME" --timeout=5m

if [ "$TYPE" != "fullnode" ]; then
  echo
  source ./scripts/set-version.sh
fi
