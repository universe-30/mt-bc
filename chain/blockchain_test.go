package chain

import (
	"fmt"
	"log"
	"testing"

	"github.com/universe-30/mt-bc/chain/types"
)

func TestSetBlockData(t *testing.T) {

	genesisBlock := CreateGenesisBlock()

	bc := CreateNewBlockChain(genesisBlock)

	l := len(bc.Blocks)
	preBlock := bc.Blocks[l-1]
	currentBlock := CreateNewBlock(preBlock)

	bc.InsertBlock(currentBlock)

	log.Printf("bc out:")
	log.Printf("detail: %v", *bc)

	fmt.Printf("detail: %+v", bc)
}

func TestBlockDataEqual(t *testing.T) {

	genesisBlock := CreateGenesisBlock()

	bc := CreateNewBlockChain(genesisBlock)
	l := len(bc.Blocks)
	preBlock := bc.Blocks[l-1]
	currentBlock := CreateNewBlock(preBlock)
	bc.InsertBlock(currentBlock)

	genesisBlock2 := CreateGenesisBlock()

	bc2 := CreateNewBlockChain(genesisBlock2)
	l2 := len(bc2.Blocks)
	preBlock2 := bc.Blocks[l2-1]
	currentBlock2 := CreateNewBlock(preBlock2)
	bc2.InsertBlock(currentBlock2)

	if currentBlock.Hash() != currentBlock2.Hash() {
		t.Errorf("Hash Not Equal %x, %x", currentBlock.Hash(), currentBlock2.Hash())
	} else {
		t.Logf("Hash Equal %x, %x", currentBlock.Hash(), currentBlock2.Hash())
	}

	log.Printf("bc out:")
	log.Printf("detail: %v", *bc)
	log.Printf("detail2: %v", bc2)
}

// 生成区块链
func CreateNewBlockChain(genesisBlock *types.Block) *BlockChain {
	blockChain := BlockChain{}
	blockChain.InsertBlock(genesisBlock)
	return &blockChain
}

func CreateNewBlock(prevBlock *types.Block) *types.Block {
	data := types.NewTxWithString("aabc")
	txs := []*types.Transaction{data}

	return types.CreateNewBlock(*prevBlock, txs)
}
