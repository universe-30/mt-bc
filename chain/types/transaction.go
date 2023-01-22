package types

import (
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/universe-30/mt-trie/common"
	"github.com/universe-30/mt-trie/rlp"
)

type Transaction struct {
	inner TxData    // Consensus contents of a transaction
	time  time.Time // Time first seen locally (spam avoidance)

	// caches
	hash atomic.Value
}

// type TxData interface {
// 	copy() TxData // creates a deep copy and initializes all fields
// 	data() []byte
// 	value() *big.Int
// 	nonce() uint64
// 	to() string
// }

// NewTx creates a new transaction.
func NewTx(inner TxData) *Transaction {
	tx := new(Transaction)
	tx.setDecoded(inner.copy(), 0)
	return tx
}

func NewTxWithString(data string) *Transaction {
	tx := new(Transaction)

	inner := &TxData{
		Data: []byte(data),
	}
	tx.setDecoded(inner.copy(), 0)
	return tx
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *Transaction) setDecoded(inner TxData, size uint64) {
	tx.inner = inner
	tx.time = time.Now()
}

// AccessListTx
type TxData struct {
	Nonce uint64          // nonce of sender account
	To    *common.Address `rlp:"nil"` // nil means contract creation
	Value *big.Int        // wei amount
	Data  []byte          // contract invocation input data
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *TxData) copy() TxData {
	cpy := &TxData{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		// These are copied below.
		Value: new(big.Int),
	}

	return *cpy
}

// copyAddressPtr copies an address.
func copyAddressPtr(a *common.Address) *common.Address {
	if a == nil {
		return nil
	}
	cpy := *a
	return &cpy
}

// Hash returns the transaction hash.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	// bytes := sha256.Sum256([]byte(tx.inner))
	// return hex.EncodeToString(bytes[:])

	h := rlpHash(tx.inner)
	tx.hash.Store(h)
	return h
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, tx.inner)
}
