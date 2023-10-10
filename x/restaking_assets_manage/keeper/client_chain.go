package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

// SetClientChainInfo todo: Temporarily use layerZeroChainId as key.
func (k Keeper) SetClientChainInfo(ctx sdk.Context, info *types2.ClientChainInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixClientChainInfo)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	store.Set([]byte(hexutil.EncodeUint64(info.LayerZeroChainId)), bz)
	return nil
}

func (k Keeper) ClientChainInfoIsExist(ctx sdk.Context, index uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixClientChainInfo)
	return store.Has([]byte(hexutil.EncodeUint64(index)))
}

func (k Keeper) GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *types2.ClientChainInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixClientChainInfo)
	ifExist := store.Has([]byte(hexutil.EncodeUint64(info.LayerZeroChainId)))
	if !ifExist {
		return nil, types2.ErrNoClientChainKey
	}

	value := store.Get([]byte(hexutil.EncodeUint64(index)))

	ret := types2.ClientChainInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*types2.ClientChainInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixClientChainInfo)
	iterator := sdk.KVStorePrefixIterator(store, types2.KeyPrefixClientChainInfo)
	defer iterator.Close()

	ret := make(map[uint64]*types2.ClientChainInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var chainInfo types2.ClientChainInfo
		k.cdc.MustUnmarshal(iterator.Value(), &chainInfo)
		ret[chainInfo.LayerZeroChainId] = &chainInfo
	}
	return ret, nil
}
