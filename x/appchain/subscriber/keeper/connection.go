package keeper

import (
	errorsmod "cosmossdk.io/errors"
	types "github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

// SetCoordinatorClientID sets the clientID of the coordinator chain
func (k Keeper) SetCoordinatorClientID(ctx sdk.Context, clientID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoordinatorClientIDKey(), []byte(clientID))
}

// GetCoordinatorClientID gets the clientID of the coordinator chain
func (k Keeper) GetCoordinatorClientID(ctx sdk.Context) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.CoordinatorClientIDKey()
	if !store.Has(key) {
		return "", false
	}
	bz := store.Get(key)
	return string(bz), true
}

// SetCoordinatorChannel sets the channelId for the channel to the coordinator.
func (k Keeper) SetCoordinatorChannel(ctx sdk.Context, channelId string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoordinatorChannelKey(), []byte(channelId))
}

// GetCoordinatorChannel gets the channelId for the channel to the coordinator.
func (k Keeper) GetCoordinatorChannel(ctx sdk.Context) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	channelIdBytes := store.Get(types.CoordinatorChannelKey())
	if len(channelIdBytes) == 0 {
		return "", false
	}
	return string(channelIdBytes), true
}

// DeleteCoordinatorChannel deletes the channelId for the channel to the coordinator.
func (k Keeper) DeleteCoordinatorChannel(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.CoordinatorChannelKey())
}

// VerifyCoordinatorChain verifies the chain trying to connect on the channel handshake.
func (k Keeper) VerifyCoordinatorChain(ctx sdk.Context, connectionHops []string) error {
	if len(connectionHops) != 1 {
		return errorsmod.Wrap(
			channeltypes.ErrTooManyConnectionHops,
			"must have direct connection to coordinator chain",
		)
	}
	connectionId := connectionHops[0]
	conn, ok := k.connectionKeeper.GetConnection(ctx, connectionId)
	if !ok {
		return errorsmod.Wrapf(
			conntypes.ErrConnectionNotFound,
			"connection not found for connection Id: %s",
			connectionId,
		)
	}
	// Verify that client id is expected clientId
	expectedClientId, ok := k.GetCoordinatorClientID(ctx)
	if !ok {
		return errorsmod.Wrapf(
			clienttypes.ErrInvalidClient,
			"could not find coordinator client id",
		)
	}
	if expectedClientId != conn.ClientId {
		return errorsmod.Wrapf(
			clienttypes.ErrInvalidClient,
			"invalid client: %s, channel must be built on top of client: %s",
			conn.ClientId,
			expectedClientId,
		)
	}

	return nil
}
