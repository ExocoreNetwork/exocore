package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) SetClientChainInfo(ctx sdk.Context, info *types2.ClientChainInfo) (exoCoreChainIndex uint64, err error) {

	/*	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixClientChainInfo)
		//key := common.HexToAddress(incentive.Contract)
		bz := k.cdc.MustMarshal(info)
		store.Set(key.Bytes(), bz)*/
	return info.ExoCoreChainIndex, nil
}

func (k Keeper) GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *types2.ClientChainInfo, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*types2.ClientChainInfo, err error) {
	//TODO implement me
	panic("implement me")
}
