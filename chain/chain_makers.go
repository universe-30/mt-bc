package chain

import (
	"errors"
	"log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/universe-30/mt-bc/chain/types"
)

func (bc *BlockChain) insertChain(chain []*types.Block) error {

	for _, b := range chain {
		if err := bc.InsertBlock(b); err != nil {
			return err
		}
	}
	return nil
}

func (bc *BlockChain) InsertBlock(block *types.Block) error {

	err := bc.validate(block)
	if err != nil {
		return err
	}

	// receipts, usedGas, err := bc.processor.Process(block, statedb, bc.vmConfig)

	if err := bc.writeBlockWithState(block, receipts, statedb); err != nil {
		return err
	}

	err = bc.blockSetHead(block)
	return err
}

func (bc *BlockChain) writeBlockWithState(block *types.Block, receipts []*types.Receipt, state *state.StateDB) error {

	// Irrelevant of the canonical status, write the block itself to the database.
	//
	// Note all the components of block(td, hash->number map, header, body, receipts)
	// should be written atomically. BlockBatch is used for containing all components.
	blockBatch := bc.db.NewBatch()
	rawdb.WriteBlock(blockBatch, block)
	rawdb.WriteReceipts(blockBatch, block.Hash(), block.NumberU64(), receipts)
	if err := blockBatch.Write(); err != nil {
		log.Fatal("Failed to write block into disk", "err", err)
	}
	// Commit all cached state changes into underlying memory database.
	root, err := state.Commit(true)
	if err != nil {
		return err
	}
	triedb := bc.db.TrieDB()

	return triedb.Commit(root, false, nil)
}

func (bc *BlockChain) validate(b *types.Block) error {

	currentBlock := bc.CurrentBlock()

	if !isValid(b, currentBlock) {
		return errors.New("Invalid Block.")
	}

	return nil
}

func isValid(newBlock *types.Block, prevBlock *types.Block) bool {

	checkNum := prevBlock.NumberU64() + 1
	if newBlock.NumberU64() != checkNum {
		return false
	}
	if newBlock.ParentHash() != prevBlock.Hash() {
		return false
	}
	return true
}

func (bc *BlockChain) blockSetHead(block *types.Block) (err error) {

	currentBlock := bc.CurrentBlock()
	reorg, err := bc.ReorgNeeded(currentBlock, block)
	if err != nil {
		return err
	}

	if reorg {
		// Reorganise the chain if the parent is not the head block
		if block.ParentHash() != currentBlock.Hash() {
			if err := bc.reorg(currentBlock, block); err != nil {
				return err
			}
		}
		// CanonStatTy

		// Set new head.
		bc.writeHeadBlock(block)
	}
	return nil
}

func (bc *BlockChain) writeHeadBlock(block *types.Block) {

	bc.currentBlock.Store(block)
}

func (bc *BlockChain) ReorgNeeded(current *types.Block, block *types.Block) (bool, error) {

	reorg := false
	return reorg, nil
}

func (bc *BlockChain) reorg(oldBlock, newBlock *types.Block) error {
	return nil

}
