package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"
	"github.com/exocore/x/restaking_assets_manage/types"
)

// UpdateStakerDelegationTotalAmount
func (k Keeper) UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerId string, assetId string, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	c := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(c.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	amount := delegationtype.ValueField{Amount: sdkmath.NewInt(0)}
	key := types.GetAssetStateKey(stakerId, assetId)
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &amount)
	}

	if opAmount.IsNegative() {
		if amount.Amount.GT(opAmount.Neg()) {
			return errorsmod.Wrap(delegationtype.ErrSubAmountIsGreaterThanOriginal, fmt.Sprintf("the OpAmount is:%s,the originalAmount is:%s", opAmount, amount.Amount))
		}
	}
	amount.Amount = amount.Amount.Add(opAmount)
	bz := k.cdc.MustMarshal(&amount)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetStakerDelegationTotalAmount(ctx sdk.Context, stakerId string, assetId string) (opAmount sdkmath.Int, err error) {
	c := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(c.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	var ret delegationtype.ValueField
	prefixKey := types.GetAssetStateKey(stakerId, assetId)
	isExit := store.Has(prefixKey)
	if !isExit {
		return sdkmath.Int{}, errorsmod.Wrap(delegationtype.ErrNoKeyInTheStore, fmt.Sprintf("GetStakerDelegationTotalAmount: key is %s", prefixKey))
	} else {
		value := store.Get(prefixKey)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}

func (k Keeper) UpdateDelegationState(ctx sdk.Context, stakerId string, assetId string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	// todo: think about the difference between init and update in future

	for opAddr, amounts := range delegationAmounts {
		if amounts == nil {
			continue
		}
		if amounts.CanUndelegationAmount.IsNil() && amounts.WaitUndelegationAmount.IsNil() {
			continue
		}
		// check operator address validation
		_, err := sdk.AccAddressFromBech32(opAddr)
		if err != nil {
			return delegationtype.OperatorAddrIsNotAccAddr
		}
		singleStateKey := delegationtype.GetDelegationStateKey(stakerId, assetId, opAddr)
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

		// save single operator delegation state
		bz := k.cdc.MustMarshal(&delegationState)
		store.Set(singleStateKey, bz)
	}
	return nil
}

func (k Keeper) GetSingleDelegationInfo(ctx sdk.Context, stakerId, assetId, operatorAddr string) (*delegationtype.DelegationAmounts, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	singleStateKey := delegationtype.GetDelegationStateKey(stakerId, assetId, operatorAddr)
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

func (k Keeper) GetDelegationInfo(ctx sdk.Context, stakerId, assetId string) (*delegationtype.QueryDelegationInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var ret delegationtype.QueryDelegationInfoResponse
	totalAmount, err := k.GetStakerDelegationTotalAmount(ctx, stakerId, assetId)
	if err != nil {
		return nil, err
	}
	ret.TotalDelegatedAmount = totalAmount

	store := prefix.NewStore(c.KVStore(k.storeKey), delegationtype.KeyPrefixRestakerDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, delegationtype.GetDelegationStateIteratorPrefix(stakerId, assetId))
	defer iterator.Close()

	ret.DelegationInfos = make(map[string]*delegationtype.DelegationAmounts, 0)
	for ; iterator.Valid(); iterator.Next() {
		var amounts delegationtype.DelegationAmounts
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := delegationtype.ParseStakerAssetIdAndOperatorAddrFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret.DelegationInfos[keys.OperatorAddr] = &amounts
	}

	return &ret, nil
}
