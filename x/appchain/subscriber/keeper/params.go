package keeper

import (
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetParams sets the appchain coordinator parameters.
func (k Keeper) SetParams(ctx sdk.Context, params commontypes.SubscriberParams) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey(), bz)
}

// GetParams gets the appchain coordinator parameters.
func (k Keeper) GetParams(ctx sdk.Context) (res commontypes.SubscriberParams) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey())
	k.cdc.MustUnmarshal(bz, &res)
	return res
}
