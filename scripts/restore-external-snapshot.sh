#!/usr/bin/env bash

set -e

# check the xmllint command is available
if ! command -v xmllint >/dev/null 2>&1; then
  echo "=> xmllint command not found, please install libxml2 and/or libxml2-utils"
  exit 1
fi

source ./scripts/core.sh

get_node_net
get_node_info_short

# trunk-ignore(shellcheck/SC2310)
if ! node_exists; then
  die "No existing SwitchlyNode found, make sure this is the correct name"
fi

# select snapshot provider
case "${NET}" in
  "mainnet")
    PROVIDER="https://snapshots.ninerealms.com"
    ;;
  "stagenet")
    PROVIDER="https://stagenet-snapshots.ninerealms.com"
    ;;
  *)
    die "Unsupported network: ${NET}"
    ;;
esac

read -r -p "=> Enter provider [${PROVIDER}]: " provider
PROVIDER=${provider:-${PROVIDER}}
echo

# get all available snapshot heights
HEIGHTS=$(
  set -o pipefail
  curl -s "${PROVIDER}/snapshots?prefix=switchlynode" |
    xmllint --xpath '//*[local-name()="Contents"]/*[local-name()="Key"]/text()' - |
    grep -oE '[0-9]+' |
    sort -nr |
    head -n 10
)
readarray -t HEIGHTS <<<"${HEIGHTS}"

echo "=> Select block height to recover"
# shellcheck disable=SC2068
menu "${HEIGHTS[0]}" ${HEIGHTS[@]}
HEIGHT=${MENU_SELECTED}

echo "=> Recovering snapshot at height ${HEIGHT} on SwitchlyNode in ${boldgreen}${NAME}${reset}"
confirm

# stop switchlynode
echo "stopping switchlynode..."
kubectl scale -n "${NAME}" --replicas=0 deploy/switchlynode --timeout=5m
kubectl wait --for=delete pods -l app.kubernetes.io/name=switchlynode -n "${NAME}" --timeout=5m >/dev/null 2>&1 || true

# create recover pod
echo "creating recover pod"
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: restore-external-switchlynode
  namespace: ${NAME}
spec:
  containers:
  - name: recover
    image: alpine:latest@sha256:4edbd2beb5f78b1014028f4fbb99f3237d9561100b6881aabbf5acce2c4f9454
    command:
      - tail
      - -f
      - /dev/null
    volumeMounts:
    - mountPath: /root
      name: data
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: switchlynode
EOF

# reset node state
echo "waiting for recover pod to be ready..."
kubectl wait --for=condition=ready pods/restore-external-switchlynode -n "${NAME}" --timeout=5m >/dev/null 2>&1

echo "clearing existing data directory..."
kubectl exec -n "${NAME}" -it restore-external-switchlynode -- rm -rf /root/.switchlynode/data/

echo "installing dependencies..."
kubectl exec -n "${NAME}" -it restore-external-switchlynode -- sh -c 'apk update && apk add aria2 pv'

echo "pulling snapshot..."
kubectl exec -n "${NAME}" -it restore-external-switchlynode -- aria2c \
  --split=16 --max-concurrent-downloads=16 --max-connection-per-server=16 \
  --continue --min-split-size=100M --out="/root/${HEIGHT}.tar.gz" \
  "${PROVIDER}/snapshots/switchlynode/${HEIGHT}.tar.gz"

echo "extracting snapshot..."
kubectl exec -n "${NAME}" -it restore-external-switchlynode -- sh -c "pv \"/root/${HEIGHT}.tar.gz\" | tar xzf - -C /root/.switchlynode/"

echo "removing snapshot..."
kubectl exec -n "${NAME}" -it restore-external-switchlynode -- rm -rf "/root/${HEIGHT}.tar.gz"

echo "=> ${boldgreen}Proceeding to clean up recovery pod and restart switchlynode${reset}"
confirm

echo "cleaning up recover pod"
kubectl -n "${NAME}" delete pod/restore-external-switchlynode

# start switchlynode
kubectl scale -n "${NAME}" --replicas=1 deploy/switchlynode --timeout=5m
