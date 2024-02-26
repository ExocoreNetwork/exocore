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
// O I (reversing the decision to opt out)
// O I R
// R O I
// The impossible list looks like:
// O R
// O R I
// Replacing the key with the same key is not possible, so it is not covered.
// The assumption is that these are enforced (allowed/denied) by the operator module.

// Cases for a fresh operator:
// I + O => KeyAdditionOrUpdate + KeyRemoval => Success + Removed => covered.
// I + R => KeyAdditionOrUpdate (old) + KeyAdditionOrUpdate (new) + KeyRemoval (old) =>
// 			Success + Success + Removed + => not considered unbonding of the old key since it
// 											 was not yet effective => covered.
// I + R + O => KeyAdditionOrUpdate (old) + KeyAdditionOrUpdate (new) + KeyRemoval (old) +
// 				KeyRemoval (new) =>
// 				Success (old) + Success (new) + Removed (old) + Removed (new) => covered

// Cases for an operator that has already opted in:
// R + O => KeyAdditionOrUpdate (new) + KeyRemoval (old) + KeyRemoval (new) =>
// 			Success (new) + Success (old) + Removed (new) =>
// 							unbonding data made (old) => covered.
// O + I
// O + I case 1 => KeyRemoval (old) + KeyAdditionOrUpdate (new) => Success (old) + Success (new)
//				=> unbonding data made (old) => covered.
// O + I case 2 => KeyRemoval (old) + KeyAdditionOrUpdate (old) => Success (old) + Removed (old)
//				=> unbonding data made (old) and then cleared => covered.
// O + I + R
// O + I + R case 1 => KeyRemoval (old) + KeyAdditionOrUpdate (old) + KeyRemoval (old) +
// 				   KeyAdditionOrUpdate (new) => Success (old) + Removed (old) + Success (old) 					   +
// Success (new) => unbonding data old made + cleared + made => covered.
// O + I + R case 2 =>
// AfterOperatorOptOut(old) => Success => unbonding data made for old
// AfterOperatorOptIn(new) => Success => no data changed
// AfterOperatorKeyReplacement(new, new2) =>
//    new2 operation KeyAdditionOrUpdate => Success => no data changed
//    new  operation KeyRemoval => Removed => no data changed
// => covered
// R + O + I => KeyAdditionOrUpdate (new) + KeyRemoval (old) + KeyRemoval (new) +
// KeyAdditionOrUpdate (X)
//              Success + Success (=> unbonding data for old) + Removed (=> no data) +
// case 1 => X == new => Success => no change => covered
// case 2 => X == old => Removed => unbonding data for old removed => covered.
// case 3 => X == abc => Success => no change and unbonding data for old is not removed =>
// covered.

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
		// remove the old key
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
		// add the new key
		// res == Removed, it means operator has added their original key again
		// res == Success, there is no additional information to store
		// res == Exists, there is no nothing to do
		if res := h.keeper.QueueOperation(
			ctx, addr, newKey, types.KeyAdditionOrUpdate,
		); res == types.QueueResultRemoved {
			// see AfterOperatorOptIn for explanation
			h.keeper.ClearUnbondingInformation(ctx, addr, newKey)
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
