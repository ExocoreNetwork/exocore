package keeper

import (
	"fmt"
	"math"
	"strings"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SetUndelegationRecords This function saves the undelegation records to be handled when the handle time expires.
// When we save the undelegation records, we save them in three kv stores which are `KeyPrefixUndelegationInfo` `KeyPrefixStakerUndelegationInfo` and `KeyPrefixWaitCompleteUndelegations`
func (k *Keeper) SetUndelegationRecords(ctx sdk.Context, records []*types.UndelegationRecord) error {
	singleRecordStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	stakerUndelegationStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	waitCompleteStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	// key := common.HexToAddress(incentive.Contract)
	for _, record := range records {
		bz := k.cdc.MustMarshal(record)
		// todo: check if the following state can only be set once?

		singleRecKey := types.GetUndelegationRecordKey(record.BlockNumber, record.LzTxNonce, record.TxHash, record.OperatorAddr)
		singleRecordStore.Set(singleRecKey, bz)

		stakerKey := types.GetStakerUndelegationRecordKey(record.StakerID, record.AssetID, record.LzTxNonce)
		stakerUndelegationStore.Set(stakerKey, singleRecKey)

		waitCompleteKey := types.GetWaitCompleteRecordKey(record.CompleteBlockNumber, record.LzTxNonce)
		waitCompleteStore.Set(waitCompleteKey, singleRecKey)
	}
	return nil
}

func (k *Keeper) DeleteUndelegationRecord(ctx sdk.Context, record *types.UndelegationRecord) error {
	singleRecordStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	stakerUndelegationStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	waitCompleteStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)

	singleRecKey := types.GetUndelegationRecordKey(record.BlockNumber, record.LzTxNonce, record.TxHash, record.OperatorAddr)
	singleRecordStore.Delete(singleRecKey)

	stakerKey := types.GetStakerUndelegationRecordKey(record.StakerID, record.AssetID, record.LzTxNonce)
	stakerUndelegationStore.Delete(stakerKey)

	waitCompleteKey := types.GetWaitCompleteRecordKey(record.CompleteBlockNumber, record.LzTxNonce)
	waitCompleteStore.Delete(waitCompleteKey)
	return nil
}

func (k *Keeper) SetSingleUndelegationRecord(ctx sdk.Context, record *types.UndelegationRecord) (recordKey []byte, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	bz := k.cdc.MustMarshal(record)
	key := types.GetUndelegationRecordKey(record.BlockNumber, record.LzTxNonce, record.TxHash, record.OperatorAddr)
	store.Set(key, bz)
	return key, nil
}

// StoreWaitCompleteRecord add it to handle the delay of completing undelegation caused by onHoldCount
// In the event that the undelegation is held by another module, this function is used within the EndBlocker to increment the scheduled completion block number by 1. Then the completion time of the undelegation will be delayed to the next block.
func (k *Keeper) StoreWaitCompleteRecord(ctx sdk.Context, singleRecKey []byte, record *types.UndelegationRecord) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	waitCompleteKey := types.GetWaitCompleteRecordKey(record.CompleteBlockNumber, record.LzTxNonce)
	store.Set(waitCompleteKey, singleRecKey)
	return nil
}

func (k *Keeper) GetUndelegationRecords(ctx sdk.Context, singleRecordKeys []string) (record []*types.UndelegationRecord, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	ret := make([]*types.UndelegationRecord, 0)
	for _, singleRecordKey := range singleRecordKeys {
		keyBytes := []byte(singleRecordKey)
		isExit := store.Has(keyBytes)
		undelegationRecord := types.UndelegationRecord{}
		if isExit {
			value := store.Get(keyBytes)
			k.cdc.MustUnmarshal(value, &undelegationRecord)
		} else {
			return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetSingleDelegationRecord: key is %s", singleRecordKey))
		}
		ret = append(ret, &undelegationRecord)
	}
	return ret, nil
}

