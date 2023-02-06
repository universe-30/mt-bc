package vm

import (
	"math/big"

	"github.com/universe-30/mt-trie/common"
)

type StateDB interface {
	CreateAccount(common.Address)
	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	ExistAccount(common.Address) bool

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	RevertToSnapshot(int)
	Snapshot() int
}
