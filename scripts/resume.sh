#!/usr/bin/env bash

set -e

source ./scripts/core.sh

get_node_info_short

echo "=> Resuming node global halt from a SwitchlyNode named $boldyellow$NAME$reset"
confirm

kubectl exec -it -n "$NAME" -c switchlynode deploy/switchlynode -- /kube-scripts/resume.sh
sleep 5
echo Switchly resumed

display_status
