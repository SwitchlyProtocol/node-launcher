#!/usr/bin/env bash

set -e

source ./scripts/core.sh

# if the existing install is for loki-stack, warn and abort
if helm -n loki-system list -o json | jq -e '.[]|select(.chart|startswith("loki-stack"))' >/dev/null 2>&1; then
  echo '=> Deprecated Loki installation detected. Please "make destroy-loki" to uninstall it first.'
  die "NOTE: Log history will be reset on new deployment."
fi

########################################################################################
# Loki
########################################################################################

if helm -n loki-system status loki >/dev/null 2>&1; then
  helm diff -C 3 upgrade loki grafana/loki --install -n loki-system -f ./loki/values.yaml
  confirm
fi

echo "=> Installing Loki Logs Management"
helm upgrade loki grafana/loki --install -n loki-system --create-namespace -f ./loki/values.yaml

echo Waiting for services to be ready...
kubectl wait --for=condition=Ready --all pods -n loki-system --timeout=5m
echo

########################################################################################
# Promtail
########################################################################################

if helm -n loki-system status promtail >/dev/null 2>&1; then
  helm diff -C 3 upgrade promtail grafana/promtail --install -n loki-system -f ./promtail/values.yaml
  confirm
fi

echo "=> Installing Promtail"
helm upgrade promtail grafana/promtail --install -n loki-system --wait -f ./promtail/values.yaml
echo Waiting for services to be ready...
kubectl wait --for=condition=Ready --all pods -n loki-system --timeout=5m
