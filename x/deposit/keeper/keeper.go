package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/assets/keeper"
	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	// other keepers
	assetsKeeper keeper.Keeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	assetsKeeper keeper.Keeper,
) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		assetsKeeper: assetsKeeper,
	}
}

// IDeposit interface will be implemented by deposit keeper
type IDeposit interface {
	// PostTxProcessing automatically call PostTxProcessing to update deposit state after receiving deposit event tx from layerZero protocol
	// PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// Deposit internal func for PostTxProcessing
	Deposit(ctx sdk.Context, event *DepositParams) error

	SetParams(ctx sdk.Context, params *deposittype.Params) error
	GetParams(ctx sdk.Context) (*deposittype.Params, error)
}
