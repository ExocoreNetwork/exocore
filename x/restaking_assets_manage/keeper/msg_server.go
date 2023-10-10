package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"strings"
)

var _ types2.MsgServer = &Keeper{}

// SetStakerExoCoreAddr don't check if the staker has existed temporarily,so users can set their ExoCoreAddr multiple times.
// It may be modified later to allow setting only once
func (k Keeper) SetStakerExoCoreAddr(ctx context.Context, addrInfo *types2.MsgSetExoCoreAddr) (*types2.MsgSetExoCoreAddrResponse, error) {
	//todo: verify client chain signature according to the client chain signature algorithm type.

	c := sdk.UnwrapSDKContext(ctx)

	//save to KeyPrefixReStakerExoCoreAddr
	store := prefix.NewStore(c.KVStore(k.storeKey), types2.KeyPrefixReStakerExoCoreAddr)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(addrInfo)

	key := strings.Join([]string{addrInfo.ClientChainAddr, hexutil.EncodeUint64(addrInfo.ClientChainIndex)}, "_")
	store.Set([]byte(key), bz)

	//save to KeyPrefixReStakerExoCoreAddrReverse

	return &types2.MsgSetExoCoreAddrResponse{}, nil
}
