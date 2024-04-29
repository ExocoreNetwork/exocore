package keeper

import (
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SetClientChainInfo todo: Temporarily use LayerZeroChainID as key.
// It provides a function to register the client chains supported by exoCore.It's called by genesis configuration now,however it will be called by the governance in the future
func (k Keeper) SetClientChainInfo(ctx sdk.Context, info *assetstype.ClientChainInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixClientChainInfo)

	bz := k.cdc.MustMarshal(info)

	store.Set([]byte(hexutil.EncodeUint64(info.LayerZeroChainID)), bz)
	return nil
}

func (k Keeper) IsExistedClientChain(ctx sdk.Context, index uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixClientChainInfo)
	return store.Has([]byte(hexutil.EncodeUint64(index)))
}

// GetClientChainInfoByIndex using LayerZeroChainID as the query index.
func (k Keeper) GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *assetstype.ClientChainInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixClientChainInfo)
	value := store.Get([]byte(hexutil.EncodeUint64(index)))
	if value == nil {
		return nil, assetstype.ErrNoClientChainKey
	}
	ret := assetstype.ClientChainInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*assetstype.ClientChainInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, assetstype.KeyPrefixClientChainInfo)
	defer iterator.Close()

	ret := make(map[uint64]*assetstype.ClientChainInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var chainInfo assetstype.ClientChainInfo
		k.cdc.MustUnmarshal(iterator.Value(), &chainInfo)
		ret[chainInfo.LayerZeroChainID] = &chainInfo
	}
	return ret, nil
}
