package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	deposittypes "github.com/ExocoreNetwork/exocore/x/deposit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DepositParams struct {
	ClientChainLzID uint64
	// The action field might need to be removed,it will be used when called from event hook.
	Action        assetstypes.CrossChainOpType
	AssetsAddress []byte
	StakerAddress []byte
	OpAmount      sdkmath.Int
}

// Deposit the deposit precompile contract will call this function to update asset state when there is a deposit.
func (k Keeper) Deposit(ctx sdk.Context, params *DepositParams) error {
	// check params parameter before executing deposit operation
	if params.OpAmount.IsNegative() {
		return errorsmod.Wrap(deposittypes.ErrInvalidDepositAmount, fmt.Sprintf("negative deposit amount:%s", params.OpAmount))
	}
	stakeID, assetID := assetstypes.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)
	assetsInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
	if err != nil {
		return err
	}
	if params.OpAmount.Add(assetsInfo.StakingTotalAmount).GT(assetsInfo.AssetBasicInfo.TotalSupply) {
		return errorsmod.Wrap(deposittypes.ErrInvalidDepositAmount, fmt.Sprintf("deposit amount will make the total staking amount greater than the total supply, amount:%s,totalStakingAmount:%s, totalSupply:%s", params.OpAmount, assetsInfo.StakingTotalAmount, assetsInfo.AssetBasicInfo.TotalSupply))
	}

	changeAmount := assetstypes.DeltaStakerSingleAsset{
		TotalDepositAmount: params.OpAmount,
		WithdrawableAmount: params.OpAmount,
	}
	// update asset state of the specified staker
	err = k.assetsKeeper.UpdateStakerAssetState(ctx, stakeID, assetID, changeAmount)
	if err != nil {
		return err
	}

	// update total amount of the deposited asset
	err = k.assetsKeeper.UpdateStakingAssetTotalAmount(ctx, assetID, params.OpAmount)
	if err != nil {
		return err
	}
	return nil
}
