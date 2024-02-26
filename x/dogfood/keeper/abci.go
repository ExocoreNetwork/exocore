package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	// start with undelegations
	undelegations := k.GetPendingUndelegations(ctx)
	for _, undelegation := range undelegations.GetList() {
		k.delegationKeeper.DecrementUndelegationHoldCount(ctx, undelegation)
	}
	k.ClearPendingUndelegations(ctx)
	// then opt outs (consensus addresses should be done after opt out)
	optOuts := k.GetPendingOptOuts(ctx)
	for _, addr := range optOuts.GetList() {
		k.operatorKeeper.CompleteOperatorOptOutFromChainId(ctx, addr, ctx.ChainID())
	}
	k.ClearPendingOptOuts(ctx)
	// then consensus addresses
	consensusAddrs := k.GetPendingConsensusAddrs(ctx)
	for _, consensusAddr := range consensusAddrs.GetList() {
		k.operatorKeeper.DeleteOperatorAddressForChainIdAndConsAddr(
			ctx, ctx.ChainID(), consensusAddr,
		)
	}
	k.ClearPendingConsensusAddrs(ctx)
	// finally, operations
	operations := k.GetPendingOperations(ctx)
	res := make([]abci.ValidatorUpdate, 0, len(operations.GetList()))
	for _, operation := range operations.GetList() {
		switch operation.OperationType {
		case types.KeyAdditionOrUpdate:
			power, err := k.restakingKeeper.GetOperatorAssetValue(
				ctx, operation.OperatorAddress,
			)
			if err != nil {
				panic(err)
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
			panic("unspecified operation type")
		}
	}
	return res
}

// SetPendingOperations sets the pending operations to be applied at the end of the block.
func (k Keeper) SetPendingOperations(ctx sdk.Context, operations types.Operations) {
	store := ctx.KVStore(k.storeKey)
	bz, err := operations.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.PendingOperationsKey(), bz)
}

// GetPendingOperations returns the pending operations to be applied at the end of the block.
func (k Keeper) GetPendingOperations(ctx sdk.Context) types.Operations {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingOperationsKey())
	if bz == nil {
		return types.Operations{}
	}
	var operations types.Operations
	if err := operations.Unmarshal(bz); err != nil {
		panic(err)
	}
	return operations
}

// ClearPendingOperations clears the pending operations to be applied at the end of the block.
func (k Keeper) ClearPendingOperations(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingOperationsKey())
}

// SetPendingOptOuts sets the pending opt-outs to be applied at the end of the block.
func (k Keeper) SetPendingOptOuts(ctx sdk.Context, addrs types.AccountAddresses) {
	store := ctx.KVStore(k.storeKey)
	bz, err := addrs.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.PendingOptOutsKey(), bz)
}

// GetPendingOptOuts returns the pending opt-outs to be applied at the end of the block.
func (k Keeper) GetPendingOptOuts(ctx sdk.Context) types.AccountAddresses {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingOptOutsKey())
	if bz == nil {
		return types.AccountAddresses{}
	}
	var addrs types.AccountAddresses
	if err := addrs.Unmarshal(bz); err != nil {
		panic(err)
	}
	return addrs
}

// ClearPendingOptOuts clears the pending opt-outs to be applied at the end of the block.
func (k Keeper) ClearPendingOptOuts(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingOptOutsKey())
}

// SetPendingConsensusAddrs sets the pending consensus addresses to be pruned at the end of the
// block.
func (k Keeper) SetPendingConsensusAddrs(ctx sdk.Context, addrs types.ConsensusAddresses) {
	store := ctx.KVStore(k.storeKey)
	bz, err := addrs.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.PendingConsensusAddrsKey(), bz)
}

// GetPendingConsensusAddrs returns the pending consensus addresses to be pruned at the end of
// the block.
func (k Keeper) GetPendingConsensusAddrs(ctx sdk.Context) types.ConsensusAddresses {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingConsensusAddrsKey())
	if bz == nil {
		return types.ConsensusAddresses{}
	}
	var addrs types.ConsensusAddresses
	if err := addrs.Unmarshal(bz); err != nil {
		panic(err)
	}
	return addrs
}

// ClearPendingConsensusAddrs clears the pending consensus addresses to be pruned at the end of
// the block.
func (k Keeper) ClearPendingConsensusAddrs(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingConsensusAddrsKey())
}

// SetPendingUndelegations sets the pending undelegations to be released at the end of the
// block.
func (k Keeper) SetPendingUndelegations(ctx sdk.Context, undelegations types.UndelegationRecordKeys) {
	store := ctx.KVStore(k.storeKey)
	bz, err := undelegations.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.PendingUndelegationsKey(), bz)
}

// GetPendingUndelegations returns the pending undelegations to be released at the end of the
// block.
func (k Keeper) GetPendingUndelegations(ctx sdk.Context) types.UndelegationRecordKeys {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingUndelegationsKey())
	if bz == nil {
		return types.UndelegationRecordKeys{}
	}
	var undelegations types.UndelegationRecordKeys
	if err := undelegations.Unmarshal(bz); err != nil {
		panic(err)
	}
	return undelegations
}

// ClearPendingUndelegations clears the pending undelegations to be released at the end of the
// block.
func (k Keeper) ClearPendingUndelegations(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingUndelegationsKey())
}
