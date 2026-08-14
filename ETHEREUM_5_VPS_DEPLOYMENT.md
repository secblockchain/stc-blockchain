# Deploy five Ethereum Mainnet full nodes on five VPSs

This runbook deploys one independent Ethereum Mainnet full node to each of five VPSs. Each full node consists of:

- one Geth execution client built from this repository's `services/go-ethereum` source;
- one Lighthouse consensus client;
- private, loopback-only HTTP and WebSocket JSON-RPC endpoints.

It does not deploy Blockscout, PostgreSQL, Redis, a frontend, a validator client, or staking keys.

The deployment uses `ethereum-mainnet-full-node.yml` from this repository. Run the same Compose file independently on all five VPSs. Because each VPS has its own public IP and Docker daemon, they can all use the standard Ethereum ports.

## 1. Prepare the five VPSs

This runbook assumes Ubuntu 24.04. Each VPS should have:

- 4 or more CPU cores;
- 16 GB or more RAM;
- a 2 TB or larger NVMe SSD;
- an unmetered connection with at least 25 Mbit/s;
- a public IPv4 address;
- accurate system time.

Plan the hosts as follows:

| Node | Example hostname | Execution RPC | Consensus API |
| --- | --- | --- | --- |
| 1 | `eth-node-1` | `127.0.0.1:8545` | `127.0.0.1:5052` |
| 2 | `eth-node-2` | `127.0.0.1:8545` | `127.0.0.1:5052` |
| 3 | `eth-node-3` | `127.0.0.1:8545` | `127.0.0.1:5052` |
| 4 | `eth-node-4` | `127.0.0.1:8545` | `127.0.0.1:5052` |
| 5 | `eth-node-5` | `127.0.0.1:8545` | `127.0.0.1:5052` |

Use SSH or the VPS provider's browser console to access each server. SSH is a management method only; Ethereum does not use SSH for peer-to-peer communication.

On each VPS, verify that the large SSD is mounted on the filesystem that contains `/var/lib/docker`:

```bash
nproc
free -h
lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINTS
df -h /var/lib/docker 2>/dev/null || df -h /
```

Do not start synchronization until Docker's data directory has at least 2 TB available. If the VPS provider supplies the NVMe disk as a separate unmounted device, follow the provider's storage-mounting instructions before installing Docker. Do not format a device that already contains data.

Enable time synchronization:

```bash
sudo timedatectl set-ntp true
timedatectl status
```

## 2. Configure the provider firewall

Create these inbound rules on every VPS:

| Port | Protocol | Source | Purpose |
| --- | --- | --- | --- |
| SSH port | TCP | Your administration IP | Server management |
| `30303` | TCP and UDP | Anywhere | Geth P2P and discovery |
| `9000` | TCP and UDP | Anywhere | Lighthouse P2P and discovery |
| `9001` | UDP | Anywhere | Lighthouse QUIC |

Do not allow public inbound access to:

- `8545` (Geth HTTP JSON-RPC);
- `8546` (Geth WebSocket JSON-RPC);
- `8551` (authenticated Engine API);
- `5052` (Lighthouse Beacon API).

Allow outbound internet traffic. The Compose file additionally binds the user-facing RPC ports to `127.0.0.1`, and it does not publish the Engine API.

If UFW is enabled on Ubuntu, configure it on each VPS:

```bash
sudo ufw allow OpenSSH
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp
sudo ufw allow 9000/tcp
sudo ufw allow 9000/udp
sudo ufw allow 9001/udp
sudo ufw enable
sudo ufw status
```

Keep the provider firewall rules even when UFW is enabled. Docker-published ports can interact with host firewall rules differently from ordinary host processes.

## 3. Install Docker on every VPS

Run on each VPS:

```bash
sudo apt update
sudo apt upgrade -y
sudo apt install -y docker.io docker-compose-v2 openssl curl ca-certificates git
sudo systemctl enable --now docker
sudo docker version
sudo docker compose version
sudo docker info --format '{{.DockerRootDir}}'
```

If `docker-compose-v2` is unavailable from the Ubuntu repository, install Docker Engine and the Compose plugin using Docker's official Ubuntu installation instructions.

## 4. Put this repository on every VPS

The whole repository is required because the standalone Compose file builds Geth from `services/go-ethereum`. Do not start the repository's default `docker-compose.yml`, because it is the Smart Energy Chain development/explorer stack rather than an Ethereum Mainnet node.

