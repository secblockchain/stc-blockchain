#!/usr/bin/env bash
# Start a node using host ./node-data (created by create-account.sh for validators).
# The published Docker image is not modified — only bind-mounted.
set -euo pipefail

if [ -z "${BASH_VERSION:-}" ]; then
  exec bash "$0" "$@"
fi

GREEN='\033[0;32m'
RED='\033[0;31m'
ORANGE='\033[0;33m'
NC='\033[0m'

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "${SCRIPT_DIR}/.." && pwd)"

IMAGE="${STC_IMAGE:-stc-geth:latest}"
CONTAINER_NAME="${STC_CONTAINER_NAME:-stc-node}"
DATADIR="${STC_DATADIR:-${REPO_ROOT}/node-data}"
GENESIS_FILE="${STC_GENESIS_FILE:-${REPO_ROOT}/genesis.json}"
GETH_DATADIR="/data"

GENESIS_INPUT=""
DATADIR_INPUT=""
MODE="rpc"
ENABLE_RPC=0
BOOTNODES=""
UNLOCK=""
NETWORK_ID=""
HTTP_PORT="${STC_HTTP_PORT:-8545}"
WS_PORT="${STC_WS_PORT:-8546}"
P2P_PORT="${STC_P2P_PORT:-30303}"
HTTP_API="eth,net,web3,txpool"
WS_API="eth,net,web3,txpool"
# Default geth cap (100 SEP) rejects many contract txs on this chain (~2 SEP/21k gas).
RPC_TX_FEE_CAP="${STC_RPC_TX_FEE_CAP:-0}"
NODISCOVER=0
FORCE=0
EXTRA_ARGS=()

