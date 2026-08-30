package vm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestPasteurKeepsSTCPrecompiles(t *testing.T) {
	for _, b := range []byte{0x66, 0x68, 0x69} {
		addr := common.BytesToAddress([]byte{b})
		if PrecompiledContractsPasteur[addr] == nil {
			t.Fatalf("%s must remain present under Pasteur", addr.Hex())
		}
	}
}
