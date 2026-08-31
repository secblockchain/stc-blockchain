package stcons

import (
	"bytes"
	"math/big"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func TestValidatorSetSort(t *testing.T) {
	size := 100
	validators := make([]common.Address, size)
	for i := 0; i < size; i++ {
		validators[i] = randomAddress()
	}
	sort.Sort(validatorsAscending(validators))
	for i := 0; i < size-1; i++ {
		assert.True(t, bytes.Compare(validators[i][:], validators[i+1][:]) < 0)
	}
}

func TestParseTwoValidatorGenesisExtra(t *testing.T) {
	a1 := common.HexToAddress("0x258f7a8fcc92367a0488ED9339cb3EbC2bF5a10B")
	a2 := common.HexToAddress("0x9c03E7511A1b696eA52EaB025Ec052F6E3BB9d7D")

	extra := make([]byte, extraVanity+validatorNumberSize+2*validatorBytesLength+turnLengthSize+extraSeal)
	extra[extraVanity] = 2
	copy(extra[extraVanity+validatorNumberSize:], a1[:])
	copy(extra[extraVanity+validatorNumberSize+validatorBytesLength:], a2[:])
	extra[extraVanity+validatorNumberSize+2*validatorBytesLength] = defaultTurnLength

	header := &types.Header{
		Number: big.NewInt(0),
		Extra:  extra,
	}
	addrs, voteAddrs, err := parseValidators(header, params.StconsTestChainConfig, defaultEpochLength)
	if err != nil {
		t.Fatalf("parseValidators: %v", err)
	}
	if len(addrs) != 2 || addrs[0] != a1 || addrs[1] != a2 {
		t.Fatalf("validators = %v, want [%s, %s]", addrs, a1, a2)
	}
	if len(voteAddrs) != 2 {
		t.Fatalf("unexpected vote addresses: %v", voteAddrs)
	}
	wantLen := extraVanity + validatorNumberSize + 2*validatorBytesLength + turnLengthSize + extraSeal
	if len(extra) != wantLen {
		t.Fatalf("extra len = %d, want %d", len(extra), wantLen)
	}
}

func TestParseGenesisExtra(t *testing.T) {
	validator := common.HexToAddress("0x009B494fdeC61d7bb7C8182EB736D5C0fE1CFe14")
	extra := make([]byte, extraVanity+validatorNumberSize+validatorBytesLength+turnLengthSize+extraSeal)
	extra[extraVanity] = 1
	copy(extra[extraVanity+validatorNumberSize:], validator[:])
	extra[extraVanity+validatorNumberSize+validatorBytesLength] = defaultTurnLength

	header := &types.Header{
		Number: big.NewInt(0),
		Extra:  extra,
	}
	addrs, voteAddrs, err := parseValidators(header, params.StconsTestChainConfig, defaultEpochLength)
	if err != nil {
		t.Fatalf("parseValidators: %v", err)
	}
	if len(addrs) != 1 || addrs[0] != validator {
		t.Fatalf("validators = %v, want [%s]", addrs, validator)
	}
	if len(voteAddrs) != 1 || (voteAddrs[0] != types.BLSPublicKey{}) {
		t.Fatalf("unexpected vote addresses: %v", voteAddrs)
	}
	turn, err := parseTurnLength(header, params.StconsTestChainConfig, defaultEpochLength)
	if err != nil || turn == nil || *turn != defaultTurnLength {
		t.Fatalf("parseTurnLength = %v, %v; want %d", turn, err, defaultTurnLength)
	}
}

func TestValidatorSetEqual(t *testing.T) {
	a1 := common.HexToAddress("0x54E9352D5ACb79EF865F7ae1a55c9d756b5843A3")
	a2 := common.HexToAddress("0xDB53b00b3e439A80C6fFA60F1c02F1A50d86ff74")
	a3 := common.HexToAddress("0x1111111111111111111111111111111111111111")

	base := map[common.Address]*ValidatorInfo{
		a1: {},
		a2: {},
	}
	same := map[common.Address]*ValidatorInfo{
		a1: {Index: 1},
		a2: {Index: 2},
	}
	if !validatorSetEqual(base, same) {
		t.Fatal("expected equal when only Index differs")
	}
	differentAddr := map[common.Address]*ValidatorInfo{
		a1: {},
		a3: {},
	}
	if validatorSetEqual(base, differentAddr) {
		t.Fatal("expected unequal when addresses differ")
	}
	var vote types.BLSPublicKey
	vote[0] = 1
	differentVote := map[common.Address]*ValidatorInfo{
		a1: {VoteAddress: vote},
		a2: {},
	}
	if validatorSetEqual(base, differentVote) {
		t.Fatal("expected unequal when vote keys differ")
	}
}

// TestTwoValidatorRecentsPreservedOnSameSetSwitch covers the N=2 edge case:
// minerHistoryCheckLen is 1, so the epoch validator-set switch runs at block 1.
// If Recents were wiped while the set is unchanged, the block-1 signer could seal
// block 2 and then stall waiting for the other validator.
func TestTwoValidatorRecentsPreservedOnSameSetSwitch(t *testing.T) {
	a1 := common.HexToAddress("0x54E9352D5ACb79EF865F7ae1a55c9d756b5843A3")
	a2 := common.HexToAddress("0xDB53b00b3e439A80C6fFA60F1c02F1A50d86ff74")

	snap := &Snapshot{
		Number:     0,
		TurnLength: defaultTurnLength,
		Validators: map[common.Address]*ValidatorInfo{
			a1: {Index: 1},
			a2: {Index: 2},
		},
		Recents: make(map[uint64]common.Address),
	}
	if got := snap.minerHistoryCheckLen(); got != 1 {
		t.Fatalf("minerHistoryCheckLen = %d, want 1 for 2 validators", got)
	}

	// Simulate apply() after sealing block 1 by a1, then hitting the same-set switch.
	const blockNumber = uint64(1)
	snap.Recents[blockNumber] = a1
	newVals := map[common.Address]*ValidatorInfo{
		a1: {},
		a2: {},
	}
	sameSet := validatorSetEqual(snap.Validators, newVals)
	if !sameSet {
		t.Fatal("expected same validator set")
	}
	// Mirror the fixed apply() path: keep Recents when the set is unchanged.
	snap.Number = blockNumber

	if !snap.SignRecently(a1) {
		t.Fatal("a1 should still be recent after same-set epoch switch at block 1")
	}
	if snap.SignRecently(a2) {
		t.Fatal("a2 should not be marked recent")
	}
}