info()  { printf "${GREEN}[INFO]${NC} %s\n" "$*"; }
warn()  { printf "${ORANGE}[WARN]${NC} %s\n" "$*"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$*" >&2; }
die()   { error "$*"; exit 1; }

usage() {
  cat <<EOF
Usage: $0 <genesis> <node-data> [OPTIONS] [-- extra-geth-flags]

Start the node with a genesis file and a host data directory.
The Docker image is not changed.

For a validator, run create-account.sh first so <node-data> already
contains the keystore.

<genesis>    genesis JSON file, or a directory that contains genesis.json
<node-data>  host data directory (keystore, password, chaindata)

Options:
  -h, --help                Show this help
  --rpc                     Enable HTTP/WS JSON-RPC (default off; combinable with --validator)
  --validator               Unlock the account in node-data and mine
  --datadir PATH            Same as the <node-data> argument
  --name NAME               Docker container name (default: ${CONTAINER_NAME})
  --image NAME[:TAG]        Docker image to use (default: ${IMAGE})
  --bootnodes ENODES        Comma-separated enode URLs of peers
  --nodiscover              Disable peer discovery
  --unlock ADDRESS          Account to unlock (default: <node-data>/address.txt)
  --networkid ID            Override network id (default: genesis config.chainId)
  --http.port PORT          Host HTTP-RPC port (default: ${HTTP_PORT})
  --ws.port PORT            Host WebSocket-RPC port (default: ${WS_PORT})
  --port PORT               Host P2P port (default: ${P2P_PORT})
  --http.api APIS           HTTP RPC APIs
  --ws.api APIS             WS RPC APIs
  --force                   Replace a running container with the same name
  --genesis                 Genesis file to use (default: ${GENESIS_FILE})
  --datadir                 Datadir to use (default: ${DATADIR})

Examples:
  sh start-node.sh
  sh start-node.sh --image registry.example.com/stc-geth:latest
  sh start-node.sh --genesis ./genesis.json --datadir ./node-data
  sh start-node.sh --validator --image stc-geth:latest
  sh start-node.sh --validator --bootnodes "enode://pubkey@1.2.3.4:30303"
EOF
}

docker_bin() {
  if docker info >/dev/null 2>&1; then
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' docker "$@"
  elif command -v sudo >/dev/null 2>&1 && sudo docker info >/dev/null 2>&1; then
    MSYS_NO_PATHCONV=1 MSYS2_ARG_CONV_EXCL='*' sudo docker "$@"
  else
    die "Docker is not available. Run: bash ${SCRIPT_DIR}/setup-node.sh"
  fi
}

docker_host_path() {
  case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*)
      cygpath -w "$1"
      ;;
    *)
      printf '%s\n' "$1"
      ;;
  esac
}

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help) usage; exit 0 ;;
      --rpc)
        ENABLE_RPC=1
        if [ "${MODE}" != "validator" ]; then
          MODE="rpc"
        fi
        ;;
      --validator)
        MODE="validator"
        HTTP_API="eth,net,web3,txpool,stcons,miner,admin,debug"
        WS_API="eth,net,web3,txpool"
        ;;
      --datadir)
        [ $# -ge 2 ] || die "--datadir requires a path"
        DATADIR_INPUT="$2"
        shift
        ;;
      --genesis)
        [ $# -ge 2 ] || die "--genesis requires a path"
        GENESIS_FILE="$2"
        shift
        ;;
      --name)
        [ $# -ge 2 ] || die "--name requires a container name"
        CONTAINER_NAME="$2"
        shift
        ;;
      --image)
        [ $# -ge 2 ] || die "--image requires NAME[:TAG]"
        IMAGE="$2"
        shift
        ;;
      --bootnodes)
        [ $# -ge 2 ] || die "--bootnodes requires an enode URL"
        BOOTNODES="$2"
        shift
        ;;
      --nodiscover) NODISCOVER=1 ;;
      --unlock)
        [ $# -ge 2 ] || die "--unlock requires an address"
        UNLOCK="$2"
        shift
        ;;
      --networkid)
        [ $# -ge 2 ] || die "--networkid requires an id"
        NETWORK_ID="$2"
        shift
        ;;
      --http.port)
        [ $# -ge 2 ] || die "--http.port requires a port"
        HTTP_PORT="$2"
        shift
        ;;
      --ws.port)
        [ $# -ge 2 ] || die "--ws.port requires a port"
        WS_PORT="$2"
        shift
        ;;
      --port)
        [ $# -ge 2 ] || die "--port requires a port"
        P2P_PORT="$2"
        shift
        ;;
      --http.api)
        [ $# -ge 2 ] || die "--http.api requires a list"
        HTTP_API="$2"
        shift
        ;;
      --ws.api)
        [ $# -ge 2 ] || die "--ws.api requires a list"
        WS_API="$2"
        shift
        ;;
      --force) FORCE=1 ;;
      --)
        shift
        EXTRA_ARGS+=("$@")
        break
        ;;
      -*) die "Unknown option: $1" ;;
      *)
        if [ -z "${GENESIS_INPUT}" ]; then
          GENESIS_INPUT="$1"
        elif [ -z "${DATADIR_INPUT}" ]; then
          DATADIR_INPUT="$1"
        else
          die "Unexpected argument: $1"
        fi
        ;;
    esac
    shift
  done
}

abspath() {
  target="$1"
  dir="$(dirname -- "${target}")"
  base="$(basename -- "${target}")"
  printf '%s/%s\n' "$(CDPATH= cd -- "${dir}" && pwd)" "${base}"
}

resolve_genesis() {
  input="$1"
  [ -n "${input}" ] || die "Genesis path is required. Example: sh start-node.sh ./genesis.json ./node-data"

  if [ -d "${input}" ]; then
    if [ -f "${input}/genesis.json" ]; then
      abspath "${input}/genesis.json"
      return
    fi
    if [ -f "${input}/genesis" ]; then
      abspath "${input}/genesis"
      return
    fi
    die "No genesis.json found in directory: ${input}"
  fi
  if [ -f "${input}" ]; then
    abspath "${input}"
    return
  fi
  if [ -f "${input}.json" ]; then
    abspath "${input}.json"
    return
  fi
  die "Genesis file not found: ${input}"
}

extract_chain_id() {
  file="$1"
  if command -v python3 >/dev/null 2>&1; then
    python3 -c 'import json,sys; print(json.load(open(sys.argv[1])).get("config",{}).get("chainId",""))' "${file}"
    return
  fi
  if command -v python >/dev/null 2>&1; then
    python -c 'import json,sys; print(json.load(open(sys.argv[1])).get("config",{}).get("chainId",""))' "${file}"
    return
  fi
  if command -v jq >/dev/null 2>&1; then
    jq -r '.config.chainId // empty' "${file}"
    return
  fi
  grep -o '"chainId"[[:space:]]*:[[:space:]]*[0-9]*' "${file}" | head -n1 | grep -o '[0-9]*$' || true
}

