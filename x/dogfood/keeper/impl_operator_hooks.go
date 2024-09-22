package keeper

import (
	exocoretypes "github.com/ExocoreNetwork/exocore/types"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OperatorHooksWrapper is the wrapper structure that implements the operator hooks for the
// dogfood keeper.
type OperatorHooksWrapper struct {
	keeper *Keeper
}

// Interface guards
var _ operatortypes.OperatorHooks = OperatorHooksWrapper{}

func (k *Keeper) OperatorHooks() OperatorHooksWrapper {
	return OperatorHooksWrapper{k}
}

// AfterOperatorKeySet is the implementation of the operator hooks.
// CONTRACT: an operator cannot set their key if they are already in the process of removing it.
func (h OperatorHooksWrapper) AfterOperatorKeySet(
	sdk.Context, sdk.AccAddress, string, exocoretypes.WrappedConsKey,
) {
	// an operator opting in does not meaningfully affect this module, since
	// this information will be fetched at the end of the epoch
	// and the operator's vote power will be calculated then.
}

// AfterOperatorKeyReplaced is the implementation of the operator hooks.
// CONTRACT: key replacement is not allowed if the operator is in the process of removing their
// key.
// CONTRACT: key replacement from newKey to oldKey is not allowed, after a replacement from
// oldKey to newKey.
func (h OperatorHooksWrapper) AfterOperatorKeyReplaced(
	ctx sdk.Context, _ sdk.AccAddress, oldKey exocoretypes.WrappedConsKey,
	_ exocoretypes.WrappedConsKey, chainID string,
) {
	// the impact of key replacement is:
	// 1. vote power of old key is 0, which happens automatically at epoch end in EndBlock. this
	// is because the key is in the previous set but not in the new one and our code will queue
	// a validator update of 0 for this.
	// 2. vote power of new key is calculated, which happens automatically at epoch end in
	// EndBlock.
	// 3. X epochs later, the reverse lookup of old cons addr + chain id -> operator addr
	// should be cleared.
	consAddr := oldKey.ToConsAddr()
	if chainID == avstypes.ChainIDWithoutRevision(ctx.ChainID()) {
		// is the oldKey already active? if not, we should not do anything.
		// this can happen if we opt in with a key, then replace it with another key
		// during the same epoch.
		_, found := h.keeper.GetExocoreValidator(ctx, consAddr)
		if found {
			unbondingEpoch := h.keeper.GetUnbondingCompletionEpoch(ctx)
			// nb: if operator sets key, it is not "at stake" till the end of the epoch.
			// before that time, any key replacement will store a superfluous entry for pruning
			// since the old key will not be in use.
			// this technically gives an operator the opportunity to spam the pruning queue
			// but it is not a security risk or a DOS vector given the cost charged.
			h.keeper.AppendConsensusAddrToPrune(ctx, unbondingEpoch, consAddr)
		} else {
			// since this consAddr isn't active, we can remove it immediately.
			h.keeper.operatorKeeper.DeleteOperatorAddressForChainIDAndConsAddr(
				ctx, chainID, consAddr,
			)
		}
	}
}

// AfterOperatorKeyRemovalInitiated is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorKeyRemovalInitiated(
	ctx sdk.Context, operator sdk.AccAddress, chainID string, key exocoretypes.WrappedConsKey,
) {
	// the impact of key removal is:
	// 1. vote power of the operator is 0, which happens automatically at epoch end in EndBlock.
	// this is because GetActiveOperatorsForChainID filters operators who are removing their
	// keys from the chain.
	// 2. X epochs later, the removal is marked complete in the operator module.
	consAddr := key.ToConsAddr()
	if chainID == avstypes.ChainIDWithoutRevision(ctx.ChainID()) {
		_, found := h.keeper.GetExocoreValidator(ctx, consAddr)
		if found {
			h.keeper.SetOptOutInformation(ctx, operator)
		} else {
			h.keeper.operatorKeeper.DeleteOperatorAddressForChainIDAndConsAddr(
				ctx, chainID, consAddr,
			)
		}
	}
}
