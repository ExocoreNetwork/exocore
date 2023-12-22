package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetStakerExoCoreAddr outdated, will be deprecated.
func (k Keeper) GetStakerExoCoreAddr(ctx sdk.Context, stakerId string) (string, error) {
	// TODO implement me
	panic("implement me")
}
