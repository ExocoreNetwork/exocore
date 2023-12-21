package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
	"strings"
)

var _ restakingtype.MsgServer = &Keeper{}

// SetStakerExoCoreAddr outdated, will be deprecated.
// don't check if the staker has existed temporarily,so users can set their ExoCoreAddr multiple times.
// It may be modified later to allow setting only once
func (k Keeper) SetStakerExoCoreAddr(ctx context.Context, addrInfo *restakingtype.MsgSetExoCoreAddr) (*restakingtype.MsgSetExoCoreAddrResponse, error) {
	//todo: verify client chain signature according to the client chain signature algorithm type.

	c := sdk.UnwrapSDKContext(ctx)

	//save to KeyPrefixReStakerExoCoreAddr
	store := prefix.NewStore(c.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerExoCoreAddr)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(addrInfo)

	key := strings.Join([]string{addrInfo.ClientChainAddr, hexutil.EncodeUint64(addrInfo.ClientChainIndex)}, "_")
	store.Set([]byte(key), bz)

	//todo: save to KeyPrefixReStakerExoCoreAddrReverse

	return &restakingtype.MsgSetExoCoreAddrResponse{}, nil
}
