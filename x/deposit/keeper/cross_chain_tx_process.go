package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Deposit(reStakerId string, assetsInfo map[string]sdkmath.Uint) error {
	//TODO implement me
	panic("implement me")
}
