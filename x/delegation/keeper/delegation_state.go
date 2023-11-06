package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerId string, assetId string, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	c := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(c.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	amount := types2.ValueField{Amount: sdkmath.NewInt(0)}
	key := types.GetAssetStateKey(stakerId, assetId)
	isExit := store.Has(key)
	if isExit {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &amount)
	}

	if opAmount.IsNegative() {
		if amount.Amount.GT(opAmount.Neg()) {
			return errorsmod.Wrap(types2.ErrSubAmountIsGreaterThanOriginal, fmt.Sprintf("the opAmount is:%s,the originalAmount is:%s", opAmount, amount.Amount))
		}
	}
	amount.Amount = amount.Amount.Add(opAmount)
	bz := k.cdc.MustMarshal(&amount)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetStakerDelegationTotalAmount(ctx sdk.Context, stakerId string, assetId string) (opAmount sdkmath.Int, err error) {
	c := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(c.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	var ret types2.ValueField
	prefixKey := types.GetAssetStateKey(stakerId, assetId)
	isExit := store.Has(prefixKey)
	if !isExit {
		return sdkmath.Int{}, errorsmod.Wrap(types2.ErrNoKeyInTheStore, fmt.Sprintf("GetStakerDelegationTotalAmount: key is %s", prefixKey))
	} else {
		value := store.Get(prefixKey)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}

func (k Keeper) UpdateDelegationState(ctx sdk.Context, stakerId string, assetId string, delegationAmounts map[string]*types2.DelegationAmounts) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	//todo: think about the difference between init and update in future

	for opAddr, amounts := range delegationAmounts {
		if amounts.CanUnDelegationAmount.Amount.IsNil() && amounts.WaitUnDelegationAmount.Amount.IsNil() {
			continue
		}
		//check operator address validation
		_, err := sdk.AccAddressFromBech32(opAddr)
		if err != nil {
			return types2.OperatorAddrIsNotAccAddr
		}
		singleStateKey := types2.GetDelegationStateKey(stakerId, assetId, opAddr)
		isExit := store.Has(singleStateKey)
		delegationState := types2.DelegationAmounts{
			CanUnDelegationAmount:  &types2.ValueField{Amount: sdkmath.NewInt(0)},
			WaitUnDelegationAmount: &types2.ValueField{Amount: sdkmath.NewInt(0)},
		}
		if isExit {
			value := store.Get(singleStateKey)
			k.cdc.MustUnmarshal(value, &delegationState)
		}

		if !amounts.CanUnDelegationAmount.Amount.IsNil() {
			if amounts.CanUnDelegationAmount.Amount.IsNegative() {
				if delegationState.CanUnDelegationAmount.Amount.LT(amounts.CanUnDelegationAmount.Amount.Neg()) {
					return errorsmod.Wrap(types2.ErrSubAmountIsGreaterThanOriginal, fmt.Sprintf("UpdateDelegationState CanUnDelegationAmount the opAmount is:%s,the originalAmount is:%s", amounts.CanUnDelegationAmount.Amount, delegationState.CanUnDelegationAmount.Amount))
				}
			}
			delegationState.CanUnDelegationAmount.Amount = delegationState.CanUnDelegationAmount.Amount.Add(amounts.CanUnDelegationAmount.Amount)
		}

		if !amounts.WaitUnDelegationAmount.Amount.IsNil() {
			if amounts.WaitUnDelegationAmount.Amount.IsNegative() {
				if delegationState.WaitUnDelegationAmount.Amount.LT(amounts.WaitUnDelegationAmount.Amount.Neg()) {
					return errorsmod.Wrap(types2.ErrSubAmountIsGreaterThanOriginal, fmt.Sprintf("UpdateDelegationState WaitUnDelegationAmount the opAmount is:%s,the originalAmount is:%s", amounts.WaitUnDelegationAmount.Amount, delegationState.WaitUnDelegationAmount.Amount))
				}
			}
			delegationState.WaitUnDelegationAmount.Amount = delegationState.WaitUnDelegationAmount.Amount.Add(amounts.WaitUnDelegationAmount.Amount)
		}

		//save single operator delegation state
		bz := k.cdc.MustMarshal(&delegationState)
		store.Set(singleStateKey, bz)
	}
	return nil
}

func (k Keeper) GetSingleDelegationInfo(ctx sdk.Context, stakerId, assetId, operatorAddr string) (*types2.DelegationAmounts, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	singleStateKey := types2.GetDelegationStateKey(stakerId, assetId, operatorAddr)
	isExit := store.Has(singleStateKey)
	delegationState := types2.DelegationAmounts{}
	if isExit {
		value := store.Get(singleStateKey)
		k.cdc.MustUnmarshal(value, &delegationState)
	} else {
		return nil, errorsmod.Wrap(types2.ErrNoKeyInTheStore, fmt.Sprintf("QuerySingleDelegationInfo: key is %s", singleStateKey))
	}
	return &delegationState, nil
}

func (k Keeper) GetDelegationInfo(ctx sdk.Context, stakerId, assetId string) (*types2.QueryDelegationInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)

	var ret types2.QueryDelegationInfoResponse
	totalAmount, err := k.GetStakerDelegationTotalAmount(ctx, stakerId, assetId)
	if err != nil {
		return nil, err
	}
	ret.TotalDelegatedAmount.Amount = totalAmount

	store := prefix.NewStore(c.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, types2.GetDelegationStateIteratorPrefix(stakerId, assetId))
	defer iterator.Close()

	ret.DelegationInfos = make(map[string]*types2.DelegationAmounts, 0)
	for ; iterator.Valid(); iterator.Next() {
		var amounts types2.DelegationAmounts
		k.cdc.MustUnmarshal(iterator.Value(), &amounts)
		keys, err := types2.ParseStakerAssetIdAndOperatorAddrFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret.DelegationInfos[keys.OperatorAddr] = &amounts
	}

	return &ret, nil
}
