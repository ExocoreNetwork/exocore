package keeper

import (
	context "context"

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
	key, err := types.StringToPubKey(req.PublicKey)
	if err != nil {
		return nil, err
	}
	err = k.SetOperatorConsKeyForChainID(
		ctx, addr, req.ChainId, key,
	)
	if err != nil {
		return nil, err
	}
	// call the basic OptIn
	avsAddr, err := k.avsKeeper.GetAvsAddrByChainID(ctx, req.ChainId)
	if err != nil {
		return nil, err
	}
	err = k.OptIn(ctx, addr, avsAddr)
	if err != nil {
		return nil, err
	}
	return &types.OptInToCosmosChainResponse{}, nil
}

// InitOptOutFromCosmosChain is a method corresponding to OptInToCosmosChain
// It provides a function to opt out from the app chain Avs for the operators.
func (k *Keeper) InitOptOutFromCosmosChain(
	goCtx context.Context,
	req *types.InitOptOutFromCosmosChainRequest,
) (*types.InitOptOutFromCosmosChainResponse, error) {
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
	return &types.InitOptOutFromCosmosChainResponse{}, nil
}
