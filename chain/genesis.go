package chain

import (
	"github.com/universe-30/mt-bc/chain/types"
)

type Genesis struct {
}

// 生成创世区块
func CreateGenesisBlock() *types.Block {

	header := &types.Header{}
	header.Number = 0

	blk := types.NewBlockWithHeader(header)
	blk.Txs = nil

	data := types.NewTxWithString("Genesis Block")
	txs := []*types.Transaction{data}

	return types.CreateNewBlock(blk, txs)
}
