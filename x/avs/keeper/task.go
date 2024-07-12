package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) SetAVSTaskInfo(ctx sdk.Context, info *types.RegisterAVSTaskReq) (err error) {
	taskAccAddr, err := sdk.AccAddressFromBech32(info.Task.TaskContractAddress)
	if err != nil {
		return errorsmod.Wrap(err, "SetTaskInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)

	bz := k.cdc.MustMarshal(info)

	store.Set(taskAccAddr, bz)
	return nil
}

func (k *Keeper) GetAVSTaskInfo(ctx sdk.Context, addr string) (info *types.TaskContractInfo, err error) {
	taskAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetAVSTaskInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	isExist := store.Has(taskAccAddr)
	if !isExist {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSTaskInfo: key is %s", taskAccAddr))
	}

	value := store.Get(taskAccAddr)

	ret := types.RegisterAVSTaskReq{}
	k.cdc.MustUnmarshal(value, &ret)
	return ret.Task, nil
}

func (k *Keeper) IsExistTask(ctx sdk.Context, addr sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	return store.Has(addr)
}

func (k *Keeper) SetOperatorPubKey(ctx sdk.Context, addr string, pub []byte) (err error) {
	opAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return errorsmod.Wrap(err, "SetOperatorPubKey: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixOperatePub)

	store.Set(opAccAddr, pub)
	return nil
}

func (k *Keeper) GetOperatorPubKey(ctx sdk.Context, addr string) (pub []byte, err error) {
	opAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetOperatorPubKey: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixOperatePub)
	// key := common.HexToAddress(incentive.Contract)
	isExist := store.Has(opAccAddr)
	if !isExist {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorPubKey: key is %s", opAccAddr))
	}

	value := store.Get(opAccAddr)

	return value, nil
}
