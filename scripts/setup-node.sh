#!/usr/bin/env bash
# Install Docker if needed and pull the published node image.
# Does not create accounts or node-data.
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
BUILD=0

info()  { printf "${GREEN}[INFO]${NC} %s\n" "$*"; }
warn()  { printf "${ORANGE}[WARN]${NC} %s\n" "$*"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$*" >&2; }
die()   { error "$*"; exit 1; }
task()  { printf "\n${ORANGE}TASK:${NC} ${GREEN}%s${NC}\n" "$*"; }

usage() {
  cat <<EOF
Usage: $0 [OPTIONS]

Install Docker if needed and pull the published node image.
This does not create accounts, genesis, or chain data.

Options:
  -h, --help          Show this help
  --image NAME[:TAG]  Image to pull (default: ${IMAGE})
  --build             Build from this repo's Dockerfile (maintainers only)

Examples:
  bash scrtips/setup-node.sh
  bash scrtips/setup-node.sh --image registry.example.com/stc-geth:latest
EOF
}

sudo_cmd() {
  if [ "$(id -u)" -eq 0 ]; then
    "$@"
  elif command -v sudo >/dev/null 2>&1; then
    sudo "$@"
  else
    die "This step needs root privileges. Re-run as root or install sudo."
  fi
}

docker_bin() {
  if docker info >/dev/null 2>&1; then
    docker "$@"
  elif command -v sudo >/dev/null 2>&1 && sudo docker info >/dev/null 2>&1; then
    sudo docker "$@"
  else
    die "Docker is installed but not usable. Start Docker, or log out and back in after joining the docker group."
  fi
}

os_family() {
  uname_s="$(uname -s 2>/dev/null || echo unknown)"
  case "${uname_s}" in
    Linux*)
      if [ -r /etc/os-release ]; then
        # shellcheck disable=SC1091
        . /etc/os-release
        printf '%s\n' "${ID:-linux}"
      else
        printf 'linux\n'
      fi
      ;;
    Darwin*) printf 'darwin\n' ;;
    MINGW*|MSYS*|CYGWIN*) printf 'windows\n' ;;
    *) printf 'unknown\n' ;;
  esac
}

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      -h|--help) usage; exit 0 ;;
      --build) BUILD=1 ;;
      --image)
        [ $# -ge 2 ] || die "--image requires NAME[:TAG]"
        IMAGE="$2"
        shift
        ;;
      *) die "Unknown option: $1" ;;
    esac
    shift
  done
}

install_prereqs_linux() {
  task "Installing host prerequisites"
  if command -v apt-get >/dev/null 2>&1; then
    sudo_cmd apt-get update -y
    sudo_cmd apt-get install -y ca-certificates curl gnupg
  elif command -v dnf >/dev/null 2>&1; then
    sudo_cmd dnf install -y ca-certificates curl gnupg2
  elif command -v yum >/dev/null 2>&1; then
    sudo_cmd yum install -y ca-certificates curl
  else
    warn "No known package manager found. Make sure curl is installed."
  fi
}

install_docker() {
  task "Checking Docker"
  if command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1; then
    info "Docker is already installed: $(docker --version 2>/dev/null | head -n1)"
    return
  fi
  if command -v docker >/dev/null 2>&1 && command -v sudo >/dev/null 2>&1 && sudo docker info >/dev/null 2>&1; then
    info "Docker is already installed (needs sudo): $(sudo docker --version 2>/dev/null | head -n1)"
    return
  fi

  family="$(os_family)"
  case "${family}" in
    ubuntu|debian|linuxmint|pop|raspbian|fedora|centos|rhel|amzn|sles|opensuse*)
      info "Installing Docker Engine via get.docker.com"
      install_prereqs_linux
      curl -fsSL https://get.docker.com | sudo_cmd sh
      if command -v systemctl >/dev/null 2>&1; then
        sudo_cmd systemctl enable --now docker || true
      fi
      if [ "$(id -u)" -ne 0 ]; then
        sudo_cmd usermod -aG docker "$(id -un)" || true
        warn "Added $(id -un) to the docker group. Log out and back in, or run: newgrp docker"
      fi
      ;;
    darwin|windows)
      die "Docker is not running. Install Docker Desktop and start it, then re-run this script."
      ;;
    *)
      die "Unsupported OS for automatic Docker install ($(uname -s)). Install Docker manually, then re-run."
      ;;
  esac

  docker_bin info >/dev/null
  info "Docker is ready: $(docker_bin --version | head -n1)"
}

pull_image() {
  task "Pulling ${IMAGE}"
  if docker_bin pull "${IMAGE}"; then
    info "Image ready: ${IMAGE}"
    return
  fi
  die "Failed to pull ${IMAGE}. Set STC_IMAGE / --image to the published image, or run with --build from this repo."
}

build_image() {
  task "Building ${IMAGE} from Dockerfile"
  [ -f "${REPO_ROOT}/Dockerfile" ] || die "Dockerfile not found at ${REPO_ROOT}/Dockerfile"
  info "This can take several minutes on the first build."
  docker_bin build -t "${IMAGE}" -f "${REPO_ROOT}/Dockerfile" "${REPO_ROOT}"
  info "Built ${IMAGE}"
}

print_next_steps() {
  printf "\n${GREEN}Setup complete.${NC} Image: ${CYAN}%s${NC}\n\n" "${IMAGE}"
  printf "${CYAN}RPC node:${NC}\n"
  printf "  bash %s/start-node.sh ./genesis.json ./node-data\n\n" "${SCRIPT_DIR}"
  printf "${CYAN}Validator:${NC}\n"
  printf "  bash %s/create-account.sh\n" "${SCRIPT_DIR}"
  printf "  # add the printed address to genesis extraData\n"
  printf "  bash %s/start-node.sh --genesis ./genesis.json --datadir ./node-data --validator\n" "${SCRIPT_DIR}"
  printf "  bash %s/start-node.sh --genesis ./genesis.json --datadir ./node-data --validator --bootnodes \"enode://...@ip:30303\"\n\n" "${SCRIPT_DIR}"
}

main() {
  parse_args "$@"

  printf "\n${GREEN}+------------------------------------------------+\n"
  printf "+  Smart Technology Chain — node setup\n"
  printf "+  Host:  %s\n" "$(uname -s)"
  printf "+  Image: %s\n" "${IMAGE}"
  printf "+------------------------------------------------+${NC}\n"

  install_docker
  if [ "${BUILD}" -eq 1 ]; then
    build_image
  else
    pull_image
  fi
  print_next_steps
}

main "$@"