package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file provides all functions about operator assets state management.

func (k Keeper) GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address, assetsFilter map[string]interface{}) (assetsInfo map[string]*restakingtype.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixOperatorAssetInfos)
	// the key is the operator address in the bech32 format
	key := []byte(operatorAddr.String())
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()

	ret := make(map[string]*restakingtype.OperatorSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo restakingtype.OperatorSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		keyList, err := restakingtype.ParseJoinedStoreKey(iterator.Key(), 2)
		if err != nil {
			return nil, err
		}
		assetId := keyList[1]
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (info *restakingtype.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixOperatorAssetInfos)
	key := restakingtype.GetJoinedStoreKey(operatorAddr.String(), assetId)
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
	key := restakingtype.GetJoinedStoreKey(operatorAddr.String(), assetId)
	assetState := restakingtype.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue:            math.NewInt(0),
		OperatorOwnAmountOrWantChangeValue:      math.NewInt(0),
		WaitUnbondingAmountOrWantChangeValue:    math.NewInt(0),
		OperatorOwnWaitUnbondingAmount:          math.NewInt(0),
		OperatorOwnCanUnbondingAmountAfterSlash: math.NewInt(0),
	}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	// update all states of the specified operator asset
	err = restakingtype.UpdateAssetValue(&assetState.TotalAmountOrWantChangeValue, &changeAmount.TotalAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState TotalAmountOrWantChangeValue error")
	}
	err = restakingtype.UpdateAssetValue(&assetState.OperatorOwnAmountOrWantChangeValue, &changeAmount.OperatorOwnAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorOwnAmountOrWantChangeValue error")
	}
	err = restakingtype.UpdateAssetValue(&assetState.WaitUnbondingAmountOrWantChangeValue, &changeAmount.WaitUnbondingAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState WaitUndelegationAmountOrWantChangeValue error")
	}
	err = restakingtype.UpdateAssetValue(&assetState.OperatorOwnWaitUnbondingAmount, &changeAmount.OperatorOwnWaitUnbondingAmount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorOwnWaitUnbondingAmount error")
	}
	err = restakingtype.UpdateAssetValue(&assetState.OperatorOwnCanUnbondingAmountAfterSlash, &changeAmount.OperatorOwnCanUnbondingAmountAfterSlash)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorOwnWaitUnbondingAmount error")
	}

	//store the updated state
	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)
	return nil
}

func (k Keeper) IteratorOperatorAssetState(ctx sdk.Context, f func(operatorAddr, assetId string, state *restakingtype.OperatorSingleAssetOrChangeInfo) error) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixOperatorAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amounts restakingtype.OperatorSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := restakingtype.ParseJoinedKey(iterator.Key())
		if err != nil {
			return err
		}
		if len(keys) == 3 {
			err = f(keys[0], keys[1], &amounts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
