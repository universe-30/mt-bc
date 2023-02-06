package chain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/universe-30/mt-bc/chain/types"
	"github.com/universe-30/mt-bc/chain/vm"
	"github.com/universe-30/mt-trie/common"
)

var TestTxOwner = common.Address{1}

type StateProcessor struct {
	bc *BlockChain
}

func NewStateProcessor(bc *BlockChain) *StateProcessor {
	return &StateProcessor{
		bc: bc,
	}
}

func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB) ([]*types.Receipt, uint64, error) {
	var (
		receipts    []*types.Receipt
		usedGas     = new(uint64)
		header      = block.Header()
		blockHash   = block.Hash()
		blockNumber = block.NumberU64()
		gp          = new(GasPool).AddGas(block.GasLimit())
	)

	blockContext := NewEVMBlockContext(header, nil)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb)
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		msg, err := tx.AsMessage(&TestTxOwner)
		if err != nil {
			return nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		statedb.Prepare(tx.Hash(), i)
		receipt, err := applyTransaction(msg, nil, gp, statedb, blockNumber, blockHash, tx, usedGas, vmenv)
		if err != nil {
			return nil, 0, fmt.Errorf("could not apply tx %d [%v]: %w", i, tx.Hash().Hex(), err)
		}
		receipts = append(receipts, receipt)
	}

	header.Root = statedb.IntermediateRoot(true)

	return receipts, *usedGas, nil
}

func applyTransaction(msg Message, author *common.Address, gp *GasPool, statedb *state.StateDB, blockNumber uint64, blockHash common.Hash, tx *types.Transaction, usedGas *uint64, evm *vm.EVM) (*types.Receipt, error) {

	// Create a new context to be used in the EVM environment.
	txContext := NewEVMTxContext(msg)
	evm.TxContext = txContext

	// Apply the transaction to the current state (included in the env).
	result, err := ApplyMessage(evm, msg, gp)
	if err != nil {
		return nil, err
	}

	// Update the state with pending changes.
	root := statedb.IntermediateRoot(true).Bytes()

	*usedGas += result.UsedGas

	// Create a new receipt for the transaction, storing the intermediate root and gas used
	// by the tx.
	receipt := &types.Receipt{PostHash: root, CumulativeGasUsed: *usedGas}
	if result.Failed() {
		receipt.Status = types.ReceiptStatusFailed
	} else {
		receipt.Status = types.ReceiptStatusSuccessful
	}
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = result.UsedGas

	// If the transaction created a contract, store the creation address in the receipt.
	if result.ContractAddr != nil {
		receipt.ContractAddress = *result.ContractAddr
	}

	// Set the receipt logs and create the bloom filter.
	// receipt.Logs = statedb.GetLogs(tx.Hash(), blockHash)
	receipt.BlockHash = blockHash
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(statedb.TxIndex())
	return receipt, err
}
