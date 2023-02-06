package vm

import (
	"math/big"

	"github.com/universe-30/mt-trie/common"
)

const (
	MaxCodeSize = 24576 // Maximum bytecode to permit for a contract

	CreateDataGas uint64 = 200 //

)

// Block information
type BlockContext struct {
	Coinbase common.Address // Provides information for COINBASE
	BaseFee  *big.Int       // Provides information for BASEFEE
}

type TxContext struct {
	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE
}

type EVM struct {
	Context   BlockContext
	TxContext TxContext
	StateDB   StateDB
}

// The returned EVM is not thread safe and should only ever be used *once*.
func NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB) *EVM {
	evm := &EVM{
		Context:   blockCtx,
		TxContext: txCtx,
		StateDB:   statedb,
	}
	return evm
}

func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64,
	value *big.Int) (ret []byte, leftOverGas uint64, err error) {

	stateDB := evm.StateDB

	// Fail if we're trying to transfer more than the available balance
	if value.Sign() > 0 && !CanTransfer(stateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	snapshot := stateDB.Snapshot()

	if !stateDB.ExistAccount(addr) {
		stateDB.CreateAccount(addr)
	}
	Transfer(stateDB, caller.Address(), addr, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	// code := stateDB.GetCode(addr)

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		stateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func CanTransfer(db StateDB, addr common.Address, amount *big.Int) bool {
	balance := db.GetBalance(addr)
	return balance.Cmp(amount) >= 0
}

func Transfer(db StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

func (evm *EVM) create(caller ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int,
	address common.Address) ([]byte, uint64, error) {

	stateDB := evm.StateDB

	if !CanTransfer(stateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	nonce := stateDB.GetNonce(caller.Address())
	stateDB.SetNonce(caller.Address(), nonce+1)

	// Ensure there's no existing contract already at the designated address
	contractHash := stateDB.GetCodeHash(address)
	if stateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := stateDB.Snapshot()
	stateDB.CreateAccount(address)
	stateDB.SetNonce(address, 1)

	Transfer(stateDB, caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	ret, err := evm.interpreter.Run(contract, nil, false)

	// Check whether the max code size has been exceeded, assign err if the case.
	if len(ret) > MaxCodeSize {
		err = ErrMaxCodeSizeExceeded
	}

	// if the contract creation ran successfully and no errors were returned
	// calculate the gas requix w xred to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil {
		createDataGas := uint64(len(ret)) * CreateDataGas
		if contract.UseGas(createDataGas) {
			stateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil && (evm.chainRules.IsHomestead || err != ErrCodeStoreOutOfGas) {
		stateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	return ret, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(contractAddr common.Address, caller ContractRef, code []byte, gas uint64,
	value *big.Int) ([]byte, uint64, error) {
	return evm.create(caller, &codeAndHash{code: code}, gas, value, contractAddr)
}
