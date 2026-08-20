# Node scripts

These scripts run a Smart Technology Chain node with a **published Docker image**.
Accounts, passwords, and chain data stay on the host in a data directory.
Nothing is written into the image.

| Script | Purpose |
| --- | --- |
| `setup-node.sh` | Install Docker if needed and pull (or build) the image |
| `create-account.sh` | Create a validator keystore on the host and print the address |
| `start-node.sh` | Init genesis (once) and start the node |
| `stop-node.sh` | Stop the container; host data is kept |

Published image: [`ghcr.io/secblockchain/stc-geth`](https://ghcr.io/secblockchain/stc-geth)

```bash
docker pull ghcr.io/secblockchain/stc-geth:latest
```
## Run without scripts

Host data stays in `./node-data`. The image is not modified.

Use the same Docker invocation as `start-node.sh`. Git Bash rewrites any argument that starts with `/` (`/root/.ethereum`, `/genesis.json`, `/tmp/geth.ipc`) unless you disable that on the `docker` command itself. `export` is not enough.

Run from the **repo root**:

```bash
IMAGE=ghcr.io/secblockchain/stc-geth:latest
mkdir -p node-data

# same as start-node.sh
docker() { MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' command docker "$@"; }
if command -v cygpath >/dev/null 2>&1; then
  DATADIR="$(cygpath -w "$PWD/node-data")"
  GENESIS="$(cygpath -w "$PWD/genesis.json")"
else
  DATADIR="$PWD/node-data"
  GENESIS="$PWD/genesis.json"
fi
```

Create a validator account, then put the printed `0x...` address into genesis Clique `extraData`:

```bash
printf 'changeme\n' > node-data/password.txt
docker run --rm -it \
  -v "${DATADIR}:/root/.ethereum" \
  "$IMAGE" account new --password /root/.ethereum/password.txt
docker run --rm \
  -v "${DATADIR}:/root/.ethereum" \
  "$IMAGE" account list
```

Init genesis **before the first start**. If you skip this, geth writes Ethereum mainnet (`Chain ID: 1`). If you already started without init, delete `node-data/geth/` (keep `keystore/` and `password.txt`) and init again:

```bash
docker run --rm \
  -v "${DATADIR}:/root/.ethereum" \
  -v "${GENESIS}:/genesis.json:ro" \
  "$IMAGE" init /genesis.json
```

RPC node:

```bash
docker run -d --name stc-node --restart unless-stopped \
  -v "${DATADIR}:/root/.ethereum" \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 -p 30303:30303/udp \
  "$IMAGE" \
  --datadir /root/.ethereum \
  --networkid 19516 \
  --http --http.addr 0.0.0.0 --http.port 8545 \
  --http.api eth,net,web3,txpool --http.vhosts "*" --http.corsdomain "*" \
  --ws --ws.addr 0.0.0.0 --ws.port 8546 \
  --ws.api eth,net,web3,txpool --ws.origins "*" \
  --port 30303 --ipcpath /tmp/geth.ipc --syncmode full
```

Validator (replace `0xYOUR_ADDRESS` with the account from `account list`):

```bash
docker run -d --name stc-node --restart unless-stopped \
  -v "${DATADIR}:/root/.ethereum" \
  -p 8545:8545 -p 8546:8546 -p 30303:30303 -p 30303:30303/udp \
  "$IMAGE" \
  --datadir /root/.ethereum \
  --networkid 19516 \
  --http --http.addr 0.0.0.0 --http.port 8545 \
  --http.api eth,net,web3,txpool,clique,miner,admin --http.vhosts "*" --http.corsdomain "*" \
  --ws --ws.addr 0.0.0.0 --ws.port 8546 \
  --ws.api eth,net,web3,txpool --ws.origins "*" \
  --port 30303 --ipcpath /tmp/geth.ipc --syncmode full \
  --mine --miner.etherbase 0xYOUR_ADDRESS \
  --unlock 0xYOUR_ADDRESS --password /root/.ethereum/password.txt \
  --allow-insecure-unlock \
  --bootnodes "enode://pubkey1@192.168.1.10:30303,enode://pubkey2@192.168.1.11:30303,enode://pubkey3@192.168.1.12:30303"
```

```bash
docker logs -f stc-node
docker exec stc-node geth attach --exec 'admin.nodeInfo.enode' /tmp/geth.ipc
docker stop stc-node
docker rm stc-node
```

## RPC node (scripts)

Use `--image ghcr.io/secblockchain/stc-geth:latest` (or `STC_IMAGE`).

```bash
bash scrtips/setup-node.sh --image ghcr.io/secblockchain/stc-geth:latest
bash scrtips/start-node.sh ./genesis.json ./node-data --image ghcr.io/secblockchain/stc-geth:latest
```

## Validator

1. Pull the image:

```bash
bash scrtips/setup-node.sh --image ghcr.io/secblockchain/stc-geth:latest
```

2. Create an account (keystore is stored in `./node-data`):

```bash
bash scrtips/create-account.sh --datadir ./node-data --image ghcr.io/secblockchain/stc-geth:latest
```

3. Put the printed `0x...` address into genesis:
   - Clique `extraData` (32-byte vanity + signer address + 65-byte seal)
   - `alloc` / `coinbase` if that account should hold funds

   This geth build uses **Clique**, not `congress`.

4. Start:

```bash
bash scrtips/start-node.sh ./genesis.json ./node-data --validator --image ghcr.io/secblockchain/stc-geth:latest
```

From `scrtips/`:

```bash
./create-account.sh --datadir ../node-data --image stc-geth:latest
./start-node.sh ../genesis.json ../node-data --validator --image stc-geth:latest
```

## Connect a peer

On node A, after start, copy the `enode://...` URL (replace `127.0.0.1` with the real IP).

On node B, use the same genesis and a **different** data directory:

```bash
bash scrtips/start-node.sh ./genesis.json ./node-data-2 \
  --image ghcr.io/secblockchain/stc-geth:latest \
  --bootnodes "enode://pubkey1@192.168.1.10:30303,enode://pubkey2@192.168.1.11:30303,enode://pubkey3@192.168.1.12:30303"
```

If the node is already running:

```bash
docker exec stc-node geth attach --exec 'admin.addPeer("enode://pubkey@192.168.1.10:30303")' /tmp/geth.ipc
docker exec stc-node geth attach --exec 'admin.peers' /tmp/geth.ipc
```

## Ports

| Port | Use |
| --- | --- |
| 8545 | HTTP JSON-RPC |
| 8546 | WebSocket RPC |
| 30303 | P2P (TCP/UDP) |

```bash
docker logs -f stc-node
bash scrtips/stop-node.sh
bash scrtips/stop-node.sh --rm
```

## Notes

- `--image` must match on `create-account` and `start-node`.
- Host data directory contains `keystore/`, `password.txt`, `address.txt`, and `geth/`.
- If you change `genesis.json`, `start-node.sh` wipes chaindata and re-inits. The keystore is kept.
- Maintainers can build locally: `bash scrtips/setup-node.sh --build`.
