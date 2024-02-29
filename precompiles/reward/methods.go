package reward

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodReward defines the ABI method name for the reward
	//  transaction.
	MethodReward = "claimReward"
)

// Reward assets to the staker, that will change the state in reward module.
func (p Precompile) Reward(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	rewardModuleParam, err := p.rewardKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	exoCoreLzAppAddr := common.HexToAddress(rewardModuleParam.ExoCoreLzAppAddress)
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
	}

	rewardParam, err := p.GetRewardParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	err = p.rewardKeeper.RewardForWithdraw(ctx, rewardParam)
	if err != nil {
		return nil, err
	}
	// get the latest asset state of staker to return.
	stakerID, assetID := types.GetStakeIDAndAssetID(rewardParam.ClientChainLzID, rewardParam.WithdrawRewardAddress, rewardParam.AssetsAddress)
	info, err := p.stakingStateKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true, info.TotalDepositAmountOrWantChangeValue.BigInt())
}
