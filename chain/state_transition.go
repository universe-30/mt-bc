package chain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/universe-30/mt-bc/chain/vm"
	"github.com/universe-30/mt-trie/common"
)

type StateTransition struct {
	gp         *GasPool
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	evm        *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	To() *common.Address

	GasPrice() *big.Int
	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	// IsFake() bool
	Data() []byte
}

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	UsedGas      uint64 // Total used gas but include the refunded gas
	Err          error  // Any error encountered during the execution(listed in core/vm/errors.go)
	ReturnData   []byte // Returned data from evm(function result or data supplied with revert opcode)
	ContractAddr *common.Address
}

func (result *ExecutionResult) Failed() bool { return result.Err != nil }

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		msg:      msg,
		gasPrice: msg.GasPrice(),
		value:    msg.Value(),
		data:     msg.Data(),
		state:    evm.StateDB,
	}
}

func ApplyMessage(evm *vm.EVM, msg Message, gp *GasPool) (*ExecutionResult, error) {
	st := NewStateTransition(evm, msg, gp)
	return st.TransitionDb()
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To() == nil /* contract creation */ {
		return common.Address{}
	}
	return *st.msg.To()
}

func (st *StateTransition) TransitionDb() (*ExecutionResult, error) {
	// First check this message satisfies all consensus rules before
	// applying the message. The rules include these clauses
	//
	// 1. the nonce of the message caller is correct
	// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
	// 3. the amount of gas required is available in the block
	// 4. the purchased gas is enough to cover intrinsic usage
	// 5. there is no overflow when calculating intrinsic gas
	// 6. caller has enough balance to cover asset transfer for **topmost** call

	if err := st.preCheck(); err != nil {
		return nil, err
	}
	if err := st.buyGas(); err != nil {
		return nil, err
	}

	var (
		statedb          = st.state
		msg              = st.msg
		sender           = vm.AccountRef(msg.From())
		contractCreation = msg.To() == nil
	)

	// Set up the initial access list.
	var (
		ret          []byte
		vmerr        error // vm errors do not effect consensus and are therefore not assigned to err
		contractAddr *common.Address
	)

	if contractCreation {
		nonce := statedb.GetNonce(sender.Address())
		*contractAddr = crypto.CreateAddress(sender.Address(), nonce)
		ret, st.gas, vmerr = st.evm.Create(*contractAddr, sender, st.data, st.gas, st.value)
	} else {
		// Increment the nonce for the next transaction
		nonce := statedb.GetNonce(sender.Address())
		statedb.SetNonce(sender.Address(), nonce+1)
		ret, st.gas, vmerr = st.evm.Call(sender, st.to(), st.data, st.gas, st.value)
	}

	gasUsed := st.gasUsed()

	// Base Fee
	effectiveTip := st.gasPrice
	fee := new(big.Int).SetUint64(gasUsed)
	fee.Mul(fee, effectiveTip)
	statedb.AddBalance(st.evm.Context.Coinbase, fee)

	return &ExecutionResult{
		UsedGas:      gasUsed,
		Err:          vmerr,
		ReturnData:   ret,
		ContractAddr: contractAddr,
	}, nil
}

func (st *StateTransition) preCheck() error {

	return nil
}

func (st *StateTransition) buyGas() error {
	statedb := st.state
	msg := st.msg
	msgGas := msg.Gas()

	mgval := new(big.Int).SetUint64(msgGas)
	mgval = mgval.Mul(mgval, st.gasPrice)
	balanceCheck := mgval

	if have, want := statedb.GetBalance(msg.From()), balanceCheck; have.Cmp(want) < 0 {
		return fmt.Errorf("%w: address %v have %v want %v", ErrInsufficientFunds, msg.From().Hex(), have, want)
	}
	if err := st.gp.SubGas(msgGas); err != nil {
		return err
	}
	st.gas += msgGas

	st.initialGas = msgGas
	st.state.SubBalance(msg.From(), mgval)
	return nil
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}
