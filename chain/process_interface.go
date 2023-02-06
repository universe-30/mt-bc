package chain

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/universe-30/mt-bc/chain/types"
	"github.com/universe-30/mt-bc/consensus"
	"github.com/universe-30/mt-trie/common"
)

type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the header corresponding to the hash/number argument pair.
	GetHeader(common.Hash, uint64) *types.Header
}

// Processor is an interface for processing blocks using a given initial state.
type Processor interface {
	Process(block *types.Block, statedb *state.StateDB) ([]*types.Receipt, uint64, error)
}
