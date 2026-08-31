#!/usr/bin/env bash
# Create a validator account using the published Docker image.
# Keystore is stored on the host in ./node-data — the image is not modified.
set -euo pipefail

if [ -z "${BASH_VERSION:-}" ]; then
  exec bash "$0" "$@"
fi

GREEN='\033[0;32m'
RED='\033[0;31m'
ORANGE='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
REPO_ROOT="$(CDPATH= cd -- "${SCRIPT_DIR}/.." && pwd)"

IMAGE="${STC_IMAGE:-stc-geth:latest}"
DATADIR="${STC_DATADIR:-${REPO_ROOT}/node-data}"
GETH_DATADIR="/data"
PASSWORD_FILE=""
FORCE=0

info()  { printf "${GREEN}[INFO]${NC} %s\n" "$*"; }
warn()  { printf "${ORANGE}[WARN]${NC} %s\n" "$*"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$*" >&2; }
die()   { error "$*"; exit 1; }

usage() {
  cat <<EOF
Usage: $0 [OPTIONS]

Create a validator account with the Docker image (geth account new).
The keystore is written to host ./node-data. The image itself is not changed.

Add the printed address to genesis, then run start-node.sh --validator.

Options:
  -h, --help          Show this help
  --datadir PATH      Host data directory (default: ${DATADIR})
  --password FILE     Password file (default: <datadir>/password.txt)
  --image NAME[:TAG]  Docker image to use (default: ${IMAGE})
  --force             Create another account even if one already exists

Examples:
  bash scrtips/create-account.sh
  bash scrtips/create-account.sh --image registry.example.com/stc-geth:latest
  bash scrtips/create-account.sh --datadir ./node-data --image stc-geth:latest
EOF
}

docker_bin() {
  # Git Bash rewrites /root/... into C:/Program Files/Git/root/... unless disabled.
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

run_geth() {
  docker_bin run --rm \
    --entrypoint geth \
    -v "$(docker_host_path "${DATADIR}"):${GETH_DATADIR}" \
    "${IMAGE}" \
    "$@"
}

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help) usage; exit 0 ;;
      --datadir)
        [ $# -ge 2 ] || die "--datadir requires a path"
        DATADIR="$2"
        shift
        ;;
      --password)
        [ $# -ge 2 ] || die "--password requires a file"
        PASSWORD_FILE="$2"
        shift
        ;;
      --image)
        [ $# -ge 2 ] || die "--image requires NAME[:TAG]"
        IMAGE="$2"
        shift
        ;;
      --force) FORCE=1 ;;
      *) die "Unknown option: $1" ;;
    esac
    shift
  done
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

list_accounts() {
  if [ -d "${DATADIR}/keystore" ]; then
    find "${DATADIR}/keystore" -maxdepth 1 -type f -name 'UTC--*' 2>/dev/null \
      | sed -n 's/.*--\([a-fA-F0-9]\{40\}\)$/\1/p' \
      | awk '{print "0x" tolower($0)}'
    return
  fi
  run_geth account list --datadir "${GETH_DATADIR}" 2>/dev/null \
    | awk '{ gsub(/[{}]/, "", $3); if ($3 != "") print $3 }'
}

first_account() {
  list_accounts | head -n1
}

write_password() {
  target="${DATADIR}/password.txt"
  if [ -n "${PASSWORD_FILE}" ] && [ -f "${PASSWORD_FILE}" ]; then
    cp -f "${PASSWORD_FILE}" "${target}"
  elif [ ! -f "${target}" ]; then
    if [ ! -t 0 ]; then
      die "No password file at ${target}. Create it first, or run this in a terminal."
    fi
    printf "Enter password for the validator account: " >&2
    stty -echo 2>/dev/null || true
    IFS= read -r password
    stty echo 2>/dev/null || true
    printf '\n' >&2
    [ -n "${password}" ] || die "Password cannot be empty"
    printf '%s\n' "${password}" > "${target}"
  fi
  chmod 600 "${target}" 2>/dev/null || true
}

save_pointer() {
  cat > "${DATADIR}/.node.env" <<EOF
IMAGE=${IMAGE}
DATADIR=${DATADIR}
EOF
}

main() {
  parse_args "$@"
  ensure_image

  mkdir -p "${DATADIR}/keystore"
  chmod 700 "${DATADIR}" 2>/dev/null || true

  existing="$(first_account || true)"
  created=0
  if [ -n "${existing}" ] && [ "${FORCE}" -eq 0 ]; then
    address="${existing}"
    info "Account already exists in ${DATADIR}"
  else
    write_password
    info "Creating account with Docker image ${IMAGE}"
    info "Keystore will be saved on the host at ${DATADIR}"
    run_geth account new --datadir "${GETH_DATADIR}" --password "${GETH_DATADIR}/password.txt"
    address="$(first_account)"
    [ -n "${address}" ] || die "Failed to create account with image ${IMAGE}"
    created=1
  fi

  printf '%s\n' "${address}" > "${DATADIR}/address.txt"
  save_pointer

  printf "\n${GREEN}+------------------------------------------------+\n"
  if [ "${created}" -eq 1 ]; then
    printf "+  Validator account created via Docker image\n"
  else
    printf "+  Validator account already exists on this host\n"
  fi
  printf "+------------------------------------------------+${NC}\n\n"
  printf "Image:   ${CYAN}%s${NC}\n" "${IMAGE}"
  printf "Address (add this to genesis extraData):\n\n"
  printf "  ${CYAN}%s${NC}\n\n" "${address}"
  printf "Host datadir: %s\n" "${DATADIR}"
  printf "Keystore:     %s/keystore\n" "${DATADIR}"
  printf "Saved to:     %s/address.txt\n\n" "${DATADIR}"
  printf "Next:\n"
  printf "  1. Put the address into genesis extraData\n"
  printf "  2. bash %s/start-node.sh ./genesis.json ./node-data --validator\n" "${SCRIPT_DIR}"
  printf "     bash %s/start-node.sh ./genesis.json ./node-data --validator --bootnodes \"enode://...@ip:30303\"\n\n" "${SCRIPT_DIR}"
}

main "$@"