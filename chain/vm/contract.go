package vm

import (
	"math/big"

	"github.com/universe-30/mt-trie/common"
)

type ContractRef interface {
	Address() common.Address
}

type AccountRef common.Address

// Address casts AccountRef to a Address
func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

type Contract struct {
	CallerAddress common.Address
	caller        ContractRef
	self          ContractRef

	Gas   uint64
	value *big.Int
}

// NewContract returns a new contract environment for the execution of EVM.
func NewContract(caller ContractRef, object ContractRef, value *big.Int, gas uint64) *Contract {
	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object}

	// Gas should be a pointer so it can safely be reduced through the run
	// This pointer will be off the state transition
	c.Gas = gas
	// ensures a value is set
	c.value = value

	return c
}
