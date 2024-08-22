package keeper

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"

	errorsmod "cosmossdk.io/errors"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls/blst"
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
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore,
			fmt.Sprintf("GetTaskInfo: key is %s", taskContractAddress))
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
		return nil, errorsmod.Wrap(err, "GetOperatorPubKey: error occurred when parse acc "+
			"address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixOperatePub)
	// key := common.HexToAddress(incentive.Contract)
	isExist := store.Has(opAccAddr)
	if !isExist {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore,
			fmt.Sprintf("GetOperatorPubKey: key is %s", opAccAddr))
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

// GetTaskID Increase the task ID by 1 each time.
func (k Keeper) GetTaskID(ctx sdk.Context, taskaddr common.Address) uint64 {
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

// SetTaskResultInfo is used to store the operator's sign task information.
func (k *Keeper) SetTaskResultInfo(
	ctx sdk.Context, addr string, info *types.TaskResultInfo,
) (err error) {
	// the operator's `addr` must match the from address.
	if addr != info.OperatorAddress {
		return errorsmod.Wrap(
			types.ErrInvalidAddr,
			"SetTaskResultInfo:from address is not equal to the operator address",
		)
	}
	opAccAddr, _ := sdk.AccAddressFromBech32(info.OperatorAddress)
	// check operator
	if !k.operatorKeeper.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(
			delegationtypes.ErrOperatorNotExist,
			fmt.Sprintf("SetTaskResultInfo:invalid operator address:%s", opAccAddr),
		)
	}
	// check operator bls pubkey
	keyInfo, _ := k.GetOperatorPubKey(ctx, info.OperatorAddress)
	pubKey, err := blst.PublicKeyFromBytes(keyInfo.PubKey)
	if err != nil || pubKey == nil {
		return errorsmod.Wrap(
			types.ErrPubKeyIsNotExists,
			fmt.Sprintf("SetTaskResultInfo:get operator address:%s", opAccAddr),
		)
	}
	//	check task contract
	task, err := k.GetTaskInfo(ctx, strconv.FormatUint(info.TaskId, 10), info.TaskContractAddress)
	if err != nil || task.TaskContractAddress == "" {
		return errorsmod.Wrap(
			types.ErrTaskIsNotExists,
			fmt.Sprintf("SetTaskResultInfo: task info not found: %s (Task ID: %d)",
				info.TaskContractAddress, info.TaskId),
		)
	}

	//  check prescribed period
	//  If submitted in the first stage, in order  to avoid plagiarism by other operators,
	//	TaskResponse and TaskResponseHash must be null values
	//	At the same time, it must be submitted within the response deadline in the first stage
	avsInfo := k.GetAVSInfoByTaskAddress(ctx, info.TaskContractAddress)
	epoch, found := k.epochsKeeper.GetEpochInfo(ctx, avsInfo.EpochIdentifier)
	if !found {
		return errorsmod.Wrap(types.ErrEpochNotFound, fmt.Sprintf("epoch info not found %s",
			avsInfo.EpochIdentifier))
	}

	switch info.Stage {
	case types.TwoPhaseCommitOne:
		if k.IsExistTaskResultInfo(ctx, info.OperatorAddress, info.TaskContractAddress, info.TaskId) {
			return errorsmod.Wrap(
				types.ErrResAlreadyExists,
				fmt.Sprintf("SetTaskResultInfo: task result is already exists, "+
					"OperatorAddress: %s (TaskContractAddress: %s),(Task ID: %d)",
					info.OperatorAddress, info.TaskContractAddress, info.TaskId),
			)
		}
		if info.TaskResponseHash != "" || info.TaskResponse != nil {
			return errorsmod.Wrap(
				types.ErrParamNotEmptyError,
				fmt.Sprintf("SetTaskResultInfo: invalid param TaskResponseHash: %s (TaskResponse: %d)",
					info.TaskResponseHash, info.TaskResponse),
			)
		}
		// check epoch，The first stage submission must be within the response window period
		if epoch.CurrentEpoch > int64(task.StartingEpoch)+int64(task.TaskResponsePeriod) {
			return errorsmod.Wrap(
				types.ErrSubmitTooLateError,
				fmt.Sprintf("SetTaskResultInfo:submit  too late, CurrentEpoch:%d", epoch.CurrentEpoch),
			)
		}
		infoKey := assetstype.GetJoinedStoreKey(info.OperatorAddress, info.TaskContractAddress,
			strconv.FormatUint(info.TaskId, 10))
		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTaskResult)
		bz := k.cdc.MustMarshal(info)
		store.Set(infoKey, bz)
		return nil

	case types.TwoPhaseCommitTwo:
		// check task response
		if info.TaskResponseHash == "" || info.TaskResponse == nil {
			return errorsmod.Wrap(
				types.ErrNotNull,
				fmt.Sprintf("SetTaskResultInfo: invalid param TaskResponseHash: %s (TaskResponse: %d)",
					info.TaskResponseHash, info.TaskResponse),
			)
		}
		// check parameters
		res, err := k.GetTaskResultInfo(ctx, info.OperatorAddress, info.TaskContractAddress, info.TaskId)
		if err != nil || res.OperatorAddress != info.OperatorAddress ||
			res.TaskContractAddress != info.TaskContractAddress ||
			res.TaskId != info.TaskId || !bytes.Equal(res.BlsSignature, info.BlsSignature) {
			return errorsmod.Wrap(
				types.ErrInconsistentParams,
				fmt.Sprintf("SetTaskResultInfo: invalid param OperatorAddress: %s ,(TaskContractAddress: %s)"+
					",(TaskId: %d),(BlsSignature: %s)",
					info.OperatorAddress, info.TaskContractAddress, info.TaskId, info.BlsSignature),
			)
		}
		//  check epoch，The second stage submission must be within the statistical window period
		if epoch.CurrentEpoch <= int64(task.StartingEpoch)+int64(task.TaskResponsePeriod) ||
			epoch.CurrentEpoch > int64(task.StartingEpoch)+int64(task.TaskResponsePeriod)+int64(task.TaskStatisticalPeriod) {
			return errorsmod.Wrap(
				types.ErrSubmitTooLateError,
				fmt.Sprintf("SetTaskResultInfo:submit  too late, CurrentEpoch:%d", epoch.CurrentEpoch),
			)
		}

		// check hash
		taskResponseDigest := crypto.Keccak256Hash(info.TaskResponse)
		if taskResponseDigest.String() != info.TaskResponseHash {
			return errorsmod.Wrap(
				types.ErrHashValue,
				"SetTaskResultInfo: task response is nil",
			)
		}
		// check bls sig
		flag, err := blst.VerifySignature(info.BlsSignature, taskResponseDigest, pubKey)
		if !flag || err != nil {
			return errorsmod.Wrap(
				types.ErrSigVerifyError,
				fmt.Sprintf("SetTaskResultInfo: invalid task address: %s (Task ID: %d)", info.TaskContractAddress, info.TaskId),
			)
		}

		infoKey := assetstype.GetJoinedStoreKey(info.OperatorAddress, info.TaskContractAddress, strconv.FormatUint(info.TaskId, 10))

		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTaskResult)
		bz := k.cdc.MustMarshal(info)
		store.Set(infoKey, bz)
		return nil
	default:
		return errorsmod.Wrap(
			types.ErrParamError,
			fmt.Sprintf("SetTaskResultInfo: invalid param value:%s", info.Stage),
		)
	}
}

