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

// GetQueuedKeyOperations returns the list of operations that are queued for execution at the
// end of the current epoch.
func (k Keeper) GetQueuedOperations(
	ctx sdk.Context,
) []types.Operation {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.QueuedOperationsKey())
	if bz == nil {
		return []types.Operation{}
	}
	var operations types.Operations
	if err := operations.Unmarshal(bz); err != nil {
		// TODO(mm): any failure to unmarshal is treated as no operations or panic?
		return []types.Operation{}
	}
	return operations.GetList()
}

// ClearQueuedOperations clears the operations to be executed at the end of the epoch.
func (k Keeper) ClearQueuedOperations(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.QueuedOperationsKey())
}

// setQueuedOperations is a private, internal function used to update the current queue of
// operations to be executed at the end of the epoch with the supplied value.
func (k Keeper) setQueuedOperations(ctx sdk.Context, operations types.Operations) {
	store := ctx.KVStore(k.storeKey)
	bz, err := operations.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.QueuedOperationsKey(), bz)
}
