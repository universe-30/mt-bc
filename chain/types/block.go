package types

import (
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/universe-30/mt-trie/common"
	"github.com/universe-30/mt-trie/rlp"
)

type Block struct {
	Number *big.Int
	Time   uint64
	// PrevBlockHash
	ParentHash common.Hash
	// Data
	txs []*Transaction

	// caches
	hash atomic.Value
}

func (b *Block) String() {
	fmt.Printf("Number: %d \n", b.Number)
	fmt.Printf("ParentHash: %s \n", b.ParentHash)
	fmt.Printf("CurrHash: %s \n", b.Hash())
	// fmt.Printf("Data: %s \n", block.Data)
	fmt.Printf("Timestamp: %d \n", b.Time)
	fmt.Println()
}

// 生成新的区块
func CreateNewBlock(prev Block, txs []*Transaction) *Block {
	newBlock := Block{}
	newBlock.Number = new(big.Int).Add(prev.Number, big.NewInt(1))
	newBlock.Time = uint64(time.Now().Unix())
	newBlock.ParentHash = prev.Hash()
	newBlock.txs = txs

	return &newBlock
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Number *big.Int
	Time   uint64
	Txs    []*Transaction
}

func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	v := calculateHash(b)
	b.hash.Store(v)
	return v
}

func calculateHash(block *Block) common.Hash {

	// data := []byte(block.Number.String()+string(block.Time)) + block.ParentHash
	// for _, tx := range block.txs {
	// 	data = data + tx.Hash()
	// }
	// bytes := sha256.Sum256([]byte(data))

	// var h common.Hash
	// copy(h[:], bytes[:])

	h := rlpHash(block)
	return h
	// return hex.EncodeToString(bytes[:])
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb Block
	if err := s.Decode(&eb); err != nil {
		return err
	}
	*b = eb
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (obj *Block) EncodeRLP(_w io.Writer) error {
	w := rlp.NewEncoderBuffer(_w)
	_tmp0 := w.List()
	w.WriteBytes(obj.ParentHash[:])
	if obj.Number == nil {
		w.Write(rlp.EmptyString)
	} else {
		if obj.Number.Sign() == -1 {
			return rlp.ErrNegativeBigInt
		}
		w.WriteBigInt(obj.Number)
	}

	w.WriteUint64(obj.Time)

	rlp.Encode(w, obj.txs)

	w.ListEnd(_tmp0)

	return w.Flush()
}
