package stcons

import (
	"container/heap"
	"context"
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/systemcontracts"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// sameDayInUTC reports whether two unix timestamps fall on the same UTC day
// as defined by BreatheBlockInterval.
func sameDayInUTC(first, second uint64) bool {
	return first/params.BreatheBlockInterval == second/params.BreatheBlockInterval
}

func isBreatheBlock(lastBlockTime, blockTime uint64) bool {
	return lastBlockTime != 0 && !sameDayInUTC(lastBlockTime, blockTime)
}

type ValidatorItem struct {
	address     common.Address
	votingPower *big.Int
	voteAddress []byte
}

// ValidatorHeap is a max-heap of validator voting power.
type ValidatorHeap []ValidatorItem

func (h *ValidatorHeap) Len() int { return len(*h) }

func (h *ValidatorHeap) Less(i, j int) bool {
	if (*h)[i].votingPower.Cmp((*h)[j].votingPower) == 0 {
		return (*h)[i].address.Hex() < (*h)[j].address.Hex()
	}
	return (*h)[i].votingPower.Cmp((*h)[j].votingPower) == 1
}

func (h *ValidatorHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *ValidatorHeap) Push(x interface{}) {
	*h = append(*h, x.(ValidatorItem))
}

func (h *ValidatorHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (p *Stcons) updateValidatorSetV2(state vm.StateDB, header *types.Header, chain core.ChainContext,
	txs *[]*types.Transaction, receipts *[]*types.Receipt, receivedTxs *[]*types.Transaction, usedGas *uint64, mode systemTxMode, tracer *tracing.Hooks,
) error {
	blockNr := rpc.BlockNumberOrHashWithHash(header.ParentHash, false)
	validatorItems, err := p.getValidatorElectionInfo(blockNr)
	if err != nil {
		return err
	}
	maxElectedValidators, err := p.getMaxElectedValidators(blockNr)
	if err != nil {
		return err
	}

	eValidators, eVotingPowers, eVoteAddrs := getTopValidatorsByVotingPower(validatorItems, maxElectedValidators)
	if len(eValidators) == 0 {
		log.Warn("skip updateValidatorSetV2: no elected validators")
		return nil
	}

	method := "updateValidatorSetV2"
	data, err := p.validatorSetABI.Pack(method, eValidators, eVotingPowers, eVoteAddrs)
	if err != nil {
		log.Error("Unable to pack tx for updateValidatorSetV2", "error", err)
		return err
	}

	msg := p.getSystemMessage(header.Coinbase, common.HexToAddress(systemcontracts.ValidatorContract), data, common.Big0)
	return p.applyTransaction(msg, state, header, chain, txs, receipts, receivedTxs, usedGas, mode, tracer)
}

func (p *Stcons) getValidatorElectionInfo(blockNr rpc.BlockNumberOrHash) ([]ValidatorItem, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	method := "getValidatorElectionInfo"
	toAddress := common.HexToAddress(systemcontracts.StakeHubContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))

	data, err := p.stakeHubABI.Pack(method, big.NewInt(0), big.NewInt(0))
	if err != nil {
		log.Error("Unable to pack tx for getValidatorElectionInfo", "error", err)
		return nil, err
	}
	msgData := (hexutil.Bytes)(data)

	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, &blockNr, nil, nil)
	if err != nil {
		return nil, err
	}

	var validators []common.Address
	var votingPowers []*big.Int
	var voteAddrs [][]byte
	var totalLength *big.Int
	if err := p.stakeHubABI.UnpackIntoInterface(&[]interface{}{&validators, &votingPowers, &voteAddrs, &totalLength}, method, result); err != nil {
		return nil, err
	}
	if totalLength.Int64() != int64(len(validators)) || totalLength.Int64() != int64(len(votingPowers)) || totalLength.Int64() != int64(len(voteAddrs)) {
		return nil, errors.New("validator length not match")
	}

	validatorItems := make([]ValidatorItem, len(validators))
	for i := 0; i < len(validators); i++ {
		validatorItems[i] = ValidatorItem{
			address:     validators[i],
			votingPower: votingPowers[i],
			voteAddress: voteAddrs[i],
		}
	}

	return validatorItems, nil
}

func (p *Stcons) getMaxElectedValidators(blockNr rpc.BlockNumberOrHash) (maxElectedValidators *big.Int, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	method := "maxElectedValidators"
	toAddress := common.HexToAddress(systemcontracts.StakeHubContract)
	gas := (hexutil.Uint64)(uint64(math.MaxUint64 / 2))

	data, err := p.stakeHubABI.Pack(method)
	if err != nil {
		log.Error("Unable to pack tx for maxElectedValidators", "error", err)
		return nil, err
	}
	msgData := (hexutil.Bytes)(data)

	result, err := p.ethAPI.Call(ctx, ethapi.TransactionArgs{
		Gas:  &gas,
		To:   &toAddress,
		Data: &msgData,
	}, &blockNr, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := p.stakeHubABI.UnpackIntoInterface(&maxElectedValidators, method, result); err != nil {
		return nil, err
	}

	return maxElectedValidators, nil
}

func getTopValidatorsByVotingPower(validatorItems []ValidatorItem, maxElectedValidators *big.Int) ([]common.Address, []uint64, [][]byte) {
	var validatorHeap ValidatorHeap
	for i := 0; i < len(validatorItems); i++ {
		if validatorItems[i].votingPower.Cmp(big.NewInt(0)) == 1 {
			validatorHeap = append(validatorHeap, validatorItems[i])
		}
	}
	hp := &validatorHeap
	heap.Init(hp)

	topN := int(maxElectedValidators.Int64())
	if topN > len(validatorHeap) {
		topN = len(validatorHeap)
	}
	eValidators := make([]common.Address, topN)
	eVotingPowers := make([]uint64, topN)
	eVoteAddrs := make([][]byte, topN)
	for i := 0; i < topN; i++ {
		item := heap.Pop(hp).(ValidatorItem)
		eValidators[i] = item.address
		eVotingPowers[i] = new(big.Int).Div(item.votingPower, big.NewInt(1e10)).Uint64()
		eVoteAddrs[i] = item.voteAddress
	}

	return eValidators, eVotingPowers, eVoteAddrs
}
