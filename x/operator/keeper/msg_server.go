package keeper

import (
	context "context"

	"github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

// RegisterOperator is an implementation of the msg server for the operator module.
func (k *Keeper) RegisterOperator(ctx context.Context, req *types.RegisterOperatorReq) (*types.RegisterOperatorResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetOperatorInfo(c, req.FromAddress, req.Info)
	if err != nil {
		return nil, err
	}
	return &types.RegisterOperatorResponse{}, nil
}

// OptIntoAVS is an implementation of the msg server for the operator module.
func (k Keeper) OptIntoAVS(ctx context.Context, req *types.OptIntoAVSReq) (*types.OptIntoAVSResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	// #nosec G703 // already validated
	accAddr, _ := sdk.AccAddressFromBech32(req.FromAddress)
	err := k.OptIn(c, accAddr, req.AvsAddress)
	if err != nil {
		return nil, err
	}
	return &types.OptIntoAVSResponse{}, nil
}

// OptOutOfAVS is an implementation of the msg server for the operator module.
func (k Keeper) OptOutOfAVS(ctx context.Context, req *types.OptOutOfAVSReq) (*types.OptOutOfAVSResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	// #nosec G703 // already validated
	accAddr, _ := sdk.AccAddressFromBech32(req.FromAddress)
	err := k.OptOut(c, accAddr, req.AvsAddress)
	if err != nil {
		return nil, err
	}
	return &types.OptOutOfAVSResponse{}, nil
}

// SetConsKey is an implementation of the msg server for the operator module.
func (k Keeper) SetConsKey(ctx context.Context, req *types.SetConsKeyReq) (*types.SetConsKeyResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	// #nosec G703 // already validated
	accAddr, _ := sdk.AccAddressFromBech32(req.Address)
	// #nosec G703 // already validated, including type
	keyString, _, _ := types.ParseConsensusKeyFromJSON(req.PublicKey)
	// #nosec G703 // already validated
	keyObj, _ := types.StringToPubKey(keyString)
	if err := k.SetOperatorConsKeyForChainID(c, accAddr, req.ChainID, keyObj); err != nil {
		return nil, err
	}
	return &types.SetConsKeyResponse{}, nil
}

// InitConsKeyRemoval is an implementation of the msg server for the operator module.
func (k Keeper) InitConsKeyRemoval(ctx context.Context, req *types.InitConsKeyRemovalReq) (*types.InitConsKeyRemovalResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	// #nosec G703 // already validated
	accAddr, _ := sdk.AccAddressFromBech32(req.Address)
	err := k.InitiateOperatorKeyRemovalForChainID(c, accAddr, req.ChainID)
	if err != nil {
		return nil, err
	}
	return &types.InitConsKeyRemovalResponse{}, nil
}
