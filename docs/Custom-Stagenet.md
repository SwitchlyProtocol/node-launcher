# Custom Stagenet

Nine Realms presently runs the official `stagenet` used for release verification and ecosystem testing. While it is necessary for Nine Realms to maintain full control over the active validators for operational usage, others may desire to experiment with new chains and features in their own sandbox `stagenet` environments against mainnet chains with real funds.

The process for deploying a custom `stagenet` is straightforward, however this guidance is not a complete step-by-step. Foundational understanding of the network will be required to succeed as different configurations and underlying Kubernetes setups will have their own nuance.

## Process

1. Build a local `switchlynode` binary by running the following in the **`switchlynode`** repository from the tag of the version to be deployed:

```bash
TAG=stagenet make install
```

2. Create a new key pair to use as the master key for the network (funding faucet and admin mimir). Example:

```bash
switchlynode keys add stagenet-master
```

3. Add the following overrides in `switchlynode-stack/stagenet.yaml`:

```yaml
switchlynode:
  net: stagenet
  chainID:
    stagenet: <your-chain-id>
  env:
    # use the "sthor" address for the key created above as the faucet and mimir admin
    FAUCET: "<your-stagenet-master-address>"
    SWITCHLY_STAGENET_ADMIN_ADDRESSES: "<your-stagenet-master-address>"

    # set seed nodes endpoint empty on the genesis node
    SWITCHLY_SEED_NODES_ENDPOINT: ""

    # set to genesis node on other validators
    # SWITCHLY_SEED_NODES_ENDPOINT: "http://switchlynode.<genesis-namespace>.svc.cluster.local:1317/switchly/nodes"

    # can re-use existing router contract deployments, or deploy your own new ones
    ETH_CONTRACT: "0xB11a1735C2e3BCC5FC8c1d147fb64629d3d0caC5"
    BSC_CONTRACT: "0x00335da4078f696b98ff619616f1c558e57b9e22"
    AVAX_CONTRACT: "0xd6a6c0b3bb4150a98a379811934e440989209db6"

bifrost:
  chainID:
    stagenet: <your-chain-id>
  peer: gateway.<genesis-namespace>.svc.cluster.local # optionally use LB IP address if in cloud
  env:
    # set seed nodes endpoint empty on the genesis node
    SWITCHLY_SEED_NODES_ENDPOINT: ""

    # set to genesis node on other validators
    # SWITCHLY_SEED_NODES_ENDPOINT: "http://switchlynode.<genesis-namespace>.svc.cluster.local:1317/switchly/nodes"
```

4. Deploy the genesis node. Example:

```bash
NAME=stagenet-genesis TYPE=genesis NET=stagenet make install
```

5. After blocks are being created, deploy other validators. This is optional - you can also run a single node network with only the genesis node if you do not need to test churn or TSS signing. Example:

```bash
NAME=stagenet-validator-1 TYPE=validator NET=stagenet make install
```

6. Bond each new validator from the master wallet and perform the standard initialization commands (https://docs.switchly.org/switchlynodes/joining).

7. Set the churn interval with admin mimir to begin churning and get more active nodes. After first churn increase the interval or set `HaltChurning` to `1` to prevent migration gas waste. Example:

```bash
switchlynode tx switchly mimir ChurnInterval --from stagenet-master --chain-id <your-chain-id> --node http://localhost:27147 -- 100
```

8. Create pools and perform whatever testing you desire. Network usage docs: https://dev.switchly.org/.

## Other Notes

- You can share a single set of daemons in a separate namespace for the genesis node and all validators. See related docs in [Multi-Validator-Cluster.md](Multi-Validator-Cluster.md).

- You can run with a subset of external chains by flagging off undesired ones in `switchlynode-stack/stagenet.yaml`. Example additions to disable BTC:

```yaml
bifrost:
  env:
    BIFROST_CHAINS_BTC_DISABLED: "true"
bitcoin-daemon:
  enabled: false
```

## Switchly Stagenet CI (DigitalOcean)

The `Deploy stagenet` GitHub workflow automates this process against a DOKS cluster using
`switchlynode-stack/stagenet.yaml` — public testnets with our own in-cluster chain daemons
(bitcoin testnet3, ethereum sepolia via geth+prysm, stellar testnet via quickstart).

**One-time setup:**

1. Create the cluster (per-node ~8GB; the chain daemons need the room):

   ```bash
   doctl kubernetes cluster create switchly-stagenet \
     --region fra1 --version latest \
     --node-pool "name=default;size=s-4vcpu-8gb;count=3;auto-scale=true;min-nodes=3;max-nodes=8"
   ```

2. Repository secrets: `DIGITALOCEAN_ACCESS_TOKEN`, `STAGENET_NODE_PASSWORD`.
   Repository variables: `DOKS_CLUSTER`, `STAGENET_FAUCET_ADDRESS`, `STAGENET_ADMIN_ADDRESSES`
   (create the master wallet locally: `switchlynode keys add stagenet-master`).

3. Install prometheus CRDs + monitoring once per cluster: `make repos tools`.

**Rollout:** run the workflow with `action=install node=genesis`, wait for blocks, then install
`validator-1` … `validator-5` one at a time (six nodes total). Every node generates a FRESH
mnemonic in-cluster on first install — back each up immediately:
`make mnemonic NAME=switchly-stagenet-<node>`. Then bond each validator from the master wallet per
the standard joining flow above.

**Upgrades (no state loss):** merges to `switchlynode` main publish a new `stagenet` image; run the
workflow with `action=upgrade`. Nodes are upgraded one at a time — helm upgrade on PVC-backed state,
wait for rollout, wait until the node reports `catching_up=false` — so the chain never halts and no
history is lost. Consensus-version activation is coordinated on-chain by native version voting once
a supermajority of active nodes runs the new version.
