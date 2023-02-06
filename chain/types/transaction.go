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
	from atomic.Value
}

func (tx *Transaction) Data() []byte { return tx.inner.data() }

func (tx *Transaction) Gas() uint64 { return tx.inner.gas() }

func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.inner.gasPrice()) }

func (tx *Transaction) Value() *big.Int { return new(big.Int).Set(tx.inner.value()) }

func (tx *Transaction) Nonce() uint64 { return tx.inner.nonce() }

func (tx *Transaction) To() *common.Address {
	return copyAddressPtr(tx.inner.to())
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
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
}

// accessors for innerTx.
func (tx *TxData) data() []byte        { return tx.Data }
func (tx *TxData) gas() uint64         { return tx.Gas }
func (tx *TxData) gasPrice() *big.Int  { return tx.GasPrice }
func (tx *TxData) value() *big.Int     { return tx.Value }
func (tx *TxData) nonce() uint64       { return tx.Nonce }
func (tx *TxData) to() *common.Address { return tx.To }

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

type TxMessage struct {
	to       *common.Address
	from     *common.Address
	nonce    uint64
	amount   *big.Int
	gasLimit uint64
	gasPrice *big.Int
	data     []byte
	isFake   bool
}

func (m TxMessage) From() common.Address { return *m.from }
func (m TxMessage) To() *common.Address  { return m.to }
func (m TxMessage) GasPrice() *big.Int   { return m.gasPrice }
func (m TxMessage) Gas() uint64          { return m.gasLimit }
func (m TxMessage) Value() *big.Int      { return m.amount }
func (m TxMessage) Nonce() uint64        { return m.nonce }
func (m TxMessage) Data() []byte         { return m.data }
func (m TxMessage) IsFake() bool         { return m.isFake }

func (tx *Transaction) AsMessage(from *common.Address) (TxMessage, error) {
	msg := TxMessage{
		nonce:    tx.Nonce(),
		gasLimit: tx.Gas(),
		gasPrice: new(big.Int).Set(tx.GasPrice()),

		// from:       tx.,

		to:     tx.To(),
		amount: tx.Value(),
		data:   tx.Data(),
		isFake: false,
	}
	msg.from = from
	return msg, nil
}
