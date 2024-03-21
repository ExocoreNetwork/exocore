package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateStakerDelegationTotalAmount The function is used to update the delegation total amount of the specified staker and asset.
// The input `opAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is delegation or undelegation related to the specified staker and asset.
func (k *Keeper) UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerID string, assetID string, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	// use stakerID+'/'+assetID as the key of total delegation amount
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	amount := delegationtype.ValueField{Amount: sdkmath.NewInt(0)}
	key := assetstype.GetJoinedStoreKey(stakerID, assetID)
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &amount)
	}

	err := assetstype.UpdateAssetValue(&amount.Amount, &opAmount)
	if err != nil {
		return err
	}

	bz := k.cdc.MustMarshal(&amount)
	store.Set(key, bz)
	return nil
}

// GetStakerDelegationTotalAmount query the total delegation amount of the specified staker and asset.
func (k *Keeper) GetStakerDelegationTotalAmount(ctx sdk.Context, stakerID string, assetID string) (opAmount sdkmath.Int, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	var ret delegationtype.ValueField
	prefixKey := assetstype.GetJoinedStoreKey(stakerID, assetID)
	if !store.Has(prefixKey) {
		return sdkmath.Int{}, errorsmod.Wrap(delegationtype.ErrNoKeyInTheStore, fmt.Sprintf("GetStakerDelegationTotalAmount: key is %s", prefixKey))
	}
	value := store.Get(prefixKey)
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}

// UpdateDelegationState The function is used to update the staker's asset amount that is delegated to a specified operator.
// Compared to `UpdateStakerDelegationTotalAmount`,they use the same kv store, but in this function the store key needs to add the operator address as a suffix.
func (k *Keeper) UpdateDelegationState(ctx sdk.Context, stakerID string, assetID string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	// todo: think about the difference between init and update in future

	for opAddr, amounts := range delegationAmounts {
		if amounts == nil {
			continue
		}
		if amounts.UndelegatableAmount.IsNil() && amounts.WaitUndelegationAmount.IsNil() {
			continue
		}
		// check operator address validation
		_, err := sdk.AccAddressFromBech32(opAddr)
		if err != nil {
			return delegationtype.OperatorAddrIsNotAccAddr
		}
		singleStateKey := assetstype.GetJoinedStoreKey(stakerID, assetID, opAddr)
		delegationState := delegationtype.DelegationAmounts{
			UndelegatableAmount:     sdkmath.NewInt(0),
			WaitUndelegationAmount:  sdkmath.NewInt(0),
			UndelegatableAfterSlash: sdkmath.NewInt(0),
		}

		if store.Has(singleStateKey) {
			value := store.Get(singleStateKey)
			k.cdc.MustUnmarshal(value, &delegationState)
		}

		err = assetstype.UpdateAssetValue(&delegationState.UndelegatableAmount, &amounts.UndelegatableAmount)
		if err != nil {
			return errorsmod.Wrap(err, "UpdateDelegationState UndelegatableAmount error")
		}

		err = assetstype.UpdateAssetValue(&delegationState.WaitUndelegationAmount, &amounts.WaitUndelegationAmount)
		if err != nil {
			return errorsmod.Wrap(err, "UpdateDelegationState WaitUndelegationAmount error")
		}

		err = assetstype.UpdateAssetValue(&delegationState.UndelegatableAfterSlash, &amounts.UndelegatableAfterSlash)
		if err != nil {
			return errorsmod.Wrap(err, "UpdateDelegationState UndelegatableAfterSlash error")
		}

		// save single operator delegation state
		bz := k.cdc.MustMarshal(&delegationState)
		store.Set(singleStateKey, bz)
	}
	return nil
}

// GetSingleDelegationInfo query the staker's asset amount that has been delegated to the specified operator.
func (k *Keeper) GetSingleDelegationInfo(ctx sdk.Context, stakerID, assetID, operatorAddr string) (*delegationtype.DelegationAmounts, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	singleStateKey := assetstype.GetJoinedStoreKey(stakerID, assetID, operatorAddr)
	isExit := store.Has(singleStateKey)
	delegationState := delegationtype.DelegationAmounts{}
	if isExit {
		value := store.Get(singleStateKey)
		k.cdc.MustUnmarshal(value, &delegationState)
	} else {
		return nil, errorsmod.Wrap(delegationtype.ErrNoKeyInTheStore, fmt.Sprintf("QuerySingleDelegationInfo: key is %s", singleStateKey))
	}
	return &delegationState, nil
}

// GetDelegationInfo query the staker's asset info that has been delegated.
func (k *Keeper) GetDelegationInfo(ctx sdk.Context, stakerID, assetID string) (*delegationtype.QueryDelegationInfoResponse, error) {
	var ret delegationtype.QueryDelegationInfoResponse
	totalAmount, err := k.GetStakerDelegationTotalAmount(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}
	ret.TotalDelegatedAmount = totalAmount

	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, delegationtype.GetDelegationStateIteratorPrefix(stakerID, assetID))
	defer iterator.Close()

	ret.DelegationInfos = make(map[string]*delegationtype.DelegationAmounts, 0)
	for ; iterator.Valid(); iterator.Next() {
		var amounts delegationtype.DelegationAmounts
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := delegationtype.ParseStakerAssetIDAndOperatorAddrFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret.DelegationInfos[keys.OperatorAddr] = &amounts
	}

	return &ret, nil
}

// DelegationStateByOperatorAssets get the specified assets state delegated to the specified operator
// assetsFilter: assetID->nil, it's used to filter the specified assets
// the first return value is a nested map, its type is: stakerID->assetID->DelegationAmounts
// It means all delegation information related to the specified operator and filtered by the specified asset IDs
func (k *Keeper) DelegationStateByOperatorAssets(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	ret := make(map[string]map[string]delegationtype.DelegationAmounts, 0)
	for ; iterator.Valid(); iterator.Next() {
		var amounts delegationtype.DelegationAmounts
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := assetstype.ParseJoinedKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		if len(keys) != 3 {
			continue
		}
		restakerID, assetID, findOperatorAddr := keys[0], keys[1], keys[2]
		if operatorAddr != findOperatorAddr {
			continue
		}
		_, assetIDExist := assetsFilter[assetID]
		_, restakerIDExist := ret[restakerID]
		if assetIDExist {
			if !restakerIDExist {
				ret[restakerID] = make(map[string]delegationtype.DelegationAmounts)
			}
			ret[restakerID][assetID] = amounts
		}
	}
	return ret, nil
}

func (k *Keeper) IterateDelegationState(ctx sdk.Context, f func(restakerID, assetID, operatorAddr string, state *delegationtype.DelegationAmounts) error) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amounts delegationtype.DelegationAmounts
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := assetstype.ParseJoinedKey(iterator.Key())
		if err != nil {
			return err
		}
		if len(keys) == 3 {
			err = f(keys[0], keys[1], keys[2], &amounts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
