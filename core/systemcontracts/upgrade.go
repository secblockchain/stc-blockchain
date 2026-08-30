package systemcontracts

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

type UpgradeConfig struct {
	BeforeUpgrade upgradeHook
	AfterUpgrade  upgradeHook
	ContractAddr  common.Address
	CommitUrl     string
	Code          string
}

type Upgrade struct {
	UpgradeName string
	Configs     []*UpgradeConfig
}

type upgradeHook func(blockNumber *big.Int, contractAddr common.Address, statedb vm.StateDB) error

// forkUpgrade is a future hard-fork system-contract rewrite. Register new
// upgrades in init(); historical BSC/STC fork bytecode is not shipped because
// a greenfield chain loads the latest contracts from genesis.
type forkUpgrade struct {
	configs map[string]*Upgrade
	isOn    func(config *params.ChainConfig, blockNumber *big.Int, lastBlockTime, blockTime uint64) bool
}

const (
	mainNet       = "Mainnet"
	stcTestnetNet = "StcTestnet"
	defaultNet    = "Default"
)

var (
	GenesisHash common.Hash

	// forkUpgrades is applied when a registered fork activates.
	// Example:
	//
	//	forkUpgrades = append(forkUpgrades, forkUpgrade{
	//		isOn: func(c *params.ChainConfig, n *big.Int, last, now uint64) bool {
	//			return c.IsOnPasteur(n, last, now)
	//		},
	//		configs: map[string]*Upgrade{
	//			mainNet: {UpgradeName: "pasteur", Configs: []*UpgradeConfig{{
	//				ContractAddr: common.HexToAddress(StakeHubContract),
	//				CommitUrl:    "https://github.com/secblockchain/stc-genesis-contract/commit/...",
	//				Code:         "<hex bytecode>",
	//			}}},
	//		},
	//	})
	forkUpgrades []forkUpgrade
)

func TryUpdateBuildInSystemContract(config *params.ChainConfig, blockNumber *big.Int, lastBlockTime uint64, blockTime uint64, statedb vm.StateDB, atBlockBegin bool) {
	if atBlockBegin {
		if !config.IsFeynman(blockNumber, lastBlockTime) {
			upgradeBuildInSystemContract(config, blockNumber, lastBlockTime, blockTime, statedb)
		}
		// HistoryStorageAddress is a special system contract in stc, which can't be upgraded
		if config.IsInSTC() && config.IsOnPrague(blockNumber, lastBlockTime, blockTime) {
			statedb.SetCode(params.HistoryStorageAddress, params.HistoryStorageCode, tracing.CodeChangeSystemContractUpgrade)
			statedb.SetNonce(params.HistoryStorageAddress, 1, tracing.NonceChangeNewContract)
			log.Info("Set code for HistoryStorageAddress", "blockNumber", blockNumber.Int64(), "blockTime", blockTime)
		}
	} else {
		if config.IsFeynman(blockNumber, lastBlockTime) {
			upgradeBuildInSystemContract(config, blockNumber, lastBlockTime, blockTime, statedb)
		}
	}
}

func upgradeBuildInSystemContract(config *params.ChainConfig, blockNumber *big.Int, lastBlockTime uint64, blockTime uint64, statedb vm.StateDB) {
	if config == nil || blockNumber == nil || statedb == nil || reflect.ValueOf(statedb).IsNil() {
		return
	}

	var network string
	switch GenesisHash {
	case params.STCGenesisHash:
		network = mainNet
	case params.StcTestnetGenesisHash:
		network = stcTestnetNet
	default:
		network = defaultNet
	}

	logger := log.New("system-contract-upgrade", network)
	for _, fork := range forkUpgrades {
		if fork.isOn != nil && fork.isOn(config, blockNumber, lastBlockTime, blockTime) {
			applySystemContractUpgrade(fork.configs[network], blockNumber, statedb, logger)
		}
	}
}

func applySystemContractUpgrade(upgrade *Upgrade, blockNumber *big.Int, statedb vm.StateDB, logger log.Logger) {
	if upgrade == nil {
		logger.Info("Empty upgrade config", "height", blockNumber.String())
		return
	}

	logger.Info(fmt.Sprintf("Apply upgrade %s at height %d", upgrade.UpgradeName, blockNumber.Int64()))
	for _, cfg := range upgrade.Configs {
		logger.Info(fmt.Sprintf("Upgrade contract %s to commit %s", cfg.ContractAddr.String(), cfg.CommitUrl))

		if cfg.BeforeUpgrade != nil {
			err := cfg.BeforeUpgrade(blockNumber, cfg.ContractAddr, statedb)
			if err != nil {
				panic(fmt.Errorf("contract address: %s, execute beforeUpgrade error: %s", cfg.ContractAddr.String(), err.Error()))
			}
		}

		newContractCode, err := hex.DecodeString(strings.TrimSpace(cfg.Code))
		if err != nil {
			panic(fmt.Errorf("failed to decode new contract code: %s", err.Error()))
		}
		statedb.SetCode(cfg.ContractAddr, newContractCode, tracing.CodeChangeSystemContractUpgrade)

		if cfg.AfterUpgrade != nil {
			err := cfg.AfterUpgrade(blockNumber, cfg.ContractAddr, statedb)
			if err != nil {
				panic(fmt.Errorf("contract address: %s, execute afterUpgrade error: %s", cfg.ContractAddr.String(), err.Error()))
			}
		}
	}
}
