package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ClearUnbondingInformation clears all information related to an operator's opt out
// or key replacement. This is done because the operator has opted back in or has
// replaced their key (again) with the original one.
func (k Keeper) ClearUnbondingInformation(
	ctx sdk.Context, addr sdk.AccAddress, pubKey tmprotocrypto.PublicKey,
) {
	optOutEpoch := k.GetOperatorOptOutFinishEpoch(ctx, addr)
	k.DeleteOperatorOptOutFinishEpoch(ctx, addr)
	k.RemoveOptOutToFinish(ctx, optOutEpoch, addr)
	consAddress, err := types.TMCryptoPublicKeyToConsAddr(pubKey)
	if err != nil {
		panic(err)
	}
	k.DeleteConsensusAddrToPrune(ctx, optOutEpoch, consAddress)
}

// SetUnbondingInformation sets information related to an operator's opt out or key replacement.
func (k Keeper) SetUnbondingInformation(
	ctx sdk.Context, addr sdk.AccAddress, pubKey tmprotocrypto.PublicKey, isOptingOut bool,
) {
	unbondingCompletionEpoch := k.GetUnbondingCompletionEpoch(ctx)
	k.AppendOptOutToFinish(ctx, unbondingCompletionEpoch, addr)
	if isOptingOut {
		k.SetOperatorOptOutFinishEpoch(ctx, addr, unbondingCompletionEpoch)
	}
	consAddress, err := types.TMCryptoPublicKeyToConsAddr(pubKey)
	if err != nil {
		panic(err)
	}
	k.AppendConsensusAddrToPrune(ctx, unbondingCompletionEpoch, consAddress)
}

// GetUnbondingCompletionEpoch returns the epoch at the end of which
// an unbonding triggered in this epoch will be completed.
func (k Keeper) GetUnbondingCompletionEpoch(
	ctx sdk.Context,
) int64 {
	unbondingEpochs := k.GetEpochsUntilUnbonded(ctx)
	epochInfo, found := k.epochsKeeper.GetEpochInfo(
		ctx, k.GetEpochIdentifier(ctx),
	)
	if !found {
		panic("current epoch not found")
	}
	// if i execute the transaction at epoch 5, the vote power change
	// goes into effect at the beginning of epoch 6. the information
	// should be held for 7 epochs, so it should be deleted at the
	// beginning of epoch 13 or the end of epoch 12.
	return epochInfo.CurrentEpoch + int64(unbondingEpochs)
}
