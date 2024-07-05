package delegation

import (
	"fmt"
	"reflect"

	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"

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

	// MethodMarkSelfDelegatedOperator defines the ABI method name for the
	// markSelfDelegatedOperator transaction.
	MethodMarkSelfDelegatedOperator = "markSelfDelegatedOperator"

	// MethodUnmarkSelfDelegatedOperator defines the ABI method name for the
	// unmarkSelfDelegatedOperator transaction.
	MethodUnmarkSelfDelegatedOperator = "unmarkSelfDelegatedOperator"

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

func (p Precompile) MarkSelfDelegatedOperator(
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
	if len(args) != 3 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	staker, ok := args[1].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", args[1])
	}
	if !common.IsHexAddress(staker) {
		return nil, xerrors.Errorf(exocmn.ErrInvalidEVMAddr, staker)
	}

	operator, ok := args[2].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", args[2])
	}
	operatorAccAddr, err := sdk.AccAddressFromBech32(operator)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse the operator address from Bech32,the operator is:%s", operator))
	}
	err = p.delegationKeeper.MarkSelfDelegatedOperator(ctx, uint64(clientChainID), operatorAccAddr, common.HexToAddress(staker))
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) UnMarkSelfDelegatedOperator(
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
	if len(args) != 2 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	staker, ok := args[1].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", args[1])
	}
	if !common.IsHexAddress(staker) {
		return nil, xerrors.Errorf(exocmn.ErrInvalidEVMAddr, staker)
	}

	err = p.delegationKeeper.UnmarkSelfDelegatedOperator(ctx, uint64(clientChainID), common.HexToAddress(staker))
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
