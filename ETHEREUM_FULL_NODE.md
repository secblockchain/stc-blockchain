# Ethereum Mainnet full node

This deployment runs only an Ethereum Mainnet full node:

- Geth execution client built from this repository's `services/go-ethereum` source
- Lighthouse consensus client
- Local HTTP and WebSocket JSON-RPC

It does not run Blockscout, a database, a frontend, or a validator client.

## VPS requirements

- Ubuntu 24.04 or another supported Linux distribution
- 4 or more CPU cores
- 16 GB or more RAM
- 2 TB or larger NVMe SSD
- Unmetered connection with at least 25 Mbit/s
- Docker Engine and Docker Compose v2

## Start the node

From the repository directory on the VPS, create the shared Engine API secret:

```bash
mkdir -p secrets
openssl rand -hex 32 > secrets/jwt.hex
chmod 600 secrets/jwt.hex
```

Build your Geth source, pull Lighthouse, and start the clients:

```bash
docker compose -f ethereum-mainnet-full-node.yml build --pull geth
docker compose -f ethereum-mainnet-full-node.yml pull lighthouse
docker compose -f ethereum-mainnet-full-node.yml up -d
docker compose -f ethereum-mainnet-full-node.yml ps
```

The standard Dockerfile builds the normal Linux Geth binary. It does not enable the optional `ziren` build tag. Confirm the custom binary before starting a long sync:

```bash
docker compose -f ethereum-mainnet-full-node.yml run --rm --no-deps geth version
```

Follow both clients' logs:

```bash
docker compose -f ethereum-mainnet-full-node.yml logs -f --tail=100
```

## Firewall

Allow the SSH port from a trusted IP and allow these Ethereum peer-to-peer ports publicly:

- `30303/tcp`
- `30303/udp`
- `9000/tcp`
- `9000/udp`
- `9001/udp`

Do not expose `8545`, `8546`, `8551`, or `5052` publicly. The Compose file binds the user-facing APIs to the VPS loopback interface and leaves the Engine API available only inside the Compose network.

## Check synchronization

Execution client:

```bash
curl -s http://127.0.0.1:8545 \
  -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}'
```

The result becomes `false` when Geth is synchronized.

Consensus client:

```bash
curl -s http://127.0.0.1:5052/eth/v1/node/syncing
```

Check the execution block number:

```bash
curl -s http://127.0.0.1:8545 \
  -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

Initial synchronization can take hours or days depending on VPS disk performance.

## Access RPC remotely

Use an SSH tunnel instead of exposing JSON-RPC publicly:

```bash
ssh -N -L 8545:127.0.0.1:8545 user@VPS_IP
```

Applications on the local computer can then use `http://127.0.0.1:8545`.

## Update clients

Review upstream Geth and Lighthouse release notes first. Merge the required Geth changes into `services/go-ethereum`, test them, and then rebuild and restart:

```bash
docker compose -f ethereum-mainnet-full-node.yml build --pull geth
docker compose -f ethereum-mainnet-full-node.yml pull lighthouse
docker compose -f ethereum-mainnet-full-node.yml up -d
```

The named Docker volumes preserve the synchronized databases during container upgrades.
