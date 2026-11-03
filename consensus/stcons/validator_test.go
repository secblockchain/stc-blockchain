package stcons

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestValidatorHeap(t *testing.T) {
	testCases := []struct {
		description string
		k           int64
		validators  []ValidatorItem
		expected    []common.Address
	}{
		{
			description: "normal case",
			k:           2,
			validators: []ValidatorItem{
				{
					address:     common.HexToAddress("0x1"),
					votingPower: new(big.Int).Mul(big.NewInt(300), big.NewInt(1e10)),
					voteAddress: []byte("0x1"),
				},
				{
					address:     common.HexToAddress("0x2"),
					votingPower: new(big.Int).Mul(big.NewInt(200), big.NewInt(1e10)),
					voteAddress: []byte("0x2"),
				},
				{
					address:     common.HexToAddress("0x3"),
					votingPower: new(big.Int).Mul(big.NewInt(100), big.NewInt(1e10)),
					voteAddress: []byte("0x3"),
				},
			},
			expected: []common.Address{
				common.HexToAddress("0x1"),
				common.HexToAddress("0x2"),
			},
		},
		{
			description: "same voting power",
			k:           2,
			validators: []ValidatorItem{
				{
					address:     common.HexToAddress("0x1"),
					votingPower: new(big.Int).Mul(big.NewInt(300), big.NewInt(1e10)),
					voteAddress: []byte("0x1"),
				},
				{
					address:     common.HexToAddress("0x2"),
					votingPower: new(big.Int).Mul(big.NewInt(100), big.NewInt(1e10)),
					voteAddress: []byte("0x2"),
				},
				{
					address:     common.HexToAddress("0x3"),
					votingPower: new(big.Int).Mul(big.NewInt(100), big.NewInt(1e10)),
					voteAddress: []byte("0x3"),
				},
			},
			expected: []common.Address{
				common.HexToAddress("0x1"),
				common.HexToAddress("0x2"),
			},
		},
		{
			description: "zero voting power and k > len(validators)",
			k:           5,
			validators: []ValidatorItem{
				{
					address:     common.HexToAddress("0x1"),
					votingPower: new(big.Int).Mul(big.NewInt(300), big.NewInt(1e10)),
					voteAddress: []byte("0x1"),
				},
				{
					address:     common.HexToAddress("0x2"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x2"),
				},
				{
					address:     common.HexToAddress("0x3"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x3"),
				},
				{
					address:     common.HexToAddress("0x4"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x4"),
				},
			},
			expected: []common.Address{
				common.HexToAddress("0x1"),
			},
		},
		{
			description: "zero voting power and k < len(validators)",
			k:           2,
			validators: []ValidatorItem{
				{
					address:     common.HexToAddress("0x1"),
					votingPower: new(big.Int).Mul(big.NewInt(300), big.NewInt(1e10)),
					voteAddress: []byte("0x1"),
				},
				{
					address:     common.HexToAddress("0x2"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x2"),
				},
				{
					address:     common.HexToAddress("0x3"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x3"),
				},
				{
					address:     common.HexToAddress("0x4"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x4"),
				},
			},
			expected: []common.Address{
				common.HexToAddress("0x1"),
			},
		},
		{
			description: "all zero voting power",
			k:           2,
			validators: []ValidatorItem{
				{
					address:     common.HexToAddress("0x1"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x1"),
				},
				{
					address:     common.HexToAddress("0x2"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x2"),
				},
				{
					address:     common.HexToAddress("0x3"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x3"),
				},
				{
					address:     common.HexToAddress("0x4"),
					votingPower: big.NewInt(0),
					voteAddress: []byte("0x4"),
				},
			},
			expected: []common.Address{},
		},
	}
	for _, tc := range testCases {
		eligibleValidators, _, _ := getTopValidatorsByVotingPower(tc.validators, big.NewInt(tc.k))

		// check
		if len(eligibleValidators) != len(tc.expected) {
			t.Errorf("expected %d, got %d", len(tc.expected), len(eligibleValidators))
		}
		for i := 0; i < len(tc.expected); i++ {
			if eligibleValidators[i] != tc.expected[i] {
				t.Errorf("expected %s, got %s", tc.expected[i].Hex(), eligibleValidators[i].Hex())
			}
		}
	}
}

func TestMergeElectedWithGenesis(t *testing.T) {
	g1 := common.HexToAddress("0xaaa1")
	g2 := common.HexToAddress("0xaaa2")
	e1 := common.HexToAddress("0xbbb1")
	e2 := common.HexToAddress("0xbbb2")

	var g1Vote, g2Vote types.BLSPublicKey
	g1Vote[0] = 0x11
	g2Vote[0] = 0x22
	genesisVotes := map[common.Address]types.BLSPublicKey{
		g1: g1Vote,
		g2: g2Vote,
	}
	genesis := []common.Address{g1, g2}

	t.Run("empty elected keeps genesis", func(t *testing.T) {
		addrs, powers, votes := mergeElectedWithGenesis(nil, nil, nil, genesis, genesisVotes)
		if len(addrs) != 2 || addrs[0] != g1 || addrs[1] != g2 {
			t.Fatalf("unexpected addrs: %v", addrs)
		}
		if powers[0] != genesisValidatorVotingPower || powers[1] != genesisValidatorVotingPower {
			t.Fatalf("unexpected powers: %v", powers)
		}
		if !bytes.Equal(votes[0], g1Vote[:]) || !bytes.Equal(votes[1], g2Vote[:]) {
			t.Fatalf("unexpected votes")
		}
	})

	t.Run("union genesis and elected", func(t *testing.T) {
		elected := []common.Address{e1, e2}
		powers := []uint64{100, 50}
		voteAddrs := [][]byte{{0x31}, {0x32}}
		addrs, outPowers, outVotes := mergeElectedWithGenesis(elected, powers, voteAddrs, genesis, genesisVotes)
		if len(addrs) != 4 {
			t.Fatalf("expected 4 validators, got %d: %v", len(addrs), addrs)
		}
		if addrs[0] != g1 || addrs[1] != g2 || addrs[2] != e1 || addrs[3] != e2 {
			t.Fatalf("unexpected order: %v", addrs)
		}
		if outPowers[2] != 100 || outPowers[3] != 50 {
			t.Fatalf("unexpected elected powers: %v", outPowers)
		}
		if !bytes.Equal(outVotes[2], []byte{0x31}) {
			t.Fatalf("unexpected elected vote")
		}
	})

	t.Run("genesis already in elected keeps elected power", func(t *testing.T) {
		elected := []common.Address{g1, e1}
		powers := []uint64{999, 10}
		voteAddrs := [][]byte{{0xaa}, {0xbb}}
		addrs, outPowers, outVotes := mergeElectedWithGenesis(elected, powers, voteAddrs, genesis, genesisVotes)
		if len(addrs) != 3 {
			t.Fatalf("expected 3 validators, got %d: %v", len(addrs), addrs)
		}
		if addrs[0] != g1 || addrs[1] != g2 || addrs[2] != e1 {
			t.Fatalf("unexpected order: %v", addrs)
		}
		if outPowers[0] != 999 {
			t.Fatalf("expected genesis validator to keep elected power, got %d", outPowers[0])
		}
		if !bytes.Equal(outVotes[0], []byte{0xaa}) {
			t.Fatalf("expected genesis validator to keep elected vote")
		}
		if outPowers[1] != genesisValidatorVotingPower {
			t.Fatalf("expected default power for g2, got %d", outPowers[1])
		}
	})

	t.Run("empty genesis returns elected unchanged", func(t *testing.T) {
		elected := []common.Address{e1}
		powers := []uint64{7}
		voteAddrs := [][]byte{{0x01}}
		addrs, outPowers, outVotes := mergeElectedWithGenesis(elected, powers, voteAddrs, nil, nil)
		if len(addrs) != 1 || addrs[0] != e1 || outPowers[0] != 7 || !bytes.Equal(outVotes[0], []byte{0x01}) {
			t.Fatalf("elected set should be unchanged")
		}
	})
}
