package stcons

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/systemcontracts"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	buildertypes "github.com/ethereum/go-ethereum/core/types/builder"
)

var signableSystemTxSelectors = map[string][4]byte{
	"distributeFinalityReward": {0x30, 0x0c, 0x35, 0x67},
	"updateValidatorSet":       {0xe6, 0x92, 0xf0, 0x6b},
}

type expectedSystemTxEntry struct {
	method   string
	selector [4]byte
}

// PrepareForBidBlock is Prepare with Coinbase set to the in-turn validator instead of p.val.
func (p *Stcons) PrepareForBidBlock(chain consensus.ChainHeaderReader, header *types.Header) error {
	// Coinbase must be set before prepare(): backOffTime and calcDifficulty depend on it.
	number := header.Number.Uint64()
	snap, err := p.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	header.Coinbase = snap.inturnValidator()
	return p.prepare(chain, header)
}

// FinalizeAndAssembleBidBlock assembles a BidBlock with unsigned system txs.
func (p *Stcons) FinalizeAndAssembleBidBlock(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB,
	body *types.Body, receipts []*types.Receipt, tracer *tracing.Hooks) (*types.Block, []*types.Receipt, error) {
	block, receipts, err := p.finalizeAndAssemble(chain, header, state, body, receipts, tracer, systemTxPacking)
	if err != nil {
		return nil, nil, err
	}
	return block, receipts, nil
}

// SignSystemTx signs a BidBlock system tx with the validator key.
func (p *Stcons) SignSystemTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.signTxFn == nil {
		return nil, errors.New("signTxFn not set")
	}
	return p.signTxFn(accounts.Account{Address: p.val}, tx, chainID)
}

// isUnsignedSystemTxCandidate reports whether tx looks like an unsigned
// BidBlock system tx. It does not recover the sender.
func (p *Stcons) isUnsignedSystemTxCandidate(tx *types.Transaction) bool {
	if tx == nil || tx.To() == nil || !isToSystemContract(*tx.To()) {
		return false
	}
	// Canonical system txs are legacy (see getSystemMessage); typed envelopes
	// carry fields the shape check does not bind, so never blind-sign them.
	if tx.Type() != types.LegacyTxType {
		return false
	}
	if tx.GasPrice() == nil || tx.EffectiveGasPriceForSTC().Sign() != 0 {
		return false
	}
	v, r, s := tx.RawSignatureValues()
	return isZeroSig(v, r, s)
}

// isSignableSystemTx reports whether tx can be bind-signed for BidBlock.
func (p *Stcons) isSignableSystemTx(tx *types.Transaction) bool {
	if !p.isUnsignedSystemTxCandidate(tx) {
		return false
	}
	if *tx.To() != common.HexToAddress(systemcontracts.ValidatorContract) {
		return false
	}
	return p.hasSignableSelector(tx.Data())
}

// expectedSystemTxShape returns the expected trailing system-tx order for accepted BidBlocks:
//
//	distributeFinalityReward (cond.) -> updateValidatorSet (cond.)
//
// Fee deposit is applied in Finalize via applySystemCall and is not a packed system tx.
func (p *Stcons) expectedSystemTxShape(header, parent *types.Header) []expectedSystemTxEntry {
	shape := make([]expectedSystemTxEntry, 0, 2)

	if header.Number.Uint64()%finalityRewardInterval == 0 {
		shape = append(shape, expectedSystemTxEntry{
			method:   "distributeFinalityReward",
			selector: p.selectorFor("distributeFinalityReward"),
		})
	}

	if isBreatheBlock(parent.Time, header.Time) {
		shape = append(shape, expectedSystemTxEntry{
			method:   "updateValidatorSet",
			selector: p.selectorFor("updateValidatorSet"),
		})
	}

	return shape
}

