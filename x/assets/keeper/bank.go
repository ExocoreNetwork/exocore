package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OpParams struct {
	ClientChainLzID uint64
	Action          assetstypes.CrossChainOpType
	AssetsAddress   []byte
	StakerAddress   []byte
	OpAmount        sdkmath.Int
}

// PerformDepositOrWithdraw the assets precompile contract will call this function to update asset state
// when there is a deposit or withdraw.
func (k Keeper) PerformDepositOrWithdraw(ctx sdk.Context, params *OpParams) error {
	// check params parameter before executing deposit operation
	if params.OpAmount.IsNegative() {
		return errorsmod.Wrap(assetstypes.ErrInvalidDepositAmount, fmt.Sprintf("negative deposit amount:%s", params.OpAmount))
	}
	stakeID, assetID := assetstypes.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)
	assetsInfo, err := k.GetStakingAssetInfo(ctx, assetID)
	if err != nil {
		return err
	}

	actualOpAmount := params.OpAmount
	switch params.Action {
	case assetstypes.Deposit:
		if params.OpAmount.Add(assetsInfo.StakingTotalAmount).GT(assetsInfo.AssetBasicInfo.TotalSupply) {
			return errorsmod.Wrapf(assetstypes.ErrInvalidDepositAmount, "deposit amount will make the total staking amount greater than the total supply, amount:%s,totalStakingAmount:%s, totalSupply:%s", params.OpAmount, assetsInfo.StakingTotalAmount, assetsInfo.AssetBasicInfo.TotalSupply)
		}
	case assetstypes.WithdrawPrinciple:
		actualOpAmount = actualOpAmount.Neg()
	default:
		return errorsmod.Wrapf(assetstypes.ErrInvalidOperationType, "the operation type is: %v", params.Action)
	}

	changeAmount := assetstypes.DeltaStakerSingleAsset{
		TotalDepositAmount: actualOpAmount,
		WithdrawableAmount: actualOpAmount,
	}
	// update asset state of the specified staker
	err = k.UpdateStakerAssetState(ctx, stakeID, assetID, changeAmount)
	if err != nil {
		return err
	}

	// update total amount of the deposited asset
	err = k.UpdateStakingAssetTotalAmount(ctx, assetID, params.OpAmount)
	if err != nil {
		return err
	}
	return nil
}
