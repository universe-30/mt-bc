package chain

import (
	"sync/atomic"

	"github.com/universe-30/mt-bc/chain/types"
	"github.com/universe-30/mt-trie/accdb"
)

type BlockChain struct {
	db accdb.Database // Low level persistent database to store final content in

	currentBlock atomic.Value // Current head of the block chain

	processor Processor // Block transaction processor interface

}

func NewBlockChain() (*BlockChain, error) {

	bc := &BlockChain{}

	bc.processor = NewStateProcessor(bc)

	return bc, nil
}

func (bc *BlockChain) CurrentBlock() *types.Block {
	return bc.currentBlock.Load().(*types.Block)
}
