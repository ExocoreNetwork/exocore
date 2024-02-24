package keeper

import (
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAppChainInfo stores the info for the app chain to the db. At the moment, it is called by
// the
// genesis process. In the future, it should be called by governance.
func (k Keeper) SetAppChainInfo(
	ctx sdk.Context,
	info restakingtype.AppChainInfo,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixAppChainInfo)
	bz := k.cdc.MustMarshal(&info)
	store.Set([]byte(info.ChainId), bz)
}

// AppChainInfoIsExist returns whether the app chain info for the specified chainId exists
func (k Keeper) AppChainInfoIsExist(ctx sdk.Context, chainId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixAppChainInfo)
	return store.Has([]byte(chainId))
}

// GetAppChainInfoByChainId gets the app chain info for the specified chainId, if it exists
func (k Keeper) GetAppChainInfoByChainId(
	ctx sdk.Context,
	chainId string,
) (info restakingtype.AppChainInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixAppChainInfo)
	ifExist := store.Has([]byte(chainId))
	if !ifExist {
		return restakingtype.AppChainInfo{}, restakingtype.ErrNoAppChainKey
	}
	value := store.Get([]byte(chainId))
	ret := restakingtype.AppChainInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return ret, nil
}

// GetAllAppChainInfo gets all the app chain info, indexed by chainId
func (k Keeper) GetAllAppChainInfo(
	ctx sdk.Context,
) (infos map[string]restakingtype.AppChainInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, restakingtype.KeyPrefixAppChainInfo)
	defer iterator.Close()

	ret := make(map[string]restakingtype.AppChainInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var chainInfo restakingtype.AppChainInfo
		k.cdc.MustUnmarshal(iterator.Value(), &chainInfo)
		ret[chainInfo.ChainId] = chainInfo
	}
	return ret
}
