package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/exocore/x/delegation/types"
	"strings"
)

type GetStakerUnDelegationRecordType uint8

const (
	PendingRecords GetStakerUnDelegationRecordType = iota
	CompletedRecords
	AllRecords
)

func (k Keeper) SetUnDelegationStates(ctx sdk.Context, records []*types.UnDelegationRecord) error {
	singleRecordStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUnDelegationInfo)
	stakerUnDelegationStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUnDelegationInfo)
	waitCompleteStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUnDelegations)
	//key := common.HexToAddress(incentive.Contract)
	for _, record := range records {
		bz := k.cdc.MustMarshal(record)
		singleRecKey := types.GetUnDelegationRecordKey(record.LzTxNonce, record.TxHash, record.OperatorAddr)
		singleRecordStore.Set(singleRecKey, bz)

		stakerKey := types.GetStakerUnDelegationRecordKey(record.StakerId, record.AssetId, record.LzTxNonce)
		stakerUnDelegationStore.Set(stakerKey, singleRecKey)

		waitCompleteKey := types.GetWaitCompleteRecordKey(record.CompleteBlockNumber, record.LzTxNonce)
		waitCompleteStore.Set(waitCompleteKey, []byte(singleRecKey))
	}
	return nil
}

func (k Keeper) SetSingleUnDelegationRecord(ctx sdk.Context, record *types.UnDelegationRecord) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUnDelegationInfo)
	bz := k.cdc.MustMarshal(record)
	key := types.GetUnDelegationRecordKey(record.LzTxNonce, record.TxHash, record.OperatorAddr)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetUnDelegationRecords(ctx sdk.Context, singleRecordKeys []string) (record []*types.UnDelegationRecord, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUnDelegationInfo)
	ret := make([]*types.UnDelegationRecord, 0)
	for _, singleRecordKey := range singleRecordKeys {
		keyBytes := []byte(singleRecordKey)
		isExit := store.Has(keyBytes)
		unDelegationRecord := types.UnDelegationRecord{}
		if isExit {
			value := store.Get(keyBytes)
			k.cdc.MustUnmarshal(value, &unDelegationRecord)
		} else {
			return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetSingleDelegationRecord: key is %s", singleRecordKey))
		}
		ret = append(ret, &unDelegationRecord)
	}

	return ret, nil
}

func (k Keeper) SetStakerUnDelegationInfo(ctx sdk.Context, stakerId, assetId, recordKey string, lzNonce uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUnDelegationInfo)
	key := types.GetStakerUnDelegationRecordKey(stakerId, assetId, lzNonce)
	store.Set(key, []byte(recordKey))
	return nil
}

func (k Keeper) GetStakerUnDelegationRecKeys(ctx sdk.Context, stakerId, assetId string) (recordKeyList []string, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUnDelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, []byte(strings.Join([]string{stakerId, assetId}, "/")))
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, string(iterator.Value()))
	}
	return ret, nil
}

func (k Keeper) GetStakerUnDelegationPendingRecords(ctx sdk.Context, stakerId, assetId string, getType GetStakerUnDelegationRecordType) (records []*types.UnDelegationRecord, err error) {
	ret := make([]*types.UnDelegationRecord, 0)
	recordKeys, err := k.GetStakerUnDelegationRecKeys(ctx, stakerId, assetId)
	if err != nil {
		return nil, err
	}

	getAllRecords, err := k.GetUnDelegationRecords(ctx, recordKeys)
	if err != nil {
		return nil, err
	}

	if getType == AllRecords {
		return getAllRecords, nil
	}

	for _, record := range getAllRecords {
		if getType == PendingRecords {
			if record.IsPending {
				ret = append(ret, record)
			}
		} else if getType == CompletedRecords {
			if !record.IsPending {
				ret = append(ret, record)
			}
		} else {
			return nil, errorsmod.Wrap(types.ErrStakerGetRecordType, fmt.Sprintf("the getType is:%v", getType))
		}
	}
	return ret, nil
}

func (k Keeper) SetWaitCompleteUnDelegationInfo(ctx sdk.Context, height, lzNonce uint64, recordKey string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUnDelegations)
	key := types.GetWaitCompleteRecordKey(height, lzNonce)
	store.Set(key, []byte(recordKey))
	return nil
}

func (k Keeper) GetWaitCompleteUnDelegationRecKeys(ctx sdk.Context, height uint64) (recordKeyList []string, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUnDelegations)
	iterator := sdk.KVStorePrefixIterator(store, []byte(hexutil.EncodeUint64(height)))
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, string(iterator.Value()))
	}
	return ret, nil
}

func (k Keeper) GetWaitCompleteUnDelegationRecords(ctx sdk.Context, height uint64) (records []*types.UnDelegationRecord, err error) {
	recordKeys, err := k.GetWaitCompleteUnDelegationRecKeys(ctx, height)
	if err != nil {
		return nil, err
	}
	return k.GetUnDelegationRecords(ctx, recordKeys)
}
