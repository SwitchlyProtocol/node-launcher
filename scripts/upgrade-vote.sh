#!/usr/bin/env bash

source ./scripts/core.sh

get_node_info_short

echo "=> Voting on SwitchlyNode upgrade proposal"

get_upgrade_proposal_name() {
  read -r -p "=> Enter SwitchlyNode upgrade proposal name: [$UPGRADE_PROPOSAL_NAME] " upgrade_proposal_name
  UPGRADE_PROPOSAL_NAME=${upgrade_proposal_name:-$UPGRADE_PROPOSAL_NAME}
  echo
}

get_upgrade_proposal_vote() {
  echo "=> Select SwitchlyNode upgrade proposal vote: "
  menu yes yes no
  UPGRADE_PROPOSAL_VOTE=$MENU_SELECTED
  echo
}

get_upgrade_proposal_name
get_upgrade_proposal_vote

kubectl exec -it -n "$NAME" -c switchlynode deploy/switchlynode -- /kube-scripts/retry.sh /kube-scripts/upgrade-vote.sh "$UPGRADE_PROPOSAL_NAME" "$UPGRADE_PROPOSAL_VOTE"
