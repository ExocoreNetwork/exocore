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

// OptInToChainID add for dogfood
func (k *Keeper) OptInToChainID(
	goCtx context.Context,
	req *types.OptInToChainIDRequest,
) (*types.OptInToChainIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	key, err := stringToPubKey(req.PublicKey)
	if err != nil {
		return nil, err
	}
	err = k.SetOperatorConsKeyForChainID(
		ctx, addr, req.ChainId, key,
	)
	if err != nil {
		return nil, err
	}
	return &types.OptInToChainIDResponse{}, nil
}

func (k *Keeper) InitiateOptOutFromChainID(
	goCtx context.Context,
	req *types.InitiateOptOutFromChainIDRequest,
) (*types.InitiateOptOutFromChainIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	if err := k.InitiateOperatorOptOutFromChainID(
		ctx, addr, req.ChainId,
	); err != nil {
		return nil, err
	}
	return &types.InitiateOptOutFromChainIDResponse{}, nil
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