// IterateUndelegationsByOperator iterate the undelegation records according to the operator
// and height filter. If the heightFilter isn't nil, only return the undelegations that the
// created height is greater than or equal to the filter height.
func (k *Keeper) IterateUndelegationsByOperator(
	ctx sdk.Context, operator string, heightFilter *uint64, isUpdate bool,
	opFunc func(undelegation *types.UndelegationRecord) error) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixUndelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, []byte(operator))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if heightFilter != nil {
			keyFields, err := types.ParseUndelegationRecordKey(iterator.Key())
			if err != nil {
				return err
			}
			if keyFields.BlockHeight < *heightFilter {
				continue
			}
		}
		undelegation := types.UndelegationRecord{}
		k.cdc.MustUnmarshal(iterator.Value(), &undelegation)
		err := opFunc(&undelegation)
		if err != nil {
			return err
		}

		if isUpdate {
			bz := k.cdc.MustMarshal(&undelegation)
			store.Set(iterator.Key(), bz)
		}
	}
	return nil
}

func (k *Keeper) SetStakerUndelegationInfo(ctx sdk.Context, stakerID, assetID string, recordKey []byte, lzNonce uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	key := types.GetStakerUndelegationRecordKey(stakerID, assetID, lzNonce)
	store.Set(key, recordKey)
	return nil
}

func (k *Keeper) GetStakerUndelegationRecKeys(ctx sdk.Context, stakerID, assetID string) (recordKeyList []string, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixStakerUndelegationInfo)
	iterator := sdk.KVStorePrefixIterator(store, []byte(strings.Join([]string{stakerID, assetID}, "/")))
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, string(iterator.Value()))
	}
	return ret, nil
}

func (k *Keeper) GetStakerUndelegationRecords(ctx sdk.Context, stakerID, assetID string) (records []*types.UndelegationRecord, err error) {
	recordKeys, err := k.GetStakerUndelegationRecKeys(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}

	return k.GetUndelegationRecords(ctx, recordKeys)
}

func (k *Keeper) SetWaitCompleteUndelegationInfo(ctx sdk.Context, height, lzNonce uint64, recordKey string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	key := types.GetWaitCompleteRecordKey(height, lzNonce)
	store.Set(key, []byte(recordKey))
	return nil
}

func (k *Keeper) GetWaitCompleteUndelegationRecKeys(ctx sdk.Context, height uint64) (recordKeyList []string, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixWaitCompleteUndelegations)
	iterator := sdk.KVStorePrefixIterator(store, []byte(hexutil.EncodeUint64(height)))
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		ret = append(ret, string(iterator.Value()))
	}
	return ret, nil
}

func (k *Keeper) GetWaitCompleteUndelegationRecords(ctx sdk.Context, height uint64) (records []*types.UndelegationRecord, err error) {
	recordKeys, err := k.GetWaitCompleteUndelegationRecKeys(ctx, height)
	if err != nil {
		return nil, err
	}
	if len(recordKeys) == 0 {
		return nil, nil
	}
	// The states of records stored by WaitCompleteUndelegations kvStore should always be IsPending,so using AllRecords as getType here is ok.
	return k.GetUndelegationRecords(ctx, recordKeys)
}

func (k *Keeper) IncrementUndelegationHoldCount(ctx sdk.Context, recordKey []byte) error {
	prev := k.GetUndelegationHoldCount(ctx, recordKey)
	if prev == math.MaxUint64 {
		return types.ErrCannotIncHoldCount
	}
	now := prev + 1
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetUndelegationOnHoldKey(recordKey), sdk.Uint64ToBigEndian(now))
	return nil
}

func (k *Keeper) GetUndelegationHoldCount(ctx sdk.Context, recordKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUndelegationOnHoldKey(recordKey))
	return sdk.BigEndianToUint64(bz)
}

func (k *Keeper) DecrementUndelegationHoldCount(ctx sdk.Context, recordKey []byte) error {
	prev := k.GetUndelegationHoldCount(ctx, recordKey)
	if prev == 0 {
		return types.ErrCannotDecHoldCount
	}
	now := prev - 1
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetUndelegationOnHoldKey(recordKey), sdk.Uint64ToBigEndian(now))
	return nil
}
