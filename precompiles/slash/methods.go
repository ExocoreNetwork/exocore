// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package slash

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodSlash defines the ABI method name for the slash
	//  transaction.
	MethodSlash = "submitSlash"
)

// SubmitSlash Slash assets to the staker, that will change the state in slash module.
func (p Precompile) SubmitSlash(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	//check the invalidation of caller contract
	slashModuleParam, err := p.slashKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	exoCoreLzAppAddr := common.HexToAddress(slashModuleParam.ExoCoreLzAppAddress)
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
	}

	slashParam, err := GetSlashParamsFromInputs(args)
	if err != nil {
		return nil, err
	}

	err = p.slashKeeper.Slash(ctx, slashParam)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}
