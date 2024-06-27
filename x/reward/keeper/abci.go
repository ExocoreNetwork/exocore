package keeper

import (
	"cosmossdk.io/api/tendermint/abci"
	"github.com/ExocoreNetwork/exocore/x/reward/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	avsAddrList, err := k.avsKeeper.GetEpochEndAVSs(ctx)
	if err != nil {
		panic(err)
	}
	if len(avsAddrList) == 0 {
		return []abci.ValidatorUpdate{}
	}

	pool := k.getPool(ctx, types.ModuleName)

	ForEach(avsAddrList, func(p string) { pool.ReleaseRewards(p) })
	return []abci.ValidatorUpdate{}
}

// ForEach apply the function on every element within the slice
func ForEach[T any](source []T, f func(T)) {
	for i := range source {
		f(source[i])
	}
}
