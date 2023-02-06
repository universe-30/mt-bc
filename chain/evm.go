package chain

import (
	"math/big"

	"github.com/universe-30/mt-bc/chain/types"
	"github.com/universe-30/mt-bc/chain/vm"
	"github.com/universe-30/mt-trie/common"
)

func NewEVMBlockContext(header *types.Header, author *common.Address) vm.BlockContext {
	var (
		beneficiary common.Address
	)

	// If we don't have an explicit author (i.e. not mining), extract from the header
	if author == nil {
		beneficiary = header.Coinbase
	} else {
		beneficiary = *author
	}

	return vm.BlockContext{
		Coinbase: beneficiary,
		BaseFee:  header.BaseFee,
	}
}

func NewEVMTxContext(msg Message) vm.TxContext {
	return vm.TxContext{
		Origin:   msg.From(),
		GasPrice: new(big.Int).Set(msg.GasPrice()),
	}
}
