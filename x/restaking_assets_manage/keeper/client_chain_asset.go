package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) UpdateStakingAssetTotalAmount(ctx sdk.Context, assetId string, changeAmount sdkmath.Int) (err error) {
	if changeAmount.IsNil() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakingAssetInfo)
	key := []byte(assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return types2.ErrNoClientChainAssetKey
	}

	value := store.Get(key)

	ret := types2.StakingAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)

	if changeAmount.IsNegative() {
		if ret.StakingTotalAmount.LT(changeAmount.Neg()) {
			return errorsmod.Wrap(types2.ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("StakingTotalAmount:%s,changeValue:%s", ret.StakingTotalAmount, changeAmount))
		}
	}
	ret.StakingTotalAmount = ret.StakingTotalAmount.Add(changeAmount)

	bz := k.cdc.MustMarshal(&ret)

	store.Set(key, bz)

	return nil
}

// SetStakingAssetInfo todo: Temporarily use clientChainAssetAddr+'_'+layerZeroChainId as the key.
func (k Keeper) SetStakingAssetInfo(ctx sdk.Context, info *types2.StakingAssetInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakingAssetInfo)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	_, assetId := types2.GetStakeIDAndAssetIdFromStr(info.AssetBasicInfo.LayerZeroChainId, "", info.AssetBasicInfo.Address)
	store.Set([]byte(assetId), bz)
	return nil
}

func (k Keeper) StakingAssetIsExist(ctx sdk.Context, assetId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakingAssetInfo)
	return store.Has([]byte(assetId))
}

func (k Keeper) GetStakingAssetInfo(ctx sdk.Context, assetId string) (info *types2.StakingAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakingAssetInfo)
	ifExist := store.Has([]byte(assetId))
	if !ifExist {
		return nil, types2.ErrNoClientChainAssetKey
	}

	value := store.Get([]byte(assetId))

	ret := types2.StakingAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*types2.StakingAssetInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types2.KeyPrefixReStakingAssetInfo)
	defer iterator.Close()

	ret := make(map[string]*types2.StakingAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var assetInfo types2.StakingAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &assetInfo)
		_, assetId := types2.GetStakeIDAndAssetIdFromStr(assetInfo.AssetBasicInfo.LayerZeroChainId, "", assetInfo.AssetBasicInfo.Address)
		ret[assetId] = &assetInfo
	}
	return ret, nil
}
