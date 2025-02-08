#!/usr/bin/env bash

set -euo pipefail

update_blockstore_hashes() {
  sed -i '/thorchain-blockstore-hashes/q' $1
  curl -s "$2" | sed -e 's/^/    /' >>$1
  echo "{{- end }}" >>$1
}

# update mainnet midgard hashes
update_blockstore_hashes \
  midgard/templates/configmap-blockstore-mainnet.yaml \
  https://snapshots.ninerealms.com/snapshots/midgard-blockstore/hashes

# update stagenet midgard hashes
update_blockstore_hashes \
  midgard/templates/configmap-blockstore-stagenet.yaml \
  https://stagenet-snapshots.ninerealms.com/snapshots/midgard-blockstore/hashes
