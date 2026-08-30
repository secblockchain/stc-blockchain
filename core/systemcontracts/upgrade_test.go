package systemcontracts

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

func TestUpgradeBuildInSystemContractNilInterface(t *testing.T) {
	var (
		config               = params.STCChainConfig
		blockNumber          = big.NewInt(1)
		lastBlockTime uint64 = 0
		blockTime     uint64 = 1
		statedb       vm.StateDB
	)

	GenesisHash = params.STCGenesisHash

	upgradeBuildInSystemContract(config, blockNumber, lastBlockTime, blockTime, statedb)
}

func TestUpgradeBuildInSystemContractNilValue(t *testing.T) {
	var (
		config                   = params.STCChainConfig
		blockNumber              = big.NewInt(1)
		lastBlockTime uint64     = 0
		blockTime     uint64     = 1
		statedb       vm.StateDB = (*state.StateDB)(nil)
	)

	GenesisHash = params.STCGenesisHash

	upgradeBuildInSystemContract(config, blockNumber, lastBlockTime, blockTime, statedb)
}
