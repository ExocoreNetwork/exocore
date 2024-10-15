package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MapHeightToChainVscID stores the height corresponding to a chainID and vscID
func (k Keeper) MapHeightToChainVscID(ctx sdk.Context, chainID string, vscID uint64, height uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.HeightToChainVscIDKey(chainID, vscID)
	store.Set(key, sdk.Uint64ToBigEndian(height))
}

// GetHeightForChainVscID gets the height corresponding to a chainID and vscID
func (k Keeper) GetHeightForChainVscID(ctx sdk.Context, chainID string, vscID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.HeightToChainVscIDKey(chainID, vscID)
	// if store.Has(key) is false will return a height of 0
	return sdk.BigEndianToUint64(store.Get(key))
}

// SetVscIDForChain stores the vscID corresponding to a chainID
func (k Keeper) SetVscIDForChain(ctx sdk.Context, chainID string, vscID uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.VscIDForChainKey(chainID)
	store.Set(key, sdk.Uint64ToBigEndian(vscID))
}

// GetVscIDForChain gets the vscID corresponding to a chainID
func (k Keeper) GetVscIDForChain(ctx sdk.Context, chainID string) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.VscIDForChainKey(chainID)
	return sdk.BigEndianToUint64(store.Get(key))
}

// IncrementVscIDForChain increments the vscID corresponding to a chainID, and
// stores/returns the new vscID
func (k Keeper) IncrementVscIDForChain(ctx sdk.Context, chainID string) uint64 {
	vscID := k.GetVscIDForChain(ctx, chainID)
	vscID++
	k.SetVscIDForChain(ctx, chainID, vscID)
	return vscID
}
