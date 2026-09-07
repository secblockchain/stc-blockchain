package stcons

const validatorSetABI = `
[
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "InvalidValue",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyCoinbase",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "systemContract",
        "type": "address"
      }
    ],
    "name": "OnlySystemContract",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyZeroGasPrice",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "UnknownParam",
    "type": "error"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "ParamChange",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "deprecatedDeposit",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "deprecatedFinalityRewardDeposit",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "feeBurned",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "finalityRewardDeposit",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "paramChange",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "systemTransfer",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "validatorDeposit",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "validatorEnterMaintenance",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "validatorEvicted",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "validatorExitMaintenance",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "amount",
        "type": "uint256"
      }
    ],
    "name": "validatorMinorSlashed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [],
    "name": "validatorSetUpdated",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "BLOCK_FEES_RATIO_SCALE",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "BURN_ADDRESS",
    "outputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "CODE_OK",
    "outputs": [
      {
        "internalType": "uint32",
        "name": "",
        "type": "uint32"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "INIT_BURN_RATIO",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "INIT_MAINTAIN_SLASH_SCALE",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "INIT_MAX_NUM_OF_MAINTAINING",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "INIT_SYSTEM_REWARD_RATIO",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "MAX_SYSTEM_REWARD_BALANCE",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "alreadyInit",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "burnRatio",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "index",
        "type": "uint256"
      }
    ],
    "name": "canEnterMaintenance",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "currentValidatorSet",
    "outputs": [
      {
        "internalType": "address",
        "name": "consensusAddress",
        "type": "address"
      },
      {
        "internalType": "uint64",
        "name": "votingPower",
        "type": "uint64"
      },
      {
        "internalType": "bool",
        "name": "jailed",
        "type": "bool"
      },
      {
        "internalType": "uint256",
        "name": "incoming",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "name": "currentValidatorSetMap",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "currentVoteAddrFullSet",
    "outputs": [
      {
        "internalType": "bytes",
        "name": "",
        "type": "bytes"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "valAddr",
        "type": "address"
      }
    ],
    "name": "deposit",
    "outputs": [],
    "stateMutability": "payable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address[]",
        "name": "valAddrs",
        "type": "address[]"
      },
      {
        "internalType": "uint256[]",
        "name": "weights",
        "type": "uint256[]"
      }
    ],
    "name": "distributeFinalityReward",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "enterMaintenance",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "evictValidator",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "exitMaintenance",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "getCurrentValidatorIndex",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "getIncoming",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getLivingValidators",
    "outputs": [
      {
        "internalType": "address[]",
        "name": "",
        "type": "address[]"
      },
      {
        "internalType": "bytes[]",
        "name": "",
        "type": "bytes[]"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getMiningValidators",
    "outputs": [
      {
        "internalType": "address[]",
        "name": "",
        "type": "address[]"
      },
      {
        "internalType": "bytes[]",
        "name": "",
        "type": "bytes[]"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getTurnLength",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getValidators",
    "outputs": [
      {
        "internalType": "address[]",
        "name": "",
        "type": "address[]"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getWorkingValidatorCount",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "workingValidatorCount",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "init",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "isCurrentValidator",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes",
        "name": "voteAddr",
        "type": "bytes"
      }
    ],
    "name": "isMonitoredForMaliciousVote",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "isSystemRewardIncluded",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "index",
        "type": "uint256"
      }
    ],
    "name": "isWorkingValidator",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "maintainSlashScale",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "maxNumOfCandidates",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "maxNumOfMaintaining",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "maxNumOfWorkingCandidates",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "minorSlash",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "numOfCabinets",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "numOfMaintaining",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "previousHeight",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "previousVoteAddrFullSet",
    "outputs": [
      {
        "internalType": "bytes",
        "name": "",
        "type": "bytes"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "stcChainID",
    "outputs": [
      {
        "internalType": "uint16",
        "name": "",
        "type": "uint16"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "systemRewardAntiMEVRatio",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "systemRewardBaseRatio",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "totalInComing",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "turnLength",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "updateParam",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address[]",
        "name": "_consensusAddrs",
        "type": "address[]"
      },
      {
        "internalType": "uint64[]",
        "name": "_votingPowers",
        "type": "uint64[]"
      },
      {
        "internalType": "bytes[]",
        "name": "_voteAddrs",
        "type": "bytes[]"
      }
    ],
    "name": "updateValidatorSet",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "validatorExtraSet",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "enterMaintenanceHeight",
        "type": "uint256"
      },
      {
        "internalType": "bool",
        "name": "isMaintaining",
        "type": "bool"
      },
      {
        "internalType": "bytes",
        "name": "voteAddress",
        "type": "bytes"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "stateMutability": "payable",
    "type": "receive"
  }
]
`

