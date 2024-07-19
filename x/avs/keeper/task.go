package keeper

import (
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetTaskInfo(ctx sdk.Context, task *types.TaskInfo) (err error) {
	taskContractAddress := task.TaskContractAddress
	// check operator address validation
	_, err = sdk.AccAddressFromBech32(taskContractAddress)
	if err != nil {
		return types.ErrInvalidAddr
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	infoKey := assetstype.GetJoinedStoreKey(taskContractAddress, task.TaskId)
	bz := k.cdc.MustMarshal(task)
	store.Set(infoKey, bz)
	return nil
}

func (k *Keeper) GetTaskInfo(ctx sdk.Context, taskID, taskContractAddress string) (info *types.TaskInfo, err error) {
	// check task address validation
	_, err = sdk.AccAddressFromBech32(taskContractAddress)
	if err != nil {
		return nil, types.ErrInvalidAddr
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	infoKey := assetstype.GetJoinedStoreKey(taskContractAddress, taskID)
	value := store.Get(infoKey)
	if value == nil {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetTaskInfo: key is %s", taskContractAddress))
	}

	ret := types.TaskInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k *Keeper) IsExistTask(ctx sdk.Context, taskID, taskContractAddress string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	infoKey := assetstype.GetJoinedStoreKey(taskContractAddress, taskID)

	return store.Has(infoKey)
}

func (k *Keeper) SetOperatorPubKey(ctx sdk.Context, pub *types.BlsPubKeyInfo) (err error) {
	operatorAddress, err := sdk.AccAddressFromBech32(pub.Operator)
	if err != nil {
		return types.ErrInvalidAddr
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixOperatePub)
	bz := k.cdc.MustMarshal(pub)
	store.Set(operatorAddress, bz)
	return nil
}

func (k *Keeper) GetOperatorPubKey(ctx sdk.Context, addr string) (pub *types.BlsPubKeyInfo, err error) {
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
	ret := types.BlsPubKeyInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k *Keeper) IsExistPubKey(ctx sdk.Context, addr string) bool {
	opAccAddr, _ := sdk.AccAddressFromBech32(addr)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixOperatePub)
	return store.Has(opAccAddr)
}
