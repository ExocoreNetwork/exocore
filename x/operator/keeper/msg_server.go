package keeper

import (
	context "context"
	"encoding/base64"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"

	"github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k *Keeper) RegisterOperator(ctx context.Context, req *types.RegisterOperatorReq) (*types.RegisterOperatorResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetOperatorInfo(c, req.FromAddress, req.Info)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// OptInToChainId add for dogfood
func (k *Keeper) OptInToChainId(
	goCtx context.Context,
	req *types.OptInToChainIdRequest,
) (*types.OptInToChainIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	key, err := stringToPubKey(req.PublicKey)
	if err != nil {
		return nil, err
	}
	err = k.SetOperatorConsKeyForChainId(
		ctx, addr, req.ChainId, key,
	)
	if err != nil {
		return nil, err
	}
	return &types.OptInToChainIdResponse{}, nil
}

func (k *Keeper) InitiateOptOutFromChainId(
	goCtx context.Context,
	req *types.InitiateOptOutFromChainIdRequest,
) (*types.InitiateOptOutFromChainIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	if err := k.InitiateOperatorOptOutFromChainId(
		ctx, addr, req.ChainId,
	); err != nil {
		return nil, err
	}
	return &types.InitiateOptOutFromChainIdResponse{}, nil
}

func stringToPubKey(pubKey string) (key tmprotocrypto.PublicKey, err error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return
	}
	subscriberTMConsKey := tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Ed25519{
			Ed25519: pubKeyBytes,
		},
	}
	return subscriberTMConsKey, nil
}
