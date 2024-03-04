package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	// start with clearing the hold on the undelegations.
	undelegations := k.GetPendingUndelegations(ctx)
	for _, undelegation := range undelegations.GetList() {
		k.delegationKeeper.DecrementUndelegationHoldCount(ctx, undelegation)
	}
	k.ClearPendingUndelegations(ctx)
	// then, let the operator module know that the opt out has finished.
	optOuts := k.GetPendingOptOuts(ctx)
	for _, addr := range optOuts.GetList() {
		k.operatorKeeper.CompleteOperatorOptOutFromChainId(ctx, addr, ctx.ChainID())
	}
	k.ClearPendingOptOuts(ctx)
	// for slashing, the operator module is required to store a mapping of chain id + cons addr
	// to operator address. this information can now be pruned, since the opt out is considered
	// complete.
	consensusAddrs := k.GetPendingConsensusAddrs(ctx)
	for _, consensusAddr := range consensusAddrs.GetList() {
		k.operatorKeeper.DeleteOperatorAddressForChainIdAndConsAddr(
			ctx, ctx.ChainID(), consensusAddr,
		)
	}
	k.ClearPendingConsensusAddrs(ctx)
	// finally, perform the actual operations of vote power changes.
	operations := k.GetPendingOperations(ctx)
	id, _ := k.getValidatorSetID(ctx, ctx.BlockHeight())
	if len(operations.GetList()) == 0 {
		// there is no validator set change, so we just increment the block height
		// and retain the same val set id mapping.
		k.setValidatorSetID(ctx, ctx.BlockHeight()+1, id)
		return []abci.ValidatorUpdate{}
	}
	res := make([]abci.ValidatorUpdate, 0, len(operations.GetList()))
	for _, operation := range operations.GetList() {
		switch operation.OperationType {
		case types.KeyAdditionOrUpdate:
			power, err := k.restakingKeeper.GetOperatorAssetValue(
				ctx, operation.OperatorAddress,
			)
			if err != nil {
				// this should never happen, but if it does, we just skip the operation.
				continue
			}
			res = append(res, abci.ValidatorUpdate{
				PubKey: operation.PubKey,
				Power:  power,
			})
		case types.KeyRemoval:
			res = append(res, abci.ValidatorUpdate{
				PubKey: operation.PubKey,
				Power:  0,
			})
		case types.KeyOpUnspecified:
			// this should never happen, but if it does, we just skip the operation.
			continue
		}
	}
	// call via wrapper function so that validator info is stored.
	// the id is incremented by 1 for the next block.
	return k.ApplyValidatorChanges(ctx, res, id+1, false)
}
