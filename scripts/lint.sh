#!/usr/bin/env bash
set -euo pipefail

get_image_versions() {
  local CONF="$1"
  local NET="$1"
  (
    pushd thornode-stack/
    helm dependency build
    popd
  ) &>/dev/null
  helm template --values thornode-stack/"$CONF".yaml \
    --set "global.net=$NET" \
    --set "midgard.enabled=true" thornode-stack/ |
    grep -E '^\s*image:\s*[^\s]+'
}

check_charts() {
  local NET="$1"

  # Check for k8s definitions that aren't using explicit hashes.
  UNCHAINED=$(get_image_versions "$NET" | grep -v sha256 || true)

  if [ -n "$UNCHAINED" ]; then
    cat <<EOF
[ERR] Some container images are specified without an explicit hash in config $NET:

$UNCHAINED

EOF
    exit 1
  fi
}

command -v kubeconform &>/dev/null || go install github.com/yannh/kubeconform/cmd/kubeconform@v0.6.7

for NET in stagenet mainnet; do
  check_charts "$NET"
done

./scripts/trunk check --no-fix --upstream origin/master

# Lint the Helm charts.
find . -type f -name 'Chart.yaml' -printf '%h\n' |
  while read -r CHART_DIR; do
    pushd "$CHART_DIR"
    printf "Helm lint of %s\n" "$CHART_DIR"
    helm lint .
    printf "Kubeconform of %s\n" "$CHART_DIR"
    # If this triggers an issue for an unknown object type, check if a JSON schema
    # already exists for it in https://github.com/datreeio/CRDs-catalog
    helm template . | kubeconform -schema-location default -schema-location https://raw.githubusercontent.com/datreeio/CRDs-catalog/refs/heads/main/monitoring.coreos.com/servicemonitor_v1.json
    popd
  done

# Check thornode-stack with the various net configs.
for NET in stagenet mainnet; do
  helm lint --values thornode-stack/"$NET".yaml thornode-stack/
done
