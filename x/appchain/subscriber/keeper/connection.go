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
func (k Keeper) SetCoordinatorChannel(ctx sdk.Context, channelID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoordinatorChannelKey(), []byte(channelID))
}

// GetCoordinatorChannel gets the channelId for the channel to the coordinator.
func (k Keeper) GetCoordinatorChannel(ctx sdk.Context) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CoordinatorChannelKey())
	if len(bz) == 0 {
		return "", false
	}
	return string(bz), true
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
	connectionID := connectionHops[0]
	conn, ok := k.connectionKeeper.GetConnection(ctx, connectionID)
	if !ok {
		return errorsmod.Wrapf(
			conntypes.ErrConnectionNotFound,
			"connection not found for connection Id: %s",
			connectionID,
		)
	}
	// Verify that client id is expected clientId
	expectedClientID, ok := k.GetCoordinatorClientID(ctx)
	if !ok {
		return errorsmod.Wrapf(
			clienttypes.ErrInvalidClient,
			"could not find coordinator client id",
		)
	}
	if expectedClientID != conn.ClientId {
		return errorsmod.Wrapf(
			clienttypes.ErrInvalidClient,
			"invalid client: %s, channel must be built on top of client: %s",
			conn.ClientId,
			expectedClientID,
		)
	}

	return nil
}
