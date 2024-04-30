package keeper

import (
	"strings"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAppChainInfo stores the info for the app chain to the db. At the moment, it is called by
// the genesis process. In the future, it should be called by governance.
func (k Keeper) SetAppChainInfo(
	ctx sdk.Context,
	info *assetstype.AppChainInfo,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixAppChainInfo)
	bz := k.cdc.MustMarshal(info)
	store.Set([]byte(info.ChainID), bz)
}

// AppChainInfoIsExist returns whether the app chain info for the specified chainID exists
func (k Keeper) AppChainInfoIsExist(ctx sdk.Context, chainID string) bool {
	// short circuit if information is for the current chain
	if strings.Compare(chainID, ctx.ChainID()) == 0 {
		return true
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixAppChainInfo)
	return store.Has([]byte(chainID))
}

// GetAppChainInfoByChainID gets the app chain info for the specified chainID, if it exists
func (k Keeper) GetAppChainInfoByChainID(
	ctx sdk.Context,
	chainID string,
) (info *assetstype.AppChainInfo, err error) {
	// short circuit if information is for the current chain
	if strings.Compare(chainID, ctx.ChainID()) == 0 {
		return &assetstype.AppChainInfo{
			ChainID: ctx.ChainID(),
		}, nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixAppChainInfo)
	value := store.Get([]byte(chainID))
	if value == nil {
		return nil, assetstype.ErrUnknownAppChainID
	}
	ret := &assetstype.AppChainInfo{}
	k.cdc.MustUnmarshal(value, ret)
	return ret, nil
}

// GetAllAppChainInfo gets all the app chain info, indexed by chainID
func (k Keeper) GetAllAppChainInfo(
	ctx sdk.Context,
) (infos map[string]*assetstype.AppChainInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, assetstype.KeyPrefixAppChainInfo)
	defer iterator.Close()

	ret := make(map[string]*assetstype.AppChainInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		chainInfo := &assetstype.AppChainInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), chainInfo)
		ret[chainInfo.ChainID] = chainInfo
	}
	ret[ctx.ChainID()] = &assetstype.AppChainInfo{
		ChainID: ctx.ChainID(),
	}
	return ret
}
