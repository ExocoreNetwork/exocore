package subscriber

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/keeper"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// IBCModule is the IBC module for the subscriber module.
type IBCModule struct {
	keeper keeper.Keeper
}

// interface guard
var _ porttypes.IBCModule = IBCModule{}

// NewIBCModule creates a new IBCModule instance
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface for the subscriber module.
// The function is called when the channel is created, typically by the relayer,
// which must be informed that the channel should be created on this chain.
// Starting the channel on the coordinator chain is not supported.
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenInit",
	)

	// ICS26 requires that it is set to the default version if empty
	if strings.TrimSpace(version) == "" {
		version = commontypes.Version
	}
	// check if channel has already been initialized
	if storedChannel, ok := im.keeper.GetCoordinatorChannel(ctx); ok {
		return "", errorsmod.Wrapf(commontypes.ErrDuplicateChannel,
			"channel already exists with ID %s", storedChannel)
	}

	// check channel params (subscriber end)
	if order != channeltypes.ORDERED {
		return "", errorsmod.Wrapf(
			channeltypes.ErrInvalidChannelOrdering,
			"expected %s channel, got %s ",
			channeltypes.ORDERED,
			order,
		)
	}
	// we set our port at genesis. check that the port of the channel is the same
	if boundPort := im.keeper.GetPort(ctx); portID != boundPort {
		return "", errorsmod.Wrapf(
			porttypes.ErrInvalidPort,
			"invalid port ID: %s, expected: %s",
			portID, boundPort,
		)
	}
	// check that the version is correct
	if version != commontypes.Version {
		return "", errorsmod.Wrapf(
			commontypes.ErrInvalidVersion,
			"invalid version: %s, expected: %s",
			version, commontypes.Version,
		)
	}
	// check channel params (coordinator end)
	if counterparty.PortId != commontypes.CoordinatorPortID {
		return "", errorsmod.Wrapf(
			porttypes.ErrInvalidPort,
			"invalid counterparty port ID: %s, expected: %s",
			counterparty.PortId, commontypes.CoordinatorPortID,
		)
	}

	// claim channel capability passed back by IBC module
	if err := im.keeper.ClaimCapability(
		ctx, chanCap,
		host.ChannelCapabilityPath(portID, channelID),
	); err != nil {
		return "", err
	}

	// check connection hops, connection, and the client id (set on genesis)
	if err := im.keeper.VerifyCoordinatorChain(ctx, connectionHops); err != nil {
		return "", err
	}

	return commontypes.Version, nil
}

// OnChanOpenTry implements the IBCModule interface. It rejects attempts by
// the counterparty chain to open a channel here, since our spec requires
// that the channel is opened by this chain.
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	_ channeltypes.Order,
	_ []string,
	_ string,
	_ string,
	_ *capabilitytypes.Capability,
	_ channeltypes.Counterparty,
	_ string,
) (string, error) {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenTry",
	)
	return "", errorsmod.Wrap(
		commontypes.ErrInvalidChannelFlow,
		"channel handshake must be initiated by subscriber chain",
	)
}

