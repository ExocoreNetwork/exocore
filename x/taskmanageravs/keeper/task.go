package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"strings"
)

func (k *Keeper) SetAvsTaskInfo(ctx sdk.Context, info *types.RegisterAVSTaskReq) (res string, err error) {

	key := strings.Join([]string{strings.ToLower(info.AVSAddress), uuid.NewString()}, "_")

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)

	bz := k.cdc.MustMarshal(info)

	store.Set([]byte(key), bz)
	return key, nil
}

func (k *Keeper) GetAvsTaskInfo(ctx sdk.Context, index string) (info *types.TaskContractInfo, err error) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	//key := common.HexToAddress(incentive.Contract)
	isExist := store.Has([]byte(index))
	if !isExist {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorInfo: key is %suite", index))
	}

	value := store.Get([]byte(index))

	ret := types.RegisterAVSTaskReq{}
	k.cdc.MustUnmarshal(value, &ret)
	return ret.Task, nil
}

func (k *Keeper) IsExistTask(ctx sdk.Context, index string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskInfo)
	return store.Has([]byte(index))
}

func (k Keeper) GetTaskIndex(ctx sdk.Context, addr string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAvsTaskIdMap)
	return store.Get([]byte(addr))
}

func (k Keeper) SetTaskIndex(ctx sdk.Context, addr string, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAvsTaskIdMap)
	store.Set([]byte(addr), []byte(id))
}

func (k Keeper) SetTaskforAvs(ctx sdk.Context, params *CreateNewTaskParams) (recordKey string, err error) {
	key := strings.Join([]string{strings.ToLower(params.ContractAddr), uuid.NewString()}, "_")

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSTaskMap)

	task := types.TaskInstance{
		NumberToBeSquared:         params.NumberToBeSquared,
		QuorumThresholdPercentage: uint64(params.QuorumThresholdPercentage),
		TaskCreatedBlock:          uint64(params.TaskCreatedBlock),
	}
	bz := k.cdc.MustMarshal(&task)

	store.Set([]byte(key), bz)
	return key, nil
}

func (k *Keeper) CreateNewTask(ctx sdk.Context, params *CreateNewTaskParams) (bool, error) {
	// create a new task struct
	k.SetTaskforAvs(ctx, params)
	return true, nil
}

type CreateNewTaskParams struct {
	TaskIndex                 uint32
	NumberToBeSquared         uint64
	QuorumThresholdPercentage uint32
	QuorumNumbers             []byte
	ContractAddr              string
	TaskCreatedBlock          int64
}
