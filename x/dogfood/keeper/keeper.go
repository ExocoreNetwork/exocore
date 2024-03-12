package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		// internal hooks to allow other modules to subscriber to our events
		dogfoodHooks types.DogfoodHooks

		// external keepers as interfaces
		epochsKeeper     types.EpochsKeeper
		operatorKeeper   types.OperatorKeeper
		delegationKeeper types.DelegationKeeper
		restakingKeeper  types.AssetsKeeper
		slashingKeeper   types.SlashingKeeper
	}
)

// NewKeeper creates a new dogfood keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	epochsKeeper types.EpochsKeeper,
	operatorKeeper types.OperatorKeeper,
	delegationKeeper types.DelegationKeeper,
	restakingKeeper types.AssetsKeeper,
	slashingKeeper types.SlashingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		paramstore:       ps,
		epochsKeeper:     epochsKeeper,
		operatorKeeper:   operatorKeeper,
		delegationKeeper: delegationKeeper,
		restakingKeeper:  restakingKeeper,
		slashingKeeper:   slashingKeeper,
	}
}

// Logger returns a logger object for use within the module.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the hooks on the keeper. It intentionally has a pointer receiver so that
// changes can be saved to the object.
func (k *Keeper) SetHooks(sh types.DogfoodHooks) *Keeper {
	if k.dogfoodHooks != nil {
		panic("cannot set dogfood hooks twice")
	}
	if sh == nil {
		panic("cannot set nil dogfood hooks")
	}
	k.dogfoodHooks = sh
	return k
}

// Hooks returns the hooks registered to the module.
func (k Keeper) Hooks() types.DogfoodHooks {
	return k.dogfoodHooks
}

// MarkEpochEnd marks the end of the epoch. It is called within the BeginBlocker to inform
// the module to apply the validator updates at the end of this block.
func (k Keeper) MarkEpochEnd(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEndKey()
	store.Set(key, []byte{1})
}

// IsEpochEnd returns true if the epoch ended in the beginning of this block, or the end of the
// previous block.
func (k Keeper) IsEpochEnd(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEndKey()
	return store.Has(key)
}

// ClearEpochEnd clears the epoch end marker. It is called after the epoch end operations are
// applied.
func (k Keeper) ClearEpochEnd(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := types.EpochEndKey()
	store.Delete(key)
}
