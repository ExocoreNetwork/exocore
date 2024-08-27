package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetClientForChain sets the ibc client id for a given chain id.
func (k Keeper) SetClientForChain(
	ctx sdk.Context, chainId string, clientId string,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ClientForChainKey(chainId), []byte(clientId))
}

// GetClientForChain gets the ibc client id for a given chain id.
func (k Keeper) GetClientForChain(
	ctx sdk.Context, chainId string,
) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.ClientForChainKey(chainId))
	if bytes == nil {
		return "", false
	}
	return string(bytes), true
}

// DeleteClientForChain deletes the ibc client id for a given chain id.
func (k Keeper) DeleteClientForChain(ctx sdk.Context, chainId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ClientForChainKey(chainId))
}