func (k *Keeper) IsExistTaskResultInfo(ctx sdk.Context, operatorAddress, taskContractAddress string, taskID uint64) bool {
	infoKey := assetstype.GetJoinedStoreKey(operatorAddress, taskContractAddress,
		strconv.FormatUint(taskID, 10))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTaskResult)
	return store.Has(infoKey)
}

func (k *Keeper) GetTaskResultInfo(ctx sdk.Context, operatorAddress, taskContractAddress string, taskID uint64) (info *types.TaskResultInfo, err error) {
	if !common.IsHexAddress(taskContractAddress) {
		return nil, types.ErrInvalidAddr
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTaskResult)
	infoKey := assetstype.GetJoinedStoreKey(operatorAddress, taskContractAddress,
		strconv.FormatUint(taskID, 10))
	value := store.Get(infoKey)
	if value == nil {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore,
			fmt.Sprintf("GetTaskResultInfo: key is %s", infoKey))
	}

	ret := types.TaskResultInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// IterateResultInfo iterate through task result info
func (k Keeper) IterateResultInfo(ctx sdk.Context, fn func(index int64, info types.TaskResultInfo) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixTaskResult)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		task := types.TaskResultInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &task)

		stop := fn(i, task)

		if stop {
			break
		}
		i++
	}
}
