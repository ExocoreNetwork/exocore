package keeper

import (
	"fmt"
	"math/big"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		paramstore    paramtypes.Subspace
		stakingKeeper stakingkeeper.Keeper
	}
)

// TODO
// add(block txs)_remove_maxDistande
var recentTxs map[uint64][]*types.PriceWithSource

type cacheItem struct {
	pSources  []*types.PriceWithSource
	validator string
	power     *big.Int
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) addCache(item *cacheItem)
func (k *Keeper) removeCache(item *cacheItem)

// persist cache into KV
func (k *Keeper) updateRecentTxs() {
	//TODO: add

	//TODO: delete, use map in the KVStore, so actually dont need to implement the exactly 'delete' at firt, just remove all blocks before maxDistance
}
