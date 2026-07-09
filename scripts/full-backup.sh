#!/usr/bin/env bash

source ./scripts/core.sh

get_node_info_short

echo "Performing snapshot and backup for SwitchlyNode and Bifrost..."
confirm

if snapshot_available; then
  make_snapshot switchlynode
  make_snapshot bifrost
else
  warn "Snapshot not available in this cluster, performing backup only..."
  echo
fi

make_backup switchlynode
make_backup bifrost
