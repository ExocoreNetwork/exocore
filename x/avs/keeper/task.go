package keeper

import (
	"fmt"
	"strconv"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ethereum/go-ethereum/common"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetTaskInfo(ctx sdk.Context, task *types.TaskInfo) (err error) {
	if !common.IsHexAddress(task.TaskContractAddress) {
		return types.ErrInvalidAddr
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	infoKey := assetstype.GetJoinedStoreKey(task.TaskContractAddress, strconv.FormatUint(task.TaskId, 10))
	bz := k.cdc.MustMarshal(task)
	store.Set(infoKey, bz)
	return nil
}

func (k *Keeper) GetTaskInfo(ctx sdk.Context, taskID, taskContractAddress string) (info *types.TaskInfo, err error) {
	if !common.IsHexAddress(taskContractAddress) {
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

// IterateTaskAVSInfo iterate through task
func (k Keeper) IterateTaskAVSInfo(ctx sdk.Context, fn func(index int64, taskInfo types.TaskInfo) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		task := types.TaskInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &task)

		stop := fn(i, task)

		if stop {
			break
		}
		i++
	}
}

// GetTaskId Increase the task ID by 1 each time.
func (k Keeper) GetTaskId(ctx sdk.Context, taskaddr common.Address) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLatestTaskNum)
	var id uint64
	if store.Has(taskaddr.Bytes()) {
		bz := store.Get(taskaddr.Bytes())
		id = sdk.BigEndianToUint64(bz)
		id++
	} else {
		id = 1
	}
	store.Set(taskaddr.Bytes(), sdk.Uint64ToBigEndian(id))
	return id
}
