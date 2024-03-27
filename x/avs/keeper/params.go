package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams()
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

type AVSParams struct {
	AVSName         string
	AVSAddress      []byte
	OperatorAddress []byte
	Action          uint64
}

const (
	RegisterAction = iota + 1
	DeRegisterAction
)
