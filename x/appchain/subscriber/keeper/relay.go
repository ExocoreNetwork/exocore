package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/utils"
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// OnRecvVSCPacket processes a validator set change packet
func (k Keeper) OnRecvVSCPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data commontypes.ValidatorSetChangePacketData,
) exported.Acknowledgement {
	coordinatorChannel, found := k.GetCoordinatorChannel(ctx)
	if found && packet.SourceChannel != coordinatorChannel {
		// should never happen
		k.Logger(ctx).Error(
			"received VSCPacket on non-coordinator channel",
			"source channel", packet.SourceChannel,
			"coordinator channel", coordinatorChannel,
		)
		return nil
	}
	if !found {
		// first message on channel
		k.SetCoordinatorChannel(ctx, packet.SourceChannel)
		k.Logger(ctx).Info(
			"channel established",
			"port", packet.DestinationPort,
			"channel", packet.DestinationChannel,
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				commontypes.EventTypeChannelEstablished,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
				sdk.NewAttribute(channeltypes.AttributeKeyChannelID, packet.DestinationChannel),
				sdk.NewAttribute(channeltypes.AttributeKeyPortID, packet.DestinationPort),
			),
		)
	}
	// the changes are received within blocks, but can only be forwarded to
	// Tendermint during EndBlock. hence, get the changes received so far, append to them
	// and save them back
	currentChanges, _ := k.GetPendingChanges(ctx)
	pendingChanges := utils.AccumulateChanges(
		currentChanges.ValidatorUpdates,
		data.ValidatorUpdates,
	)

	k.SetPendingChanges(ctx, &commontypes.ValidatorSetChangePacketData{
		ValidatorUpdates: pendingChanges,
	})

	// Save maturity time and packet
	maturityTime := ctx.BlockTime().Add(k.GetUnbondingPeriod(ctx))
	k.SetPacketMaturityTime(ctx, data.ValsetUpdateID, maturityTime)
	k.Logger(ctx).Info(
		"packet maturity time was set",
		"vscID", data.ValsetUpdateID,
		"maturity time (utc)", maturityTime.UTC(),
		"maturity time (nano)", uint64(maturityTime.UnixNano()),
	)

	// set height to VSC id mapping; it is effective as of the next block
	k.SetValsetUpdateIDForHeight(
		ctx, ctx.BlockHeight()+1, data.ValsetUpdateID,
	)
	k.Logger(ctx).Info(
		"block height was mapped to vscID",
		"height", ctx.BlockHeight()+1,
		"vscID", data.ValsetUpdateID,
	)

	// remove outstanding slashing flags of the validators
	// for which the slashing was acknowledged by the coordinator chain
	// TODO(mm): since this packet is only received when there are validator set changes
	// there is some additional lag between the slashing occurrence on the coordinator
	// and deletion of this flag on the subscriber. does it matter?
	for _, addr := range data.GetSlashAcks() {
		consAddr, err := sdk.ConsAddressFromBech32(addr)
		if err != nil {
			k.Logger(ctx).Error(
				"failed to parse consensus address",
				"address", addr, "error", err,
			)
			// returning an error will cause the coordinator to drop us
			continue
		}
		k.DeleteOutstandingDowntime(ctx, consAddr)
	}

	k.Logger(ctx).Info(
		"finished receiving/handling VSCPacket",
		"vscID", data.ValsetUpdateID,
		"len updates", len(data.ValidatorUpdates),
		"len slash acks", len(data.SlashAcks),
	)
	// Acknowledge the packet
	return channeltypes.NewResultAcknowledgement([]byte{byte(1)})
}

// OnAcknowledgementPacket processes an acknowledgement packet
func (k Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	ack channeltypes.Acknowledgement,
) error {
	// the ack can only be error when packet parsing is failed
	// or the packet type is wrong, or when the slash packet
	// has incorrect data. none of these should happen
	if err := ack.GetError(); err != "" {
		k.Logger(ctx).Error(
			"recv ErrorAcknowledgement",
			"channel", packet.SourceChannel,
			"error", err,
		)
		// Initiate ChanCloseInit using packet source (non-counterparty) port and channel
		err := k.ChanCloseInit(ctx, packet.SourcePort, packet.SourceChannel)
		if err != nil {
			return fmt.Errorf("ChanCloseInit(%s) failed: %s", packet.SourceChannel, err.Error())
		}
		// check if there is an established channel to coordinator
		channelID, found := k.GetCoordinatorChannel(ctx)
		if !found {
			return errorsmod.Wrapf(
				types.ErrNoProposerChannelID,
				"recv ErrorAcknowledgement on non-established channel %s",
				packet.SourceChannel,
			)
		}
		if channelID != packet.SourceChannel {
			// Close the established channel as well
			return k.ChanCloseInit(ctx, commontypes.SubscriberPortID, channelID)
		}
	}
	return nil
}
