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

func (k *Keeper) UpdateStateForAsset(ctx sdk.Context, assetID, avsAddr, operatorAddr string, changeState operatortypes.OptedInAssetStateChange) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	if changeState.ChangeForAmount.IsNil() && changeState.ChangeForValue.IsNil() {
		return nil
	}
	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return assetstype.ErrOperatorAddr
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

	err = assetstype.UpdateAssetValue(&optedInAssetState.Amount, &changeState.ChangeForAmount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStateForAsset OptedInAssetState.Amount error")
	}

	err = assetstype.UpdateAssetDecValue(&optedInAssetState.Value, &changeState.ChangeForValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStateForAsset OptedInAssetState.Value error")
	}

	// save single operator delegation state
	bz := k.cdc.MustMarshal(&optedInAssetState)
	store.Set(stateKey, bz)
	return nil
}

func (k *Keeper) DeleteAssetState(ctx sdk.Context, assetID, avsAddr, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return assetstype.ErrOperatorAddr
	}
	stateKey := assetstype.GetJoinedStoreKey(assetID, avsAddr, operatorAddr)
	store.Delete(stateKey)
	return nil
}

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

func (k *Keeper) DeleteStakerShare(ctx sdk.Context, avsAddr, stakerID, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	key := assetstype.GetJoinedStoreKey(avsAddr, stakerID, operatorAddr)
	store.Delete(key)
	return nil
}

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

func (k *Keeper) GetStakerByAVSOperator(ctx sdk.Context, _, _ string) (map[string]interface{}, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	stakers := make(map[string]interface{}, 0)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keys, err := assetstype.ParseJoinedStoreKey(iterator.Key(), 3)
		if err != nil {
			return nil, err
		}
		if keys[1] != "" {
			stakers[keys[1]] = nil
		}
	}
	return stakers, nil
}
