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
