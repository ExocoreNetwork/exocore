package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file provides all functions about operator assets state management.

func (k Keeper) GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address, _ map[string]interface{}) (assetsInfo map[string]*assetstype.OperatorSingleAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixOperatorAssetInfos)
	// the key is the operator address in the bech32 format
	key := []byte(operatorAddr.String())
	iterator := sdk.KVStorePrefixIterator(store, key)
	defer iterator.Close()

	ret := make(map[string]*assetstype.OperatorSingleAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo assetstype.OperatorSingleAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		keyList, err := assetstype.ParseJoinedStoreKey(iterator.Key(), 2)
		if err != nil {
			return nil, err
		}
		assetID := keyList[1]
		ret[assetID] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetID string) (info *assetstype.OperatorSingleAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixOperatorAssetInfos)
	key := assetstype.GetJoinedStoreKey(operatorAddr.String(), assetID)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, assetstype.ErrNoOperatorAssetKey
	}

	value := store.Get(key)

	ret := assetstype.OperatorSingleAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// UpdateOperatorAssetState It's used to update the operator states that include TotalAmount OperatorOwnAmount and WaitUndelegationAmount
// The input `changeAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is delegation or undelegation related to the operator. In the future,it will also be called when the operator deposit their own assets.

func (k Keeper) UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetID string, changeAmount assetstype.OperatorSingleAssetChangeInfo) (err error) {
	// get the latest state,use the default initial state if the state hasn't been stored
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixOperatorAssetInfos)
	key := assetstype.GetJoinedStoreKey(operatorAddr.String(), assetID)
	assetState := assetstype.OperatorSingleAssetInfo{
		TotalAmount:                        math.NewInt(0),
		OperatorOwnAmount:                  math.NewInt(0),
		WaitUnbondingAmount:                math.NewInt(0),
		OperatorUnbondingAmount:            math.NewInt(0),
		OperatorUnbondableAmountAfterSlash: math.NewInt(0),
	}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	// update all states of the specified operator asset
	err = assetstype.UpdateAssetValue(&assetState.TotalAmount, &changeAmount.ChangeForTotalAmount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState TotalAmountOrWantChangeValue error")
	}
	err = assetstype.UpdateAssetValue(&assetState.OperatorOwnAmount, &changeAmount.ChangeForOperatorOwn)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorOwnAmountOrWantChangeValue error")
	}
	err = assetstype.UpdateAssetValue(&assetState.WaitUnbondingAmount, &changeAmount.ChangeForWaitUnbonding)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState WaitUndelegationAmountOrWantChangeValue error")
	}
	err = assetstype.UpdateAssetValue(&assetState.OperatorUnbondingAmount, &changeAmount.ChangeForOperatorUnbonding)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorUnbondingAmount error")
	}
	err = assetstype.UpdateAssetValue(&assetState.OperatorUnbondableAmountAfterSlash, &changeAmount.ChangeForUnbondableAfterSlash)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAssetState OperatorUnbondingAmount error")
	}

	// store the updated state
	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)
	return nil
}

func (k Keeper) IteratorOperatorAssetState(ctx sdk.Context, f func(operatorAddr, assetID string, state *assetstype.OperatorSingleAssetInfo) error) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixOperatorAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amounts assetstype.OperatorSingleAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := assetstype.ParseJoinedKey(iterator.Key())
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
