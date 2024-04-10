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

// UpdateOperatorShare is a function to update the USD share for specified operator and Avs,
// The key and value that will be changed is:
// AVSAddr + '/' + operatorAddr -> types.DecValueField (the total USD share of specified operator and Avs)
// This function will be called when some assets supported by Avs are delegated/undelegated or slashed.
func (k *Keeper) UpdateOperatorShare(ctx sdk.Context, avsAddr, operatorAddr string, opAmount sdkmath.LegacyDec) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var key []byte
	if operatorAddr == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, "UpdateOperatorShare the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)

	totalValue := operatortypes.DecValueField{Amount: sdkmath.LegacyNewDec(0)}
	if store.Has(key) {
		value := store.Get(key)
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

// DeleteOperatorShare is a function to delete the USD share related to specified operator and Avs,
// The key and value that will be deleted is:
// AVSAddr + '/' + operatorAddr -> types.DecValueField (the total USD share of specified operator and Avs)
// This function will be called when the operator opts out of the AVS, because the USD share
// doesn't need to be stored.
func (k *Keeper) DeleteOperatorShare(ctx sdk.Context, avsAddr, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var key []byte
	if operatorAddr == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, "UpdateOperatorShare the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)

	store.Delete(key)
	return nil
}

// GetOperatorShare is a function to retrieve the USD share of specified operator and Avs,
// The key and value to retrieve is:
// AVSAddr + '/' + operatorAddr -> types.DecValueField (the total USD share of specified operator and Avs)
// This function will be called when the operator opts out of the AVS, because the total USD share
// of Avs should decrease the USD share of the opted-out operator
// This function can also serve as an RPC in the future.
func (k *Keeper) GetOperatorShare(ctx sdk.Context, avsAddr, operatorAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var ret operatortypes.DecValueField
	var key []byte
	if operatorAddr == "" {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrParameterInvalid, "GetOperatorShare the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)

	isExist := store.Has(key)
	if !isExist {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorShare: key is %suite", key))
	}
	value := store.Get(key)
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}

// UpdateAVSShare is a function to update the total USD share of an Avs,
// The key and value that will be changed is:
// AVSAddr -> types.DecValueField（the total USD share of specified Avs）
// This function will be called when some assets of operator supported by the specified Avs
// are delegated/undelegated or slashed. Additionally, when an operator opts out of
// the Avs, this function also will be called.
func (k *Keeper) UpdateAVSShare(ctx sdk.Context, avsAddr string, opAmount sdkmath.LegacyDec) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	key := []byte(avsAddr)
	totalValue := operatortypes.DecValueField{Amount: sdkmath.LegacyNewDec(0)}
	if store.Has(key) {
		value := store.Get(key)
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

// BatchUpdateShareForAVSAndOperator is a function to update the USD share for operator and Avs in bulk,
// The key and value that will be changed is:
// AVSAddr -> types.DecValueField（the total USD share of specified Avs）
// AVSAddr + '/' + operatorAddr -> types.DecValueField (the total USD share of specified operator and Avs)
// This function will be called when the prices of assets supported by Avs are changed.
func (k *Keeper) BatchUpdateShareForAVSAndOperator(ctx sdk.Context, avsOperatorChange map[string]sdkmath.LegacyDec) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	for avs, opAmount := range avsOperatorChange {
		key := []byte(avs)
		totalValue := operatortypes.DecValueField{Amount: sdkmath.LegacyNewDec(0)}
		if store.Has(key) {
			value := store.Get(key)
			k.cdc.MustUnmarshal(value, &totalValue)
		}
		tmpOpAmount := opAmount
		err := assetstype.UpdateAssetDecValue(&totalValue.Amount, &tmpOpAmount)
		if err != nil {
			return err
		}
		bz := k.cdc.MustMarshal(&totalValue)
		store.Set(key, bz)
	}
	return nil
}

// GetAVSShare is a function to retrieve the USD share of specified Avs,
// The key and value to retrieve is:
// AVSAddr -> types.DecValueField（the total USD share of specified Avs）
// It hasn't been used now. but it can serve as an RPC in the future.
func (k *Keeper) GetAVSShare(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var ret operatortypes.DecValueField
	key := []byte(avsAddr)
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSShare: key is %suite", key))
	}
	value := store.Get(key)
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}

// UpdateStateForAsset is a function to update the opted-in amount and USD share for
// the specified asset
// The key and value that will be changed is:
// assetID + '/' + AVSAddr + '/' + operatorAddr -> types.OptedInAssetState
// This function will be called when the amount of a specified asset opted-in by the operator
// changes, such as: opt-in, delegation, undelegation and slash.
func (k *Keeper) UpdateStateForAsset(ctx sdk.Context, assetID, avsAddr, operatorAddr string, changeState operatortypes.DeltaOptedInAssetState) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	if changeState.Amount.IsNil() && changeState.Value.IsNil() {
		return nil
	}
	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return assetstype.ErrInvalidOperatorAddr
	}
	stateKey := assetstype.GetJoinedStoreKey(assetID, avsAddr, operatorAddr)
	optedInAssetState := operatortypes.OptedInAssetState{
		Amount: sdkmath.NewInt(0),
		Value:  sdkmath.LegacyNewDec(0),
	}

	if store.Has(stateKey) {
		value := store.Get(stateKey)
		k.cdc.MustUnmarshal(value, &optedInAssetState)
	}

	err = assetstype.UpdateAssetValue(&optedInAssetState.Amount, &changeState.Amount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStateForAsset OptedInAssetState.Amount error")
	}

	err = assetstype.UpdateAssetDecValue(&optedInAssetState.Value, &changeState.Value)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStateForAsset OptedInAssetState.Value error")
	}

	// save single operator delegation state
	bz := k.cdc.MustMarshal(&optedInAssetState)
	store.Set(stateKey, bz)
	return nil
}