const slashABI = `
[
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "InvalidValue",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyCoinbase",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "systemContract",
        "type": "address"
      }
    ],
    "name": "OnlySystemContract",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyZeroGasPrice",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "UnknownParam",
    "type": "error"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "ParamChange",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "slashCount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "failReason",
        "type": "bytes"
      }
    ],
    "name": "failedEvictionSlash",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [],
    "name": "indicatorCleaned",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "paramChange",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "validatorSlashed",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "CODE_OK",
    "outputs": [
      {
        "internalType": "uint32",
        "name": "",
        "type": "uint32"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "DECREASE_RATE",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "EVICTION_SLASH_THRESHOLD",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "INIT_EVIDENCE_REPORTER_REWARD_RATIO",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "INIT_SLASH_EVIDENCE_WINDOW",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "MINOR_SLASH_THRESHOLD",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "alreadyInit",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "clean",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "internalType": "uint256",
        "name": "count",
        "type": "uint256"
      },
      {
        "internalType": "bool",
        "name": "shouldRevert",
        "type": "bool"
      }
    ],
    "name": "downtimeSlash",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "enableMaliciousVoteSlash",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "evictionSlashThreshold",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "evidenceReporterRewardRatio",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "getSlashIndicator",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "getSlashThresholds",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "name": "indicators",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "height",
        "type": "uint256"
      },
      {
        "internalType": "uint256",
        "name": "count",
        "type": "uint256"
      },
      {
        "internalType": "bool",
        "name": "exist",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "init",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "minorSlashThreshold",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "previousHeight",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "sendEvictionSlashPackage",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "validator",
        "type": "address"
      }
    ],
    "name": "slash",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "slashEvidenceWindow",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "stcChainID",
    "outputs": [
      {
        "internalType": "uint16",
        "name": "",
        "type": "uint16"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes",
        "name": "header1",
        "type": "bytes"
      },
      {
        "internalType": "bytes",
        "name": "header2",
        "type": "bytes"
      }
    ],
    "name": "submitDoubleSignEvidence",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "components": [
          {
            "components": [
              {
                "internalType": "uint256",
                "name": "srcNum",
                "type": "uint256"
              },
              {
                "internalType": "bytes32",
                "name": "srcHash",
                "type": "bytes32"
              },
              {
                "internalType": "uint256",
                "name": "tarNum",
                "type": "uint256"
              },
              {
                "internalType": "bytes32",
                "name": "tarHash",
                "type": "bytes32"
              },
              {
                "internalType": "bytes",
                "name": "sig",
                "type": "bytes"
              }
            ],
            "internalType": "struct SlashIndicator.VoteData",
            "name": "voteA",
            "type": "tuple"
          },
          {
            "components": [
              {
                "internalType": "uint256",
                "name": "srcNum",
                "type": "uint256"
              },
              {
                "internalType": "bytes32",
                "name": "srcHash",
                "type": "bytes32"
              },
              {
                "internalType": "uint256",
                "name": "tarNum",
                "type": "uint256"
              },
              {
                "internalType": "bytes32",
                "name": "tarHash",
                "type": "bytes32"
              },
              {
                "internalType": "bytes",
                "name": "sig",
                "type": "bytes"
              }
            ],
            "internalType": "struct SlashIndicator.VoteData",
            "name": "voteB",
            "type": "tuple"
          },
          {
            "internalType": "bytes",
            "name": "voteAddr",
            "type": "bytes"
          }
        ],
        "internalType": "struct SlashIndicator.FinalityEvidence",
        "name": "_evidence",
        "type": "tuple"
      }
    ],
    "name": "submitFinalityViolationEvidence",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "updateParam",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "validators",
    "outputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]
`

