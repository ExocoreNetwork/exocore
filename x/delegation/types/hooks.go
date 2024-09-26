package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
) error {
	for _, hook := range hooks {
		err := hook.AfterUndelegationStarted(ctx, addr, recordKey)
		if err != nil {
			return err
		}
	}
	return nil
}