Commit and push `ethereum-mainnet-full-node.yml` and this runbook before cloning. On each VPS, clone the same reviewed commit:

```bash
sudo install -d -o "$USER" -g "$USER" /opt/core-blockchain-v2
git clone --branch develop YOUR_REPOSITORY_URL /opt/core-blockchain-v2
cd /opt/core-blockchain-v2
git rev-parse HEAD
```

Use the same commit hash on all five hosts. If the repository is private, configure a read-only deploy key or use your approved HTTPS credential. As an alternative, copy the complete repository directory to each VPS; copying only the Compose file is not sufficient.

This checkout identifies its Geth source as `1.17.2-unstable`. Execution-client code is consensus-critical: test the exact commit on node 1 before rolling it to nodes 2-5, and keep it current for Ethereum network upgrades. The repository's standard Dockerfile builds the normal Linux binary and does not enable the optional `ziren` build tag.

## 5. Generate a unique JWT secret on every VPS

The JWT secret authenticates communication between that VPS's Geth and Lighthouse containers. It is not a wallet key and does not hold ETH, but it should remain private.

Run independently on every VPS:

```bash
sudo install -d -m 0750 /opt/core-blockchain-v2/secrets
sudo sh -c 'umask 077; openssl rand -hex 32 > /opt/core-blockchain-v2/secrets/jwt.hex'
sudo test -s /opt/core-blockchain-v2/secrets/jwt.hex
sudo stat /opt/core-blockchain-v2/secrets/jwt.hex
```

Do not copy one JWT secret to all five VPSs. Generate a separate value on every host.

## 6. Validate and start each node

Run on each VPS:

```bash
cd /opt/core-blockchain-v2
sudo docker compose -f ethereum-mainnet-full-node.yml config --quiet
sudo docker compose -f ethereum-mainnet-full-node.yml build --pull geth
sudo docker compose -f ethereum-mainnet-full-node.yml pull lighthouse
sudo docker compose -f ethereum-mainnet-full-node.yml run --rm --no-deps geth version
sudo docker compose -f ethereum-mainnet-full-node.yml up -d
sudo docker compose -f ethereum-mainnet-full-node.yml ps
```

The service list should contain only:

```text
geth
lighthouse
```

Both should show a running state. The clients restart automatically after a VPS reboot because the Compose services use `restart: unless-stopped` and Docker is enabled at boot.

## 7. Inspect startup logs

Geth logs:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml logs -f --tail=100 geth
```

Lighthouse logs:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml logs -f --tail=100 lighthouse
```

Press `Ctrl+C` to stop following logs. This does not stop the containers.

Temporary execution-endpoint errors from Lighthouse can occur while Geth is starting. The clients should reconnect automatically. Persistent JWT authentication errors usually mean the containers are not reading the same `/opt/core-blockchain-v2/secrets/jwt.hex` file.

## 8. Check synchronization

Check Geth on each VPS:

```bash
curl -s http://127.0.0.1:8545 \
  -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}'
```

During synchronization, `result` is an object containing progress values. When execution synchronization is complete, `result` becomes `false`.

Check Lighthouse:

```bash
curl -s http://127.0.0.1:5052/eth/v1/node/syncing
```

The node is ready when Lighthouse reports `is_syncing: false`, `el_offline: false`, and a small or zero `sync_distance`.

Check the current execution block number:

```bash
curl -s http://127.0.0.1:8545 \
  -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

Check Geth peer count:

```bash
curl -s http://127.0.0.1:8545 \
  -H 'Content-Type: application/json' \
  --data '{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}'
```

The peer count is hexadecimal. For example, `0x19` is 25 peers.

Check Lighthouse peer count:

```bash
curl -s http://127.0.0.1:5052/eth/v1/node/peer_count
```

Initial synchronization can take hours or days depending on SSD performance and bandwidth. All five VPSs may synchronize concurrently.

## 9. Understand how the five nodes connect

Geth and Lighthouse automatically discover Ethereum Mainnet peers. The five nodes do not need to be manually connected to each other in order to synchronize or serve RPC.

Optionally, create direct Geth connections among the five nodes. First retrieve the enode URL from each VPS:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml \
  exec -T geth geth attach \
  --exec 'admin.nodeInfo.enode' \
  /root/.ethereum/geth.ipc
```

