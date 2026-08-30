package stcons

import (
	"time"

	cmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

const millisecondsUnit = 50 // not enforced at the consensus level

func blockIntervalMs(config *params.StconsConfig) uint64 {
	if config != nil && config.Period > 0 {
		return config.Period * 1000
	}
	return defaultBlockInterval
}

func (p *Stcons) sealingDelay(header *types.Header) time.Duration {
	return time.Until(time.UnixMilli(int64(header.MilliTimestamp())))
}

func (p *Stcons) nextBlockMilliTime(snap *Snapshot, header, parent *types.Header) uint64 {
	blockTime := parent.MilliTimestamp() + snap.BlockInterval + p.backOffTime(snap, parent, header, header.Coinbase)
	if now := uint64(time.Now().UnixMilli()); blockTime < now {
		// Align the millisecond part of the timestamp.
		blockTime = uint64(cmath.CeilDiv(int(now), millisecondsUnit)) * millisecondsUnit
	}
	return blockTime
}

func (p *Stcons) verifyBlockTime(snap *Snapshot, header, parent *types.Header) error {
	if header.MilliTimestamp() < parent.MilliTimestamp()+snap.BlockInterval+p.backOffTime(snap, parent, header, header.Coinbase) {
		return consensus.ErrFutureBlock
	}
	return nil
}
