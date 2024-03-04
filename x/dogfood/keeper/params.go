package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetEpochsUntilUnbonded returns the number of epochs after which an unbonding that is made
// during the current epoch will be released. It is a parameter of the dogfood module.
func (k Keeper) GetEpochsUntilUnbonded(ctx sdk.Context) uint32 {
	var epochsUntilUnbonded uint32
	k.paramstore.Get(ctx, types.KeyEpochsUntilUnbonded, &epochsUntilUnbonded)
	return epochsUntilUnbonded
}

// GetEpochIdentifier returns the epoch identifier used to measure an epoch. It is a parameter
// of the dogfood module.
func (k Keeper) GetEpochIdentifier(ctx sdk.Context) string {
	var epochIdentifier string
	k.paramstore.Get(ctx, types.KeyEpochIdentifier, &epochIdentifier)
	return epochIdentifier
}

// GetMaxValidators returns the maximum number of validators that can be asked to validate for
// the chain. It is a parameter of the dogfood module.
func (k Keeper) GetMaxValidators(ctx sdk.Context) uint32 {
	var maxValidators uint32
	k.paramstore.Get(ctx, stakingtypes.KeyMaxValidators, &maxValidators)
	return maxValidators
}

// GetHistorialEntries is the number of historical info entries to persist in the store. These
// entries are used by the IBC module. The return value is a parameter of the dogfood module.
func (k Keeper) GetHistoricalEntries(ctx sdk.Context) uint32 {
	var historicalEntries uint32
	k.paramstore.Get(ctx, stakingtypes.KeyHistoricalEntries, &historicalEntries)
	return historicalEntries
}

// GetAssetIDs returns the asset IDs that are accepted by the dogfood module. It is a parameter
// of the dogfood module.
func (k Keeper) GetAssetIDs(ctx sdk.Context) []string {
	var assetIDs []string
	k.paramstore.Get(ctx, types.KeyAssetIDs, &assetIDs)
	return assetIDs
}

// SetParams sets the params for the dogfood module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetDogfoodParams returns the parameters for the dogfood module. Note that this function is
// intentionally called GetDogfoodParams and not GetParams, since the GetParams function is used
// to implement the slashingtypes.StakingKeeper interface `GetParams(sdk.Context)
// stakingtypes.Params`.
func (k Keeper) GetDogfoodParams(ctx sdk.Context) (params types.Params) {
	return types.NewParams(
		k.GetEpochsUntilUnbonded(ctx),
		k.GetEpochIdentifier(ctx),
		k.GetMaxValidators(ctx),
		k.GetHistoricalEntries(ctx),
		k.GetAssetIDs(ctx),
	)
}
