package core

import (
	"math/big"
	"os"
	"regexp"
	"slices"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var removedSystemContracts = []common.Address{
	common.HexToAddress("0x0000000000000000000000000000000000001003"),
	common.HexToAddress("0x0000000000000000000000000000000000001004"),
	common.HexToAddress("0x0000000000000000000000000000000000001005"),
	common.HexToAddress("0x0000000000000000000000000000000000001006"),
	common.HexToAddress("0x0000000000000000000000000000000000002000"),
	common.HexToAddress("0x0000000000000000000000000000000000001008"),
	common.HexToAddress("0x0000000000000000000000000000000000003000"),
}

type allocItem struct {
	Addr    *big.Int
	Balance *big.Int
	Misc    *allocItemMisc `rlp:"optional"`
}

type allocItemMisc struct {
	Nonce uint64
	Code  []byte
	Slots []allocItemStorageItem
}

type allocItemStorageItem struct {
	Key common.Hash
	Val common.Hash
}

func filterEmbeddedAlloc(data string) string {
	alloc := decodePrealloc(data)
	remove := make(map[common.Address]bool, len(removedSystemContracts))
	for _, addr := range removedSystemContracts {
		remove[addr] = true
	}
	for addr := range alloc {
		if remove[addr] {
			delete(alloc, addr)
		}
	}
	items := make([]allocItem, 0, len(alloc))
	for addr, account := range alloc {
		var misc *allocItemMisc
		if len(account.Storage) > 0 || len(account.Code) > 0 || account.Nonce != 0 {
			misc = &allocItemMisc{
				Nonce: account.Nonce,
				Code:  account.Code,
				Slots: make([]allocItemStorageItem, 0, len(account.Storage)),
			}
			for key, val := range account.Storage {
				misc.Slots = append(misc.Slots, allocItemStorageItem{key, val})
			}
			slices.SortFunc(misc.Slots, func(a, b allocItemStorageItem) int {
				return a.Key.Cmp(b.Key)
			})
		}
		bigAddr := new(big.Int).SetBytes(addr.Bytes())
		items = append(items, allocItem{bigAddr, account.Balance, misc})
	}
	slices.SortFunc(items, func(a, b allocItem) int {
		return a.Addr.Cmp(b.Addr)
	})
	encoded, err := rlp.EncodeToBytes(items)
	if err != nil {
		panic(err)
	}
	return strconv.QuoteToASCII(string(encoded))
}

func TestRewriteGenesisAlloc(t *testing.T) {
	if os.Getenv("REWRITE_GENESIS_ALLOC") == "" {
		t.Skip("set REWRITE_GENESIS_ALLOC=1 to rewrite genesis_alloc.go")
	}
	path := "genesis_alloc.go"
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	re := regexp.MustCompile(`(?m)^(const (stcMainnetAllocData|stcTestnetAllocData) = )("(\\.|[^"\\])*")`)
	out := re.ReplaceAllStringFunc(text, func(match string) string {
		sub := re.FindStringSubmatch(match)
		if sub == nil {
			return match
		}
		prefix, quoted := sub[1], sub[3]
		raw, err := strconv.Unquote(quoted)
		if err != nil {
			t.Fatal(err)
		}
		return prefix + filterEmbeddedAlloc(raw)
	})
	if out == text {
		t.Fatal("no STC alloc constants rewritten")
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRemovedSystemContractsNotInDefaultSTCGenesis(t *testing.T) {
	for _, addr := range removedSystemContracts {
		if _, ok := DefaultSTCGenesisBlock().Alloc[addr]; ok {
			t.Fatalf("STC mainnet alloc still contains %s", addr.Hex())
		}
		if _, ok := DefaultStcTestnetGenesisBlock().Alloc[addr]; ok {
			t.Fatalf("STC testnet alloc still contains %s", addr.Hex())
		}
	}
}
