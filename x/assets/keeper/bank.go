package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DepositWithdrawParams struct {
	ClientChainLzID uint64
	Action          assetstypes.CrossChainOpType
	AssetsAddress   []byte
	StakerAddress   []byte
	OpAmount        sdkmath.Int
}

// PerformDepositOrWithdraw the assets precompile contract will call this function to update asset state
// when there is a deposit or withdraw.
func (k Keeper) PerformDepositOrWithdraw(ctx sdk.Context, params *DepositWithdrawParams) error {
	// check params parameter before executing deposit operation
	if params.OpAmount.IsNegative() {
		return errorsmod.Wrap(assetstypes.ErrInvalidDepositAmount, fmt.Sprintf("negative deposit amount:%s", params.OpAmount))
	}
	stakeID, assetID := assetstypes.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)

	actualOpAmount := params.OpAmount
	switch params.Action {
	case assetstypes.Deposit:
	case assetstypes.WithdrawPrincipal:
		actualOpAmount = actualOpAmount.Neg()
	default:
		return errorsmod.Wrapf(assetstypes.ErrInvalidOperationType, "the operation type is: %v", params.Action)
	}

	changeAmount := assetstypes.DeltaStakerSingleAsset{
		TotalDepositAmount: actualOpAmount,
		WithdrawableAmount: actualOpAmount,
	}
	// don't update staker info for exo-native-token
	// TODO: do we need additional process for exo-native-token ?
	if assetID != assetstypes.NativeAssetID {
		// update asset state of the specified staker
		err := k.UpdateStakerAssetState(ctx, stakeID, assetID, changeAmount)
		if err != nil {
			return errorsmod.Wrapf(err, "stakeID:%s assetID:%s", stakeID, assetID)
		}

		// update total amount of the deposited asset
		err = k.UpdateStakingAssetTotalAmount(ctx, assetID, actualOpAmount)
		if err != nil {
			return errorsmod.Wrapf(err, "assetID:%s", assetID)
		}
	}
	return nil
}
