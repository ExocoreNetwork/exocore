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

type DelegationOpFunc func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) error

// IterateDelegationsForStakerAndAsset processes all operations
// that require iterating over delegations for a specified staker and asset.
func (k Keeper) IterateDelegationsForStakerAndAsset(ctx sdk.Context, stakerID string, assetID string, opFunc DelegationOpFunc) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, delegationtype.GetDelegationStateIteratorPrefix(stakerID, assetID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amounts delegationtype.DelegationAmounts
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := delegationtype.ParseStakerAssetIDAndOperatorAddrFromKey(iterator.Key())
		if err != nil {
			return err
		}
		err = opFunc(keys, &amounts)
		if err != nil {
			return err
		}
	}
	return nil
}

// StakerDelegatedTotalAmount query the total delegation amount of the specified staker and asset.
// It needs to be calculated from the share and amount of the asset pool.
func (k Keeper) StakerDelegatedTotalAmount(ctx sdk.Context, stakerID string, assetID string) (amount sdkmath.Int, err error) {
	amount = sdkmath.NewInt(0)
	opFunc := func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) error {
		if amounts.UndelegatableShare.IsZero() {
			return nil
		}
		opAccAddr := sdk.MustAccAddressFromBech32(keys.GetOperatorAddr())
		// get the asset state of operator
		operatorAsset, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, opAccAddr, assetID)
		if err != nil {
			return err
		}
		singleAmount, err := TokensFromShares(amounts.UndelegatableShare, operatorAsset.TotalShare, operatorAsset.TotalAmount)
		if err != nil {
			return err
		}
		amount = amount.Add(singleAmount)
		return nil
	}
	err = k.IterateDelegationsForStakerAndAsset(ctx, stakerID, assetID, opFunc)
	if err != nil {
		return amount, err
	}
	return amount, nil
}

// AllDelegatedAmountForStakerAsset returns all delegated amount of the specified staker and asset
// the key of return value is the operator address, and the value is the asset amount.
func (k *Keeper) AllDelegatedAmountForStakerAsset(ctx sdk.Context, stakerID string, assetID string) (map[string]sdkmath.Int, error) {
	ret := make(map[string]sdkmath.Int)
	opFunc := func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) error {
		opAccAddr := sdk.MustAccAddressFromBech32(keys.GetOperatorAddr())
		// get the asset state of operator
		operatorAsset, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, opAccAddr, assetID)
		if err != nil {
			return err
		}
		singleAmount, err := TokensFromShares(amounts.UndelegatableShare, operatorAsset.TotalShare, operatorAsset.TotalAmount)
		if err != nil {
			return err
		}
		ret[keys.OperatorAddr] = singleAmount
		return nil
	}
	err := k.IterateDelegationsForStakerAndAsset(ctx, stakerID, assetID, opFunc)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpdateDelegationState is used to update the staker's asset amount that is delegated to a specified operator.
