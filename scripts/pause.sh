#!/usr/bin/env bash

set -e

source ./scripts/core.sh

get_node_info_short

echo "=> Pausing node global halt from a SwitchlyNode named $boldyellow$NAME$reset"
confirm

kubectl exec -it -n "$NAME" -c switchlynode deploy/switchlynode -- /kube-scripts/pause.sh
sleep 5
echo Switchly paused

display_status
