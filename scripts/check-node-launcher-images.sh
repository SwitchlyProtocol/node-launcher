#!/bin/bash
set -euo pipefail

. ./scripts/gitlab.sh

# get auth for container registry
REGISTRY="thorchain/devops/node-launcher"
TOKEN=$(gitlab_registry_token $REGISTRY)

check_image() {
  NAME="$1"
  TAG="$2"
  HASH="$3"

  if [[ $NAME == *"$REGISTRY"* ]]; then
    echo "  Checking $NAME:$TAG..."

    DIGEST=$(gitlab_registry_digest "$REGISTRY" "$TAG" "$TOKEN")

    if [[ "sha256:$HASH" != "$DIGEST" ]]; then
      echo "Hash mismatch!"
      echo "Expected: ${DIGEST}"
      echo "  Actual: sha256:${HASH}"
      exit 1
    fi
  else
    echo "  Not checking $NAME:$TAG - outside of node-launcher registry"
  fi
}

# check <chart>
check() {
  echo "Checking $1..."

  VALS="$1/values.yaml"
  NAME=$(yq -r '.image.name' "$VALS")
  TAG=$(yq -r '.image.tag' "$VALS")
  HASH=$(yq -r '.image.hash' "$VALS")

  if [[ $NAME == "null" ]]; then
    # Retry as nested maps.
    for K in $(yq -r '.image | keys[]' "$VALS"); do
      NAME=$(yq -r '.image.'$K'.name' "$VALS")
      TAG=$(yq -r '.image.'$K'.tag' "$VALS")
      HASH=$(yq -r '.image.'$K'.hash' "$VALS")
      check_image "$NAME" "$TAG" "$HASH"
    done
  else
    check_image "$NAME" "$TAG" "$HASH"
  fi
}

# check all images hosted in node-launcher registry
for CHART in *-daemon; do
  check "${CHART%/}"
done
