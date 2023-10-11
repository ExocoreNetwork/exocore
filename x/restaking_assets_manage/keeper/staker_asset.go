package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"strings"
)

func GetStakerAssetStateKey(stakerId, assetId string) []byte {
	return []byte(strings.Join([]string{stakerId, assetId}, "_"))
}

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]*types2.StakerSingleAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakerAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, types2.KeyPrefixReStakerAssetInfos)
	defer iterator.Close()

	ret := make(map[string]*types2.StakerSingleAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo types2.StakerSingleAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		stringList := strings.SplitAfter(string(iterator.Key()), "_")
		assetId := stringList[len(stringList)-1]
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerId string, assetId string) (info *types2.StakerSingleAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakerAssetInfos)
	key := GetStakerAssetStateKey(stakerId, assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, types2.ErrNoStakerAssetKey
	}

	value := store.Get(key)

	ret := types2.StakerSingleAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) UpdateStakerAssetsState(ctx sdk.Context, stakerId string, assetsUpdate map[string]types2.StakerSingleAssetInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakerAssetInfos)
	for assetId, changeAmount := range assetsUpdate {
		key := GetStakerAssetStateKey(stakerId, assetId)
		isExit := store.Has(key)
		assetState := types2.StakerSingleAssetInfo{
			TotalDepositAmount: math.NewInt(0),
			CanWithdrawAmount:  math.NewInt(0),
		}
		if isExit {
			value := store.Get(key)
			k.cdc.MustUnmarshal(value, &assetState)
		}

		if changeAmount.TotalDepositAmount.IsZero() && changeAmount.CanWithdrawAmount.IsZero() {
			return types2.ErrInputUpdateStateIsZero
		}

		if changeAmount.TotalDepositAmount.IsNegative() {
			if assetState.TotalDepositAmount.LT(changeAmount.TotalDepositAmount.Abs()) {
				return types2.ErrSubDepositAmountIsMoreThanOrigin
			}
		}
		if changeAmount.CanWithdrawAmount.IsNegative() {
			if assetState.CanWithdrawAmount.LT(changeAmount.CanWithdrawAmount.Abs()) {
				return types2.ErrSubCanWithdrawAmountIsMoreThanOrigin
			}
		}

		bz := k.cdc.MustMarshal(&assetState)
		store.Set(key, bz)
	}
	return nil
}
