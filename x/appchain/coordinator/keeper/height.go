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
