package keeper

import (
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAppChainInfo stores the info for the app chain to the db. At the moment, it is called by
// the
// genesis process. In the future, it should be called by governance.
func (k Keeper) SetAppChainInfo(
	ctx sdk.Context,
	info assetstype.AppChainInfo,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixAppChainInfo)
	bz := k.cdc.MustMarshal(&info)
	store.Set([]byte(info.ChainId), bz)
}

// AppChainInfoIsExist returns whether the app chain info for the specified chainId exists
func (k Keeper) AppChainInfoIsExist(ctx sdk.Context, chainID string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixAppChainInfo)
	return store.Has([]byte(chainID))
}

// GetAppChainInfoByChainID gets the app chain info for the specified chainId, if it exists
func (k Keeper) GetAppChainInfoByChainID(
	ctx sdk.Context,
	chainID string,
) (info assetstype.AppChainInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixAppChainInfo)
	ifExist := store.Has([]byte(chainID))
	if !ifExist {
		return assetstype.AppChainInfo{}, assetstype.ErrNoAppChainKey
	}
	value := store.Get([]byte(chainID))
	ret := assetstype.AppChainInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return ret, nil
}

// GetAllAppChainInfo gets all the app chain info, indexed by chainId
func (k Keeper) GetAllAppChainInfo(
	ctx sdk.Context,
) (infos map[string]assetstype.AppChainInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, assetstype.KeyPrefixAppChainInfo)
	defer iterator.Close()

	ret := make(map[string]assetstype.AppChainInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var chainInfo assetstype.AppChainInfo
		k.cdc.MustUnmarshal(iterator.Value(), &chainInfo)
		ret[chainInfo.ChainId] = chainInfo
	}
	return ret
}
