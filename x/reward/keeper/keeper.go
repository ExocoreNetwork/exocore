package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"google.golang.org/protobuf/proto"

	"github.com/ExocoreNetwork/exocore/utils"
	"github.com/ExocoreNetwork/exocore/utils/key"
	"github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/ExocoreNetwork/exocore/x/reward/types"
)

var (
	poolNamePrefix = "pool"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	// other keepers
	assetsKeeper keeper.Keeper
	banker       bankkeeper.Keeper
	distributor  types.Distributor
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	assetsKeeper keeper.Keeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		assetsKeeper: assetsKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) setPool(ctx sdk.Context, pool types.Pool) {
	poolKey := key.FromStr(poolNamePrefix).Append(key.FromStr(pool.Name))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRewardInfo)
	store.Set(poolKey.Bytes(), k.cdc.MustMarshal(&pool))
}

func (k Keeper) getPools(ctx sdk.Context) ([]types.Pool, error) {
	var pools []types.Pool

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRewardInfo)
	iter := store.Iterator(utils.LowerCaseKey(poolNamePrefix))
	defer utils.CloseLogError(iter, k.Logger(ctx))

	for ; iter.Valid(); iter.Next() {
		var pool types.Pool
		err := proto.Unmarshal(iter.Value(), pool)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode the pool")
		}
		pools = append(pools, pool)
	}

	return pools, nil
}

func (k Keeper) getPool(ctx sdk.Context, name string) *rewardPool {
	var pool types.Pool
	poolKey := key.FromStr(poolNamePrefix).Append(key.FromStr(name))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixRewardInfo)
	if !store.Has(poolKey.Bytes()) {
		return newRewardPool(ctx, k, k.banker, k.distributor, types.NewPool(name))
	}
	return newRewardPool(ctx, k, k.banker, k.distributor, pool)

}
