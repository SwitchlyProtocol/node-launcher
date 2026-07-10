#!/usr/bin/env bash

source ./scripts/core.sh

get_node_info_short

if [[ $TYPE != "daemons" ]]; then
  node_exists || die "No existing SwitchlyNode found, make sure this is the correct name"
fi

display_status

echo -e "=> Destroying a $boldgreen$TYPE$reset SwitchlyNode on $boldgreen$NET$reset named $boldgreen$NAME$reset"
echo
echo

if [[ $TYPE == "daemons" ]]; then
  warn "!!! Make sure your daemons are not being used !!!"
  confirm
  echo "=> Deleting Daemons"
else
  warn "!!! Make sure your got your BOND back before destroying your SwitchlyNode !!!"
  confirm
  echo "=> Deleting SwitchlyNode"
fi

helm delete "$NAME" -n "$NAME"
kubectl delete namespace "$NAME"
