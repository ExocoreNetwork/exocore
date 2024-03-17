package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ DelegationHooks = &MultiDelegationHooks{}

type MultiDelegationHooks []DelegationHooks

func NewMultiDelegationHooks(hooks ...DelegationHooks) MultiDelegationHooks {
	return hooks
}

func (hooks MultiDelegationHooks) AfterDelegation(ctx sdk.Context, operator sdk.AccAddress) {
	for _, hook := range hooks {
		hook.AfterDelegation(ctx, operator)
	}
}

func (hooks MultiDelegationHooks) AfterUndelegationStarted(
	ctx sdk.Context,
	addr sdk.AccAddress,
	recordKey []byte,
) {
	for _, hook := range hooks {
		hook.AfterUndelegationStarted(ctx, addr, recordKey)
	}
}

func (hooks MultiDelegationHooks) AfterUndelegationCompleted(ctx sdk.Context, addr sdk.AccAddress) {
	for _, hook := range hooks {
		hook.AfterUndelegationCompleted(ctx, addr)
	}
}
