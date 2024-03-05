package keeper

import (
	"sort"

	abci "github.com/cometbft/cometbft/abci/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	id, _ := k.getValidatorSetID(ctx, ctx.BlockHeight())
	if !k.IsEpochEnd(ctx) {
		// save the same id for the next block height.
		k.setValidatorSetID(ctx, ctx.BlockHeight()+1, id)
		return []abci.ValidatorUpdate{}
	}
	defer k.ClearEpochEnd(ctx)
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
	// 1. find all operator keys for the chain.
	// 2. find last stored operator keys + their powers.
	// 3. find newest vote power for the operator keys, and sort them.
	// 4. loop through #1 and see if anything has changed.
	//    if it hasn't, do nothing for that operator key.
	//    if it has, queue an update.
	prev := k.getKeyPowerMapping(ctx).List
	res := make([]abci.ValidatorUpdate, 0, len(prev))
	operators, keys := k.operatorKeeper.GetActiveOperatorsForChainId(ctx, ctx.ChainID())
	powers, err := k.restakingKeeper.GetAvgDelegatedValue(
		ctx, operators, k.GetAssetIDs(ctx), k.GetEpochIdentifier(ctx),
	)
	if err != nil {
		return []abci.ValidatorUpdate{}
	}
	operators, keys, powers = sortByPower(operators, keys, powers)
	maxVals := k.GetMaxValidators(ctx)
	for i := range operators {
		// #nosec G701 // ok if 64-bit.
		if i >= int(maxVals) {
			// we have reached the maximum number of validators.
			break
		}
		key := keys[i]
		power := powers[i]
		if power < 1 {
			// we have reached the bottom of the rung.
			break
		}
		// find the previous power.
		prevPower, found := prev[key.String()]
		if found && prevPower == power {
			delete(prev, key.String())
			continue
		}
		// either the key was not in the previous set,
		// or the power has changed.
		res = append(res, abci.ValidatorUpdate{
			PubKey: key,
			Power:  power,
		})
	}
	// the remaining keys in prev have lost their power.
	// gosec does not like `for key := range prev` while others do not like
	// `for key, _ := range prev`
	// #nosec G705
	for key := range prev {
		bz := []byte(key)
		var keyObj tmprotocrypto.PublicKey
		k.cdc.MustUnmarshal(bz, &keyObj)
		res = append(res, abci.ValidatorUpdate{
			PubKey: keyObj,
			Power:  0,
		})
	}
	// call via wrapper function so that validator info is stored.
	// the id is incremented by 1 for the next block.
	return k.ApplyValidatorChanges(ctx, res, id+1)
}

func sortByPower(
	operatorAddrs []sdk.AccAddress,
	pubKeys []tmprotocrypto.PublicKey,
	powers []int64,
) ([]sdk.AccAddress, []tmprotocrypto.PublicKey, []int64) {
	// Create a slice of indices
	indices := make([]int, len(powers))
	for i := range indices {
		indices[i] = i
	}

	// Sort the indices slice based on the powers slice
	sort.SliceStable(indices, func(i, j int) bool {
		return powers[indices[i]] > powers[indices[j]]
	})

	// Reorder all slices using the sorted indices
	sortedOperatorAddrs := make([]sdk.AccAddress, len(operatorAddrs))
	sortedPubKeys := make([]tmprotocrypto.PublicKey, len(pubKeys))
	sortedPowers := make([]int64, len(powers))
	for i, idx := range indices {
		sortedOperatorAddrs[i] = operatorAddrs[idx]
		sortedPubKeys[i] = pubKeys[idx]
		sortedPowers[i] = powers[idx]
	}

	return sortedOperatorAddrs, sortedPubKeys, sortedPowers
}
