package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ OperatorHooks = &MultiOperatorHooks{}

type MultiOperatorHooks []OperatorHooks

func NewMultiOperatorHooks(hooks ...OperatorHooks) MultiOperatorHooks {
	return hooks
}

func (hooks MultiOperatorHooks) AfterOperatorKeySet(
	ctx sdk.Context,
	addr sdk.AccAddress,
	chainID string,
	pubKey WrappedConsKey,
) {
	for _, hook := range hooks {
		hook.AfterOperatorKeySet(ctx, addr, chainID, pubKey)
	}
}

func (hooks MultiOperatorHooks) AfterOperatorKeyReplaced(
	ctx sdk.Context,
	addr sdk.AccAddress,
	oldKey WrappedConsKey,
	newAddr WrappedConsKey,
	chainID string,
) {
	for _, hook := range hooks {
		hook.AfterOperatorKeyReplaced(ctx, addr, oldKey, newAddr, chainID)
	}
}

func (hooks MultiOperatorHooks) AfterOperatorKeyRemovalInitiated(
	ctx sdk.Context, addr sdk.AccAddress, chainID string, key WrappedConsKey,
) {
	for _, hook := range hooks {
		hook.AfterOperatorKeyRemovalInitiated(ctx, addr, chainID, key)
	}
}
