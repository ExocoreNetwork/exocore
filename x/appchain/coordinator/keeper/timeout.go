package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendChainToInitTimeout appends a chain to the list of chains which will timeout (if not
// initialized by then) at the end of the epoch.
func (k Keeper) AppendChainToInitTimeout(
	ctx sdk.Context, epoch epochstypes.Epoch, chainID string,
) {
	prev := k.GetChainsToInitTimeout(ctx, epoch)
	prev.List = append(prev.List, chainID)
	k.SetChainsToInitTimeout(ctx, epoch, prev)
}

// GetChainsToInitTimeout returns the list of chains which will timeout (if not initialized by then)
// at the end of the epoch.
func (k Keeper) GetChainsToInitTimeout(
	ctx sdk.Context, epoch epochstypes.Epoch,
) (res types.ChainIDs) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InitTimeoutEpochKey(epoch))
	k.cdc.MustUnmarshal(bz, &res)
	return res
}

// SetChainsToInitTimeout sets the list of chains which will timeout (if not initialized by then)
// at the end of the epoch.
func (k Keeper) SetChainsToInitTimeout(
	ctx sdk.Context, epoch epochstypes.Epoch, chains types.ChainIDs,
) {
	store := ctx.KVStore(k.storeKey)
	if len(chains.List) == 0 {
		store.Delete(types.InitTimeoutEpochKey(epoch))
		return
	}
	bz := k.cdc.MustMarshal(&chains)
	store.Set(types.InitTimeoutEpochKey(epoch), bz)
}

// RemoveChainFromInitTimeout removes a chain from the list of chains which will timeout (if not
// initialized by then) at the end of the epoch.
func (k Keeper) RemoveChainFromInitTimeout(
	ctx sdk.Context, epoch epochstypes.Epoch, chainID string,
) {
	prev := k.GetChainsToInitTimeout(ctx, epoch)
	for i, id := range prev.List {
		if id == chainID {
			prev.List = append(prev.List[:i], prev.List[i+1:]...)
			break
		}
	}
	k.SetChainsToInitTimeout(ctx, epoch, prev)
}

// SetChainInitTimeout stores a lookup from chain to the epoch by the end of which
// it must be initialized.
func (k Keeper) SetChainInitTimeout(
	ctx sdk.Context, chainID string, epoch epochstypes.Epoch,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChainInitTimeoutKey(chainID), k.cdc.MustMarshal(&epoch))
}

// GetChainInitTimeout returns the epoch by the end of which the chain must be initialized.
func (k Keeper) GetChainInitTimeout(
	ctx sdk.Context, chainID string,
) (epoch epochstypes.Epoch, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ChainInitTimeoutKey(chainID))
	if bz == nil {
		return epoch, false
	}
	k.cdc.MustUnmarshal(bz, &epoch)
	return epoch, true
}

// DeleteChainInitTimeout deletes the lookup from chain to the epoch by the end of which
// it must be initialized.
func (k Keeper) DeleteChainInitTimeout(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ChainInitTimeoutKey(chainID))
}
