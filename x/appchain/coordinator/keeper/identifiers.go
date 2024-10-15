package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetClientForChain sets the ibc client id for a given chain id.
func (k Keeper) SetClientForChain(
	ctx sdk.Context, chainID string, clientID string,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ClientForChainKey(chainID), []byte(clientID))
}

// GetClientForChain gets the ibc client id for a given chain id.
func (k Keeper) GetClientForChain(
	ctx sdk.Context, chainID string,
) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.ClientForChainKey(chainID))
	if bytes == nil {
		return "", false
	}
	return string(bytes), true
}

// DeleteClientForChain deletes the ibc client id for a given chain id.
func (k Keeper) DeleteClientForChain(ctx sdk.Context, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ClientForChainKey(chainID))
}

// GetAllChainsWithClients gets all chain ids that have an ibc client id.
func (k Keeper) GetAllChainsWithClients(ctx sdk.Context) []string {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ClientForChainBytePrefix})
	defer iterator.Close()

	var chains []string
	for ; iterator.Valid(); iterator.Next() {
		chainID := string(iterator.Key()[1:])
		chains = append(chains, chainID)
	}

	return chains
}

// SetChannelForChain sets the ibc channel id for a given chain id.
func (k Keeper) SetChannelForChain(ctx sdk.Context, chainID string, channelID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChannelForChainKey(chainID), []byte(channelID))
}

// GetChannelForChain gets the ibc channel id for a given chain id.
func (k Keeper) GetChannelForChain(ctx sdk.Context, chainID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.ChannelForChainKey(chainID))
	if bytes == nil {
		return "", false
	}
	return string(bytes), true
}

// GetAllChainsWithChannels gets all chain ids that have an ibc channel id, on top of the
// client id.
func (k Keeper) GetAllChainsWithChannels(ctx sdk.Context) []string {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ChannelForChainBytePrefix})
	defer iterator.Close()

	var chains []string
	for ; iterator.Valid(); iterator.Next() {
		chainID := string(iterator.Key()[1:])
		chains = append(chains, chainID)
	}

	return chains
}

// SetChainForChannel sets the chain id for a given channel id.
func (k Keeper) SetChainForChannel(ctx sdk.Context, channelID string, chainID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ChainForChannelKey(channelID), []byte(chainID))
}

// GetChainForChannel gets the chain id for a given channel id.
func (k Keeper) GetChainForChannel(ctx sdk.Context, channelID string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.ChainForChannelKey(channelID))
	if bytes == nil {
		return "", false
	}
	return string(bytes), true
}
