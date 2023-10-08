package keeper

import types2 "github.com/exocore/x/restaking_assets_manage/types"

func (k Keeper) SetClientChainInfo(info *types2.ClientChainInfo) (exoCoreChainIndex uint64, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetClientChainInfoByIndex(index uint64) (info types2.ClientChainInfo, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetAllClientChainInfo() (infos map[uint64]types2.ClientChainInfo, err error) {
	//TODO implement me
	panic("implement me")
}