const stakeABI = `
[
  {
    "inputs": [],
    "name": "AlreadyPaused",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "AlreadySlashed",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "ConsensusAddressExpired",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "DelegationAmountTooSmall",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "DuplicateConsensusAddress",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "DuplicateMoniker",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "DuplicateNodeID",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "DuplicateVoteAddress",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "ExceedsMaxNodeIDs",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InBlackList",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidAgent",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidCommission",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidConsensusAddress",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidMoniker",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidNodeID",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidRequest",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidSynPackage",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidValidator",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "InvalidValue",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "InvalidVoteAddress",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "JailTimeNotExpired",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "MaxEvictionsPerEpochReached",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "NotPaused",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyCoinbase",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyProtector",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlySelfDelegation",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "address",
        "name": "systemContract",
        "type": "address"
      }
    ],
    "name": "OnlySystemContract",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "OnlyZeroGasPrice",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "SameValidator",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "SelfDelegationNotEnough",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "TransferFailed",
    "type": "error"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "UnknownParam",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "UpdateTooFrequently",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "ValidatorExisted",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "ValidatorNotExisted",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "ValidatorNotJailed",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "VoteAddressExpired",
    "type": "error"
  },
  {
    "inputs": [],
    "name": "ZeroShares",
    "type": "error"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "oldAgent",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "newAgent",
        "type": "address"
      }
    ],
    "name": "AgentChanged",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "target",
        "type": "address"
      }
    ],
    "name": "BlackListed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "delegator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "sepAmount",
        "type": "uint256"
      }
    ],
    "name": "Claimed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint64",
        "name": "newCommissionRate",
        "type": "uint64"
      }
    ],
    "name": "CommissionRateEdited",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "newConsensusAddress",
        "type": "address"
      }
    ],
    "name": "ConsensusAddressEdited",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "delegator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "shares",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "sepAmount",
        "type": "uint256"
      }
    ],
    "name": "Delegated",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      }
    ],
    "name": "DescriptionEdited",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint8",
        "name": "version",
        "type": "uint8"
      }
    ],
    "name": "Initialized",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "nodeID",
        "type": "bytes32"
      }
    ],
    "name": "NodeIDAdded",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "validator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "nodeID",
        "type": "bytes32"
      }
    ],
    "name": "NodeIDRemoved",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "key",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "value",
        "type": "bytes"
      }
    ],
    "name": "ParamChange",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [],
    "name": "Paused",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "oldProtector",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "newProtector",
        "type": "address"
      }
    ],
    "name": "ProtectorChanged",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "srcValidator",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "dstValidator",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "delegator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "oldShares",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "newShares",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "sepAmount",
        "type": "uint256"
      }
    ],
    "name": "Redelegated",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [],
    "name": "Resumed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "failReason",
        "type": "bytes"
      }
    ],
    "name": "RewardDistributeFailed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "reward",
        "type": "uint256"
      }
    ],
    "name": "RewardDistributed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "creditContract",
        "type": "address"
      }
    ],
    "name": "StakeCreditInitialized",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "target",
        "type": "address"
      }
    ],
    "name": "UnBlackListed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "delegator",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "shares",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "sepAmount",
        "type": "uint256"
      }
    ],
    "name": "Undelegated",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "consensusAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "creditContract",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "voteAddress",
        "type": "bytes"
      }
    ],
    "name": "ValidatorCreated",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      }
    ],
    "name": "ValidatorEmptyJailed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      }
    ],
    "name": "ValidatorJailed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "jailUntil",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "slashAmount",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "enum StakeHub.SlashType",
        "name": "slashType",
        "type": "uint8"
      }
    ],
    "name": "ValidatorSlashed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      }
    ],
    "name": "ValidatorUnjailed",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "operatorAddress",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes",
        "name": "newVoteAddress",
        "type": "bytes"
      }
    ],
    "name": "VoteAddressEdited",
    "type": "event"
  },
  {
    "inputs": [],
        "name": "BREATHE_BLOCK_INTERVAL",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
        "inputs": [],
    "name": "CODE_OK",
        "outputs": [
            {
        "internalType": "uint32",
                "name": "",
        "type": "uint32"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "DEAD_ADDRESS",
        "outputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "INIT_MAX_NUMBER_NODE_ID",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "LOCK_AMOUNT",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "REDELEGATE_FEE_RATE_BASE",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "bytes32[]",
                "name": "nodeIDs",
        "type": "bytes32[]"
            }
        ],
    "name": "addNodeIDs",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "account",
        "type": "address"
            }
        ],
    "name": "addToBlackList",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "name": "agentToOperator",
        "outputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "alreadyInit",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
        "inputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "name": "blackList",
        "outputs": [
            {
        "internalType": "bool",
                "name": "",
        "type": "bool"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            },
            {
        "internalType": "uint256",
                "name": "requestNumber",
        "type": "uint256"
            }
        ],
    "name": "claim",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address[]",
                "name": "operatorAddresses",
        "type": "address[]"
            },
            {
        "internalType": "uint256[]",
                "name": "requestNumbers",
        "type": "uint256[]"
            }
        ],
    "name": "claimBatch",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "name": "consensusExpiration",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "name": "consensusToOperator",
        "outputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "consensusAddress",
        "type": "address"
            },
            {
        "internalType": "bytes",
                "name": "voteAddress",
        "type": "bytes"
            },
            {
        "internalType": "bytes",
                "name": "blsProof",
        "type": "bytes"
            },
            {
                "components": [
                    {
            "internalType": "uint64",
                        "name": "rate",
            "type": "uint64"
                    },
                    {
            "internalType": "uint64",
                        "name": "maxRate",
            "type": "uint64"
                    },
                    {
            "internalType": "uint64",
                        "name": "maxChangeRate",
            "type": "uint64"
                    }
        ],
        "internalType": "struct StakeHub.Commission",
        "name": "commission",
        "type": "tuple"
            },
            {
                "components": [
                    {
            "internalType": "string",
                        "name": "moniker",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "identity",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "website",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "details",
            "type": "string"
                    }
        ],
        "internalType": "struct StakeHub.Description",
        "name": "description",
        "type": "tuple"
            }
        ],
    "name": "createValidator",
        "outputs": [],
    "stateMutability": "payable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            },
            {
        "internalType": "bool",
                "name": "delegateVotePower",
        "type": "bool"
            }
        ],
    "name": "delegate",
        "outputs": [],
    "stateMutability": "payable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "consensusAddress",
        "type": "address"
            }
        ],
    "name": "distributeReward",
        "outputs": [],
    "stateMutability": "payable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "consensusAddress",
        "type": "address"
            }
        ],
    "name": "doubleSignSlash",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "downtimeJailTime",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "consensusAddress",
        "type": "address"
            }
        ],
    "name": "downtimeSlash",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "downtimeSlashAmount",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "uint64",
                "name": "commissionRate",
        "type": "uint64"
            }
        ],
    "name": "editCommissionRate",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "newConsensusAddress",
        "type": "address"
            }
        ],
    "name": "editConsensusAddress",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
                "components": [
                    {
            "internalType": "string",
                        "name": "moniker",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "identity",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "website",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "details",
            "type": "string"
                    }
        ],
        "internalType": "struct StakeHub.Description",
        "name": "description",
        "type": "tuple"
            }
        ],
    "name": "editDescription",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "bytes",
                "name": "newVoteAddress",
        "type": "bytes"
            },
            {
        "internalType": "bytes",
                "name": "blsProof",
        "type": "bytes"
            }
        ],
    "name": "editVoteAddress",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "evictionJailTime",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "evictionSlashAmount",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address[]",
        "name": "validatorsToQuery",
        "type": "address[]"
      }
    ],
    "name": "getNodeIDs",
    "outputs": [
      {
        "internalType": "address[]",
        "name": "consensusAddresses",
        "type": "address[]"
      },
      {
        "internalType": "bytes32[][]",
        "name": "nodeIDsList",
        "type": "bytes32[][]"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
        "inputs": [],
    "name": "getProtector",
        "outputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorAgent",
        "outputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorBasicInfo",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "createdTime",
        "type": "uint256"
            },
            {
        "internalType": "bool",
                "name": "jailed",
        "type": "bool"
            },
            {
        "internalType": "uint256",
                "name": "jailUntil",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorCommission",
        "outputs": [
            {
                "components": [
                    {
            "internalType": "uint64",
                        "name": "rate",
            "type": "uint64"
                    },
                    {
            "internalType": "uint64",
                        "name": "maxRate",
            "type": "uint64"
                    },
                    {
            "internalType": "uint64",
                        "name": "maxChangeRate",
            "type": "uint64"
                    }
        ],
        "internalType": "struct StakeHub.Commission",
        "name": "",
        "type": "tuple"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorConsensusAddress",
        "outputs": [
            {
        "internalType": "address",
                "name": "consensusAddress",
        "type": "address"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorCreditContract",
        "outputs": [
            {
        "internalType": "address",
                "name": "creditContract",
        "type": "address"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorDescription",
        "outputs": [
            {
                "components": [
                    {
            "internalType": "string",
                        "name": "moniker",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "identity",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "website",
            "type": "string"
                    },
                    {
            "internalType": "string",
                        "name": "details",
            "type": "string"
                    }
        ],
        "internalType": "struct StakeHub.Description",
        "name": "",
        "type": "tuple"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "uint256",
                "name": "offset",
        "type": "uint256"
            },
            {
        "internalType": "uint256",
                "name": "limit",
        "type": "uint256"
            }
        ],
    "name": "getValidatorElectionInfo",
        "outputs": [
            {
        "internalType": "address[]",
                "name": "consensusAddrs",
        "type": "address[]"
            },
            {
        "internalType": "uint256[]",
                "name": "votingPowers",
        "type": "uint256[]"
            },
            {
        "internalType": "bytes[]",
                "name": "voteAddrs",
        "type": "bytes[]"
            },
            {
        "internalType": "uint256",
                "name": "totalLength",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            },
            {
        "internalType": "uint256",
                "name": "index",
        "type": "uint256"
            }
        ],
    "name": "getValidatorRewardRecord",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            },
            {
        "internalType": "uint256",
                "name": "index",
        "type": "uint256"
            }
        ],
    "name": "getValidatorTotalPooledSEPRecord",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorUpdateTime",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "getValidatorVoteAddress",
        "outputs": [
            {
        "internalType": "bytes",
                "name": "voteAddress",
        "type": "bytes"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "uint256",
                "name": "offset",
        "type": "uint256"
            },
            {
        "internalType": "uint256",
                "name": "limit",
        "type": "uint256"
            }
        ],
    "name": "getValidators",
        "outputs": [
            {
        "internalType": "address[]",
                "name": "operatorAddrs",
        "type": "address[]"
            },
            {
        "internalType": "address[]",
                "name": "creditAddrs",
        "type": "address[]"
            },
            {
        "internalType": "uint256",
                "name": "totalLength",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
        "inputs": [],
    "name": "initialize",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "isPaused",
        "outputs": [
            {
        "internalType": "bool",
                "name": "",
        "type": "bool"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "bytes",
                "name": "voteAddress",
        "type": "bytes"
            }
        ],
    "name": "maliciousVoteSlash",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "maxElectedValidators",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "maxEvictionsPerEpoch",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "maxNodeIDs",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "minDelegationSEPChange",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "minSelfDelegationSEP",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "numOfJailed",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "pause",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "srcValidator",
        "type": "address"
            },
            {
        "internalType": "address",
                "name": "dstValidator",
        "type": "address"
            },
            {
        "internalType": "uint256",
                "name": "shares",
        "type": "uint256"
            },
            {
        "internalType": "bool",
                "name": "delegateVotePower",
        "type": "bool"
            }
        ],
    "name": "redelegate",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "redelegateFeeRate",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "account",
        "type": "address"
            }
        ],
    "name": "removeFromBlackList",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "bytes32[]",
                "name": "targetNodeIDs",
        "type": "bytes32[]"
            }
        ],
    "name": "removeNodeIDs",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
    "name": "resume",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "stcChainID",
    "outputs": [
      {
        "internalType": "uint16",
        "name": "",
        "type": "uint16"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
        "inputs": [
            {
        "internalType": "address[]",
                "name": "operatorAddresses",
        "type": "address[]"
            },
            {
        "internalType": "address",
                "name": "account",
        "type": "address"
            }
        ],
    "name": "syncGovToken",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [],
        "name": "unbondPeriod",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            },
            {
        "internalType": "uint256",
                "name": "shares",
        "type": "uint256"
            }
        ],
    "name": "undelegate",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "operatorAddress",
        "type": "address"
            }
        ],
    "name": "unjail",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "address",
                "name": "newAgent",
        "type": "address"
            }
        ],
    "name": "updateAgent",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "string",
                "name": "key",
        "type": "string"
            },
            {
        "internalType": "bytes",
                "name": "value",
        "type": "bytes"
            }
        ],
    "name": "updateParam",
        "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "bytes",
                "name": "",
        "type": "bytes"
            }
        ],
    "name": "voteExpiration",
        "outputs": [
            {
        "internalType": "uint256",
                "name": "",
        "type": "uint256"
            }
        ],
    "stateMutability": "view",
    "type": "function"
    },
    {
        "inputs": [
            {
        "internalType": "bytes",
                "name": "",
        "type": "bytes"
            }
        ],
    "name": "voteToOperator",
        "outputs": [
            {
        "internalType": "address",
                "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "stateMutability": "payable",
    "type": "receive"
    }
]
`
