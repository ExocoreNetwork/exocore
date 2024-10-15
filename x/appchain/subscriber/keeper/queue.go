package keeper

import (
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// getAndIncrementPendingPacketsIdx returns the current pending packets index and increments it.
func (k Keeper) getAndIncrementPendingPacketsIdx(ctx sdk.Context) (toReturn uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingPacketsIndexKey())
	if bz != nil {
		toReturn = sdk.BigEndianToUint64(bz)
	}
	toStore := toReturn + 1
	store.Set(types.PendingPacketsIndexKey(), sdk.Uint64ToBigEndian(toStore))
	return toReturn
}

// AppendPendingPacket appends a packet to the pending packets queue, indexed by the current index.
func (k Keeper) AppendPendingPacket(
	ctx sdk.Context,
	packetType commontypes.SubscriberPacketDataType,
	packet commontypes.WrappedSubscriberPacketData,
) {
	store := ctx.KVStore(k.storeKey)
	// it is appended to a key with idx value, and not an overall array
	idx := k.getAndIncrementPendingPacketsIdx(ctx)
	key := types.PendingDataPacketsKey(idx)
	wrapped := commontypes.NewSubscriberPacketData(packetType, packet)
	bz := k.cdc.MustMarshal(&wrapped)
	store.Set(key, bz)
}

// GetPendingPackets returns ALL the pending packets from the store without indexes.
func (k Keeper) GetPendingPackets(ctx sdk.Context) []commontypes.SubscriberPacketData {
	ppWithIndexes := k.GetAllPendingPacketsWithIdx(ctx)
	ppList := make([]commontypes.SubscriberPacketData, 0)
	for _, ppWithIndex := range ppWithIndexes {
		ppList = append(ppList, ppWithIndex.SubscriberPacketData)
	}
	return ppList
}

// GetAllPendingPacketsWithIdx returns ALL pending packet data from the store
// with indexes relevant to the pending packets queue.
func (k Keeper) GetAllPendingPacketsWithIdx(ctx sdk.Context) []types.SubscriberPacketDataWithIdx {
	packets := []types.SubscriberPacketDataWithIdx{}
	store := ctx.KVStore(k.storeKey)
	// Note: PendingDataPacketsBytePrefix is the correct prefix, NOT PendingDataPacketsByteKey.
	// See consistency with PendingDataPacketsKey().
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.PendingDataPacketsBytePrefix})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var packet commontypes.SubscriberPacketData
		bz := iterator.Value()
		k.cdc.MustUnmarshal(bz, &packet)
		packetWithIdx := types.SubscriberPacketDataWithIdx{
			SubscriberPacketData: packet,
			// index stored in key after prefix, see PendingDataPacketsKey()
			Idx: sdk.BigEndianToUint64(iterator.Key()[1:]),
		}
		packets = append(packets, packetWithIdx)
	}
	return packets
}

// DeletePendingDataPackets deletes pending data packets with given indexes
func (k Keeper) DeletePendingDataPackets(ctx sdk.Context, idxs ...uint64) {
	store := ctx.KVStore(k.storeKey)
	for _, idx := range idxs {
		store.Delete(types.PendingDataPacketsKey(idx))
	}
}

// DeleteAllPendingDataPackets deletes all pending data packets from the store.
func (k Keeper) DeleteAllPendingDataPackets(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// Note: PendingDataPacketsBytePrefix is the correct prefix, NOT PendingDataPacketsByteKey.
	// See consistency with PendingDataPacketsKey().
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.PendingDataPacketsBytePrefix})
	keysToDel := [][]byte{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keysToDel = append(keysToDel, iterator.Key())
	}
	for _, key := range keysToDel {
		store.Delete(key)
	}
}