// OnChanOpenAck implements the IBCModule interface. It is ran after `OnChanOpenTry`
// is run on the counterparty chain.
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	_ string, // unused as per spec
	counterpartyMetadata string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenAck",
	)

	// ensure coordinator channel has not already been created
	if coordinatorChannel, ok := im.keeper.GetCoordinatorChannel(ctx); ok {
		return errorsmod.Wrapf(commontypes.ErrDuplicateChannel,
			"coordinator channel: %s already established", coordinatorChannel)
	}

	var md commontypes.HandshakeMetadata
	// no k.cdc.MustUnmarshal available here, so we use this way.
	if err := (&md).Unmarshal([]byte(counterpartyMetadata)); err != nil {
		return errorsmod.Wrapf(
			commontypes.ErrInvalidHandshakeMetadata,
			"error unmarshalling ibc-ack metadata: \n%v", err,
		)
	}

	if md.Version != commontypes.Version {
		return errorsmod.Wrapf(
			commontypes.ErrInvalidVersion,
			"invalid counterparty version: %s, expected %s",
			md.Version,
			commontypes.Version,
		)
	}

	// This address is not required to be supplied at the time of chain registration.
	// Rather, it is set later by the coordinator chain.
	im.keeper.SetCoordinatorFeePoolAddrStr(ctx, md.CoordinatorFeePoolAddr)

	///////////////////////////////////////////////////
	// Initialize distribution token transfer channel

	// First check if an existing transfer channel already exists.
	transChannelID := im.keeper.GetDistributionTransmissionChannel(ctx)
	if found := im.keeper.TransferChannelExists(ctx, transChannelID); found {
		return nil
	}

	// NOTE The handshake for this channel is handled by the ibc-go/transfer
	// module. If the transfer-channel fails here (unlikely) then the transfer
	// channel should be manually created and parameters set accordingly.

	// reuse the connection hops for this channel for the
	// transfer channel being created.
	connHops, err := im.keeper.GetConnectionHops(ctx, portID, channelID)
	if err != nil {
		return err
	}

	distrTransferMsg := channeltypes.NewMsgChannelOpenInit(
		transfertypes.PortID,
		transfertypes.Version,
		channeltypes.UNORDERED,
		connHops,
		transfertypes.PortID,
		"", // signer unused
	)

	resp, err := im.keeper.ChannelOpenInit(ctx, distrTransferMsg)
	if err != nil {
		return err
	}
	im.keeper.SetDistributionTransmissionChannel(ctx, resp.ChannelId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeFeeTransferChannelOpened,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(channeltypes.AttributeKeyChannelID, channelID),
			sdk.NewAttribute(channeltypes.AttributeKeyPortID, types.PortID),
		),
	)

	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	_ string,
	_ string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenConfirm",
	)

	return errorsmod.Wrap(
		commontypes.ErrInvalidChannelFlow,
		"channel handshake must be initiated by subscriber chain",
	)
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	_ string,
	channelID string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanCloseInit",
	)

	// allow relayers to close duplicate OPEN channels, if the coordinator channel has already
	// been established
	if coordinatorChannel, ok := im.keeper.GetCoordinatorChannel(ctx); ok &&
		coordinatorChannel != channelID {
		return nil
	}
	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	_ string,
	_ string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanCloseConfirm",
	)
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	im.keeper.Logger(ctx).Debug(
		"OnRecvPacket",
	)
	var (
		ack  ibcexported.Acknowledgement
		data commontypes.ValidatorSetChangePacketData
	)
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		errAck := commontypes.NewErrorAcknowledgementWithLog(
			ctx, fmt.Errorf("cannot unmarshal packet data"),
		)
		ack = &errAck
	} else {
		im.keeper.Logger(ctx).Debug(
			"OnRecvPacket",
			"packet data", data,
		)
		ack = im.keeper.OnRecvVSCPacket(ctx, packet, data)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			commontypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(commontypes.AttributeKeyAckSuccess, fmt.Sprintf("%t", ack != nil)),
		),
	)

	return ack
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	_ sdk.AccAddress,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnAcknowledgementPacket",
	)
	var ack channeltypes.Acknowledgement
	if err := commontypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal subscriber packet acknowledgement: %v",
			err,
		)
	}

	if err := im.keeper.OnAcknowledgementPacket(ctx, packet, ack); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			commontypes.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(commontypes.AttributeKeyAck, ack.String()),
		),
	)
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				commontypes.EventTypePacket,
				sdk.NewAttribute(commontypes.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				commontypes.EventTypePacket,
				sdk.NewAttribute(commontypes.AttributeKeyAckError, resp.Error),
			),
		)
	}
	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	_ channeltypes.Packet,
	_ sdk.AccAddress,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnTimeoutPacket",
	)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			commontypes.EventTypeTimeout,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return nil
}
