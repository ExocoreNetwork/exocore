package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock : completed task create events according to the canCompleted blockHeight
func (k Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	ctx.Logger().Info("the blockHeight is:", "height", ctx.BlockHeight())
	task := CreateNewTaskParams{
		TaskCreatedBlock: ctx.BlockHeight(),
	}
	records, err := k.SetTaskforAvs(ctx, &task)
	if err != nil {
		panic(err)
	}
	if len(records) == 0 {
		return []abci.ValidatorUpdate{}
	}
	return []abci.ValidatorUpdate{}
}
