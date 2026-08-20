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

Default image: `stc-geth:latest` (override with `--image` or `STC_IMAGE`).

## RPC node

```bash
bash scrtips/setup-node.sh --image registry.example.com/stc-geth:latest
bash scrtips/start-node.sh ./genesis.json ./node-data --image registry.example.com/stc-geth:latest
```

## Validator

1. Pull the image:

```bash
bash scrtips/setup-node.sh --image registry.example.com/stc-geth:latest
```

2. Create an account (keystore is stored in `./node-data`):

```bash
bash scrtips/create-account.sh --datadir ./node-data --image registry.example.com/stc-geth:latest
```

3. Put the printed `0x...` address into genesis:
   - Clique `extraData` (32-byte vanity + signer address + 65-byte seal)
   - `alloc` / `coinbase` if that account should hold funds

   This geth build uses **Clique**, not `congress`.

4. Start:

```bash
bash scrtips/start-node.sh ./genesis.json ./node-data --validator --image registry.example.com/stc-geth:latest
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
  --image registry.example.com/stc-geth:latest \
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