// Compared to `UpdateStakerDelegationTotalAmount`,they use the same kv store, but in this function the store key needs to add the operator address as a suffix.
func (k Keeper) UpdateDelegationState(ctx sdk.Context, stakerID, assetID, opAddr string, deltaAmounts *delegationtype.DeltaDelegationAmounts) (bool, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	// todo: think about the difference between init and update in future
	shareIsZero := false
	if deltaAmounts == nil {
		return false, errorsmod.Wrap(
			assetstype.ErrInputPointerIsNil,
			fmt.Sprintf("UpdateDelegationState opAddr:%v,deltaAmounts:%v", opAddr, deltaAmounts),
		)
	}
	// check operator address validation
	_, err := sdk.AccAddressFromBech32(opAddr)
	if err != nil {
		return shareIsZero, delegationtype.OperatorAddrIsNotAccAddr
	}
	singleStateKey := assetstype.GetJoinedStoreKey(stakerID, assetID, opAddr)
	delegationState := delegationtype.DelegationAmounts{
		WaitUndelegationAmount: sdkmath.NewInt(0),
		UndelegatableShare:     sdkmath.LegacyNewDec(0),
	}

	if store.Has(singleStateKey) {
		value := store.Get(singleStateKey)
		k.cdc.MustUnmarshal(value, &delegationState)
	}

	err = assetstype.UpdateAssetValue(&delegationState.WaitUndelegationAmount, &deltaAmounts.WaitUndelegationAmount)
	if err != nil {
		return shareIsZero, errorsmod.Wrap(err, "UpdateDelegationState WaitUndelegationAmount error")
	}

	err = assetstype.UpdateAssetDecValue(&delegationState.UndelegatableShare, &deltaAmounts.UndelegatableShare)
	if err != nil {
		return shareIsZero, errorsmod.Wrap(err, "UpdateDelegationState UndelegatableShare error")
	}

	if delegationState.UndelegatableShare.IsZero() {
		shareIsZero = true
	}

	// todo: should we delete the delegation state if both the share and the WaitUndelegationAmount are zero
	// to reduce the state storage?

	// save single operator delegation state
	bz := k.cdc.MustMarshal(&delegationState)
	store.Set(singleStateKey, bz)

	return shareIsZero, nil
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
	ret.DelegationInfos = make(map[string]*delegationtype.DelegationAmounts, 0)
	opFunc := func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) error {
		ret.DelegationInfos[keys.OperatorAddr] = amounts
		return nil
	}
	err := k.IterateDelegationsForStakerAndAsset(ctx, stakerID, assetID, opFunc)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (k *Keeper) SetStakersForOperator(ctx sdk.Context, operator, assetID, stakerID string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixStakersByOperator)
	Key := assetstype.GetJoinedStoreKey(operator, assetID)
	stakerMap := delegationtype.StakerMap{
		Stakers: make(map[string]bool),
	}
	if store.Has(Key) {
		value := store.Get(Key)
		k.cdc.MustUnmarshal(value, &stakerMap)
	}
	if _, ok := stakerMap.Stakers[stakerID]; !ok {
		stakerMap.Stakers[stakerID] = true
		bz := k.cdc.MustMarshal(&stakerMap)
		store.Set(Key, bz)
	}
	return nil
}

func (k *Keeper) DeleteStakerForOperator(ctx sdk.Context, operator, assetID, stakerID string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixStakersByOperator)
	Key := assetstype.GetJoinedStoreKey(operator, assetID)
	stakerMap := delegationtype.StakerMap{
		Stakers: make(map[string]bool),
	}
	if !store.Has(Key) {
		return delegationtype.ErrNoKeyInTheStore
	}
	delete(stakerMap.Stakers, stakerID)
	bz := k.cdc.MustMarshal(&stakerMap)
	store.Set(Key, bz)
	return nil
}

func (k *Keeper) DeleteStakerMapForOperator(ctx sdk.Context, operator, assetID string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixStakersByOperator)
	Key := assetstype.GetJoinedStoreKey(operator, assetID)
	if !store.Has(Key) {
		return delegationtype.ErrNoKeyInTheStore
	}
	store.Delete(Key)
	return nil
}

func (k *Keeper) GetStakersByOperator(ctx sdk.Context, operator, assetID string) (delegationtype.StakerMap, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixStakersByOperator)
	Key := assetstype.GetJoinedStoreKey(operator, assetID)
	if !store.Has(Key) {
		return delegationtype.StakerMap{}, delegationtype.ErrNoKeyInTheStore
	}
	value := store.Get(Key)
	stakerMap := delegationtype.StakerMap{}
	k.cdc.MustUnmarshal(value, &stakerMap)
	return stakerMap, nil
}

func (k *Keeper) SetStakerShareToZero(ctx sdk.Context, operator, assetID string, stakerMap delegationtype.StakerMap) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	for stakerID := range stakerMap.Stakers {
		singleStateKey := assetstype.GetJoinedStoreKey(stakerID, assetID, operator)
		delegationState := delegationtype.DelegationAmounts{}
		if store.Has(singleStateKey) {
			value := store.Get(singleStateKey)
			k.cdc.MustUnmarshal(value, &delegationState)
			delegationState.UndelegatableShare = sdkmath.LegacyNewDec(0)
			bz := k.cdc.MustMarshal(&delegationState)
			store.Set(singleStateKey, bz)
		}
	}
	return nil
}

// DelegationStateByOperatorAssets get the specified assets state delegated to the specified operator
// assetsFilter: assetID->nil, it's used to filter the specified assets
// the first return value is a nested map, its type is: stakerID->assetID->DelegationAmounts
// It means all delegation information related to the specified operator and filtered by the specified asset IDs
func (k Keeper) DelegationStateByOperatorAssets(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error) {
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
