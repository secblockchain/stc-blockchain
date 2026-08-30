#!/bin/bash
set -e

STC_CONFIG=${STC_HOME}/config/config.toml
STC_GENESIS=${STC_HOME}/config/genesis.json

# Init genesis state if geth not exist
DATA_DIR=$(cat ${STC_CONFIG} | grep -A1 '\[Node\]' | grep -oP '\"\K.*?(?=\")')

GETH_DIR=${DATA_DIR}/geth
if [ ! -d "$GETH_DIR" ]; then
  geth --datadir ${DATA_DIR} init ${STC_GENESIS}
fi

exec "geth" "--config" ${STC_CONFIG} "$@"
