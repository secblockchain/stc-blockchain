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
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/systemcontracts"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// genesisValidatorVotingPower is the compressed voting power written into
// updateValidatorSetV2 for genesis validators that are not present in StakeHub
// election results. It is intentionally small but non-zero so they remain in the set.
const genesisValidatorVotingPower = uint64(1)

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

// loadGenesisValidators parses and caches the validator list from the genesis
// header extraData. Safe to call repeatedly.
func (p *Stcons) loadGenesisValidators() {
	p.genesisValsOnce.Do(func() {
		header := rawdb.ReadHeader(p.db, p.genesisHash, 0)
		if header == nil {
			log.Error("Failed to read genesis header for validator floor", "hash", p.genesisHash)
			return
		}
		vals, votes, err := parseValidators(header, p.chainConfig, defaultEpochLength)
		if err != nil {
			log.Error("Failed to parse genesis validators", "hash", p.genesisHash, "err", err)
			return
		}
		p.genesisVals = vals
		p.genesisVoteKeys = make(map[common.Address]types.BLSPublicKey, len(vals))
		for i, v := range vals {
			p.genesisVoteKeys[v] = votes[i]
		}
		log.Info("Loaded genesis validator floor", "count", len(vals))
	})
}

// unionGenesisValidators returns contract validators union genesis validators.
// Genesis addresses missing from the contract set are appended with their genesis BLS keys.
func (p *Stcons) unionGenesisValidators(valSet []common.Address, voteAddrMap map[common.Address]*types.BLSPublicKey) ([]common.Address, map[common.Address]*types.BLSPublicKey) {
	p.loadGenesisValidators()
	if len(p.genesisVals) == 0 {
		return valSet, voteAddrMap
	}
	if voteAddrMap == nil {
		voteAddrMap = make(map[common.Address]*types.BLSPublicKey)
	}
	seen := make(map[common.Address]struct{}, len(valSet)+len(p.genesisVals))
	for _, v := range valSet {
		seen[v] = struct{}{}
	}
	for _, g := range p.genesisVals {
		if _, ok := seen[g]; ok {
			continue
		}
		valSet = append(valSet, g)
		seen[g] = struct{}{}
		if _, ok := voteAddrMap[g]; !ok {
			k := new(types.BLSPublicKey)
			*k = p.genesisVoteKeys[g]
			voteAddrMap[g] = k
		}
	}
	return valSet, voteAddrMap
}

// mergeElectedWithGenesis ensures genesis validators are always present in the
// set written by updateValidatorSetV2. Genesis validators come first; elected
// validators from StakeHub that are not already genesis entries follow.
func mergeElectedWithGenesis(
	elected []common.Address,
	powers []uint64,
	voteAddrs [][]byte,
	genesis []common.Address,
	genesisVotes map[common.Address]types.BLSPublicKey,
) ([]common.Address, []uint64, [][]byte) {
	if len(genesis) == 0 {
		return elected, powers, voteAddrs
	}

	type electedInfo struct {
		power uint64
		vote  []byte
	}
	electedLookup := make(map[common.Address]electedInfo, len(elected))
	for i, addr := range elected {
		info := electedInfo{power: powers[i]}
		if i < len(voteAddrs) {
			info.vote = voteAddrs[i]
		}
		electedLookup[addr] = info
	}

	outAddrs := make([]common.Address, 0, len(genesis)+len(elected))
	outPowers := make([]uint64, 0, len(genesis)+len(elected))
	outVotes := make([][]byte, 0, len(genesis)+len(elected))
	seen := make(map[common.Address]struct{}, len(genesis)+len(elected))

	for _, g := range genesis {
		seen[g] = struct{}{}
		outAddrs = append(outAddrs, g)
		if info, ok := electedLookup[g]; ok {
			outPowers = append(outPowers, info.power)
			if len(info.vote) > 0 {
				outVotes = append(outVotes, info.vote)
			} else {
				key := genesisVotes[g]
				outVotes = append(outVotes, key[:])
			}
		} else {
			outPowers = append(outPowers, genesisValidatorVotingPower)
			key := genesisVotes[g]
			outVotes = append(outVotes, key[:])
		}
	}
	for i, addr := range elected {
		if _, ok := seen[addr]; ok {
			continue
		}
		outAddrs = append(outAddrs, addr)
		outPowers = append(outPowers, powers[i])
		outVotes = append(outVotes, voteAddrs[i])
	}
	return outAddrs, outPowers, outVotes
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

	p.loadGenesisValidators()
	eValidators, eVotingPowers, eVoteAddrs = mergeElectedWithGenesis(
		eValidators, eVotingPowers, eVoteAddrs, p.genesisVals, p.genesisVoteKeys,
	)

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
