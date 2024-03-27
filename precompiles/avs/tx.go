package avs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// AVSRegister defines the ABI method name for the avs
	// related transaction.
	MethodAVSAction = "avsAction"
)

// AVSInfoRegister register the avs related information and change the state in avs keeper module.
func (p Precompile) AVSInfoRegisterOrDeregister(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// parse the avs input params first.
	avsParams, err := p.GetAVSParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
