package types

import (
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/universe-30/mt-trie/common"
	"github.com/universe-30/mt-trie/rlp"
)

type BlockNonce uint64

type Header struct {
	ParentHash common.Hash    `json:"parentHash"`
	Coinbase   common.Address `json:"miner"`
	Root       common.Hash    `json:"stateRoot"`
	TxHash     common.Hash    `json:"transactionsRoot"`

	GasLimit uint64     `json:"gasLimit"`
	Number   uint64     `json:"number"`
	Time     uint64     `json:"timestamp"`
	Nonce    BlockNonce `json:"nonce"`

	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`
}

type Block struct {
	header *Header

	// PrevBlockHash
	Txs []*Transaction

	// caches
	hash atomic.Value
}

func (b *Block) Transactions() []*Transaction { return b.Txs }

func (b *Block) NumberU64() uint64       { return b.header.Number }
func (b *Block) GasLimit() uint64        { return b.header.GasLimit }
func (b *Block) ParentHash() common.Hash { return b.header.ParentHash }

func (b *Block) Header() *Header { return b.header }

func (b *Block) SetFinal(txhash common.Hash, nonce BlockNonce) {
	b.header.TxHash = txhash
	b.header.Nonce = nonce
}

// 生成新的区块
func CreateNewBlock(prev *Block, txs []*Transaction) *Block {

	header := &Header{}
	header.Number = prev.NumberU64() + 1
	header.Time = uint64(time.Now().Unix())
	header.ParentHash = prev.Hash()

	blk := &Block{header: header}
	blk.Txs = txs

	return blk
}

func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: header}
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Number uint64
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

// DecodeRLP decodes
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb Block
	if err := s.Decode(&eb); err != nil {
		return err
	}
	*b = eb
	return nil
}

// EncodeRLP serializes b into the RLP block format.
func (obj *Block) EncodeRLP(_w io.Writer) error {
	w := rlp.NewEncoderBuffer(_w)
	_tmp0 := w.List()
	w.WriteBytes(obj.header.ParentHash[:])
	if obj.header.Number < 0 {
		return rlp.ErrNegativeBigInt
	}
	w.WriteUint64(obj.header.Number)

	w.WriteUint64(obj.header.Time)

	rlp.Encode(w, obj.Txs)

	w.ListEnd(_tmp0)

	return w.Flush()
}
