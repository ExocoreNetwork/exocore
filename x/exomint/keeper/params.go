package keeper

import (
	"cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/exomint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetMintDenom gets the mint denomination.
func (k Keeper) GetMintDenom(ctx sdk.Context) string {
	var mintDenom string
	k.paramstore.Get(ctx, types.KeyMintDenom, &mintDenom)
	return mintDenom
}

// SetMintDenom sets the mint denomination.
func (k Keeper) SetMintDenom(ctx sdk.Context, mintDenom string) {
	k.paramstore.Set(ctx, types.KeyMintDenom, mintDenom)
}

// GetEpochReward gets the reward minted per epoch.
func (k Keeper) GetEpochReward(ctx sdk.Context) math.Int {
	var epochReward math.Int
	k.paramstore.Get(ctx, types.KeyEpochReward, &epochReward)
	return epochReward
}

// SetEpochReward sets the reward minted per epoch.
func (k Keeper) SetEpochReward(ctx sdk.Context, epochReward math.Int) {
	k.paramstore.Set(ctx, types.KeyEpochReward, epochReward)
}

// GetEpochIdentifier gets the epoch identifier at the end of which
// the epoch reward is minted.
func (k Keeper) GetEpochIdentifier(ctx sdk.Context) string {
	var epochIdentifier string
	k.paramstore.Get(ctx, types.KeyEpochIdentifier, &epochIdentifier)
	return epochIdentifier
}

// SetEpochIdentifier sets the epoch identifier at the end of which
// the epoch reward is minted.
func (k Keeper) SetEpochIdentifier(ctx sdk.Context, epochIdentifier string) {
	k.paramstore.Set(ctx, types.KeyEpochIdentifier, epochIdentifier)
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetMintDenom(ctx),
		k.GetEpochReward(ctx),
		k.GetEpochIdentifier(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