func (p *Stcons) verifySystemTxShape(txs []*types.Transaction, shape []expectedSystemTxEntry) error {
	if len(txs) < len(shape) {
		return fmt.Errorf("missing required system tx %q", shape[len(txs)].method)
	}
	if len(txs) > len(shape) {
		return fmt.Errorf("unexpected extra system tx at position %d (selector 0x%x)",
			len(shape), txSelector(txs[len(shape)]))
	}
	for i, exp := range shape {
		if !bytes.HasPrefix(txs[i].Data(), exp.selector[:]) {
			return fmt.Errorf("expected system tx %q at position %d, got selector 0x%x",
				exp.method, i, txSelector(txs[i]))
		}
	}
	return nil
}

// ExtractBidBlockSystemTxStart locates the trailing unsigned system-tx region and
// returns its start index. Fee deposit is no longer a packed system tx; GasFee is
// supplied on BidBlock.GasFee instead.
func (p *Stcons) ExtractBidBlockSystemTxStart(txs []*types.Transaction) int {
	systemTxStart := len(txs)
	for i := len(txs) - 1; i >= 0; i-- {
		if !p.isUnsignedSystemTxCandidate(txs[i]) {
			break
		}
		systemTxStart = i
	}
	return systemTxStart
}

// ExtractBidBlockDepositValue is kept as a thin wrapper for older callers.
// GasFee is always zero here; use BidBlock.GasFee for fee ranking.
func (p *Stcons) ExtractBidBlockDepositValue(txs []*types.Transaction) (int, *big.Int) {
	return p.ExtractBidBlockSystemTxStart(txs), new(big.Int)
}

// VerifyBidBlockSystemTxs validates the trailing unsigned system-tx region starting at systemTxStart.
//
//	Stage 1 — each trailing unsigned tx must be on the signable whitelist.
//	Stage 2 — selectors & order must match expectedSystemTxShape for this header.
func (p *Stcons) VerifyBidBlockSystemTxs(decoded *buildertypes.DecodedBidBlock, parent *types.Header, systemTxStart int) error {
	for i := systemTxStart; i < len(decoded.Txs); i++ {
		if !p.isSignableSystemTx(decoded.Txs[i]) {
			toAddr := "<nil>"
			if decoded.Txs[i].To() != nil {
				toAddr = decoded.Txs[i].To().Hex()
			}
			return fmt.Errorf("unsigned system tx at position %d (to=%s) is not on the signable whitelist", i, toAddr)
		}
	}
	shape := p.expectedSystemTxShape(decoded.Header, parent)
	return p.verifySystemTxShape(decoded.Txs[systemTxStart:], shape)
}

func (p *Stcons) hasSignableSelector(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	selector := data[:4]
	for _, methodSelector := range signableSystemTxSelectors {
		if bytes.Equal(selector, methodSelector[:]) {
			return true
		}
	}
	return false
}

func (p *Stcons) selectorFor(methodName string) [4]byte {
	selector, ok := signableSystemTxSelectors[methodName]
	if !ok {
		panic(fmt.Sprintf("missing fixed system tx selector %s", methodName))
	}
	return selector
}

func (p *Stcons) BlockTimeUpperCheck(chain consensus.ChainHeaderReader, header *types.Header) error {
	number := header.Number.Uint64()
	snap, err := p.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	maxAllowed := p.nextBlockMilliTime(snap, header, parent)
	if header.MilliTimestamp() > maxAllowed {
		return fmt.Errorf("BidBlock time too far in future: headerTime=%d, maxAllowed=%d",
			header.MilliTimestamp(), maxAllowed)
	}
	return nil
}

func txSelector(tx *types.Transaction) []byte {
	data := tx.Data()
	if len(data) < 4 {
		return data
	}
	return data[:4]
}

func isZeroSig(v, r, s *big.Int) bool {
	isZero := func(x *big.Int) bool { return x == nil || x.Sign() == 0 }
	return isZero(v) && isZero(r) && isZero(s)
}