If the returned enode URL contains a Docker address such as `172.18.0.2`, replace only that address with the VPS's reachable public IP or private-network IP while retaining the public key and port `30303`:

```text
enode://PUBLIC_KEY@REACHABLE_VPS_IP:30303
```

On another VPS, open the local Geth console:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml \
  exec geth geth attach /root/.ethereum/geth.ipc
```

Add the remote node:

```javascript
admin.addPeer("enode://PUBLIC_KEY@REACHABLE_VPS_IP:30303")
```

Verify the direct connection:

```javascript
admin.peers
```

Exit the console:

```javascript
exit
```

A simple optional ring is sufficient: node 1 tracks node 2, node 2 tracks node 3, node 3 tracks node 4, node 4 tracks node 5, and node 5 tracks node 1. All five continue to connect to normal Mainnet peers. Lighthouse manages its own consensus-layer peer discovery; direct Lighthouse peering is not required.

## 10. Access RPC remotely without making it public

The HTTP RPC endpoint is bound to `127.0.0.1:8545` on every VPS. From an administration computer, create an SSH tunnel to a specific node:

```bash
ssh -N -L 18545:127.0.0.1:8545 USER@NODE_1_IP
```

Then use `http://127.0.0.1:18545` locally. Use different local ports for simultaneous tunnels:

| Node | Suggested local tunnel endpoint |
| --- | --- |
| 1 | `http://127.0.0.1:18545` |
| 2 | `http://127.0.0.1:28545` |
| 3 | `http://127.0.0.1:38545` |
| 4 | `http://127.0.0.1:48545` |
| 5 | `http://127.0.0.1:58545` |

If an application server must access these RPC endpoints continuously, use a private VPS network or VPN and apply explicit firewall restrictions. Do not change the binding to a public `0.0.0.0:8545` endpoint without an authenticated, rate-limited RPC gateway.

## 11. Monitor disk space and service state

Run periodically on each VPS:

```bash
df -h /var/lib/docker 2>/dev/null || df -h /
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml ps
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml logs --tail=100
```

Do not allow the disk to become completely full. Forced termination during database writes can require a lengthy recovery or resynchronization.

## 12. Restart, stop, and update safely

Restart one node:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml restart
```

Stop it gracefully:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml stop
```

Start it again:

```bash
sudo docker compose -f /opt/core-blockchain-v2/ethereum-mainnet-full-node.yml start
```

`docker compose down` removes the containers and network but preserves the named database volumes. Never add `--volumes` or `-v` unless intentionally deleting the synchronized databases and accepting a complete resync.

Lighthouse is pinned in `ethereum-mainnet-full-node.yml`, while Geth is built from the checked-out repository commit. Before Ethereum network upgrades and after security releases, review both clients' release notes, merge and test the required upstream Geth changes, then roll out the reviewed repository commit. On the first VPS run:

```bash
cd /opt/core-blockchain-v2
git pull --ff-only
sudo docker compose -f ethereum-mainnet-full-node.yml build --pull geth
sudo docker compose -f ethereum-mainnet-full-node.yml pull lighthouse
sudo docker compose -f ethereum-mainnet-full-node.yml up -d
```

Upgrade and verify one VPS first, then roll the same tested version across the other four nodes.

## 13. Final checklist for every VPS

- Docker and Compose are installed and enabled at boot.
- `/var/lib/docker` is backed by a 2 TB or larger NVMe filesystem.
- `/opt/core-blockchain-v2/services/go-ethereum` contains the reviewed Geth source.
- All five VPSs use the same repository commit.
- `/opt/core-blockchain-v2/ethereum-mainnet-full-node.yml` exists.
- `/opt/core-blockchain-v2/secrets/jwt.hex` exists and is unique to the VPS.
- Geth and Lighthouse containers are running.
- P2P firewall ports are open.
- RPC, Engine API, and Beacon API ports are not public.
- Geth has execution peers.
- Lighthouse has consensus peers.
- Geth reports `eth_syncing: false` after initial synchronization.
- Lighthouse reports `is_syncing: false` and `el_offline: false`.

This setup provides five independent, non-validating, pruned Ethereum Mainnet full nodes. It does not provide complete archive-state history and it does not earn staking rewards.
