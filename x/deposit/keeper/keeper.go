// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	types2 "github.com/exocore/x/deposit/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	//other keepers
	retakingStateKeeper keeper.Keeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	retakingStateKeeper keeper.Keeper,
) Keeper {
	return Keeper{
		storeKey:            storeKey,
		cdc:                 cdc,
		retakingStateKeeper: retakingStateKeeper,
	}
}

// IDeposit interface will be implemented by deposit keeper
type IDeposit interface {
	// PostTxProcessing automatically call PostTxProcessing to update deposit state after receiving deposit event tx from layerZero protocol
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// Deposit internal func for PostTxProcessing
	Deposit(ctx sdk.Context, event *DepositParams) error

	SetParams(ctx sdk.Context, params *types2.Params) error
	GetParams(ctx sdk.Context) (*types2.Params, error)
}
