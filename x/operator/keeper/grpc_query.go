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

// QueryOperatorConsKeyForChainId add for dogfood
func (k *Keeper) QueryOperatorConsKeyForChainId(
	goCtx context.Context,
	req *operatortypes.QueryOperatorConsKeyForChainIdRequest,
) (*operatortypes.QueryOperatorConsKeyForChainIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Addr)
	if err != nil {
		return nil, err
	}
	found, key, err := k.GetOperatorConsKeyForChainId(
		ctx, addr, req.ChainId,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("no key assigned")
	}
	return &operatortypes.QueryOperatorConsKeyForChainIdResponse{
		PublicKey: key,
	}, nil
}
