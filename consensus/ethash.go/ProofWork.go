package ethash

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/big"

	"github.com/universe-30/mt-bc/chain/types"
	"github.com/universe-30/mt-trie/common"
)

/**
target计算方式  假设：Hash为8位，targetBit为2位
eg:0000 0001(8位的Hash)
1.8-2 = 6 将上值左移6位
2.0000 0001 << 6 = 0100 0000 = target
3.只要计算的Hash满足 ：hash < target，便是符合POW的哈希值
*/

// 16 个 0
const targetBits = 16

type ProofOfWork struct {
	//工作量难度 big.Int大数存储
	target *big.Int
}

func NewProofOfWork() *ProofOfWork {

	//1.创建一个初始值为1的target
	target := big.NewInt(1)
	//2.左移bits(Hash) - targetBit 位
	target = target.Lsh(target, 256-targetBits)

	// target  = new(big.Int).Div(two256, header.Difficulty)

	return &ProofOfWork{target}
}

func (pow *ProofOfWork) prepareData(block *types.Block) []byte {

	data := bytes.Join(
		[][]byte{
			block.Hash().Bytes(),
			IntToHex(int64(targetBits)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Seal(block *types.Block) (common.Hash, types.BlockNonce, error) {

	// Start generating random nonces until we abort or find a good one
	var (
		attempts  = int64(0)
		powBuffer = new(big.Int)
		target    = pow.target
	)

	var nonce uint64
	nonce = 0
	// nonce    = seed

	//准备数据
	dataBytes := pow.prepareData(block)

	// logger.Trace("Started ethash search for new nonces", "seed", seed)

	for {
		attempts++

		// Compute the PoW value of this nonce
		hash := hashimotoFull(dataBytes, nonce)
		if powBuffer.SetBytes(hash[:]).Cmp(target) <= 0 {
			return common.BytesToHash(hash[:]), types.BlockNonce(nonce), nil
		}

		nonce++
	}

	return common.Hash{}, 0, nil
}

// hashimotoFull aggregates data from the full dataset (using the full in-memory
// dataset) in order to produce our final value for a particular header hash and
// nonce.
func hashimotoFull(dataset []byte, nonce uint64) [32]byte {
	// Compute the PoW value of this nonce

	pad := IntToHex(int64(nonce))

	data := bytes.Join(
		[][]byte{
			dataset,
			pad,
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	return hash
}

func UintToHex(num uint64) []byte {

	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func IntToHex(num int64) []byte {

	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {

		log.Panic(err)
	}

	return buff.Bytes()
}
