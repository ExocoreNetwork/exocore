// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package deposit

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodDepositTo defines the ABI method name for the deposit
	// DepositTo transaction.
	MethodDepositTo = "DepositTo"
)

// DepositTo deposit the client chain assets to the staker, that will change the state in deposit module.
func (p Precompile) DepositTo(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	//check the invalidation of caller contract
	depositModuleParam, err := p.depositKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	exoCoreLzAppAddr := common.HexToAddress(depositModuleParam.ExoCoreLzAppAddress)
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
	}

	depositParams, err := GetDepositToParamsFromInputs(args)
	if err != nil {
		return nil, err
	}

	err = p.depositKeeper.Deposit(ctx, depositParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
