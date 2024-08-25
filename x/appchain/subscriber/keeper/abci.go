package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlock(sdk.Context) {}

func (k Keeper) EndBlock(sdk.Context) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