ensure_image() {
  if docker_bin image inspect "${IMAGE}" >/dev/null 2>&1; then
    info "Using Docker image ${IMAGE}"
    return
  fi
  info "Image ${IMAGE} not found locally — pulling"
  docker_bin pull "${IMAGE}" || die "Failed to pull ${IMAGE}. Pass --image NAME[:TAG] or run setup-node.sh."
  info "Using Docker image ${IMAGE}"
}

container_state() {
  docker_bin inspect -f '{{.State.Status}}' "${CONTAINER_NAME}" 2>/dev/null || true
}

remove_existing_container() {
  state="$(container_state)"
  [ -n "${state}" ] || return 0
  if [ "${state}" = "running" ] && [ "${FORCE}" -eq 0 ]; then
    die "Container ${CONTAINER_NAME} is already running. Stop it with stop-node.sh, or pass --force."
  fi
  info "Removing existing container ${CONTAINER_NAME} (${state})"
  docker_bin rm -f "${CONTAINER_NAME}" >/dev/null
}

normalize_address() {
  addr="$(printf '%s' "$1" | tr -d '[:space:]')"
  case "${addr}" in
    0x|0X) addr="" ;;
    0x*) addr="0x$(printf '%s' "${addr#0x}" | tr '[:upper:]' '[:lower:]')" ;;
    *) addr="0x$(printf '%s' "${addr}" | tr '[:upper:]' '[:lower:]')" ;;
  esac
  printf '%s\n' "${addr}"
}

list_keystore_addresses() {
  find "${DATADIR}/keystore" -maxdepth 1 -type f -name 'UTC--*' 2>/dev/null \
    | sed -n 's/.*--\([a-fA-F0-9]\{40\}\)$/\1/p' \
    | awk '{ print "0x" tolower($0) }'
}

keystore_address() {
  list_keystore_addresses | head -n1
}

keystore_has_address() {
  target="$(normalize_address "$1")"
  while IFS= read -r addr; do
    [ -n "${addr}" ] || continue
    if [ "${addr}" = "${target}" ]; then
      return 0
    fi
  done <<EOF
$(list_keystore_addresses)
EOF
  return 1
}

load_account_from_datadir() {
  if [ -n "${UNLOCK}" ]; then
    UNLOCK="$(normalize_address "${UNLOCK}")"
    return
  fi

  ks_addr="$(keystore_address || true)"
  if [ -n "${ks_addr}" ]; then
    UNLOCK="${ks_addr}"
  fi

  if [ -f "${DATADIR}/address.txt" ]; then
    file_addr="$(normalize_address "$(tr -d '[:space:]' < "${DATADIR}/address.txt")")"
    if [ -n "${file_addr}" ]; then
      if [ -z "${UNLOCK}" ]; then
        UNLOCK="${file_addr}"
      elif [ "${file_addr}" != "${UNLOCK}" ]; then
        warn "address.txt (${file_addr}) does not match keystore (${UNLOCK}); using keystore address"
      fi
    fi
  fi

  if [ -z "${UNLOCK}" ]; then
    UNLOCK="$(
      docker_bin run --rm \
        --entrypoint geth \
        -v "$(docker_host_path "${DATADIR}"):${GETH_DATADIR}" \
        "${IMAGE}" \
        account list --datadir "${GETH_DATADIR}" 2>/dev/null \
        | awk '{ gsub(/[{}]/, "", $3); if ($3 != "") print $3 }' \
        | head -n1
    )"
    UNLOCK="$(normalize_address "${UNLOCK}")"
  fi
}

require_validator_datadir() {
  if [ ! -d "${DATADIR}/keystore" ] || [ -z "$(ls -A "${DATADIR}/keystore" 2>/dev/null || true)" ]; then
    die "No validator keystore in ${DATADIR}. Run: bash ${SCRIPT_DIR}/create-account.sh --datadir ${DATADIR}"
  fi
  [ -f "${DATADIR}/password.txt" ] || die "No password file at ${DATADIR}/password.txt. Run create-account.sh first."
  load_account_from_datadir
  [ -n "${UNLOCK}" ] || die "Could not read the validator address from ${DATADIR}. Run create-account.sh first."
  if ! keystore_has_address "${UNLOCK}"; then
    die "No keystore key for ${UNLOCK} in ${DATADIR}/keystore. Re-run create-account.sh --datadir ${DATADIR}, or pass --unlock with an address that exists in the keystore."
  fi
  info "Using validator account ${UNLOCK} from ${DATADIR}"
}

