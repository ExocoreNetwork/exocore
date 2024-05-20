package keeper

import (
	"context"
	"errors"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ operatortypes.QueryServer = &Keeper{}

func (k *Keeper) GetOperatorInfo(ctx context.Context, req *operatortypes.GetOperatorInfoReq) (*operatortypes.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.OperatorInfo(c, req.OperatorAddr)
}

// QueryOperatorConsKeyForChainID add for dogfood
func (k *Keeper) QueryOperatorConsKeyForChainID(
	goCtx context.Context,
	req *operatortypes.QueryOperatorConsKeyRequest,
) (*operatortypes.QueryOperatorConsKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Addr)
	if err != nil {
		return nil, err
	}
	found, key, err := k.GetOperatorConsKeyForChainID(
		ctx, addr, req.ChainId,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("no key assigned")
	}
	return &operatortypes.QueryOperatorConsKeyResponse{
		PublicKey: *key,
	}, nil
}

func (k *Keeper) QueryOperatorUSDValue(ctx context.Context, req *operatortypes.QueryOperatorUSDValueRequest) (*operatortypes.DecValueField, error) {
	c := sdk.UnwrapSDKContext(ctx)
	usdValue, err := k.GetOperatorUSDValue(c, req.OperatorAddr, req.AVSAddress)
	if err != nil {
		return nil, err
	}
	return &operatortypes.DecValueField{
		Amount: usdValue,
	}, nil
}

func (k *Keeper) QueryAVSUSDValue(ctx context.Context, req *operatortypes.QueryAVSUSDValueRequest) (*operatortypes.DecValueField, error) {
	c := sdk.UnwrapSDKContext(ctx)
	usdValue, err := k.GetAVSUSDValue(c, req.AVSAddress)
	if err != nil {
		return nil, err
	}
	return &operatortypes.DecValueField{
		Amount: usdValue,
	}, nil
}

func (k *Keeper) QueryOperatorSlashInfo(ctx context.Context, req *operatortypes.QueryOperatorSlashInfoRequest) (*operatortypes.QueryOperatorSlashInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	slashInfo, err := k.AllOperatorSlashInfo(c, req.OperatorAddr, req.AVSAddress)
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryOperatorSlashInfoResponse{
		AllSlashInfo: slashInfo,
	}, nil
}
