#!/usr/bin/env bash

set -e

source ./scripts/core.sh

get_node_info_short

echo "=> Setting SwitchlyNode keys"
kubectl exec -it -n "$NAME" -c switchlynode deploy/switchlynode -- /kube-scripts/set-node-keys.sh
sleep 5
echo SwitchlyNode Keys updated
