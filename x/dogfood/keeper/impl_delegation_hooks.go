package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationHooksWrapper is the wrapper structure that implements the delegation hooks for the
// dogfood keeper.
type DelegationHooksWrapper struct {
	keeper *Keeper
}

// Interface guard
var _ types.DelegationHooks = DelegationHooksWrapper{}

// DelegationHooks returns the delegation hooks wrapper. It follows the "accept interfaces,
// return concretes" pattern.
func (k *Keeper) DelegationHooks() DelegationHooksWrapper {
	return DelegationHooksWrapper{k}
}

// AfterDelegation is called after a delegation is made.
func (wrapper DelegationHooksWrapper) AfterDelegation(
	ctx sdk.Context, operator sdk.AccAddress,
) {
	found, pubKey, err := wrapper.keeper.operatorKeeper.GetOperatorConsKeyForChainId(
		ctx, operator, ctx.ChainID(),
	)
	if err != nil {
		// the operator keeper can offer two errors: not an operator and not a chain.
		// both of these should not happen here because the dogfooding genesis will
		// register the chain, and the operator must be known to the delegation module
		// when it calls this hook.
		panic(err)
	}
	if found {
		if !wrapper.keeper.operatorKeeper.IsOperatorOptingOutFromChainId(
			ctx, operator, ctx.ChainID(),
		) {
			// only queue the operation if operator is still opted into the chain.
			res := wrapper.keeper.QueueOperation(
				ctx, operator, pubKey, types.KeyAdditionOrUpdate,
			)
			switch res {
			case types.QueueResultExists:
				// nothing to do because the operation is in the queue already.
			case types.QueueResultRemoved:
				// a KeyRemoval was in the queue which has now been cleared from the queue.
				// the KeyRemoval can only be in the queue if the operator is opting out from
				// the chain, or has replaced their key. if it is the former, it means that
				// there is some inconsistency. if it is the latter, it means that the operator
				// module just reported the old key in `GetOperatorConsKeyForChainId`, which
				// should not happen.
				panic("unexpected removal of operation from queue")
			case types.QueueResultSuccess:
				// best case, nothing to do.
			case types.QueueResultUnspecified:
				panic("unspecified queue result")
			}
		}
	}
}

// AfterUndelegationStarted is called after an undelegation is started.
func (wrapper DelegationHooksWrapper) AfterUndelegationStarted(
	ctx sdk.Context, operator sdk.AccAddress, recordKey []byte,
) {
	found, pubKey, err := wrapper.keeper.operatorKeeper.GetOperatorConsKeyForChainId(
		ctx, operator, ctx.ChainID(),
	)
	if err != nil {
		panic(err)
	}
	if found {
		// note that this is still key addition or update because undelegation does not remove
		// the operator from the list. it only decreases their vote power.
		if !wrapper.keeper.operatorKeeper.IsOperatorOptingOutFromChainId(
			ctx, operator, ctx.ChainID(),
		) {
			// only queue the operation if operator is still opted into the chain.
			res := wrapper.keeper.QueueOperation(
				ctx, operator, pubKey, types.KeyAdditionOrUpdate,
			)
			switch res {
			case types.QueueResultExists:
				// nothing to do
			case types.QueueResultRemoved:
				// KeyRemoval + KeyAdditionOrUpdate => Removed
				// KeyRemoval can happen
				// 1. if the operator is opting out from the chain,which is inconsistent.
				// 2. if the operator is replacing their old key, which should not be returned
				//    by `GetOperatorConsKeyForChainId`.
				panic("unexpected removal of operation from queue")
			case types.QueueResultSuccess:
				// best case, nothing to do.
			case types.QueueResultUnspecified:
				panic("unspecified queue result")
			}
		}
		// now handle the unbonding timeline.
		wrapper.keeper.delegationKeeper.IncrementUndelegationHoldCount(ctx, recordKey)
		// mark for unbonding release.
		// note that we aren't supporting redelegation yet, so this undelegated amount will be
		// held until the end of the unbonding period or the operator opt out period, whichever
		// is first.
		var unbondingCompletionEpoch int64
		if wrapper.keeper.operatorKeeper.IsOperatorOptingOutFromChainId(
			ctx, operator, ctx.ChainID(),
		) {
			unbondingCompletionEpoch = wrapper.keeper.GetOperatorOptOutFinishEpoch(
				ctx, operator,
			)
		} else {
			unbondingCompletionEpoch = wrapper.keeper.GetUnbondingCompletionEpoch(ctx)
		}
		wrapper.keeper.AppendUndelegationToMature(ctx, unbondingCompletionEpoch, recordKey)
	}
}

// AfterUndelegationCompleted is called after an undelegation is completed.
func (DelegationHooksWrapper) AfterUndelegationCompleted(
	sdk.Context, sdk.AccAddress,
) {
	// no-op
}
