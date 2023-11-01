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

func (k Keeper) UpdateDelegationState(ctx sdk.Context, stakerId string, assetId string, operatorAndAmounts map[string]sdkmath.Int) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	//todo: think about the difference between init and update in future
	//key := common.HexToAddress(incentive.Contract)
	assetTotalAmountKey := types.GetAssetStateKey(stakerId, assetId)
	assetTotalAmount := types2.ValueField{Amount: sdkmath.NewInt(0)}
	isExit := store.Has(assetTotalAmountKey)
	if isExit {
		value := store.Get(assetTotalAmountKey)
		k.cdc.MustUnmarshal(value, &assetTotalAmount)
	}

	for opAddr, amount := range operatorAndAmounts {
		if amount.IsNil() {
			continue
		}
		//check operator address validation
		_, err := sdk.AccAddressFromBech32(opAddr)
		if err != nil {
			return types2.OperatorAddrIsNotAccAddr
		}
		singleStateKey := types2.GetDelegationStateKey(stakerId, assetId, opAddr)
		isExit := store.Has(singleStateKey)
		delegationState := types2.ValueField{Amount: sdkmath.NewInt(0)}
		if isExit {
			value := store.Get(singleStateKey)
			k.cdc.MustUnmarshal(value, &delegationState)
		}
		if amount.IsNegative() {
			if delegationState.Amount.LT(amount.Neg()) {
				return types2.ErrSubAmountIsGreaterThanOriginal
			}
		}
		delegationState.Amount = delegationState.Amount.Add(amount)
		//save single operator delegation state
		bz := k.cdc.MustMarshal(&delegationState)
		store.Set(singleStateKey, bz)

		//add amount to total delegation amount of the same asset
		assetTotalAmount.Amount = assetTotalAmount.Amount.Add(amount)
	}

	bz := k.cdc.MustMarshal(&assetTotalAmount)
	store.Set(assetTotalAmountKey, bz)
	return nil
}

func (k Keeper) GetSingleDelegationInfo(ctx sdk.Context, stakerId, assetId, operatorAddr string) (*types2.ValueField, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	singleStateKey := types2.GetDelegationStateKey(stakerId, assetId, operatorAddr)
	isExit := store.Has(singleStateKey)
	delegationState := types2.ValueField{}
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
	store := prefix.NewStore(c.KVStore(k.storeKey), types2.KeyPrefixRestakerDelegationInfo)
	var ret types2.QueryDelegationInfoResponse
	prefixKey := types.GetAssetStateKey(stakerId, assetId)
	isExit := store.Has(prefixKey)
	if !isExit {
		return nil, errorsmod.Wrap(types2.ErrNoKeyInTheStore, fmt.Sprintf("QueryDelegationInfo: key is %s", prefixKey))
	} else {
		value := store.Get(prefixKey)
		k.cdc.MustUnmarshal(value, ret.TotalDelegatedAmount)
	}

	iterator := sdk.KVStorePrefixIterator(store, prefixKey)
	defer iterator.Close()

	ret.DelegationInfos = make(map[string]*types2.ValueField, 0)
	for ; iterator.Valid(); iterator.Next() {
		var amount types2.ValueField
		k.cdc.MustUnmarshal(iterator.Value(), &amount)
		keys, err := types2.ParseStakerAssetIdAndOperatorAddrFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret.DelegationInfos[keys.OperatorAddr] = &amount
	}

	return &ret, nil
}
