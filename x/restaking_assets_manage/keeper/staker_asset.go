package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateAssetValue It's used to update asset state,negative or positive `changeValue` represents a decrease or increase in the asset state
// newValue = valueToUpdate + changeVale
func UpdateAssetValue(valueToUpdate *math.Int, changeValue *math.Int) error {
	if valueToUpdate == nil || changeValue == nil {
		return errorsmod.Wrap(restakingtype.ErrInputPointerIsNil, fmt.Sprintf("valueToUpdate:%v,changeValue:%v", valueToUpdate, changeValue))
	}

	if !changeValue.IsNil() {
		if changeValue.IsNegative() {
			if valueToUpdate.LT(changeValue.Neg()) {
				return errorsmod.Wrap(restakingtype.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("valueToUpdate:%s,changeValue:%s", *valueToUpdate, *changeValue))
			}
		}
		if !changeValue.IsZero() {
			*valueToUpdate = valueToUpdate.Add(*changeValue)
		}
	}
	return nil
}

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerID string) (assetsInfo map[string]*restakingtype.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, []byte(stakerID))
	defer iterator.Close()

	ret := make(map[string]*restakingtype.StakerSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo restakingtype.StakerSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		_, assetID, err := restakingtype.ParseStakerAndAssetIDFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret[assetID] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerID string, assetID string) (info *restakingtype.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	key := restakingtype.GetAssetStateKey(stakerID, assetID)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, errorsmod.Wrap(restakingtype.ErrNoStakerAssetKey, fmt.Sprintf("the key is:%s", key))
	}

	value := store.Get(key)

	ret := restakingtype.StakerSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// UpdateStakerAssetState It's used to update the staker asset state
// The input `changeAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is deposit or withdraw related to the specified staker.
func (k Keeper) UpdateStakerAssetState(ctx sdk.Context, stakerID string, assetID string, changeAmount restakingtype.StakerSingleAssetOrChangeInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	key := restakingtype.GetAssetStateKey(stakerID, assetID)
	assetState := restakingtype.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     math.NewInt(0),
		CanWithdrawAmountOrWantChangeValue:      math.NewInt(0),
		WaitUndelegationAmountOrWantChangeValue: math.NewInt(0),
	}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	// update all states of the specified restaker asset
	err = UpdateAssetValue(&assetState.TotalDepositAmountOrWantChangeValue, &changeAmount.TotalDepositAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState TotalDepositAmountOrWantChangeValue error")
	}
	err = UpdateAssetValue(&assetState.CanWithdrawAmountOrWantChangeValue, &changeAmount.CanWithdrawAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState CanWithdrawAmountOrWantChangeValue error")
	}
	err = UpdateAssetValue(&assetState.WaitUndelegationAmountOrWantChangeValue, &changeAmount.WaitUndelegationAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState WaitUndelegationAmountOrWantChangeValue error")
	}

	// store the updated state
	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)

	return nil
}
