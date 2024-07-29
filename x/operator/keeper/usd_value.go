package keeper

import (
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpdateOperatorUSDValue is a function to update the USD share for specified operator and Avs,
// The key and value that will be changed is:
// AVSAddr + '/' + operatorAddr -> types.OperatorOptedUSDValue (the total USD share of specified operator and Avs)
// This function will be called when some assets supported by Avs are delegated/undelegated or slashed.
func (k *Keeper) UpdateOperatorUSDValue(ctx sdk.Context, avsAddr, operatorAddr string, delta operatortypes.DeltaOperatorUSDInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForOperator)
	var key []byte
	if operatorAddr == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, "UpdateOperatorUSDValue the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)

	usdInfo := operatortypes.OperatorOptedUSDValue{
		SelfUSDValue:   sdkmath.LegacyNewDec(0),
		TotalUSDValue:  sdkmath.LegacyNewDec(0),
		ActiveUSDValue: sdkmath.LegacyNewDec(0),
	}
	value := store.Get(key)
	if value != nil {
		k.cdc.MustUnmarshal(value, &usdInfo)
	}

	err := assetstype.UpdateAssetDecValue(&usdInfo.SelfUSDValue, &delta.SelfUSDValue)
	if err != nil {
		return err
	}
	err = assetstype.UpdateAssetDecValue(&usdInfo.TotalUSDValue, &delta.TotalUSDValue)
	if err != nil {
		return err
	}
	err = assetstype.UpdateAssetDecValue(&usdInfo.ActiveUSDValue, &delta.ActiveUSDValue)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&usdInfo)
	store.Set(key, bz)
	return nil
}

func (k *Keeper) InitOperatorUSDValue(ctx sdk.Context, avsAddr, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForOperator)
	var key []byte
	if operatorAddr == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, "UpdateOperatorUSDValue the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)
	if store.Has(key) {
		return errorsmod.Wrap(operatortypes.ErrKeyAlreadyExist, fmt.Sprintf("avsAddr operatorAddr is: %s, %s", avsAddr, operatorAddr))
	}
	initValue := operatortypes.OperatorOptedUSDValue{
		SelfUSDValue:   sdkmath.LegacyNewDec(0),
		TotalUSDValue:  sdkmath.LegacyNewDec(0),
		ActiveUSDValue: sdkmath.LegacyNewDec(0),
	}
	bz := k.cdc.MustMarshal(&initValue)
	store.Set(key, bz)
	return nil
}

// DeleteOperatorUSDValue is a function to delete the USD share related to specified operator and Avs,
// The key and value that will be deleted is:
// AVSAddr + '/' + operatorAddr -> types.OperatorOptedUSDValue (the total USD share of specified operator and Avs)
// This function will be called when the operator opts out of the AVS, because the USD share
// doesn't need to be stored.
func (k *Keeper) DeleteOperatorUSDValue(ctx sdk.Context, avsAddr, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForOperator)
	var key []byte
	if operatorAddr == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, "UpdateOperatorUSDValue the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)
	store.Delete(key)

	return nil
}

