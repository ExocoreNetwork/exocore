package delegation

import (
	"fmt"
	"reflect"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodDelegateToThroughClientChain defines the ABI method name for the
	// DelegateToThroughClientChain transaction.
	MethodDelegateToThroughClientChain = "delegateToThroughClientChain"

	// MethodUndelegateFromThroughClientChain defines the ABI method name for the
	// UndelegateFromThroughClientChain transaction.
	MethodUndelegateFromThroughClientChain = "undelegateFromThroughClientChain"

	CtxKeyTxHash = "TxHash"
)

// DelegateToThroughClientChain delegate the client chain assets to the operator through client chain, that will change the states in delegation and assets module
func (p Precompile) DelegateToThroughClientChain(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	err := p.assetsKeeper.CheckExocoreGatewayAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, exocmn.ErrContractCaller)
	}

	delegationParams, err := p.GetDelegationParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	err = p.delegationKeeper.DelegateTo(ctx, delegationParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// UndelegateFromThroughClientChain Undelegation the client chain assets from the operator through client chain, that will change the states in delegation and assets module
func (p Precompile) UndelegateFromThroughClientChain(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	err := p.assetsKeeper.CheckExocoreGatewayAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, exocmn.ErrContractCaller)
	}

	undelegationParams, err := p.GetDelegationParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	txHash, ok := ctx.Value(CtxKeyTxHash).(common.Hash)
	if !ok || txHash.Bytes() == nil {
		return nil, fmt.Errorf(ErrCtxTxHash, reflect.TypeOf(ctx.Value(CtxKeyTxHash)), txHash)
	}
	undelegationParams.TxHash = txHash

	err = p.delegationKeeper.UndelegateFrom(ctx, undelegationParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
