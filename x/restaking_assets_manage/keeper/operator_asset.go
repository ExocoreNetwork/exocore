package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]*types2.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, operatorAddr.Bytes())
	defer iterator.Close()

	ret := make(map[string]*types2.OperatorSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo types2.OperatorSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		_, assetId, err := types2.ParseStakerAndAssetIdFromKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (info *types2.OperatorSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorAssetInfos)
	key := types2.GetAssetStateKey(operatorAddr.String(), assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, types2.ErrNoOperatorAssetKey
	}

	value := store.Get(key)

	ret := types2.OperatorSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetId string, changeAmount types2.OperatorSingleAssetOrChangeInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorAssetInfos)

	key := types2.GetAssetStateKey(operatorAddr.String(), assetId)
	isExit := store.Has(key)
	assetState := types2.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue:            math.NewInt(0),
		OperatorOwnAmountOrWantChangeValue:      math.NewInt(0),
		WaitUnDelegationAmountOrWantChangeValue: math.NewInt(0),
	}
	if isExit {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	if !changeAmount.TotalAmountOrWantChangeValue.IsNil() {
		if changeAmount.TotalAmountOrWantChangeValue.IsNegative() {
			if assetState.TotalAmountOrWantChangeValue.LT(changeAmount.TotalAmountOrWantChangeValue.Abs()) {
				return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("TotalAmount:%s,changeValue:%s", assetState.TotalAmountOrWantChangeValue, changeAmount.TotalAmountOrWantChangeValue))
			}
		}
		if !changeAmount.TotalAmountOrWantChangeValue.IsZero() {
			assetState.TotalAmountOrWantChangeValue = assetState.TotalAmountOrWantChangeValue.Add(changeAmount.TotalAmountOrWantChangeValue)
		}
	}

	if !changeAmount.OperatorOwnAmountOrWantChangeValue.IsNil() {
		if changeAmount.OperatorOwnAmountOrWantChangeValue.IsNegative() {
			if assetState.OperatorOwnAmountOrWantChangeValue.LT(changeAmount.OperatorOwnAmountOrWantChangeValue.Abs()) {
				return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("OperatorOwnAmount:%s,changeValue:%s", assetState.OperatorOwnAmountOrWantChangeValue, changeAmount.OperatorOwnAmountOrWantChangeValue))
			}
		}
		if !changeAmount.OperatorOwnAmountOrWantChangeValue.IsZero() {
			assetState.OperatorOwnAmountOrWantChangeValue = assetState.OperatorOwnAmountOrWantChangeValue.Add(changeAmount.OperatorOwnAmountOrWantChangeValue)
		}
	}

	if !changeAmount.WaitUnDelegationAmountOrWantChangeValue.IsNil() {
		if changeAmount.WaitUnDelegationAmountOrWantChangeValue.IsNegative() {
			if assetState.WaitUnDelegationAmountOrWantChangeValue.LT(changeAmount.WaitUnDelegationAmountOrWantChangeValue.Abs()) {
				return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("WaitUnDelegationAmount:%s,changeValue:%s", assetState.WaitUnDelegationAmountOrWantChangeValue, changeAmount.WaitUnDelegationAmountOrWantChangeValue))
			}
		}
		if !changeAmount.WaitUnDelegationAmountOrWantChangeValue.IsZero() {
			assetState.WaitUnDelegationAmountOrWantChangeValue = assetState.WaitUnDelegationAmountOrWantChangeValue.Add(changeAmount.WaitUnDelegationAmountOrWantChangeValue)
		}
	}

	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}
