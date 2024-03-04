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

// AfterOperatorOptIn is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorOptIn(
	ctx sdk.Context, addr sdk.AccAddress, chainID string, key tmprotocrypto.PublicKey,
) {
	// an operator opting in does not meaningfully affect this module, since
	// this information will be fetched at the end of the epoch
	// and the operator's vote power will be calculated then.
	// however, we will still clear the unbonding information, if it exists.
	h.keeper.ClearUnbondingInformation(ctx, addr, key)
}

// AfterOperatorKeyReplacement is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorKeyReplacement(
	ctx sdk.Context, operator sdk.AccAddress, oldKey tmprotocrypto.PublicKey,
	newKey tmprotocrypto.PublicKey, chainID string,
) {
	if strings.Compare(chainID, ctx.ChainID()) == 0 {
		// a key replacement means that the old key needs to be pruned upon maturity.
		h.keeper.SetUnbondingInformation(ctx, operator, oldKey, false)
		h.keeper.ClearUnbondingInformation(ctx, operator, newKey)
	}
	// 	// remove the old key
	// 	// res == Removed, it means operator had added this key and is now removing it.
	// 	// no additional information to clear.
	// 	// res == Success, the old key should be pruned from the operator module.
	// 	// res == Exists, there is nothing to do.
	// 	if res := h.keeper.QueueOperation(
	// 		ctx, addr, oldKey, types.KeyRemoval,
	// 	); res == types.QueueResultSuccess {
	// 		// the old key can be marked for pruning
	// 		h.keeper.SetUnbondingInformation(ctx, addr, oldKey, false)
	// 	}
	// 	// add the new key
	// 	// res == Removed, it means operator has added their original key again
	// 	// res == Success, there is no additional information to store
	// 	// res == Exists, there is no nothing to do
	// 	if res := h.keeper.QueueOperation(
	// 		ctx, addr, newKey, types.KeyAdditionOrUpdate,
	// 	); res == types.QueueResultRemoved {
	// 		// see AfterOperatorOptIn for explanation
	// 		h.keeper.ClearUnbondingInformation(ctx, addr, newKey)
	// 	}
	// }
}

// AfterOperatorOptOutInitiated is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorOptOutInitiated(
	ctx sdk.Context, operator sdk.AccAddress, chainID string, key tmprotocrypto.PublicKey,
) {
	if strings.Compare(chainID, ctx.ChainID()) == 0 {
		h.keeper.SetUnbondingInformation(ctx, operator, key, true)
	}
}
