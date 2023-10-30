// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	types2 "github.com/exocore/x/delegation/types"
	keeper2 "github.com/exocore/x/deposit/keeper"
	"github.com/exocore/x/restaking_assets_manage/keeper"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	//other keepers
	retakingStateKeeper keeper.Keeper
	depositKeeper       keeper2.Keeper
}

func (k Keeper) CompleteUnDelegateAssetFromOperator() error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) SetOperatorInfo(ctx sdk.Context, addr string, info *types2.OperatorInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorInfo)
	//todo: think about the difference between init and update in future
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	store.Set([]byte(addr), bz)
	return nil
}

func (k Keeper) GetOperatorInfo(ctx sdk.Context, addr string) (info *types2.OperatorInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixOperatorInfo)
	//key := common.HexToAddress(incentive.Contract)
	ifExist := store.Has([]byte(addr))
	if !ifExist {
		return nil, types2.ErrNoOperatorInfoKey
	}

	value := store.Get([]byte(addr))

	ret := types2.OperatorInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// IDelegation interface will be implemented by deposit keeper
type IDelegation interface {
	// PostTxProcessing automatically call PostTxProcessing to update delegation state after receiving delegation event tx from layerZero protocol
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// RegisterOperator handle the registerOperator txs from msg service
	RegisterOperator(ctx context.Context, req *types2.RegisterOperatorReq) (*types2.RegisterOperatorResponse, error)
	// DelegateAssetToOperator handle the DelegateAssetToOperator txs from msg service
	DelegateAssetToOperator(ctx context.Context, delegation *types2.MsgDelegation) (*types2.DelegationResponse, error)
	// UnDelegateAssetFromOperator handle the UnDelegateAssetFromOperator txs from msg service
	UnDelegateAssetFromOperator(ctx context.Context, delegation *types2.MsgUnDelegation) (*types2.UnDelegationResponse, error)

	//GetDelegationInfo grpc_query interface
	//GetDelegationInfo(context.Context, *types2.DelegationInfoReq) (*types2.QueryDelegationInfoResponse, error)

	// CompleteUnDelegateAssetFromOperator scheduled execute to handle UnDelegateAssetFromOperator through two steps
	CompleteUnDelegateAssetFromOperator() error
}
