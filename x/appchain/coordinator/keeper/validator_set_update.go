package keeper

import (
	"fmt"

	exocoretypes "github.com/ExocoreNetwork/exocore/types/keys"
	"github.com/ExocoreNetwork/exocore/utils"
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

// QueueValidatorUpdatesForEpochID queues all the validator updates to be sent to the subscriber
// chains at the end of the epoch. After this function, call SendQueuedValidatorUpdates, which
// will actually send the updates.
func (k Keeper) QueueValidatorUpdatesForEpochID(
	ctx sdk.Context, epochID string, epochNumber int64,
) {
	// Get all the chains that need to be updated
	chainIDs := k.avsKeeper.GetEpochEndChainIDs(ctx, epochID, epochNumber)
	for _, chainID := range chainIDs {
		cctx, writeCache, err := k.QueueValidatorUpdatesForChainIDInCachedCtx(ctx, chainID)
		if err != nil {
			k.Logger(ctx).Error(
				"error queuing validator updates for chain",
				"chainID", chainID,
				"error", err,
			)
			continue
		}
		// copy over the events from the cached ctx
		ctx.EventManager().EmitEvents(cctx.EventManager().Events())
		writeCache()
	}
}

// QueueValidatorUpdatesForChainIDInCachedCtx is a wrapper function around QueueValidatorUpdatesForChainID.
func (k Keeper) QueueValidatorUpdatesForChainIDInCachedCtx(
	ctx sdk.Context, chainID string,
) (cctx sdk.Context, writeCache func(), err error) {
	cctx, writeCache = ctx.CacheContext()
	err = k.QueueValidatorUpdatesForChainID(cctx, chainID)
	return
}

// QueueValidatorUpdatesForChainID queues all the validator updates to be sent to the subscriber, saving the
// updates as individual validators as well.
func (k Keeper) QueueValidatorUpdatesForChainID(
	ctx sdk.Context, chainID string,
) error {
	// Get the current validator set for the chain, which is sorted
	// by the consensus address (bytes). This sorting is okay to use
	prevList := k.GetAllSubscriberValidatorsForChain(ctx, chainID)
	// to check whether the new set has a changed vote power, convert to map.
	prevMap := make(map[string]int64, len(prevList))
	for _, val := range prevList {
		// we are okay to use ConsAddress here even though the bech32 prefix
		// is different, because we don't print the address.
		prevMap[sdk.ConsAddress(val.ConsAddress).String()] = val.Power
	}
	operators, keys := k.operatorKeeper.GetActiveOperatorsForChainID(ctx, chainID)
	powers, err := k.operatorKeeper.GetVotePowerForChainID(
		ctx, operators, chainID,
	)
	if err != nil {
		k.Logger(ctx).Error(
			"error getting vote power for chain",
			"chainID", chainID,
			"error", err,
		)
		// skip this chain, if consecutive failures are reported, it will eventually be
		// timed out and then dropped.
		return err
	}
	operators, keys, powers = utils.SortByPower(operators, keys, powers)
	maxVals := k.GetMaxValidatorsForChain(ctx, chainID)
	// double the capacity assuming that all validators are removed and an entirely new
	// set of validators is added.
	validatorUpdates := make([]abci.ValidatorUpdate, 0, maxVals*2)
	for i := range operators {
		if i >= int(maxVals) {
			break
		}
		power := powers[i]
		if power < 1 {
			break
		}
		wrappedKey := keys[i]
		addressString := wrappedKey.ToConsAddr().String()
		prevPower, found := prevMap[addressString]
		if found {
			if prevPower != power {
				validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{
					PubKey: *wrappedKey.ToTmProtoKey(),
					Power:  power,
				})
			}
			delete(prevMap, addressString)
			validator, err := commontypes.NewSubscriberValidator(
				wrappedKey.ToConsAddr(), power, wrappedKey.ToSdkKey(),
			)
			if err != nil {
				// should never happen, but just in case.
				// don't skip the chain though, instead, skip the validator.
				continue
			}
			k.SetSubscriberValidatorForChain(ctx, chainID, validator)
		} else {
			// new key, add it to the list.
			validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{
				PubKey: *wrappedKey.ToTmProtoKey(),
				Power:  power,
			})
			validator, err := commontypes.NewSubscriberValidator(
				wrappedKey.ToConsAddr(), power, wrappedKey.ToSdkKey(),
			)
			if err != nil {
				// should never happen, but just in case.
				// don't skip the chain though, instead, skip the validator.
				continue
			}
			k.SetSubscriberValidatorForChain(ctx, chainID, validator)
		}
	}
	// if there is any element in the prevList, which is still in prevMap, that element
	// needs to have a vote power of 0 queued.
	for _, validator := range prevList {
		pubKey, err := validator.ConsPubKey()
		if err != nil {
			k.Logger(ctx).Error(
				"error deserializing consensus public key",
				"chainID", chainID,
				"error", err,
			)
			return err
		}
		wrappedKey := exocoretypes.NewWrappedConsKeyFromSdkKey(pubKey)
		// alternatively, the below could be replaced by wrappedKey.ToConsAddr(), but
		// since we generated this address when saving it, we can use it directly.
		consAddress := sdk.ConsAddress(validator.ConsAddress)
		if _, found := prevMap[consAddress.String()]; found {
			validatorUpdates = append(validatorUpdates, abci.ValidatorUpdate{
				PubKey: *wrappedKey.ToTmProtoKey(),
				Power:  0,
			})
			k.DeleteSubscriberValidatorForChain(ctx, chainID, consAddress)
		}
	}
	// default is 0 for the subscriber genesis. any updates will start with 1.
	// increment gets the value of 0, increments it to 1, stores it and returns it.
	vscID := k.IncrementVscIDForChain(ctx, chainID)
	data := commontypes.NewVscPacketData(
		validatorUpdates, vscID, k.ConsumeSlashAcks(ctx, chainID),
	)
	k.AppendPendingVscPacket(ctx, chainID, data)
	return nil
}

