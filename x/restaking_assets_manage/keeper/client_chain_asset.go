package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"strings"
)

// SetStakingAssetInfo todo: Temporarily use clientChainAssetAddr+'_'+layerZeroChainId as the key.
func (k Keeper) SetStakingAssetInfo(ctx sdk.Context, info *types2.StakingAssetInfo) (err error) {
	//TODO implement me
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakingAssetInfo)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	key := strings.Join([]string{info.AssetBasicInfo.Address, hexutil.EncodeUint64(info.AssetBasicInfo.LayerZeroChainId)}, "_")
	store.Set([]byte(key), bz)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixReStakingAssetInfo)
	iterator := sdk.KVStorePrefixIterator(store, types2.KeyPrefixReStakingAssetInfo)
	defer iterator.Close()

	ret := make(map[string]*types2.StakingAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var assetInfo types2.StakingAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &assetInfo)
		assetId := strings.Join([]string{assetInfo.AssetBasicInfo.Address, hexutil.EncodeUint64(assetInfo.AssetBasicInfo.LayerZeroChainId)}, "_")
		ret[assetId] = &assetInfo
	}
	return ret, nil
}
