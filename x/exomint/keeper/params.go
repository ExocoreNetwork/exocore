package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/exomint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetMintDenom(ctx sdk.Context) string {
	var mintDenom string
	k.paramstore.Get(ctx, types.KeyMintDenom, &mintDenom)
	return mintDenom
}

func (k Keeper) SetMintDenom(ctx sdk.Context, mintDenom string) {
	k.paramstore.Set(ctx, types.KeyMintDenom, mintDenom)
}

func (k Keeper) GetBlockReward(ctx sdk.Context) sdk.Int {
	var blockReward sdk.Int
	k.paramstore.Get(ctx, types.KeyBlockReward, &blockReward)
	return blockReward
}

func (k Keeper) SetBlockReward(ctx sdk.Context, blockReward sdk.Int) {
	k.paramstore.Set(ctx, types.KeyBlockReward, blockReward)
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetMintDenom(ctx),
		k.GetBlockReward(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
