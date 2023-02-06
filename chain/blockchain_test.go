package chain

import (
	"fmt"
	"log"
	"testing"

	"github.com/universe-30/mt-bc/chain/types"
	"github.com/universe-30/mt-bc/consensus/ethash.go"
)

func TestSetBlockData(t *testing.T) {

	genesisBlock := CreateGenesisBlock()

	bc := CreateNewBlockChain(genesisBlock)

	currentBlock := CreateNewBlock(genesisBlock)

	bc.InsertBlock(currentBlock)

	log.Printf("bc out:")
	log.Printf("detail: %v", *bc)

	fmt.Printf("detail: %+v", bc)
}

func TestBlockDataEqual(t *testing.T) {

	genesisBlock := CreateGenesisBlock()

	bc := CreateNewBlockChain(genesisBlock)
	currentBlock := CreateNewBlock(genesisBlock)
	bc.InsertBlock(currentBlock)

	genesisBlock2 := CreateGenesisBlock()

	bc2 := CreateNewBlockChain(genesisBlock2)
	currentBlock2 := CreateNewBlock(genesisBlock2)
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
	blockChain, _ := NewBlockChain()
	blockChain.InsertBlock(genesisBlock)
	return blockChain
}

func CreateNewBlock(prevBlock *types.Block) *types.Block {
	data := types.NewTxWithString("aabc")
	txs := []*types.Transaction{data}

	block := types.CreateNewBlock(prevBlock, txs)

	pow := ethash.NewProofOfWork()
	hash, nonce, err := pow.Seal(block)
	if err != nil {
		log.Panic(err)
		return nil
	}

	block.SetFinal(hash, nonce)

	return block
}
