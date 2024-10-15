package keeper

import (
	"encoding/binary"

	errorsmod "cosmossdk.io/errors"
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	subscribertypes "github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibchost "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

// VerifySubscriberChain verifies the chain trying to connect on the channel handshake.
// The verification includes the number of connection hops, the presence of a light client,
// as well as the channel's presence.
func (k Keeper) VerifySubscriberChain(
	ctx sdk.Context,
	_ string,
	connectionHops []string,
) error {
	if len(connectionHops) != 1 {
		return errorsmod.Wrap(
			channeltypes.ErrTooManyConnectionHops,
			"must have direct connection to coordinator chain",
		)
	}
	connectionID := connectionHops[0]
	clientID, tmClient, err := k.getUnderlyingClient(ctx, connectionID)
	if err != nil {
		return err
	}
	storedClientID, found := k.GetClientForChain(ctx, tmClient.ChainId)
	if !found {
		return errorsmod.Wrapf(
			commontypes.ErrClientNotFound,
			"cannot find client for subscriber chain %s",
			tmClient.ChainId,
		)
	}
	if storedClientID != clientID {
		return errorsmod.Wrapf(
			types.ErrInvalidSubscriberClient,
			"channel must be built on top of client. expected %s, got %s",
			storedClientID, clientID,
		)
	}

	// Verify that there isn't already a stored channel
	if prevChannel, ok := k.GetChannelForChain(ctx, tmClient.ChainId); ok {
		return errorsmod.Wrapf(
			commontypes.ErrDuplicateChannel,
			"channel with ID: %s already created for subscriber chain %s",
			prevChannel, tmClient.ChainId,
		)
	}
	return nil
}

// getUnderlyingClient gets the client state of the subscriber chain,
// as deployed on the coordinator chain.
func (k Keeper) getUnderlyingClient(ctx sdk.Context, connectionID string) (
	clientID string, tmClient *ibctmtypes.ClientState, err error,
) {
	conn, ok := k.connectionKeeper.GetConnection(ctx, connectionID)
	if !ok {
		return "", nil, errorsmod.Wrapf(conntypes.ErrConnectionNotFound,
			"connection not found for connection ID: %s", connectionID)
	}
	clientID = conn.ClientId
	clientState, ok := k.clientKeeper.GetClientState(ctx, clientID)
	if !ok {
		return "", nil, errorsmod.Wrapf(clienttypes.ErrClientNotFound,
			"client not found for client ID: %s", clientID)
	}
	tmClient, ok = clientState.(*ibctmtypes.ClientState)
	if !ok {
		return "", nil, errorsmod.Wrapf(
			clienttypes.ErrInvalidClientType,
			"invalid client type. expected %s, got %s",
			ibchost.Tendermint,
			clientState.ClientType(),
		)
	}
	return clientID, tmClient, nil
}

// SetSubscriberChain sets the subscriber chain for the given channel ID.
// It is called when the connection handshake is complete.
func (k Keeper) SetSubscriberChain(ctx sdk.Context, channelID string) error {
	channel, ok := k.channelKeeper.GetChannel(ctx, commontypes.CoordinatorPortID, channelID)
	if !ok {
		return errorsmod.Wrapf(
			channeltypes.ErrChannelNotFound,
			"channel not found for channel ID: %s", channelID,
		)
	}
	if len(channel.ConnectionHops) != 1 {
		return errorsmod.Wrap(
			channeltypes.ErrTooManyConnectionHops,
			"must have direct connection to subscriber chain",
		)
	}
	connectionID := channel.ConnectionHops[0]
	clientID, tmClient, err := k.getUnderlyingClient(ctx, connectionID)
	if err != nil {
		return err
	}
	// Verify that there isn't already a channel for the subscriber chain
	chainID := tmClient.ChainId
	if prevChannelID, ok := k.GetChannelForChain(ctx, chainID); ok {
		return errorsmod.Wrapf(
			commontypes.ErrDuplicateChannel,
			"channel with Id: %s already created for subscriber chain %s",
			prevChannelID, chainID,
		)
	}

	// the channel is established:
	// - set channel mappings
	k.SetChannelForChain(ctx, chainID, channelID)
	k.SetChainForChannel(ctx, channelID, chainID)
	// - set current block height for the subscriber chain initialization
	k.SetInitChainHeight(ctx, chainID, uint64(ctx.BlockHeight()))
	// remove init timeout timestamp
	timeout, exists := k.GetChainInitTimeout(ctx, chainID)
	if exists {
		k.DeleteChainInitTimeout(ctx, chainID)
		k.RemoveChainFromInitTimeout(ctx, timeout, chainID)
	} else {
		k.Logger(ctx).Error("timeout not found for chain", "chainID", chainID)
	}

	// emit event on successful addition
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			commontypes.EventTypeChannelEstablished,
			sdk.NewAttribute(sdk.AttributeKeyModule, subscribertypes.ModuleName),
			sdk.NewAttribute(commontypes.AttributeChainID, chainID),
			sdk.NewAttribute(conntypes.AttributeKeyClientID, clientID),
			sdk.NewAttribute(channeltypes.AttributeKeyChannelID, channelID),
			sdk.NewAttribute(conntypes.AttributeKeyConnectionID, connectionID),
		),
	)
	return nil
}

// SetInitChainHeight sets the Exocore block height when the given app chain was initiated
func (k Keeper) SetInitChainHeight(ctx sdk.Context, chainID string, height uint64) {
	store := ctx.KVStore(k.storeKey)
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	store.Set(types.InitChainHeightKey(chainID), heightBytes)
}

// GetInitChainHeight returns the Exocore block height when the given app chain was initiated
func (k Keeper) GetInitChainHeight(ctx sdk.Context, chainID string) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InitChainHeightKey(chainID))
	if bz == nil {
		return 0, false
	}

	return binary.BigEndian.Uint64(bz), true
}
