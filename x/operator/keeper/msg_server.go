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

// OptInToCosmosChain this is an RPC for the operators
// that want to service as a validator for the app chain Avs
// The operator can opt in the cosmos app chain through this RPC
// In this function, the basic function `OptIn` need to be called
func (k *Keeper) OptInToCosmosChain(
	goCtx context.Context,
	req *types.OptInToCosmosChainRequest,
) (*types.OptInToCosmosChainResponse, error) {
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
	return &types.OptInToCosmosChainResponse{}, nil
}

// InitiateOptOutFromCosmosChain is a method corresponding to OptInToCosmosChain
// It provides a function to opt out from the app chain Avs for the operators.
func (k *Keeper) InitiateOptOutFromCosmosChain(
	goCtx context.Context,
	req *types.InitiateOptOutFromCosmosChainRequest,
) (*types.InitiateOptOutFromCosmosChainResponse, error) {
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
	return &types.InitiateOptOutFromCosmosChainResponse{}, nil
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