// DeleteAssetState is a function to delete the opted-in amount and USD share for
// the specified asset
// The key and value that will be deleted is:
// assetID + '/' + AVSAddr + '/' + operatorAddr -> types.OptedInAssetState
// This function will be called when the specified operator opts out of the Avs.
func (k *Keeper) DeleteAssetState(ctx sdk.Context, assetID, avsAddr, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return assetstype.ErrInvalidOperatorAddr
	}
	stateKey := assetstype.GetJoinedStoreKey(assetID, avsAddr, operatorAddr)
	store.Delete(stateKey)
	return nil
}

// GetAssetState is a function to retrieve the opted-in amount and USD share for the specified asset
// The key and value to retrieve is:
// assetID + '/' + AVSAddr + '/' + operatorAddr -> types.OptedInAssetState
// It hasn't been used now. but it can serve as an RPC in the future.
func (k *Keeper) GetAssetState(ctx sdk.Context, assetID, avsAddr, operatorAddr string) (changeState *operatortypes.OptedInAssetState, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	stateKey := assetstype.GetJoinedStoreKey(assetID, avsAddr, operatorAddr)
	isExit := store.Has(stateKey)
	optedInAssetState := operatortypes.OptedInAssetState{}
	if isExit {
		value := store.Get(stateKey)
		k.cdc.MustUnmarshal(value, &optedInAssetState)
	} else {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAssetState: key is %suite", stateKey))
	}
	return &optedInAssetState, nil
}

// IterateUpdateAssetState is a function to iteratively update the opted-in amount and USD share for
// the specified asset
// The key and value that will be changed is:
// assetID + '/' + AVSAddr + '/' + operatorAddr -> types.OptedInAssetState
// This function will be called when the prices of opted-in assets are changed.
func (k *Keeper) IterateUpdateAssetState(ctx sdk.Context, assetID string, f func(assetID string, keys []string, state *operatortypes.OptedInAssetState) error) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	iterator := sdk.KVStorePrefixIterator(store, []byte(assetID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keys, err := assetstype.ParseJoinedStoreKey(iterator.Key(), 3)
		if err != nil {
			return err
		}
		optedInAssetState := &operatortypes.OptedInAssetState{}
		k.cdc.MustUnmarshal(iterator.Value(), optedInAssetState)
		err = f(assetID, keys, optedInAssetState)
		if err != nil {
			return err
		}
		bz := k.cdc.MustMarshal(optedInAssetState)
		store.Set(iterator.Key(), bz)
	}
	return nil
}

