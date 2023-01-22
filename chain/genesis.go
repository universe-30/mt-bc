package chain

import (
	"math/big"

	"github.com/universe-30/mt-bc/chain/types"
)

type Genesis struct {
}

// 生成创世区块
func CreateGenesisBlock() *types.Block {
	block := types.Block{}
	block.Number = big.NewInt(0)

	data := types.NewTxWithString("Genesis Block")
	txs := []*types.Transaction{data}

	return types.CreateNewBlock(block, txs)
}
