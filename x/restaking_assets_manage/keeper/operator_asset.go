package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

// This file provides all functions about operator assets state management.

func (k Keeper) GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]*restakingtype.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixOperatorAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, operatorAddr.Bytes())
	defer iterator.Close()

	ret := make(map[string]*restakingtype.OperatorSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo restakingtype.OperatorSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		_, assetId, err := restakingtype.ParseStakerAndAssetIdFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (info *restakingtype.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixOperatorAssetInfos)
	key := restakingtype.GetAssetStateKey(operatorAddr.String(), assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, restakingtype.ErrNoOperatorAssetKey
	}

	value := store.Get(key)

	ret := restakingtype.OperatorSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// UpdateOperatorAssetState It's used to update the operator states that include TotalAmount OperatorOwnAmount and WaitUndelegationAmount
// The input `changeAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is delegation or undelegation related to the operator. In the future,it will also be called when the operator deposit their own assets.

func (k Keeper) UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetId string, changeAmount restakingtype.OperatorSingleAssetOrChangeInfo) (err error) {
	//get the latest state,use the default initial state if the state hasn't been stored
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixOperatorAssetInfos)
	key := restakingtype.GetAssetStateKey(operatorAddr.String(), assetId)
	assetState := restakingtype.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue:            math.NewInt(0),
		OperatorOwnAmountOrWantChangeValue:      math.NewInt(0),
		WaitUndelegationAmountOrWantChangeValue: math.NewInt(0),
	}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	// update all states of the specified operator asset
	err = UpdateAssetValue(&assetState.TotalAmountOrWantChangeValue, &changeAmount.TotalAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState TotalAmountOrWantChangeValue error")
	}
	err = UpdateAssetValue(&assetState.OperatorOwnAmountOrWantChangeValue, &changeAmount.OperatorOwnAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorOwnAmountOrWantChangeValue error")
	}
	err = UpdateAssetValue(&assetState.WaitUndelegationAmountOrWantChangeValue, &changeAmount.WaitUndelegationAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState WaitUndelegationAmountOrWantChangeValue error")
	}

	//store the updated state
	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)
	return nil
}

// GetOperatorAssetOptedInMiddleWare This function should be implemented in the operator opt-in module
func (k Keeper) GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}
