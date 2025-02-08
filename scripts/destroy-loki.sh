#!/usr/bin/env bash

set -e

source ./scripts/core.sh

if helm -n loki-system status promtail >/dev/null 2>&1; then
  helm get manifest promtail -n loki-system
fi
helm get manifest loki -n loki-system
echo -n "The above resources will be deleted "
confirm

if helm -n loki-system status promtail >/dev/null 2>&1; then
  echo "=> Deleting Promtail"
  helm delete promtail -n loki-system
fi
echo "=> Deleting Loki Logs Management"
helm delete loki -n loki-system
kubectl delete namespace loki-system
echo
