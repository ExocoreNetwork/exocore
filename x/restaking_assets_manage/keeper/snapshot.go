package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// we need to ensure that a snapshot exists for the block heights at which
// a delegation/undelegation happened. these are the heights at which
// the validator set update is broadcasted to the chain(s), and they
// therefore represent the mapped infraction heights.

// updateOperatorLastSnapshotHeight stores the last snapshot height for an operator
// equal to the current block height.
func (k Keeper) updateOperatorLastSnapshotHeight(
	ctx sdk.Context, operator sdk.AccAddress,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.OperatorLastSnapshotHeightKey(operator),
		sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight())),
	)
}

// GetOperatorLastSnapshotHeight returns the last snapshot height for an operator
func (k Keeper) getOperatorLastSnapshotHeight(
	ctx sdk.Context, operator sdk.AccAddress,
) uint64 {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.OperatorLastSnapshotHeightKey(operator))
	return sdk.BigEndianToUint64(val)
}

// GetOperatorLastSnapshot returns the last snapshot for an operator
func (k Keeper) getOperatorLastSnapshot(
	ctx sdk.Context, operator sdk.AccAddress,
) (types.Snapshot, bool) {
	lastHeight := k.getOperatorLastSnapshotHeight(ctx, operator)
	if lastHeight == 0 {
		return types.Snapshot{}, false
	}
	return k.GetOperatorSnapshotAtHeight(ctx, operator, lastHeight)
}

// GetOperatorSnapshotAtHeight returns the snapshot for an operator at a given height.
// If no snapshot exists at that height, an empty snapshot is returned.
func (k Keeper) GetOperatorSnapshotAtHeight(
	ctx sdk.Context, operator sdk.AccAddress, height uint64,
) (types.Snapshot, bool) {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.OperatorSnapshotKey(operator, height))
	if val == nil {
		return types.Snapshot{}, false
	}
	var snapshot types.Snapshot
	k.cdc.MustUnmarshal(val, &snapshot)
	return snapshot, true
}

// SetOperatorSnapshot stores a snapshot for an operator at the current height
func (k Keeper) setOperatorSnapshot(
	ctx sdk.Context, operator sdk.AccAddress, snapshot types.Snapshot,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.OperatorSnapshotKey(operator, uint64(ctx.BlockHeight())),
		k.cdc.MustMarshal(&snapshot),
	)
	k.updateOperatorLastSnapshotHeight(ctx, operator)
}

// How is a slashing request sent to the coordinator?
// The subscriber chain finds the id of the val set update
// corresponding to the infraction height on the subcriber.
// It then sends a slashing request to the coordinator with
// the val set update id and the infraction type.
// The coordinator then finds the Exocore height for which
// the val set update was broadcasted to the subscriber.
// It then loads the operator snapshot at that height and
// calculates the amount to slash.
// The amount that is liable to be slashed is the amount
// that was delegated to the operator at the infraction height.
// It will be burnt from either the current amount that is pending
// undelegation first (since we keep it pending for exactly situations
// like these), and then the delegation (which will happen if there
// were no undelegations). If the operator is slashed multiple times
// for the same infraction height, it may happen that both the
// amounts run out.

func (k Keeper) UpdateAndStoreSnapshotForDelegation(
	ctx sdk.Context, operator sdk.AccAddress, assetId string,
	stakerId string, changeAmount sdkmath.Int,
) {
	snapshot, _ := k.getOperatorLastSnapshot(ctx, operator)
	snapshot.UpdateForDelegation(assetId, stakerId, changeAmount)
	k.setOperatorSnapshot(ctx, operator, snapshot)
}

func (k Keeper) UpdateAndStoreSnapshotForUndelegation(
	ctx sdk.Context, operator sdk.AccAddress, assetId string,
	stakerId string, changeAmount sdkmath.Int,
) {
	snapshot, found := k.getOperatorLastSnapshot(ctx, operator)
	if !found {
		panic("snapshot not found, cnanot update for undelegation")
	}
	snapshot.UpdateForUndelegation(assetId, stakerId, changeAmount)
	k.setOperatorSnapshot(ctx, operator, snapshot)
}
