package chain

import (
	"fmt"

	"github.com/universe-30/mt-bc/chain/types"
)

type BlockChain struct {
	Blocks []*types.Block
}

func (bc *BlockChain) String() {
	for _, block := range bc.Blocks {
		fmt.Printf("Number: %d \n", block.Number)
		fmt.Printf("ParentHash: %s \n", block.ParentHash)
		fmt.Printf("CurrHash: %s \n", block.Hash())
		// fmt.Printf("Data: %s \n", block.Data)
		fmt.Printf("Timestamp: %d \n", block.Time)
		fmt.Println()
	}
}
