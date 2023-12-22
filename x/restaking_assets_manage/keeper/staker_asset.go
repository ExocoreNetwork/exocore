package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

// UpdateAssetValue is used to update asset state,negative or positive `changeValue` represents a decrease or increase in the asset state
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

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]*restakingtype.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, []byte(stakerId))
	defer iterator.Close()

	ret := make(map[string]*restakingtype.StakerSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo restakingtype.StakerSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		_, assetId, err := restakingtype.ParseStakerAndAssetIdFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerId string, assetId string) (info *restakingtype.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	key := restakingtype.GetAssetStateKey(stakerId, assetId)
	isExist := store.Has(key)
	if !isExist {
		return nil, errorsmod.Wrap(restakingtype.ErrNoStakerAssetKey, fmt.Sprintf("the key is:%s", key))
	}

	value := store.Get(key)

	ret := restakingtype.StakerSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// UpdateStakerAssetState is used to update the staker asset state
func (k Keeper) UpdateStakerAssetState(ctx sdk.Context, stakerId string, assetId string, changeAmount restakingtype.StakerSingleAssetOrChangeInfo) (err error) {
	// get the latest state,use the default initial state if the state hasn't been stored
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	key := restakingtype.GetAssetStateKey(stakerId, assetId)
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
