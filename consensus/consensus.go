package consensus

import "github.com/universe-30/mt-bc/chain/types"

type Engine interface {
	Seal(block *types.Block) error
}
