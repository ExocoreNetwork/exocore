package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateStakerDelegationTotalAmount The function is used to update the delegation total amount of the specified staker and asset.
// The input `opAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is delegation or undelegation related to the specified staker and asset.
func (k Keeper) UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerID string, assetID string, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	c := sdk.UnwrapSDKContext(ctx)
	// use stakerID+'/'+assetID as the key of total delegation amount
	store := prefix.NewStore(c.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	amount := delegationtype.ValueField{Amount: sdkmath.NewInt(0)}
	key := types.GetAssetStateKey(stakerID, assetID)
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &amount)
	}

	err := keeper.UpdateAssetValue(&amount.Amount, &opAmount)
	if err != nil {
		return err
	}

	bz := k.cdc.MustMarshal(&amount)
	store.Set(key, bz)
	return nil
}

// GetStakerDelegationTotalAmount query the total delegation amount of the specified staker and asset.
func (k Keeper) GetStakerDelegationTotalAmount(ctx sdk.Context, stakerID string, assetID string) (opAmount sdkmath.Int, err error) {
	c := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(c.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	var ret delegationtype.ValueField
	prefixKey := types.GetAssetStateKey(stakerID, assetID)
	if !store.Has(prefixKey) {
		return sdkmath.Int{}, errorsmod.Wrap(delegationtype.ErrNoKeyInTheStore, fmt.Sprintf("GetStakerDelegationTotalAmount: key is %s", prefixKey))
	}
	value := store.Get(prefixKey)
	k.cdc.MustUnmarshal(value, &ret)
	return ret.Amount, nil
}

// UpdateDelegationState The function is used to update the staker's asset amount that is delegated to a specified operator.
// Compared to `UpdateStakerDelegationTotalAmount`,they use the same kv store, but in this function the store key needs to add the operator address as a suffix.
func (k Keeper) UpdateDelegationState(ctx sdk.Context, stakerID string, assetID string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	// todo: think about the difference between init and update in future

	for opAddr, amounts := range delegationAmounts {
		if amounts == nil {
			continue
		}
		if amounts.CanUndelegationAmount.IsNil() && amounts.WaitUndelegationAmount.IsNil() {
			continue
		}

		_, err := sdk.AccAddressFromBech32(opAddr)
		if err != nil {
			return delegationtype.OperatorAddrIsNotAccAddr
		}
		singleStateKey := delegationtype.GetDelegationStateKey(stakerID, assetID, opAddr)
		delegationState := delegationtype.DelegationAmounts{
			CanUndelegationAmount:  sdkmath.NewInt(0),
			WaitUndelegationAmount: sdkmath.NewInt(0),
		}

		if store.Has(singleStateKey) {
			value := store.Get(singleStateKey)
			k.cdc.MustUnmarshal(value, &delegationState)
		}

		err = keeper.UpdateAssetValue(&delegationState.CanUndelegationAmount, &amounts.CanUndelegationAmount)
		if err != nil {
			return errorsmod.Wrap(err, "UpdateDelegationState CanUndelegationAmount error")
		}

		err = keeper.UpdateAssetValue(&delegationState.WaitUndelegationAmount, &amounts.WaitUndelegationAmount)
		if err != nil {
			return errorsmod.Wrap(err, "UpdateDelegationState WaitUndelegationAmount error")
		}

		bz := k.cdc.MustMarshal(&delegationState)
		store.Set(singleStateKey, bz)
	}
	return nil
}

// GetSingleDelegationInfo query the staker's asset amount that has been delegated to the specified operator.
func (k Keeper) GetSingleDelegationInfo(ctx sdk.Context, stakerID, assetID, operatorAddr string) (*delegationtype.DelegationAmounts, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	singleStateKey := delegationtype.GetDelegationStateKey(stakerID, assetID, operatorAddr)
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
func (k Keeper) GetDelegationInfo(ctx sdk.Context, stakerID, assetID string) (*delegationtype.QueryDelegationInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var ret delegationtype.QueryDelegationInfoResponse
	totalAmount, err := k.GetStakerDelegationTotalAmount(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}
	ret.TotalDelegatedAmount = totalAmount

	store := prefix.NewStore(c.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
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
