package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ assetstype.MsgServer = &Keeper{}

// SetStakerExoCoreAddr outdated, will be deprecated.
// don't check if the staker has existed temporarily,so users can set their ExoCoreAddr multiple times.
// It may be modified later to allow setting only once
func (k Keeper) SetStakerExoCoreAddr(ctx context.Context, addrInfo *assetstype.MsgSetExoCoreAddr) (*assetstype.MsgSetExoCoreAddrResponse, error) {
	// todo: verify client chain signature according to the client chain signature algorithm type.

	c := sdk.UnwrapSDKContext(ctx)

	store := prefix.NewStore(c.KVStore(k.storeKey), assetstype.KeyPrefixReStakerExoCoreAddr)

	bz := k.cdc.MustMarshal(addrInfo)

	key := strings.Join([]string{addrInfo.ClientChainAddr, hexutil.EncodeUint64(addrInfo.ClientChainIndex)}, "_")
	store.Set([]byte(key), bz)

	// todo: save to KeyPrefixReStakerExoCoreAddrReverse

	return &assetstype.MsgSetExoCoreAddrResponse{}, nil
}

func (k Keeper) RegisterClientChain(ctx context.Context, req *assetstype.RegisterClientChainReq) (*assetstype.RegisterClientChainResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetClientChainInfo(c, req.Info)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (k Keeper) RegisterAsset(ctx context.Context, req *assetstype.RegisterAssetReq) (*assetstype.RegisterAssetResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetStakingAssetInfo(c, &assetstype.StakingAssetInfo{
		AssetBasicInfo:     req.Info,
		StakingTotalAmount: math.NewInt(0),
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
