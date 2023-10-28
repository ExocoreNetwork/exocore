package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/exocore/utils"
	"github.com/exocore/utils/key"
	"github.com/exocore/x/reward/exported"
	"github.com/exocore/x/reward/types"
)

var poolNamePrefix = "pool"

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	bankKeeper types.BankKeeper,
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

		bankKeeper: bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetPool returns the reward pool of the given name, or returns an empty reward pool if not found
func (k Keeper) GetPool(ctx sdk.Context, name string) exported.RewardPool {
	var pool types.Pool
	ok := k.getStore(ctx).GetNew(key.FromStr(poolNamePrefix).Append(key.FromStr(name)), &pool)
	if !ok {
		return newPool(ctx, k, k.bankKeeper, types.NewPool(name))
	}

	return newPool(ctx, k, k.bankKeeper, pool)
}

func (k Keeper) getPools(ctx sdk.Context) []types.Pool {
	var pools []types.Pool

	store := k.getStore(ctx)
	iter := store.Iterator(utils.LowerCaseKey(poolNamePrefix))
	defer utils.CloseLogError(iter, k.Logger(ctx))

	for ; iter.Valid(); iter.Next() {
		var pool types.Pool
		iter.UnmarshalValue(&pool)

		pools = append(pools, pool)
	}

	return pools
}

func (k Keeper) setPool(ctx sdk.Context, pool types.Pool) {
	k.getStore(ctx).SetNewValidated(key.FromStr(poolNamePrefix).Append(key.FromStr(pool.Name)), &pool)
}

func (k Keeper) getStore(ctx sdk.Context) utils.KVStore {
	return utils.NewNormalizedStore(ctx.KVStore(k.storeKey), k.cdc)
}
