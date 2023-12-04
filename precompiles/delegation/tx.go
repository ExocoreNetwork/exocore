// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package delegation

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"reflect"
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

// DelegateToThroughClientChain delegate the client chain assets to the operator through client chain, that will change the states in delegation and restaking_assets_manage module
func (p Precompile) DelegateToThroughClientChain(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	//check the invalidation of caller contract
	exoCoreLzAppAddr, err := p.delegationKeeper.GetExoCoreLzAppAddress(ctx)
	if err != nil {
		return nil, err
	}
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
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

// UndelegateFromThroughClientChain Undelegation the client chain assets from the operator through client chain, that will change the states in delegation and restaking_assets_manage module
func (p Precompile) UndelegateFromThroughClientChain(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	//check the invalidation of caller contract
	exoCoreLzAppAddr, err := p.delegationKeeper.GetExoCoreLzAppAddress(ctx)
	if err != nil {
		return nil, err
	}
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
	}

	UndelegationParams, err := p.GetDelegationParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	txHash, ok := ctx.Value(CtxKeyTxHash).(common.Hash)
	if !ok || txHash.Bytes() == nil {
		return nil, fmt.Errorf(ErrCtxTxHash, reflect.TypeOf(ctx.Value(CtxKeyTxHash)), txHash)
	}
	UndelegationParams.TxHash = txHash

	err = p.delegationKeeper.UndelegateFrom(ctx, UndelegationParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
