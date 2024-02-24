package types

import (
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ OperatorConsentHooks = &MultiOperatorConsentHooks{}

type MultiOperatorConsentHooks []OperatorConsentHooks

func NewMultiOperatorConsentHooks(hooks ...OperatorConsentHooks) MultiOperatorConsentHooks {
	return hooks
}

func (hooks MultiOperatorConsentHooks) AfterOperatorOptIn(
	ctx sdk.Context,
	addr sdk.AccAddress,
	chainId string,
	pubKey tmprotocrypto.PublicKey,
) {
	for _, hook := range hooks {
		hook.AfterOperatorOptIn(ctx, addr, chainId, pubKey)
	}
}

func (hooks MultiOperatorConsentHooks) AfterOperatorKeyReplacement(
	ctx sdk.Context,
	addr sdk.AccAddress,
	oldKey tmprotocrypto.PublicKey,
	newAddr tmprotocrypto.PublicKey,
	chainId string,
) {
	for _, hook := range hooks {
		hook.AfterOperatorKeyReplacement(ctx, addr, oldKey, newAddr, chainId)
	}
}

func (hooks MultiOperatorConsentHooks) AfterOperatorOptOutInitiated(
	ctx sdk.Context, addr sdk.AccAddress, chainId string, key tmprotocrypto.PublicKey,
) {
	for _, hook := range hooks {
		hook.AfterOperatorOptOutInitiated(ctx, addr, chainId, key)
	}
}
