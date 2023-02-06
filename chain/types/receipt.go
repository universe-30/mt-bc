package types

import (
	"github.com/universe-30/mt-trie/common"
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint64(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint64(1)
)

type Receipt struct {
	PostHash          []byte `json:"root"`
	Status            uint64 `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`

	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	GasUsed         uint64         `json:"gasUsed" gencodec:"required"`

	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      uint64      `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}
