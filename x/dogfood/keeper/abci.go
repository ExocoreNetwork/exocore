package keeper

import (
	"sort"

	"cosmossdk.io/math"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) {
	// for IBC, track historical validator set
	k.TrackHistoricalInfo(ctx)
}

func (k Keeper) EndBlock(ctx sdk.Context) []abci.ValidatorUpdate {
	if !k.IsEpochEnd(ctx) {
		k.SetValidatorUpdates(ctx, []abci.ValidatorUpdate{})
		return []abci.ValidatorUpdate{}
	}
	defer k.ClearEpochEnd(ctx)
	// start with clearing the hold on the undelegations.
	undelegations := k.GetPendingUndelegations(ctx)
	for _, undelegation := range undelegations.GetList() {
		err := k.delegationKeeper.DecrementUndelegationHoldCount(ctx, undelegation)
		if err != nil {
			panic(err)
		}
		k.ClearUndelegationMaturityEpoch(ctx, undelegation)
	}
	k.ClearPendingUndelegations(ctx)
	// then, let the operator module know that the opt out has finished.
	optOuts := k.GetPendingOptOuts(ctx)
	for _, addr := range optOuts.GetList() {
		// TODO log the error
		_ = k.operatorKeeper.CompleteOperatorKeyRemovalForChainID(ctx, addr, ctx.ChainID())
	}
	k.ClearPendingOptOuts(ctx)
	// for slashing, the operator module is required to store a mapping of chain id + cons addr
	// to operator address. this information can now be pruned, since the opt out is considered
	// complete.
	consensusAddrs := k.GetPendingConsensusAddrs(ctx)
	for _, consensusAddr := range consensusAddrs.GetList() {
		k.operatorKeeper.DeleteOperatorAddressForChainIDAndConsAddr(
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
	// 5. keep in mind the total vote power.
	totalPower := math.ZeroInt()
	prevList := k.GetAllExocoreValidators(ctx)
	// prevMap is a map of the previous validators, indexed by the consensus address
	// and the value being the vote power.
	prevMap := make(map[string]int64, len(prevList))
	for _, validator := range prevList {
		pubKey, err := validator.ConsPubKey()
		if err != nil {
			// indicates an error in deserialization, and should never happen.
			continue
		}
		addressString := sdk.GetConsAddress(pubKey).String()
		prevMap[addressString] = validator.Power
	}
	operators, keys := k.operatorKeeper.GetActiveOperatorsForChainID(ctx, ctx.ChainID())
	powers, err := k.operatorKeeper.GetAvgDelegatedValue(
		ctx, operators, ctx.ChainID(), k.GetEpochIdentifier(ctx),
	)
	if err != nil {
		return []abci.ValidatorUpdate{}
	}
	operators, keys, powers = sortByPower(operators, keys, powers)
	maxVals := k.GetMaxValidators(ctx)
	// the capacity of this list is twice the maximum number of validators.
	// this is because we can have a maximum of maxVals validators, and we can also have
	// a maximum of maxVals validators that are removed.
	res := make([]abci.ValidatorUpdate, 0, maxVals*2)
	for i := range operators {
		// #nosec G701 // ok on 64-bit systems.
		if i > int(maxVals) {
			// we have reached the maximum number of validators, amongst all the validators.
			// even if there are intersections with the previous validator set, this will
			// only be reached if we exceed the threshold.
			// if there are no intersections, this case is glaringly obvious.
			break
		}
		power := powers[i]
		if power < 1 {
			// we have reached the bottom of the rung.
			// assumption is that negative vote power isn't provided by the module.
			// the consensus engine will reject it anyway and panic.
			break
		}
		// find the previous power.
		key := keys[i]
		address, err := operatortypes.TMCryptoPublicKeyToConsAddr(key)
		if err != nil {
			// indicates an error in deserialization, and should never happen.
			continue
		}
		addressString := address.String()
		prevPower, found := prevMap[addressString]
		if found {
			// if the power has changed, queue an update.
			if prevPower != power {
				res = append(res, abci.ValidatorUpdate{
					PubKey: *key,
					Power:  power,
				})
			}
			// remove the validator from the previous map, so that 0 power
			// is not queued for it.
			delete(prevMap, addressString)
		} else {
			// new consensus key, queue an update.
			res = append(res, abci.ValidatorUpdate{
				PubKey: *key,
				Power:  power,
			})
		}
		totalPower = totalPower.Add(sdk.NewInt(power))
	}
	// the remaining validators in prevMap have been removed.
	// we need to queue a change in power to 0 for them.
	for _, validator := range prevList { // O(N)
		// #nosec G703 // already checked in the previous iteration over prevList.
		pubKey, _ := validator.ConsPubKey()
		addressString := sdk.GetConsAddress(pubKey).String()
		// Check if this validator is still in prevMap (i.e., hasn't been deleted)
		if _, exists := prevMap[addressString]; exists { // O(1) since hash map
			tmprotoKey, err := cryptocodec.ToTmProtoPublicKey(pubKey)
			if err != nil {
				continue
			}
			res = append(res, abci.ValidatorUpdate{
				PubKey: tmprotoKey,
				Power:  0,
			})
			// while calculating total power, we started with 0 and not previous power.
			// so the previous power of these validators does not need to be subtracted.
		}
	}
	// if there are any updates, set total power on lookup index.
	if len(res) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}

	// call via wrapper function so that validator info is stored.
	return k.ApplyValidatorChanges(ctx, res)
}

// sortByPower sorts operators, their pubkeys, and their powers by the powers.
// the sorting is descending, so the highest power is first.
func sortByPower(
	operatorAddrs []sdk.AccAddress,
	pubKeys []*tmprotocrypto.PublicKey,
	powers []int64,
) ([]sdk.AccAddress, []*tmprotocrypto.PublicKey, []int64) {
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
	sortedPubKeys := make([]*tmprotocrypto.PublicKey, len(pubKeys))
	sortedPowers := make([]int64, len(powers))
	for i, idx := range indices {
		sortedOperatorAddrs[i] = operatorAddrs[idx]
		sortedPubKeys[i] = pubKeys[idx]
		sortedPowers[i] = powers[idx]
	}

	return sortedOperatorAddrs, sortedPubKeys, sortedPowers
}
