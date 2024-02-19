package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type GetUndelegationRecordType uint8

const (
	PendingRecords GetUndelegationRecordType = iota
	CompletedRecords
	AllRecords
)

// SetUndelegationRecords This function saves the undelegation records to be handled when the handle time expires.
// When we save the undelegation records, we save them in three kv stores which are `KeyPrefixUndelegationInfo` `KeyPrefixStakerUndelegationInfo` and `KeyPrefixWaitCompleteUndelegations`
func (k Keeper) SetUndelegationRecords(ctx sdk.Context, records []*types.UndelegationRecord) error {
	singleRecordStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	stakerUndelegationStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	waitCompleteStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	// key := common.HexToAddress(incentive.Contract)
	for _, record := range records {
		bz := k.cdc.MustMarshal(record)
		// todo: check if the following state can only be set once?

		singleRecKey := types.GetUndelegationRecordKey(record.LzTxNonce, record.TxHash, record.OperatorAddr)
		singleRecordStore.Set(singleRecKey, bz)

		stakerKey := types.GetStakerUndelegationRecordKey(record.StakerID, record.AssetID, record.LzTxNonce)
		stakerUndelegationStore.Set(stakerKey, singleRecKey)

		waitCompleteKey := types.GetWaitCompleteRecordKey(record.CompleteBlockNumber, record.LzTxNonce)
		waitCompleteStore.Set(waitCompleteKey, singleRecKey)
	}
	return nil
}

func (k Keeper) SetSingleUndelegationRecord(ctx sdk.Context, record *types.UndelegationRecord) (recordKey []byte, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	bz := k.cdc.MustMarshal(record)
	key := types.GetUndelegationRecordKey(record.LzTxNonce, record.TxHash, record.OperatorAddr)
	store.Set(key, bz)
	return key, nil
}

func (k Keeper) GetUndelegationRecords(ctx sdk.Context, singleRecordKeys []string, getType GetUndelegationRecordType) (record []*types.UndelegationRecord, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	ret := make([]*types.UndelegationRecord, 0)
	for _, singleRecordKey := range singleRecordKeys {
		keyBytes := []byte(singleRecordKey)
		isExit := store.Has(keyBytes)
		UndelegationRecord := types.UndelegationRecord{}
		if isExit {
			value := store.Get(keyBytes)
			k.cdc.MustUnmarshal(value, &UndelegationRecord)
		} else {
			return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetSingleDelegationRecord: key is %s", singleRecordKey))
		}

		if getType == PendingRecords {
			if UndelegationRecord.IsPending {
				ret = append(ret, &UndelegationRecord)
			}
		} else if getType == CompletedRecords {
			if !UndelegationRecord.IsPending {
				ret = append(ret, &UndelegationRecord)
			}
		} else if getType == AllRecords {
			ret = append(ret, &UndelegationRecord)
		} else {
			return nil, errorsmod.Wrap(types.ErrStakerGetRecordType, fmt.Sprintf("the getType is:%v", getType))
		}
	}
	return ret, nil
}

func (k Keeper) SetStakerUndelegationInfo(ctx sdk.Context, stakerId, assetId string, recordKey []byte, lzNonce uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	key := types.GetStakerUndelegationRecordKey(stakerId, assetId, lzNonce)
	store.Set(key, recordKey)
	return nil
}

func (k Keeper) GetStakerUndelegationRecKeys(ctx sdk.Context, stakerId, assetId string) (recordKeyList []string, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, []byte(strings.Join([]string{stakerId, assetId}, "/")))
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, string(iterator.Value()))
	}
	return ret, nil
}

func (k Keeper) GetStakerUndelegationRecords(ctx sdk.Context, stakerId, assetId string, getType GetUndelegationRecordType) (records []*types.UndelegationRecord, err error) {
	recordKeys, err := k.GetStakerUndelegationRecKeys(ctx, stakerId, assetId)
	if err != nil {
		return nil, err
	}

	return k.GetUndelegationRecords(ctx, recordKeys, getType)
}

func (k Keeper) SetWaitCompleteUndelegationInfo(ctx sdk.Context, height, lzNonce uint64, recordKey string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	key := types.GetWaitCompleteRecordKey(height, lzNonce)
	store.Set(key, []byte(recordKey))
	return nil
}

func (k Keeper) GetWaitCompleteUndelegationRecKeys(ctx sdk.Context, height uint64) (recordKeyList []string, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	iterator := sdk.KVStorePrefixIterator(store, []byte(hexutil.EncodeUint64(height)))
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, string(iterator.Value()))
	}
	return ret, nil
}

func (k Keeper) GetWaitCompleteUndelegationRecords(ctx sdk.Context, height uint64) (records []*types.UndelegationRecord, err error) {
	recordKeys, err := k.GetWaitCompleteUndelegationRecKeys(ctx, height)
	if err != nil {
		return nil, err
	}
	if len(recordKeys) == 0 {
		return nil, nil
	}
	// The states of records stored by WaitCompleteUndelegations kvStore should always be IsPending,so using AllRecords as getType here is ok.
	return k.GetUndelegationRecords(ctx, recordKeys, AllRecords)
}
