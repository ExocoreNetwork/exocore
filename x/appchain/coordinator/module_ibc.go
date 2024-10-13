package coordinator

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/keeper"
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
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

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	_ channeltypes.Order,
	_ []string,
	_ string,
	_ string,
	_ *capabilitytypes.Capability,
	_ channeltypes.Counterparty,
	version string,
) (string, error) {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenInit",
	)
	return version, errorsmod.Wrap(
		commontypes.ErrInvalidChannelFlow,
		"channel handshake must be initiated by subscriber chain",
	)
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenTry",
	)
	// channel ordering
	if order != channeltypes.ORDERED {
		return "", errorsmod.Wrapf(
			channeltypes.ErrInvalidChannelOrdering,
			"expected %s channel, got %s", channeltypes.ORDERED, order,
		)
	}

	// the channel's portId should match the module's
	boundPort := im.keeper.GetPort(ctx)
	if boundPort != portID {
		return "", errorsmod.Wrapf(
			porttypes.ErrInvalidPort,
			"invalid port: %s, expected %s", portID, boundPort,
		)
	}

	if counterpartyVersion != commontypes.Version {
		return "", errorsmod.Wrapf(
			commontypes.ErrInvalidVersion,
			"invalid counterparty version: got: %s, expected %s",
			counterpartyVersion,
			commontypes.Version,
		)
	}

	if counterparty.PortId != commontypes.SubscriberPortID {
		return "", errorsmod.Wrapf(
			porttypes.ErrInvalidPort,
			"invalid counterparty port Id: got %s, expected %s",
			counterparty.PortId,
			commontypes.SubscriberPortID,
		)
	}

	// Claim channel capability
	if err := im.keeper.ClaimCapability(
		ctx, chanCap, host.ChannelCapabilityPath(portID, channelID),
	); err != nil {
		return "", err
	}

	if err := im.keeper.VerifySubscriberChain(
		ctx, channelID, connectionHops,
	); err != nil {
		return "", err
	}

	md := commontypes.HandshakeMetadata{
		CoordinatorFeePoolAddr: im.keeper.GetSubscriberRewardsPoolAddressStr(ctx),
		Version:                commontypes.Version,
	}
	// we can use `MustMarshal` for data that we create
	mdBz := commontypes.ModuleCdc.MustMarshal(&md)
	return string(mdBz), nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	_ string,
	_ string,
	_ string,
	_ string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenAck",
	)
	return errorsmod.Wrap(
		commontypes.ErrInvalidChannelFlow,
		"channel handshake must be initiated by subscriber chain",
	)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context, _ string, dstChannelID string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanOpenConfirm",
	)
	err := im.keeper.SetSubscriberChain(ctx, dstChannelID)
	if err != nil {
		return err
	}
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context, _ string, _ string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanCloseInit",
	)
	// Disallow user-initiated channel closing for channels
	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context, _ string, _ string,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnChanCloseConfirm",
	)
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context, packet channeltypes.Packet, _ sdk.AccAddress,
) ibcexported.Acknowledgement {
	im.keeper.Logger(ctx).Debug(
		"OnRecvPacket",
	)

	var (
		ack  ibcexported.Acknowledgement
		data commontypes.SubscriberPacketData
		err  error
		res  []byte
	)

	// (1) Since this is a packet originating from the subscriber, we cannot use MustUnmarshal,
	// because such packets are not guaranteed to be correctly formed.
	// (2) When the subscriber chain marshals the data, it should use MarshalJSON.
	if unmarshalErr := commontypes.ModuleCdc.UnmarshalJSON(
		packet.GetData(), &data,
	); unmarshalErr != nil {
		im.keeper.Logger(ctx).Error(
			"cannot unmarshal subscriber packet data",
			"error", unmarshalErr,
		)
		err = sdkerrors.ErrInvalidType.Wrapf(
			"cannot unmarshal coordinator packet data: %s", unmarshalErr,
		)
	} else {
		switch data.Type {
		case commontypes.SlashPacket:
			im.keeper.Logger(ctx).Debug(
				"OnRecvSlashPacket",
				"packet data", data,
			)
			res, err = im.keeper.OnRecvSlashPacket(ctx, packet, *data.GetSlashPacketData())
		case commontypes.VscMaturedPacket:
			im.keeper.Logger(ctx).Debug(
				"OnRecvVscMaturedPacket",
				"packet data", data,
			)
			// no need to send an ack for this packet type
			err = im.keeper.OnRecvVscMaturedPacket(ctx, packet, *data.GetVscMaturedPacketData())
		default:
			err = sdkerrors.ErrInvalidType.Wrapf("unknown packet type: %s", data.Type)
		}
	}
	switch {
	case err != nil:
		ack = commontypes.NewErrorAcknowledgementWithLog(ctx, err)
	case res != nil:
		ack = commontypes.NewResultAcknowledgementWithLog(ctx, res)
	default:
		ack = commontypes.NewResultAcknowledgementWithLog(ctx, nil)
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
	// same as before, this packet is sent by the subscriber, so we cannot use MustUnmarshal
	var ack channeltypes.Acknowledgement
	if err := commontypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal packet acknowledgement: %s", err,
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
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) error {
	im.keeper.Logger(ctx).Debug(
		"OnTimeoutPacket",
	)
	if err := im.keeper.OnTimeoutPacket(ctx, packet); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			commontypes.EventTypeTimeout,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	return nil
}
