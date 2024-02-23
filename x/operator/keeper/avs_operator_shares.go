package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
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
	} else {
		key = restakingtype.GetJoinedStoreKey(avsAddr, operatorAddr)
	}

	totalValue := operatortypes.ValueField{Amount: sdkmath.LegacyNewDec(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &totalValue)
	}
	err := restakingtype.UpdateAssetDecValue(&totalValue.Amount, &opAmount)
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
	} else {
		key = restakingtype.GetJoinedStoreKey(avsAddr, operatorAddr)
	}
	store.Delete(key)
	return nil
}

func (k *Keeper) GetOperatorShare(ctx sdk.Context, avsAddr, operatorAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var ret operatortypes.ValueField
	var key []byte
	if operatorAddr == "" {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrParameterInvalid, "GetOperatorShare the operatorAddr is empty")
	} else {
		key = restakingtype.GetJoinedStoreKey(avsAddr, operatorAddr)
	}
	isExist := store.Has(key)
	if !isExist {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorShare: key is %suite", key))
	} else {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}

func (k *Keeper) UpdateAVSShare(ctx sdk.Context, avsAddr string, opAmount sdkmath.LegacyDec) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	key := []byte(avsAddr)
	totalValue := operatortypes.ValueField{Amount: sdkmath.LegacyNewDec(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &totalValue)
	}
	err := restakingtype.UpdateAssetDecValue(&totalValue.Amount, &opAmount)
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
		totalValue := operatortypes.ValueField{Amount: sdkmath.LegacyNewDec(0)}
		if store.Has(key) {
			value := store.Get(key)
			k.cdc.MustUnmarshal(value, &totalValue)
		}
		err := restakingtype.UpdateAssetDecValue(&totalValue.Amount, &opAmount)
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
	var ret operatortypes.ValueField
	key := []byte(avsAddr)
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSShare: key is %suite", key))
	} else {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}

func (k *Keeper) UpdateStateForAsset(ctx sdk.Context, assetId, avsAddr, operatorAddr string, changeState operatortypes.AssetOptedInState) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	if changeState.Amount.IsNil() && changeState.Value.IsNil() {
		return nil
	}
	//check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return restakingtype.ErrOperatorAddr
	}
	stateKey := restakingtype.GetJoinedStoreKey(assetId, avsAddr, operatorAddr)
	assetOptedInState := operatortypes.AssetOptedInState{
		Amount: sdkmath.NewInt(0),
		Value:  sdkmath.LegacyNewDec(0),
	}

	if store.Has(stateKey) {
		value := store.Get(stateKey)
		k.cdc.MustUnmarshal(value, &assetOptedInState)
	}

	err = restakingtype.UpdateAssetValue(&assetOptedInState.Amount, &changeState.Amount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStateForAsset assetOptedInState.Amount error")
	}

	err = restakingtype.UpdateAssetDecValue(&assetOptedInState.Value, &changeState.Value)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStateForAsset assetOptedInState.Value error")
	}

	//save single operator delegation state
	bz := k.cdc.MustMarshal(&assetOptedInState)
	store.Set(stateKey, bz)
	return nil
}

func (k *Keeper) DeleteAssetState(ctx sdk.Context, assetId, avsAddr, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	//check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return restakingtype.ErrOperatorAddr
	}
	stateKey := restakingtype.GetJoinedStoreKey(assetId, avsAddr, operatorAddr)
	store.Delete(stateKey)
	return nil
}

func (k *Keeper) GetAssetState(ctx sdk.Context, assetId, avsAddr, operatorAddr string) (changeState *operatortypes.AssetOptedInState, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	stateKey := restakingtype.GetJoinedStoreKey(assetId, avsAddr, operatorAddr)
	isExit := store.Has(stateKey)
	assetOptedInState := operatortypes.AssetOptedInState{}
	if isExit {
		value := store.Get(stateKey)
		k.cdc.MustUnmarshal(value, &assetOptedInState)
	} else {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAssetState: key is %suite", stateKey))
	}
	return &assetOptedInState, nil
}

func (k *Keeper) IterateUpdateAssetState(ctx sdk.Context, assetId string, f func(assetId string, keys []string, state *operatortypes.AssetOptedInState) error) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	iterator := sdk.KVStorePrefixIterator(store, []byte(assetId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keys, err := restakingtype.ParseJoinedStoreKey(iterator.Key(), 3)
		if err != nil {
			return err
		}
		assetOptedInState := &operatortypes.AssetOptedInState{}
		k.cdc.MustUnmarshal(iterator.Value(), assetOptedInState)
		err = f(assetId, keys, assetOptedInState)
		if err != nil {
			return err
		}
		bz := k.cdc.MustMarshal(assetOptedInState)
		store.Set(iterator.Key(), bz)
	}
	return nil
}

func (k *Keeper) UpdateStakerShare(ctx sdk.Context, avsAddr, stakerId, operatorAddr string, opAmount sdkmath.LegacyDec) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	key := restakingtype.GetJoinedStoreKey(avsAddr, stakerId, operatorAddr)

	optedInValue := operatortypes.ValueField{Amount: sdkmath.LegacyNewDec(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &optedInValue)
	}
	err := restakingtype.UpdateAssetDecValue(&optedInValue.Amount, &opAmount)
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
		optedInValue := operatortypes.ValueField{Amount: value}
		if store.Has([]byte(key)) {
			value := store.Get([]byte(key))
			k.cdc.MustUnmarshal(value, &optedInValue)
		}

		bz := k.cdc.MustMarshal(&optedInValue)
		store.Set([]byte(key), bz)
	}
	return nil
}

func (k *Keeper) DeleteStakerShare(ctx sdk.Context, avsAddr, stakerId, operatorAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	key := restakingtype.GetJoinedStoreKey(avsAddr, stakerId, operatorAddr)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetStakerShare(ctx sdk.Context, avsAddr, stakerId, operatorAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	var ret operatortypes.ValueField
	key := restakingtype.GetJoinedStoreKey(avsAddr, stakerId, operatorAddr)
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetStakerShare: key is %s", key))
	} else {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}

func (k *Keeper) GetStakerByAVSOperator(ctx sdk.Context, avsAddr, operatorAddr string) (map[string]interface{}, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	stakers := make(map[string]interface{}, 0)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		keys, err := restakingtype.ParseJoinedStoreKey(iterator.Key(), 3)
		if err != nil {
			return nil, err
		}
		if keys[1] != "" {
			stakers[keys[1]] = nil
		}
	}
	return stakers, nil
}
