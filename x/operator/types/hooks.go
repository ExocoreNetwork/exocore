package types

import (
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ OperatorHooks = &MultiOperatorHooks{}

type MultiOperatorHooks []OperatorHooks

func NewMultiOperatorHooks(hooks ...OperatorHooks) MultiOperatorHooks {
	return hooks
}

func (hooks MultiOperatorHooks) AfterOperatorOptIn(
	ctx sdk.Context,
	addr sdk.AccAddress,
	chainID string,
	pubKey *tmprotocrypto.PublicKey,
) {
	for _, hook := range hooks {
		hook.AfterOperatorOptIn(ctx, addr, chainID, pubKey)
	}
}

func (hooks MultiOperatorHooks) AfterOperatorKeyReplacement(
	ctx sdk.Context,
	addr sdk.AccAddress,
	oldKey *tmprotocrypto.PublicKey,
	newAddr *tmprotocrypto.PublicKey,
	chainID string,
) {
	for _, hook := range hooks {
		hook.AfterOperatorKeyReplacement(ctx, addr, oldKey, newAddr, chainID)
	}
}

func (hooks MultiOperatorHooks) AfterOperatorOptOutInitiated(
	ctx sdk.Context, addr sdk.AccAddress, chainID string, key *tmprotocrypto.PublicKey,
) {
	for _, hook := range hooks {
		hook.AfterOperatorOptOutInitiated(ctx, addr, chainID, key)
	}
}
