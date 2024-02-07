package keeper

import (
	"strings"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OperatorHooksWrapper is the wrapper structure that implements the operator hooks for the
// dogfood keeper.
type OperatorHooksWrapper struct {
	keeper *Keeper
}

// Interface guards
var _ types.OperatorHooks = OperatorHooksWrapper{}

func (k *Keeper) OperatorHooks() OperatorHooksWrapper {
	return OperatorHooksWrapper{k}
}

// Hooks assumptions: Assuming I is opt-in, O is opt-out and R is key replacement, these are all
// possible within the same epoch, for a fresh operator.
// I O
// I R
// I R O
// This is not possible for a fresh operator to do:
// I O R
// R I O
// R I
// For an operator that is already opted in, the list looks like follows:
// R O
// O I
// O I R
// R O I
// The impossible list looks like:
// O R
// O R I
// TODO: list out operation results for each of these, and make sure everything is covered below

// AfterOperatorOptIn is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorOptIn(
	ctx sdk.Context, addr sdk.AccAddress,
	chainID string, pubKey tmprotocrypto.PublicKey,
) {
	if strings.Compare(ctx.ChainID(), chainID) == 0 {
		// res == Removed, it means operator has opted back in
		// res == Success, there is no additional information to store
		// res == Exists, there is nothing to do
		if res := h.keeper.QueueOperation(
			ctx, addr, pubKey, types.KeyAdditionOrUpdate,
		); res == types.QueueResultRemoved {
			// the old operation was key removal, which is now removed from the queue.
			// so all of the changes that were associated with it need to be undone.
			h.keeper.ClearUnbondingInformation(ctx, addr, pubKey)
		}
	}
}

// AfterOperatorKeyReplacement is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorKeyReplacement(
	ctx sdk.Context, addr sdk.AccAddress,
	newKey tmprotocrypto.PublicKey, oldKey tmprotocrypto.PublicKey,
	chainID string,
) {
	if strings.Compare(chainID, ctx.ChainID()) == 0 {
		// res == Removed, it means operator has added their original key again
		// res == Success, there is no additional information to store
		// res == Exists, there is no nothing to do
		if res := h.keeper.QueueOperation(
			ctx, addr, newKey, types.KeyAdditionOrUpdate,
		); res == types.QueueResultRemoved {
			// see AfterOperatorOptIn for explanation
			h.keeper.ClearUnbondingInformation(ctx, addr, newKey)
		}
		// res == Removed, it means operator had added this key and is now removing it.
		// no additional information to clear.
		// res == Success, the old key should be pruned from the operator module.
		// res == Exists, there is nothing to do.
		if res := h.keeper.QueueOperation(
			ctx, addr, oldKey, types.KeyRemoval,
		); res == types.QueueResultSuccess {
			// the old key can be marked for pruning
			h.keeper.SetUnbondingInformation(ctx, addr, oldKey, false)
		}
	}
}

// AfterOperatorOptOutInitiated is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorOptOutInitiated(
	ctx sdk.Context, addr sdk.AccAddress,
	chainID string, pubKey tmprotocrypto.PublicKey,
) {
	if strings.Compare(chainID, ctx.ChainID()) == 0 {
		// res == Removed means operator had opted in and is now opting out. nothing to do if
		// it is within the same epoch.
		// res == Success, set up pruning deadline and opt out completion deadline
		// res == Exists, there is nothing to do (should never happen)
		if res := h.keeper.QueueOperation(
			ctx, addr, pubKey, types.KeyRemoval,
		); res == types.QueueResultSuccess {
			h.keeper.SetUnbondingInformation(ctx, addr, pubKey, true)
		}
	}
}
