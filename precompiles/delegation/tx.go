package delegation

import (
	"fmt"
	"reflect"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	cmn "github.com/evmos/evmos/v14/precompiles/common"

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

	// MethodAssociateOperatorWithStaker defines the ABI method name for the
	// associateOperatorWithStaker transaction.
	MethodAssociateOperatorWithStaker = "associateOperatorWithStaker"

	// MethodDissociateOperatorFromStaker defines the ABI method name for the
	// dissociateOperatorFromStaker transaction.
	MethodDissociateOperatorFromStaker = "dissociateOperatorFromStaker"

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
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
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
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
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

func (p Precompile) AssociateOperatorWithStaker(
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
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
	}

	inputsLen := len(p.ABI.Methods[MethodAssociateOperatorWithStaker].Inputs)
	if len(args) != inputsLen {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	staker, ok := args[1].([]byte)
	if !ok || staker == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", args[1])
	}
	// TODO: In the future, the check should be the same as it is in delegation if using LayerZero to route the
	// message for non-EVM client chains, such as Solana.
	if len(staker) != common.AddressLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidEVMAddr, staker)
	}

	operator, ok := args[2].([]byte)
	if !ok || operator == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", args[2])
	}
	operatorAccAddr, err := sdk.AccAddressFromBech32(string(operator))
	if err != nil {
		return nil, fmt.Errorf("error occurred when parse the operator address from Bech32,the operator is:%s, error:%s ", operator, err.Error())
	}
	err = p.delegationKeeper.AssociateOperatorWithStaker(ctx, uint64(clientChainID), operatorAccAddr, common.Address(staker))
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) DissociateOperatorFromStaker(
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
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
	}
	inputsLen := len(p.ABI.Methods[MethodDissociateOperatorFromStaker].Inputs)
	if len(args) != inputsLen {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	staker, ok := args[1].([]byte)
	if !ok || staker == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", args[1])
	}
	// TODO: In the future, the check should be the same as it is in delegation if using LayerZero to route the
	// message for non-EVM client chains, such as Solana.
	if len(staker) != common.AddressLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidEVMAddr, staker)
	}

	err = p.delegationKeeper.DissociateOperatorFromStaker(ctx, uint64(clientChainID), common.Address(staker))
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