resolve_paths() {
  if [ -n "${DATADIR_INPUT}" ]; then
    DATADIR="$(abspath "${DATADIR_INPUT}")"
  else
    DATADIR="$(abspath "${DATADIR}")"
  fi

  if [ -n "${GENESIS_INPUT}" ]; then
    GENESIS_FILE="$(resolve_genesis "${GENESIS_INPUT}")"
  else
    GENESIS_FILE="$(resolve_genesis "${GENESIS_FILE}")"
  fi
}

init_genesis() {
  genesis="$1"
  mkdir -p "${DATADIR}"

  if [ -d "${DATADIR}/geth/chaindata" ] && [ "$(ls -A "${DATADIR}/geth/chaindata" 2>/dev/null || true)" ]; then
    if [ -f "${DATADIR}/genesis.json" ] && cmp -s "${genesis}" "${DATADIR}/genesis.json"; then
      info "Chaindata already exists — skipping geth init"
      return
    fi
    warn "Genesis changed. Removing old chaindata (keystore is kept) and re-initializing."
    rm -rf "${DATADIR}/geth/chaindata" "${DATADIR}/geth/lightchaindata" "${DATADIR}/geth/nodes" "${DATADIR}/geth/LOCK"
  fi

  cp -f "${genesis}" "${DATADIR}/genesis.json"
  info "Initializing host datadir from $(basename "${genesis}")"
  docker_bin run --rm \
    --name "${CONTAINER_NAME}-init" \
    --entrypoint geth \
    -v "$(docker_host_path "${DATADIR}"):${GETH_DATADIR}" \
    -v "$(docker_host_path "${genesis}"):/genesis.json:ro" \
    "${IMAGE}" \
    init --datadir "${GETH_DATADIR}" /genesis.json
}

wait_for_rpc() {
  url="http://127.0.0.1:${HTTP_PORT}"
  info "Waiting for HTTP-RPC on ${url}"
  i=0
  while [ "${i}" -lt 30 ]; do
    if command -v curl >/dev/null 2>&1; then
      if curl -fsS -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "${url}" >/dev/null 2>&1; then
        info "RPC is up"
        return
      fi
    fi
    i=$((i + 1))
    sleep 1
  done
  warn "RPC did not respond yet. Check logs: docker logs -f ${CONTAINER_NAME}"
}

wait_for_ipc() {
  info "Waiting for IPC at /tmp/geth.ipc"
  i=0
  while [ "${i}" -lt 30 ]; do
    if ! docker_bin inspect -f '{{.State.Running}}' "${CONTAINER_NAME}" 2>/dev/null | grep -qx true; then
      warn "Container ${CONTAINER_NAME} is not running. Check logs: docker logs -f ${CONTAINER_NAME}"
      return 1
    fi
    if docker_bin exec "${CONTAINER_NAME}" test -S /tmp/geth.ipc 2>/dev/null; then
      info "IPC is up"
      return 0
    fi
    i=$((i + 1))
    sleep 1
  done
  warn "IPC did not appear yet. Check logs: docker logs -f ${CONTAINER_NAME}"
  return 1
}

print_enode() {
  enode="$(
    docker_bin exec "${CONTAINER_NAME}" \
      geth attach --exec 'admin.nodeInfo.enode' /tmp/geth.ipc 2>/dev/null || true
  )"
  # geth may print Fatal on stdout when the socket is not ready yet.
  enode="$(printf '%s' "${enode}" | tr -d '\r' | sed -n 's/^"\(enode:\/\/.*\)"$/\1/p; t; s/^\(enode:\/\/.*\)$/\1/p')"
  if [ -n "${enode}" ]; then
    info "enode: ${enode}"
  else
    warn "Could not read enode yet. Try: docker exec ${CONTAINER_NAME} geth attach --exec 'admin.nodeInfo.enode' /tmp/geth.ipc"
  fi
}

