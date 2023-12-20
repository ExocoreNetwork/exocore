package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

// SetClientChainInfo todo: Temporarily use layerZeroChainId as key.
func (k Keeper) SetClientChainInfo(ctx sdk.Context, info *restakingtype.ClientChainInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixClientChainInfo)
	// key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	store.Set([]byte(hexutil.EncodeUint64(info.LayerZeroChainId)), bz)
	return nil
}

func (k Keeper) ClientChainInfoIsExist(ctx sdk.Context, index uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixClientChainInfo)
	return store.Has([]byte(hexutil.EncodeUint64(index)))
}

func (k Keeper) GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *restakingtype.ClientChainInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixClientChainInfo)
	ifExist := store.Has([]byte(hexutil.EncodeUint64(index)))
	if !ifExist {
		return nil, restakingtype.ErrNoClientChainKey
	}

	value := store.Get([]byte(hexutil.EncodeUint64(index)))

	ret := restakingtype.ClientChainInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*restakingtype.ClientChainInfo, err error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, restakingtype.KeyPrefixClientChainInfo)
	defer iterator.Close()

	ret := make(map[uint64]*restakingtype.ClientChainInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var chainInfo restakingtype.ClientChainInfo
		k.cdc.MustUnmarshal(iterator.Value(), &chainInfo)
		ret[chainInfo.LayerZeroChainId] = &chainInfo
	}
	return ret, nil
}
