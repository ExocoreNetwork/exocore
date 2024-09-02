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
	order channeltypes.Order,
	connectionHops []string,
	portId string,
	channelId string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
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
	// no k.cdc here so use this
	mdBz, err := (&md).Marshal()
	if err != nil {
		return "", errorsmod.Wrapf(commontypes.ErrInvalidHandshakeMetadata,
			"error marshalling ibc-try metadata: %v", err)
	}
	return string(mdBz), nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portId,
	channelId string,
	_,
	counterpartyVersion string,
) error {
	im.keeper.Logger(ctx).Error(
		"OnChanOpenAck",
	)
	return errorsmod.Wrap(
		commontypes.ErrInvalidChannelFlow,
		"channel handshake must be initiated by subscriber chain",
	)
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portId string,
	channelId string,
) error {
	im.keeper.Logger(ctx).Error(
		"OnChanOpenConfirm",
	)
	err := im.keeper.SetSubscriberChain(ctx, channelId)
	if err != nil {
		return err
	}
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portId string,
	channelId string,
) error {
	im.keeper.Logger(ctx).Error(
		"OnChanCloseInit",
	)
	// Disallow user-initiated channel closing for channels
	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portId,
	channelId string,
) error {
	im.keeper.Logger(ctx).Error(
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
	im.keeper.Logger(ctx).Error(
		"OnRecvPacket",
	)
	var (
		ack  ibcexported.Acknowledgement
		data commontypes.SubscriberPacketData
	)

	if err := commontypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		errAck := utils.NewErrorAcknowledgementWithLog(ctx, err)
		ack = &errAck
	} else {
		switch data.Type {
		case commontypes.SlashPacket:
			im.keeper.Logger(ctx).Error(
				"OnRecvSlashPacket",
				"packet data", data,
			)
			ack = im.keeper.OnRecvSlashPacket(ctx, packet, *data.GetSlashPacketData())
		case commontypes.VscMaturedPacket:
			im.keeper.Logger(ctx).Error(
				"OnRecvVscMaturedPacket",
				"packet data", data,
			)
			ack = im.keeper.OnRecvVscMaturedPacket(ctx, packet, *data.GetVscMaturedPacketData())
		default:
			errAck := utils.NewErrorAcknowledgementWithLog(ctx, fmt.Errorf("unknown packet type: %s", data.Type))
			ack = &errAck
		}
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
	im.keeper.Logger(ctx).Error(
		"OnAcknowledgementPacket",
	)
	var ack channeltypes.Acknowledgement
	if err := commontypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrUnknownRequest,
			"cannot unmarshal coordinator packet acknowledgement: %v",
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
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) error {
	im.keeper.Logger(ctx).Error(
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