start_container() {
  genesis="$1"
  geth_args=(
    --datadir "${GETH_DATADIR}"
    --networkid "${NETWORK_ID}"
    --port 30303
    --ipcpath /tmp/geth.ipc
    --syncmode full
    --rpc.txfeecap "${RPC_TX_FEE_CAP}"
    --verbosity 3
  )

  if [ "${ENABLE_RPC}" -eq 1 ]; then
    geth_args+=(
      --http
      --http.addr 0.0.0.0
      --http.port 8545
      --http.api "${HTTP_API}"
      --http.vhosts "*"
      --http.corsdomain "*"
      --ws
      --ws.addr 0.0.0.0
      --ws.port 8546
      --ws.api "${WS_API}"
      --ws.origins "*"
    )
  fi

  if [ "${NODISCOVER}" -eq 1 ]; then
    geth_args+=(--nodiscover)
  fi
  if [ -n "${BOOTNODES}" ]; then
    geth_args+=(--bootnodes "${BOOTNODES}")
  fi

  docker_args=(
    run -d
    --name "${CONTAINER_NAME}"
    --restart unless-stopped
    --entrypoint geth
    -v "$(docker_host_path "${DATADIR}"):${GETH_DATADIR}"
    -v "$(docker_host_path "${genesis}"):/genesis.json:ro"
    -p "${P2P_PORT}:30303"
    -p "${P2P_PORT}:30303/udp"
  )

  if [ "${ENABLE_RPC}" -eq 1 ]; then
    docker_args+=(
      -p "${HTTP_PORT}:8545"
      -p "${WS_PORT}:8546"
    )
  fi

  if [ "${MODE}" = "validator" ]; then
    geth_args+=(
      --mine
      --miner.etherbase "${UNLOCK}"
      --unlock "${UNLOCK}"
      --password "${GETH_DATADIR}/password.txt"
      --allow-insecure-unlock
    )
  fi

  if [ ${#EXTRA_ARGS[@]} -gt 0 ]; then
    geth_args+=("${EXTRA_ARGS[@]}")
  fi

  docker_bin "${docker_args[@]}" "${IMAGE}" "${geth_args[@]}" >/dev/null
}

save_runtime() {
  cat > "${DATADIR}/.node.env" <<EOF
IMAGE=${IMAGE}
CONTAINER_NAME=${CONTAINER_NAME}
DATADIR=${DATADIR}
MODE=${MODE}
ENABLE_RPC=${ENABLE_RPC}
NETWORK_ID=${NETWORK_ID}
HTTP_PORT=${HTTP_PORT}
WS_PORT=${WS_PORT}
P2P_PORT=${P2P_PORT}
GENESIS=${GENESIS_FILE}
EOF
}

main() {
  parse_args "$@"
  resolve_paths

  mkdir -p "${DATADIR}"

  printf "\n${GREEN}+------------------------------------------------+\n"
  printf "+  Starting Smart Technology Chain node\n"
  printf "+  Mode:     %s\n" "${MODE}"
  printf "+  Image:    %s\n" "${IMAGE}"
  printf "+  Genesis:  %s\n" "${GENESIS_FILE}"
  printf "+  Datadir:  %s\n" "${DATADIR}"
  printf "+------------------------------------------------+${NC}\n"

  ensure_image
  remove_existing_container

  if [ "${MODE}" = "validator" ]; then
    require_validator_datadir
  fi

  init_genesis "${GENESIS_FILE}"

  if [ -z "${NETWORK_ID}" ]; then
    NETWORK_ID="$(extract_chain_id "${GENESIS_FILE}")"
  fi
  [ -n "${NETWORK_ID}" ] || die "Could not read config.chainId from genesis. Pass --networkid."

  start_container "${GENESIS_FILE}"
  save_runtime

  info "Container ${CONTAINER_NAME} started"
  if [ "${ENABLE_RPC}" -eq 1 ]; then
    info "HTTP-RPC: http://127.0.0.1:${HTTP_PORT}"
    info "WS-RPC:   ws://127.0.0.1:${WS_PORT}"
  fi
  info "P2P:      0.0.0.0:${P2P_PORT}"
  info "Logs:     docker logs -f ${CONTAINER_NAME}"

  if [ "${ENABLE_RPC}" -eq 1 ]; then
    wait_for_rpc
  else
    wait_for_ipc || true
  fi
  print_enode
}

main "$@"
