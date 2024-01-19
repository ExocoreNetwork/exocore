package keeper

import (
	sdkmath "cosmossdk.io/math"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateStakingAssetTotalAmount updating the total deposited amount of a specified asset in exoCore chain
// The function will be called when stakers deposit and withdraw their assets
func (k Keeper) UpdateStakingAssetTotalAmount(ctx sdk.Context, assetId string, changeAmount sdkmath.Int) (err error) {
	if changeAmount.IsNil() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakingAssetInfo)
	key := []byte(assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return restakingtype.ErrNoClientChainAssetKey
	}

	value := store.Get(key)

	ret := restakingtype.StakingAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)

	//calculate and set new amount
	err = restakingtype.UpdateAssetValue(&ret.StakingTotalAmount, &changeAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&ret)

	store.Set(key, bz)

	return nil
}

// SetStakingAssetInfo todo: Temporarily use clientChainAssetAddr+'_'+layerZeroChainId as the key.
// It provides a function to register the client chain assets supported by exoCore.It's called by genesis configuration now,however it will be called by the governance in the future
func (k Keeper) SetStakingAssetInfo(ctx sdk.Context, info *restakingtype.StakingAssetInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakingAssetInfo)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	_, assetId := restakingtype.GetStakeIDAndAssetIdFromStr(info.AssetBasicInfo.LayerZeroChainId, "", info.AssetBasicInfo.Address)
	store.Set([]byte(assetId), bz)
	return nil
}

func (k Keeper) IsStakingAsset(ctx sdk.Context, assetId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakingAssetInfo)
	return store.Has([]byte(assetId))
}

func (k Keeper) GetStakingAssetInfo(ctx sdk.Context, assetId string) (info *restakingtype.StakingAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakingAssetInfo)
	ifExist := store.Has([]byte(assetId))
	if !ifExist {
		return nil, restakingtype.ErrNoClientChainAssetKey
	}

	value := store.Get([]byte(assetId))

	ret := restakingtype.StakingAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*restakingtype.StakingAssetInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, restakingtype.KeyPrefixReStakingAssetInfo)
	defer iterator.Close()

	ret := make(map[string]*restakingtype.StakingAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var assetInfo restakingtype.StakingAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &assetInfo)
		_, assetId := restakingtype.GetStakeIDAndAssetIdFromStr(assetInfo.AssetBasicInfo.LayerZeroChainId, "", assetInfo.AssetBasicInfo.Address)
		ret[assetId] = &assetInfo
	}
	return ret, nil
}
