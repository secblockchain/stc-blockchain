#!/usr/bin/env bash
# Stop the node container. Host node-data is kept.
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

CONTAINER_NAME="${STC_CONTAINER_NAME:-stc-node}"
DATADIR="${STC_DATADIR:-${REPO_ROOT}/node-data}"
NAME_FROM_CLI=""
DATADIR_FROM_CLI=""
REMOVE=0
TIMEOUT=30

info()  { printf "${GREEN}[INFO]${NC} %s\n" "$*"; }
warn()  { printf "${ORANGE}[WARN]${NC} %s\n" "$*"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$*" >&2; }
die()   { error "$*"; exit 1; }

usage() {
  cat <<EOF
Usage: $0 [OPTIONS]

Stop the node container started by start-node.sh.
Host chain data in node-data is kept.

Options:
  -h, --help         Show this help
  --name NAME        Container name (default: ${CONTAINER_NAME})
  --datadir PATH     Data directory used to read saved runtime settings
  --rm               Also remove the stopped container
  --timeout SECONDS  Graceful stop timeout (default: ${TIMEOUT})

Examples:
  bash scrtips/stop-node.sh
  bash scrtips/stop-node.sh --rm
EOF
}

docker_bin() {
  if docker info >/dev/null 2>&1; then
    docker "$@"
  elif command -v sudo >/dev/null 2>&1 && sudo docker info >/dev/null 2>&1; then
    sudo docker "$@"
  else
    die "Docker is not available."
  fi
}

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help) usage; exit 0 ;;
      --name)
        [ $# -ge 2 ] || die "--name requires a container name"
        NAME_FROM_CLI="$2"
        CONTAINER_NAME="$2"
        shift
        ;;
      --datadir)
        [ $# -ge 2 ] || die "--datadir requires a path"
        DATADIR_FROM_CLI="$2"
        DATADIR="$2"
        shift
        ;;
      --rm) REMOVE=1 ;;
      --timeout)
        [ $# -ge 2 ] || die "--timeout requires seconds"
        TIMEOUT="$2"
        shift
        ;;
      *) die "Unknown option: $1" ;;
    esac
    shift
  done
}

load_saved_env() {
  if [ -f "${DATADIR}/.node.env" ]; then
    # shellcheck disable=SC1091
    . "${DATADIR}/.node.env"
  fi
}

container_state() {
  docker_bin inspect -f '{{.State.Status}}' "${CONTAINER_NAME}" 2>/dev/null || true
}

main() {
  parse_args "$@"
  load_saved_env
  [ -z "${DATADIR_FROM_CLI}" ] || DATADIR="${DATADIR_FROM_CLI}"
  [ -z "${NAME_FROM_CLI}" ] || CONTAINER_NAME="${NAME_FROM_CLI}"

  state="$(container_state)"
  if [ -z "${state}" ]; then
    warn "Container ${CONTAINER_NAME} does not exist."
    exit 0
  fi

  if [ "${state}" = "running" ] || [ "${state}" = "paused" ] || [ "${state}" = "restarting" ]; then
    info "Stopping ${CONTAINER_NAME} (timeout ${TIMEOUT}s)"
    docker_bin stop -t "${TIMEOUT}" "${CONTAINER_NAME}" >/dev/null
    info "Node stopped"
  else
    info "Container ${CONTAINER_NAME} is already ${state}"
  fi

  if [ "${REMOVE}" -eq 1 ]; then
    docker_bin rm "${CONTAINER_NAME}" >/dev/null
    info "Container ${CONTAINER_NAME} removed"
  fi

  info "Chain data kept at ${DATADIR}"
}

main "$@"