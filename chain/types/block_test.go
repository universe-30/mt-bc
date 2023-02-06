package types

import (
	"fmt"
	"testing"
)

func TestHashBlock(t *testing.T) {

	header := &Header{}
	header.Number = 100
	header.Time = 42
	// header.ParentHash = prev.Hash()

	b := &Block{header: header}
	b.Txs = nil

	hash := b.Hash()
	fmt.Println(hash)
	if len(hash) != 32 {
		t.Fatal("Hashing block failed.")
	}
}