// SendQueuedValidatorUpdates sends the queued validator set updates to the subscriber chains.
// It only sends them if a client + channel for that chain are set up. Otherwise, no action
// is taken. Since it is called immediately after queuing the updates, it is guaranteed that
// only the updates from the queue (or prior) are sent. In other words, there is no possibility
// for updates from a different epoch will be sent. Hence, we simply iterate over all (active)
// chains.
func (k Keeper) SendQueuedValidatorUpdates(ctx sdk.Context, epochNumber int64) {
	chainIDs := k.GetAllChainsWithChannels(ctx)
	for _, chainID := range chainIDs {
		// a channel is guaranteed to exist.
		channelID, _ := k.GetChannelForChain(ctx, chainID)
		packets := k.GetPendingVscPackets(ctx, chainID)
		k.SendVscPacketsToChain(ctx, chainID, channelID, packets.List, epochNumber)
	}
}

// SendVscPacketsToChain sends the validator set change packets to the subscriber chain.
func (k Keeper) SendVscPacketsToChain(
	ctx sdk.Context, chainID string, channelID string,
	packets []commontypes.ValidatorSetChangePacketData,
	epochNumber int64,
) {
	params := k.GetParams(ctx)
	for i := range packets {
		data := packets[i]
		// send packet over IBC
		err := commontypes.SendIBCPacket(
			ctx,
			k.scopedKeeper,
			k.channelKeeper,
			channelID,                     // source channel id
			commontypes.CoordinatorPortID, // source port id
			commontypes.ModuleCdc.MustMarshalJSON(&data), // packet data
			params.IBCTimeoutPeriod,
		)
		if err != nil {
			if clienttypes.ErrClientNotActive.Is(err) {
				// IBC client is expired!
				// leave the packet data stored to be sent once the client is upgraded
				// the client cannot expire during iteration (in the middle of a block)
				k.Logger(ctx).Info(
					"IBC client is expired, cannot send VSC, leaving packet data stored:",
					"chainID", chainID,
					"vscID", data.ValsetUpdateID,
				)
				return
			}
			// Not able to send packet over IBC!
			k.Logger(ctx).Error(
				"cannot send VSC, removing subscriber",
				"chainID", chainID,
				"vscID", data.ValsetUpdateID,
				"err", err.Error(),
			)
			// If this happens, most likely the subscriber is malicious; remove it
			err := k.StopSubscriberChain(ctx, chainID, true)
			if err != nil {
				panic(fmt.Errorf("subscriber chain failed to stop: %w", err))
			}
			return
		}
		// even when the epoch identifier is `minute` and that of the `timeoutPeriod` is hour
		// the latter is used. this is because the `timeout` runs on a different schedule.
		timeoutPeriod := params.VSCTimeoutPeriod
		timeoutPeriod.EpochNumber += uint64(epochNumber) + 1 // 1 extra for the ended epoch
		k.SetVscTimeout(ctx, chainID, data.ValsetUpdateID, timeoutPeriod)
	}
	k.SetPendingVscPackets(ctx, chainID, types.ValidatorSetChangePackets{})
}

// AppendPendingVscPacket appends a validator set change packet to the pending list, indexed by the chainID.
func (k Keeper) AppendPendingVscPacket(ctx sdk.Context, chainID string, data commontypes.ValidatorSetChangePacketData) {
	prev := k.GetPendingVscPackets(ctx, chainID)
	prev.List = append(prev.List, data)
	k.SetPendingVscPackets(ctx, chainID, prev)
}

// GetPendingVscPackets gets the pending validator set change packets for a chain.
func (k Keeper) GetPendingVscPackets(ctx sdk.Context, chainID string) types.ValidatorSetChangePackets {
	store := ctx.KVStore(k.storeKey)
	var data types.ValidatorSetChangePackets
	key := types.ChainIDToVscPacketsKey(chainID)
	value := store.Get(key)
	k.cdc.MustUnmarshal(value, &data)
	return data
}

// SetPendingVscPackets sets the pending validator set change packets for a chain.
func (k Keeper) SetPendingVscPackets(ctx sdk.Context, chainID string, data types.ValidatorSetChangePackets) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChainIDToVscPacketsKey(chainID)
	if len(data.List) == 0 {
		store.Delete(key)
	} else {
		store.Set(key, k.cdc.MustMarshal(&data))
	}
}