// GetOperatorOptedUSDValue is a function to retrieve the USD share of specified operator and Avs,
// The key and value to retrieve is:
// AVSAddr + '/' + operatorAddr -> types.OperatorOptedUSDValue (the total USD share of specified operator and Avs)
// This function will be called when the operator opts out of the AVS, because the total USD share
// of Avs should decrease the USD share of the opted-out operator
// This function can also serve as an RPC in the future.
func (k *Keeper) GetOperatorOptedUSDValue(ctx sdk.Context, avsAddr, operatorAddr string) (operatortypes.OperatorOptedUSDValue, error) {
	// return zero if the operator has opted-out of the AVS
	if !k.IsOptedIn(ctx, operatorAddr, avsAddr) {
		return operatortypes.OperatorOptedUSDValue{
			SelfUSDValue:   sdkmath.LegacyNewDec(0),
			TotalUSDValue:  sdkmath.LegacyNewDec(0),
			ActiveUSDValue: sdkmath.LegacyNewDec(0),
		}, nil
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForOperator)
	var ret operatortypes.OperatorOptedUSDValue
	var key []byte
	if operatorAddr == "" {
		return operatortypes.OperatorOptedUSDValue{}, errorsmod.Wrap(operatortypes.ErrParameterInvalid, "GetOperatorOptedUSDValue the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)
	value := store.Get(key)
	if value == nil {
		return operatortypes.OperatorOptedUSDValue{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorOptedUSDValue: key is %s", key))
	}
	k.cdc.MustUnmarshal(value, &ret)

	return ret, nil
}

// UpdateAVSUSDValue is a function to update the total USD share of an Avs,
// The key and value that will be changed is:
// AVSAddr -> types.DecValueField（the total USD share of specified Avs）
// This function will be called when some assets of operator supported by the specified Avs
// are delegated/undelegated or slashed. Additionally, when an operator opts out of
// the Avs, this function also will be called.
func (k *Keeper) UpdateAVSUSDValue(ctx sdk.Context, avsAddr string, opAmount sdkmath.LegacyDec) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return errorsmod.Wrap(operatortypes.ErrValueIsNilOrZero, fmt.Sprintf("UpdateAVSUSDValue the opAmount is:%v", opAmount))
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForAVS)
	key := []byte(avsAddr)
	totalValue := operatortypes.DecValueField{Amount: sdkmath.LegacyNewDec(0)}
	value := store.Get(key)
	if value != nil {
		k.cdc.MustUnmarshal(value, &totalValue)
	}

	err := assetstype.UpdateAssetDecValue(&totalValue.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&totalValue)
	store.Set(key, bz)
	return nil
}

// SetAVSUSDValue is a function to set the total USD share of an Avs,
func (k *Keeper) SetAVSUSDValue(ctx sdk.Context, avsAddr string, amount sdkmath.LegacyDec) error {
	if amount.IsNil() {
		return errorsmod.Wrap(operatortypes.ErrValueIsNilOrZero, fmt.Sprintf("SetAVSUSDValue the amount is:%v", amount))
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForAVS)
	key := []byte(avsAddr)
	setValue := operatortypes.DecValueField{Amount: amount}
	bz := k.cdc.MustMarshal(&setValue)
	store.Set(key, bz)
	return nil
}

// GetAVSUSDValue is a function to retrieve the USD share of specified Avs,
// The key and value to retrieve is:
// AVSAddr -> types.DecValueField（the total USD share of specified Avs）
func (k *Keeper) GetAVSUSDValue(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		operatortypes.KeyPrefixVotingPowerForAVS,
	)
	var ret operatortypes.DecValueField
	key := []byte(avsAddr)
	value := store.Get(key)
	if value == nil {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSUSDValue: key is %s", key))
	}
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}

// IterateOperatorsForAVS is used to iterate the operators of a specified AVS and do some external operations
// `isUpdate` is a flag to indicate whether the change of the state should be set to the store.
func (k *Keeper) IterateOperatorsForAVS(ctx sdk.Context, avsAddr string, isUpdate bool, opFunc func(operator string, optedUSDValues *operatortypes.OperatorOptedUSDValue) error) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixVotingPowerForOperator)
	iterator := sdk.KVStorePrefixIterator(store, operatortypes.IterateOperatorsForAVSPrefix(avsAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keys, err := assetstype.ParseJoinedKey(iterator.Key())
		if err != nil {
			return err
		}
		var optedUSDValues operatortypes.OperatorOptedUSDValue
		k.cdc.MustUnmarshal(iterator.Value(), &optedUSDValues)
		err = opFunc(keys[1], &optedUSDValues)
		if err != nil {
			return err
		}
		if isUpdate {
			bz := k.cdc.MustMarshal(&optedUSDValues)
			store.Set(iterator.Key(), bz)
		}
	}
	return nil
}

func (k Keeper) GetVotePowerForChainID(
	ctx sdk.Context, operators []sdk.AccAddress, chainID string,
) ([]int64, error) {
	avsAddr := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
	ret := make([]int64, 0)
	for _, operator := range operators {
		// this already filters by the required assetIDs
		optedUSDValues, err := k.GetOperatorOptedUSDValue(ctx, avsAddr, operator.String())
		if err != nil {
			return nil, err
		}
		// truncate the USD value to int64, so if the usd value is smaller than 1U,
		// the returned value is 0.
		ret = append(ret, optedUSDValues.ActiveUSDValue.TruncateInt64())
	}
	return ret, nil
}

func (k *Keeper) GetOperatorAssetValue(ctx sdk.Context, operator sdk.AccAddress, chainID string) (int64, error) {
	avsAddr := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
	optedUSDValues, err := k.GetOperatorOptedUSDValue(ctx, operator.String(), avsAddr)
	if err != nil {
		return 0, err
	}
	// truncate the USD value to int64
	return optedUSDValues.ActiveUSDValue.TruncateInt64(), nil
}
