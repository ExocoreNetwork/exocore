package withdraw

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodWithdraw defines the ABI method name for the withdrawal transaction.
	MethodWithdraw = "withdrawPrinciple"
)

// Withdraw assets to the staker, that will change the state in withdraw module.
func (p Precompile) Withdraw(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	withdrawModuleParam, err := p.withdrawKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	exoCoreLzAppAddr := common.HexToAddress(withdrawModuleParam.ExoCoreLzAppAddress)
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
	}

	withdrawParam, err := p.GetWithdrawParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	err = p.withdrawKeeper.Withdraw(ctx, withdrawParam)
	if err != nil {
		return nil, err
	}
	// get the latest asset state of staker to return.
	stakerId, assetId := types.GetStakeIDAndAssetId(withdrawParam.ClientChainLzId, withdrawParam.WithdrawAddress, withdrawParam.AssetsAddress)
	info, err := p.stakingStateKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerId, assetId)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true, info.TotalDepositAmountOrWantChangeValue.BigInt())
}