// UpdateStakerShare is a function to update the opted-in USD share for the specified staker and operator ,
// The key and value that will be changed is:
// AVSAddr + '/' + ” + '/' +  operatorAddr -> types.DecValueField（the opted-in USD share owned by the operator itself）
// AVSAddr + '/' + stakerID + '/' + operatorAddr -> types.DecValueField (the opted-in USD share of the staker)
// This function will be called when the opted-in assets of operator and staker
// are delegated/undelegated or slashed. Additionally, when an operator opts in, this function also will be called.
func (k *Keeper) UpdateStakerShare(ctx sdk.Context, avsAddr, stakerID, operatorAddr string, opAmount sdkmath.LegacyDec) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	key := assetstype.GetJoinedStoreKey(avsAddr, stakerID, operatorAddr)

	optedInValue := operatortypes.DecValueField{Amount: sdkmath.LegacyNewDec(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &optedInValue)
	}
	err := assetstype.UpdateAssetDecValue(&optedInValue.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&optedInValue)
	store.Set(key, bz)
	return nil
}

// BatchSetStakerShare is a function to set the opted-in USD share for the specified staker and operator in bulk,
// The key and value that will be set is:
// AVSAddr + '/' + ” + '/' +  operatorAddr -> types.DecValueField（the opted-in USD share owned by the operator itself）
// AVSAddr + '/' + stakerID + '/' + operatorAddr -> types.DecValueField (the opted-in USD share of the staker)
// This function will be called when the prices of opted-in assets are changed.
func (k *Keeper) BatchSetStakerShare(ctx sdk.Context, newValues map[string]sdkmath.LegacyDec) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	for key, value := range newValues {
		optedInValue := operatortypes.DecValueField{Amount: value}
		if store.Has([]byte(key)) {
			value := store.Get([]byte(key))
			k.cdc.MustUnmarshal(value, &optedInValue)
		}

		bz := k.cdc.MustMarshal(&optedInValue)
		store.Set([]byte(key), bz)
	}
	return nil
}

// DeleteStakerShare is a function to delete the opted-in USD share for the specified staker and operator,
// The key and value that will be set is:
// AVSAddr + '/' + ” + '/' +  operatorAddr -> types.DecValueField（the opted-in USD share owned by the operator itself）
// AVSAddr + '/' + stakerID + '/' + operatorAddr -> types.DecValueField (the opted-in USD share of the staker)
// This function will be called when the operator opts out of the Avs.
func (k *Keeper) DeleteStakerShare(ctx sdk.Context, avsAddr, stakerID, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	key := assetstype.GetJoinedStoreKey(avsAddr, stakerID, operatorAddr)
	store.Delete(key)
	return nil
}

// GetStakerShare is a function to retrieve the opted-in USD share for the specified staker and operator,
// The key and value that will be set is:
// AVSAddr + '/' + ” + '/' +  operatorAddr -> types.DecValueField（the opted-in USD share owned by the operator itself）
// AVSAddr + '/' + stakerID + '/' + operatorAddr -> types.DecValueField (the opted-in USD share of the staker)
// It hasn't been used now. but it can serve as an RPC in the future.
func (k *Keeper) GetStakerShare(ctx sdk.Context, avsAddr, stakerID, operatorAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	var ret operatortypes.DecValueField
	key := assetstype.GetJoinedStoreKey(avsAddr, stakerID, operatorAddr)
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetStakerShare: key is %s", key))
	}
	value := store.Get(key)
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}

func (k *Keeper) GetAvgDelegatedValue(
	ctx sdk.Context, operators []sdk.AccAddress, chainID, _ string,
) ([]int64, error) {
	avsAddr, err := k.avsKeeper.GetAvsAddrByChainID(ctx, chainID)
	if err != nil {
		return nil, err
	}
	ret := make([]int64, 0)
	for _, operator := range operators {
		share, err := k.GetOperatorShare(ctx, operator.String(), avsAddr)
		if err != nil {
			return nil, err
		}
		// truncate the USD value to int64
		ret = append(ret, share.TruncateInt64())
	}
	return ret, nil
}

func (k *Keeper) GetOperatorAssetValue(ctx sdk.Context, operator sdk.AccAddress, chainID string) (int64, error) {
	avsAddr, err := k.avsKeeper.GetAvsAddrByChainID(ctx, chainID)
	if err != nil {
		return 0, err
	}
	share, err := k.GetOperatorShare(ctx, operator.String(), avsAddr)
	if err != nil {
		return 0, err
	}
	// truncate the USD value to int64
	return share.TruncateInt64(), nil
}
