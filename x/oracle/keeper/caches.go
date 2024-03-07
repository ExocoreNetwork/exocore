package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type cacheItem struct {
	feederId  int32
	pSources  []*types.PriceWithSource
	validator string
	// power     *big.Int
}

// memory cache
func (k *Keeper) addCache(item *cacheItem) {

}

// memory cache
func (k *Keeper) removeCache(item *cacheItem)

// persist cache into KV
func (k *Keeper) updateRecentTxs() {
	//TODO: add

	//TODO: delete, use map in the KVStore, so actually dont need to implement the exactly 'delete' at firt, just remove all blocks before maxDistance
}

// commit memory cache to KVStore
func (k *Keeper) commitCache() {

}

// from KVStore
func (k *Keeper) getCaches(ctx sdk.Context) map[uint64]*cacheItem {
	return nil
}

func (k *Keeper) clearCaches(ctx sdk.Context) {

}

// from KVStore
func (k *Keeper) getCacheValidators(ctx sdk.Context) map[string]*big.Int {
	return nil
}
