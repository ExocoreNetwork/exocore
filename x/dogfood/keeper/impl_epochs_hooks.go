package keeper

import (
	"strings"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ = sdk.NewCoin("stake", sdk.NewInt(1))
)

// EpochsHooksWrapper is the wrapper structure that implements the epochs hooks for the dogfood
// keeper.
type EpochsHooksWrapper struct {
	keeper *Keeper
}

// Interface guard
var _ types.EpochsHooks = EpochsHooksWrapper{}

// EpochsHooks returns the epochs hooks wrapper. It follows the "accept interfaces, return
// concretes" pattern.
func (k *Keeper) EpochsHooks() EpochsHooksWrapper {
	return EpochsHooksWrapper{k}
}

// AfterEpochEnd is called after an epoch ends.
func (wrapper EpochsHooksWrapper) AfterEpochEnd(
	ctx sdk.Context, identifier string, epoch int64,
) {
	if strings.Compare(identifier, wrapper.keeper.GetEpochIdentifier(ctx)) == 0 {
		// we will upgrade all of the queued information to "pending", which will be applied at
		// the end of the block.
		// note that this hook is called during BeginBlock, and the "pending" operations will be
		// applied within this block. however, for clarity, it is highlighted that unbonding
		// takes N epochs + 1 block to complete.
		operations := wrapper.keeper.GetQueuedOperations(ctx)
		wrapper.keeper.SetPendingOperations(ctx, types.Operations{List: operations})
		wrapper.keeper.ClearQueuedOperations(ctx)
		optOuts := wrapper.keeper.GetOptOutsToFinish(ctx, epoch)
		wrapper.keeper.SetPendingOptOuts(ctx, types.AccountAddresses{List: optOuts})
		wrapper.keeper.ClearOptOutsToFinish(ctx, epoch)
		consAddresses := wrapper.keeper.GetConsensusAddrsToPrune(ctx, epoch)
		wrapper.keeper.SetPendingConsensusAddrs(
			ctx, types.ConsensusAddresses{List: consAddresses},
		)
		wrapper.keeper.ClearConsensusAddrsToPrune(ctx, epoch)
		undelegations := wrapper.keeper.GetUndelegationsToMature(ctx, epoch)
		wrapper.keeper.SetPendingUndelegations(ctx, types.RecordKeys{List: undelegations})
		wrapper.keeper.ClearUndelegationsToMature(ctx, epoch)
	}
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	ctx sdk.Context, identifier string, epoch int64,
) {
	// nothing to do
}
