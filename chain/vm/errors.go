package vm

import "errors"

var (
	ErrCodeStoreOutOfGas = errors.New("contract creation code storage out of gas")

	ErrInsufficientBalance      = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision = errors.New("contract address collision")

	ErrExecutionReverted = errors.New("execution reverted")

	ErrMaxCodeSizeExceeded = errors.New("max code size exceeded")
)
