# Smart Energy Pay (docker-compose)

Local Blockscout stack for **Smart Energy Chain** using Docker Compose (backend + frontend + proxy + databases + optional microservices).

## Prerequisites

- Docker Engine (20.10+ recommended)
- Docker Compose v2 (`docker compose`)

## Quickstart (default: local Geth dev chain + Blockscout)

From this folder:

```bash
docker compose up -d --build
```

Then open:

- Explorer UI: `http://localhost/`
- Backend API: `http://localhost/api` (proxied)
- Postgres (Blockscout): `localhost:7432`
- Ethereum JSON-RPC (local dev geth): `http://localhost:8545` (WS: `ws://localhost:8546`)

To stop everything:

```bash
docker compose down
```

## What runs (default `docker-compose.yml`)

Core services:

- `eth`: local **geth** dev node (builds from `services/go-ethereum`)
- `db`: Postgres (published as `7432` on host)
- `redis-db`: Redis
- `backend`: Blockscout backend (port `4000` inside / also published as `4000` on host)
- `frontend`: Blockscout frontend (served via proxy)
- `proxy`: nginx reverse proxy (publishes `80`, `443`, plus extra ports used by templates)

Microservices included in this repo (enabled by default in the main compose):

- `stats` (+ `stats-db`): chain statistics service
- `visualizer`: Sol2UML visualizer
- `sig-provider`: signature provider
- `user-ops-indexer`: user ops / account abstraction indexer (service present; feature toggles live in env)

## Configuration

Environment files (edit these to customize the stack):

- Backend: `envs/common-blockscout.env`
- Frontend: `envs/common-frontend.env`
- Ethereum node: `envs/common-eth.env`
- Stats: `envs/common-stats.env`
- Visualizer: `envs/common-visualizer.env`
- User-ops-indexer: `envs/common-user-ops-indexer.env`
- Smart-contract verifier: `envs/common-smart-contract-verifier.env`
- NFT media handler: `envs/common-nft-media-handler.env`

Useful docs (upstream):

- Backend envs: `https://docs.blockscout.com/setup/env-variables`
- Frontend envs: `https://github.com/blockscout/frontend/blob/main/docs/ENVS.md`

## Using another JSON-RPC client / external node

This repo includes alternative compose files (examples: `geth.yml`, `erigon.yml`, `ganache.yml`, `hardhat-network.yml`).

When switching clients, make sure the backend is pointing to the right RPC URLs:

- `ETHEREUM_JSONRPC_HTTP_URL`
- `ETHEREUM_JSONRPC_TRACE_URL`
- `ETHEREUM_JSONRPC_WS_URL`

By default they are set in `envs/common-blockscout.env` to `http://eth:8545` / `ws://eth:8546` (the internal `eth` service).
If you want to use an RPC running on your host machine, change them to `http://host.docker.internal:8545` (and the matching WS URL).

Example:

```bash
docker compose -f geth.yml up -d
```

To stop a specific config:

```bash
docker compose -f geth.yml down
```

## Other provided compose configs

- External DB (no DB container): `docker compose -f external-db.yml up -d`
- External backend: `docker compose -f external-backend.yml up -d`
- External frontend: `docker compose -f external-frontend.yml up -d`
- Microservices only: `docker compose -f microservices.yml up -d`
- Explorer without microservices: `docker compose -f no-services.yml up -d`

## Troubleshooting

- If the explorer loads but shows no blocks/txs, the RPC URL is usually wrong. Check `envs/common-blockscout.env`.
- If you’re on Linux and connecting to a host RPC, prefer `http://host.docker.internal:8545` (requires Docker support) or expose your host RPC on `0.0.0.0`.
- If DB fails to start because of permissions, remove the `volumes/` data directory and start again (you’ll lose local data):
  - `../volumes/blockscout-db-data`
