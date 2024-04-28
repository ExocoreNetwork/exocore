package keeper

import (
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
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
	ctx sdk.Context, addr sdk.AccAddress, chainID string, key *tmprotocrypto.PublicKey,
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
	ctx sdk.Context, operator sdk.AccAddress, oldKey *tmprotocrypto.PublicKey,
	newKey *tmprotocrypto.PublicKey, chainID string,
) {
	// the impact of key replacement is:
	// 1. vote power of old key is 0, which happens automatically at epoch end in EndBlock. this
	// is because the key is in the previous set but not in the new one and our code will queue
	// a validator update of 0 fot this.
	// 2. vote power of new key is calculated, which happens automatically at epoch end in
	// EndBlock.
	// 3. X epochs later, the reverse lookup of old cons addr + chain id -> operator addr
	// should be cleared. however, if the key is replaced again with the old one, this reverse
	// lookup should be retained.
}

// AfterOperatorKeyRemovalInitiated is the implementation of the operator hooks.
func (h OperatorHooksWrapper) AfterOperatorKeyRemovalInitiated(
	ctx sdk.Context, operator sdk.AccAddress, chainID string, key *tmprotocrypto.PublicKey,
) {
	// the impact of key removal is:
	// 1. vote power of the operator is 0, which happens automatically at epoch end in EndBlock.
	// this is because GetActiveOperatorsForChainID filters operators who are removing their
	// keys from the chain.
	// 2. X epochs later, the removal is marked complete in the operator module.
	if chainID == ctx.ChainID() {
		h.keeper.SetUnbondingInformation(ctx, operator, key)
	}
}
