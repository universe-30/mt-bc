package chain

import (
	"log"
	"math/big"

	"github.com/universe-30/mt-bc/chain/types"
)

func (bc *BlockChain) InsertBlock(b *types.Block) {

	l := len(bc.Blocks)
	if l == 0 {
		bc.Blocks = append(bc.Blocks, b)
		return
	}

	prev_block := bc.Blocks[l-1]
	if !isValid(*b, *prev_block) {
		log.Fatal("Invalid Block.")
		return
	}

	bc.Blocks = append(bc.Blocks, b)
	return
}

func isValid(newBlock types.Block, prevBlock types.Block) bool {

	checkNum := new(big.Int).Add(prevBlock.Number, big.NewInt(1))
	if newBlock.Number.Cmp(checkNum) != 0 {
		return false
	}
	if newBlock.ParentHash != prevBlock.Hash() {
		return false
	}
	return true
}
