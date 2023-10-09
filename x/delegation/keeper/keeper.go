// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	"context"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	types2 "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	//other keepers
	retakingStateKeeper keeper.Keeper
}

func (k Keeper) InternalDelegateAssetToOperator(reStakerId string, operatorAssetsInfo map[string]map[string]math.Uint, approvedInfo map[string]*types2.DelegationApproveInfo) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) InternalUnDelegateAssetFromOperator(reStakerId string, operatorAssetsInfo map[string]map[string]math.Uint) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) CompleteUnDelegateAssetFromOperator() error {
	//TODO implement me
	panic("implement me")
}

// IDelegation interface will be implemented by deposit keeper
type IDelegation interface {
	// PostTxProcessing automatically call PostTxProcessing to update delegation state after receiving delegation event tx from layerZero protocol
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// InternalDelegateAssetToOperator internal func for PostTxProcessing
	InternalDelegateAssetToOperator(reStakerId string, operatorAssetsInfo map[string]map[string]math.Uint, approvedInfo map[string]*types2.DelegationApproveInfo) error
	// InternalUnDelegateAssetFromOperator internal func for PostTxProcessing
	InternalUnDelegateAssetFromOperator(reStakerId string, operatorAssetsInfo map[string]map[string]math.Uint) error

	// RegisterOperator handle the registerOperator txs from msg service
	RegisterOperator(context.Context, *types2.OperatorInfo) (*types2.RegisterOperatorResponse, error)
	// DelegateAssetToOperator handle the DelegateAssetToOperator txs from msg service
	DelegateAssetToOperator(context.Context, *types2.MsgDelegation) (*types2.DelegationResponse, error)
	// UnDelegateAssetFromOperator handle the UnDelegateAssetFromOperator txs from msg service
	UnDelegateAssetFromOperator(context.Context, *types2.MsgUnDelegation) (*types2.UnDelegationResponse, error)

	//GetDelegationInfo grpc_query interface
	GetDelegationInfo(context.Context, *types2.QueryDelegationInfo) (*types2.QueryDelegationInfoResponse, error)

	// CompleteUnDelegateAssetFromOperator scheduled execute to handle UnDelegateAssetFromOperator through two steps
	CompleteUnDelegateAssetFromOperator() error
}
