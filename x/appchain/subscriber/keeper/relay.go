package keeper

import (
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/utils"
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
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
	currentChanges := k.GetPendingChanges(ctx)
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
	for _, consAddr := range data.GetSlashAcks() {
		k.DeleteOutstandingDowntime(ctx, consAddr)
	}

	k.Logger(ctx).Info(
		"finished receiving/handling VSCPacket",
		"vscID", data.ValsetUpdateID,
		"len updates", len(data.ValidatorUpdates),
		"len slash acks", len(data.SlashAcks),
	)
	// Acknowledge the packet
	return commontypes.NewResultAcknowledgementWithLog(ctx, commontypes.VscPacketHandledResult)
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

// QueueVscMaturedPackets queues all VSC packets that have matured as of the current block time,
// to be sent to the coordinator chain at the end of the block.
func (k Keeper) QueueVscMaturedPackets(
	ctx sdk.Context,
) {
	for _, packet := range k.GetElapsedVscPackets(ctx) {
		vscPacket := commontypes.NewVscMaturedPacketData(packet.ID)
		k.AppendPendingPacket(
			ctx, commontypes.VscMaturedPacket,
			&commontypes.SubscriberPacketData_VscMaturedPacketData{
				VscMaturedPacketData: vscPacket,
			},
		)
		k.DeletePacketMaturityTime(ctx, packet.ID, packet.MaturityTime)

		k.Logger(ctx).Info("VSCMaturedPacket enqueued", "vscID", vscPacket.ValsetUpdateID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeVSCMatured,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
				sdk.NewAttribute(commontypes.AttributeChainID, ctx.ChainID()),
				sdk.NewAttribute(
					types.AttributeSubscriberHeight,
					strconv.Itoa(int(ctx.BlockHeight())),
				),
				sdk.NewAttribute(
					commontypes.AttributeValSetUpdateID,
					strconv.Itoa(int(packet.ID)),
				),
				sdk.NewAttribute(types.AttributeTimestamp, ctx.BlockTime().String()),
			),
		)
	}
}

// QueueSlashPacket queues a slashing request to be sent to the coordinator chain at the end of the block.
func (k Keeper) QueueSlashPacket(
	ctx sdk.Context,
	validator abci.Validator,
	valsetUpdateID uint64,
	infraction stakingtypes.Infraction,
) {
	consAddr := sdk.ConsAddress(validator.Address)
	downtime := infraction == stakingtypes.Infraction_INFRACTION_DOWNTIME

	// return if an outstanding downtime request is set for the validator
	if downtime && k.HasOutstandingDowntime(ctx, consAddr) {
		return
	}

	if downtime {
		// set outstanding downtime to not send multiple
		// slashing requests for the same downtime infraction
		k.SetOutstandingDowntime(ctx, consAddr)
	}

	// construct slash packet data
	slashPacket := commontypes.NewSlashPacketData(validator, valsetUpdateID, infraction)

	// append the Slash packet data to pending data packets
	k.AppendPendingPacket(
		ctx,
		commontypes.SlashPacket,
		&commontypes.SubscriberPacketData_SlashPacketData{
			SlashPacketData: slashPacket,
		},
	)

	k.Logger(ctx).Info(
		"SlashPacket enqueued",
		"vscID", slashPacket.ValsetUpdateID,
		"validator cons addr", fmt.Sprintf("%X", slashPacket.Validator.Address),
		"infraction", slashPacket.Infraction,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubscriberSlashRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(
				commontypes.AttributeValidatorAddress,
				fmt.Sprintf("%X", slashPacket.Validator.Address),
			),
			sdk.NewAttribute(
				commontypes.AttributeValSetUpdateID,
				strconv.Itoa(int(valsetUpdateID)),
			),
			sdk.NewAttribute(commontypes.AttributeInfractionType, infraction.String()),
		),
	)
}

// IsChannelClosed returns a boolean whether a given channel is in the CLOSED state
func (k Keeper) IsChannelClosed(ctx sdk.Context, channelID string) bool {
	channel, found := k.channelKeeper.GetChannel(ctx, commontypes.SubscriberPortID, channelID)
	return !found || channel.State == channeltypes.CLOSED
}

// SendPackets sends all pending packets to the coordinator chain
func (k Keeper) SendPackets(ctx sdk.Context) {
	// find destination
	channelID, ok := k.GetCoordinatorChannel(ctx)
	if !ok {
		return
	}
	// find packets, which will be returned sorted by index
	pending := k.GetAllPendingPacketsWithIdx(ctx)
	idxsForDeletion := []uint64{}
	timeoutPeriod := k.GetSubscriberParams(ctx).IBCTimeoutPeriod
	for _, p := range pending {
		// Send packet over IBC
		err := commontypes.SendIBCPacket(
			ctx,
			k.scopedKeeper,
			k.channelKeeper,
			channelID,                    // source channel id
			commontypes.SubscriberPortID, // source port id
			commontypes.ModuleCdc.MustMarshalJSON(&p.SubscriberPacketData),
			timeoutPeriod,
		)
		if err != nil {
			if clienttypes.ErrClientNotActive.Is(err) {
				// IBC client is expired!
				// leave the packet data stored to be sent once the client is upgraded
				k.Logger(ctx).Info(
					"IBC client is expired, cannot send IBC packet; leaving packet data stored:",
					"type", p.Type.String(),
				)
				return
			}
			// Not able to send packet over IBC!
			// Leave the packet data stored for the sent to be retried in the next block.
			// Note that if VSCMaturedPackets are not sent for long enough, the coordinator
			// will remove the subscriber anyway.
			k.Logger(ctx).Error(
				"cannot send IBC packet; leaving packet data stored:",
				"type", p.Type.String(), "err", err.Error(),
			)
			return
		} else {
			if p.Type == commontypes.VscMaturedPacket {
				id := p.GetVscMaturedPacketData().ValsetUpdateID
				k.Logger(ctx).Info(
					"IBC packet sent",
					"type", p.Type.String(),
					"id", id,
				)
			} else {
				data := p.GetSlashPacketData()
				addr := data.Validator.Address
				k.Logger(ctx).Info(
					"IBC packet sent",
					"type", p.Type.String(),
					"addr", addr,
				)
			}
		}
		// Otherwise the vsc matured will be deleted
		idxsForDeletion = append(idxsForDeletion, p.Idx)
	}
	// Delete pending packets that were successfully sent and did not return an error from SendIBCPacket
	k.DeletePendingDataPackets(ctx, idxsForDeletion...)
}
