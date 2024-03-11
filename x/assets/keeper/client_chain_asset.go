package keeper

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateStakingAssetTotalAmount updating the total deposited amount of a specified asset in exoCore chain
// The function will be called when stakers deposit and withdraw their assets
func (k Keeper) UpdateStakingAssetTotalAmount(ctx sdk.Context, assetID string, changeAmount sdkmath.Int) (err error) {
	if changeAmount.IsNil() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakingAssetInfo)
	key := []byte(assetID)
	ifExist := store.Has(key)
	if !ifExist {
		return assetstype.ErrNoClientChainAssetKey
	}

	value := store.Get(key)

	ret := assetstype.StakingAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)

	// calculate and set new amount
	err = assetstype.UpdateAssetValue(&ret.StakingTotalAmount, &changeAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&ret)

	store.Set(key, bz)

	return nil
}

// SetStakingAssetInfo todo: Temporarily use clientChainAssetAddr+'_'+LayerZeroChainID as the key.
// It provides a function to register the client chain assets supported by exoCore.It's called by genesis configuration now,however it will be called by the governance in the future
func (k Keeper) SetStakingAssetInfo(ctx sdk.Context, info *assetstype.StakingAssetInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakingAssetInfo)
	// key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	_, assetID := assetstype.GetStakeIDAndAssetIDFromStr(info.AssetBasicInfo.LayerZeroChainID, "", info.AssetBasicInfo.Address)
	store.Set([]byte(assetID), bz)
	return nil
}

func (k Keeper) IsStakingAsset(ctx sdk.Context, assetID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakingAssetInfo)
	return store.Has([]byte(assetID))
}

func (k Keeper) GetStakingAssetInfo(ctx sdk.Context, assetID string) (info *assetstype.StakingAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakingAssetInfo)
	ifExist := store.Has([]byte(assetID))
	if !ifExist {
		return nil, assetstype.ErrNoClientChainAssetKey
	}

	value := store.Get([]byte(assetID))

	ret := assetstype.StakingAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*assetstype.StakingAssetInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, assetstype.KeyPrefixReStakingAssetInfo)
	defer iterator.Close()

	ret := make(map[string]*assetstype.StakingAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var assetInfo assetstype.StakingAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &assetInfo)
		_, assetID := assetstype.GetStakeIDAndAssetIDFromStr(assetInfo.AssetBasicInfo.LayerZeroChainID, "", assetInfo.AssetBasicInfo.Address)
		ret[assetID] = &assetInfo
	}
	return ret, nil
}
