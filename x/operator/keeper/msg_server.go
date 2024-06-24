package keeper

import (
	context "context"

	errorsmod "cosmossdk.io/errors"
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
func (k Keeper) OptIntoAVS(ctx context.Context, req *types.OptIntoAVSReq) (res *types.OptIntoAVSResponse, err error) {
	uncachedCtx := sdk.UnwrapSDKContext(ctx)
	// only write if both calls succeed
	c, writeFunc := uncachedCtx.CacheContext()
	defer func() {
		if err == nil {
			writeFunc()
		}
	}()
	// TODO: use some form of an AVS to key-type registry here, possibly from within the AVS module to determine
	// if a key is required and that it is appropriately supplied.
	if req.AvsAddress == c.ChainID() {
		if len(req.PublicKey) == 0 {
			return nil, errorsmod.Wrap(types.ErrInvalidPubKey, "a key is required but was not supplied")
		}
	} else {
		if len(req.PublicKey) > 0 {
			return nil, errorsmod.Wrap(types.ErrInvalidPubKey, "a key is not required but was supplied")
		}
	}
	// #nosec G703 // already validated
	accAddr, _ := sdk.AccAddressFromBech32(req.FromAddress)
	err = k.OptIn(c, accAddr, req.AvsAddress)
	if err != nil {
		return nil, err
	}
	if len(req.PublicKey) > 0 {
		// we have to validate just the key; we previously validated that ctx.ChainID() == req.AvsAddress
		keyObj, _ := types.ValidateConsensusKeyJSON(req.PublicKey)
		if err := k.SetOperatorConsKeyForChainID(c, accAddr, req.AvsAddress, keyObj); err != nil {
			return nil, err
		}
	}
	return &types.OptIntoAVSResponse{}, nil
}

// OptOutOfAVS is an implementation of the msg server for the operator module.
func (k Keeper) OptOutOfAVS(ctx context.Context, req *types.OptOutOfAVSReq) (res *types.OptOutOfAVSResponse, err error) {
	uncachedCtx := sdk.UnwrapSDKContext(ctx)
	// only write if both calls succeed
	c, writeFunc := uncachedCtx.CacheContext()
	defer func() {
		if err == nil {
			writeFunc()
		}
	}()
	// #nosec G703 // already validated
	accAddr, _ := sdk.AccAddressFromBech32(req.FromAddress)
	err = k.OptOut(c, accAddr, req.AvsAddress)
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
	// #nosec G703 // already validated
	keyObj, _ := types.ValidateConsensusKeyJSON(req.PublicKey)
	if err := k.SetOperatorConsKeyForChainID(c, accAddr, req.ChainID, keyObj); err != nil {
		return nil, err
	}
	return &types.SetConsKeyResponse{}, nil
}
