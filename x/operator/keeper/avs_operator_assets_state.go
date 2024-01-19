package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	operatortypes "github.com/exocore/x/operator/types"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) UpdateAVSOperatorTotalValue(ctx sdk.Context, avsAddr, operatorAddr string, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var key []byte
	if operatorAddr == "" {
		key = []byte(avsAddr)
	} else {
		key = restakingtype.GetJoinedStoreKey(avsAddr, operatorAddr)
	}

	totalValue := restakingtype.ValueField{Amount: sdkmath.NewInt(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &totalValue)
	}
	err := restakingtype.UpdateAssetValue(&totalValue.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&totalValue)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetAVSOperatorTotalValue(ctx sdk.Context, avsAddr, operatorAddr string) (sdkmath.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorAssetsTotalValue)
	var ret restakingtype.ValueField
	var key []byte
	if operatorAddr == "" {
		key = []byte(avsAddr)
	} else {
		key = restakingtype.GetJoinedStoreKey(avsAddr, operatorAddr)
	}
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.Int{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSOperatorTotalValue: key is %s", key))
	} else {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}

func (k Keeper) UpdateOperatorAVSAssetsState(ctx sdk.Context, assetId, avsAddr, operatorAddr string, changeState operatortypes.AssetOptedInState) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)

	if changeState.Amount.IsNil() && changeState.Value.IsNil() {
		return nil
	}
	//check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return restakingtype.OperatorAddrIsNotAccAddr
	}
	stateKey := restakingtype.GetJoinedStoreKey(assetId, avsAddr, operatorAddr)
	assetOptedInState := operatortypes.AssetOptedInState{
		Amount: sdkmath.NewInt(0),
		Value:  sdkmath.NewInt(0),
	}

	if store.Has(stateKey) {
		value := store.Get(stateKey)
		k.cdc.MustUnmarshal(value, &assetOptedInState)
	}

	err = restakingtype.UpdateAssetValue(&assetOptedInState.Amount, &changeState.Amount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAVSAssetsState assetOptedInState.Amount error")
	}

	err = restakingtype.UpdateAssetValue(&assetOptedInState.Value, &changeState.Value)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateOperatorAVSAssetsState assetOptedInState.Value error")
	}

	//save single operator delegation state
	bz := k.cdc.MustMarshal(&assetOptedInState)
	store.Set(stateKey, bz)
	return nil
}

func (k Keeper) GetOperatorAVSAssetsState(ctx sdk.Context, assetId, avsAddr, operatorAddr string) (changeState *operatortypes.AssetOptedInState, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorAVSSingleAssetState)
	stateKey := restakingtype.GetJoinedStoreKey(assetId, avsAddr, operatorAddr)
	isExit := store.Has(stateKey)
	assetOptedInState := operatortypes.AssetOptedInState{}
	if isExit {
		value := store.Get(stateKey)
		k.cdc.MustUnmarshal(value, &assetOptedInState)
	} else {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorAVSAssetsState: key is %s", stateKey))
	}
	return &assetOptedInState, nil
}

func (k Keeper) UpdateAVSOperatorStakerShareValue(ctx sdk.Context, avsAddr, stakerId, operatorAddr string, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	key := restakingtype.GetJoinedStoreKey(avsAddr, stakerId, operatorAddr)

	optedInValue := restakingtype.ValueField{Amount: sdkmath.NewInt(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &optedInValue)
	}
	err := restakingtype.UpdateAssetValue(&optedInValue.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&optedInValue)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetAVSOperatorStakerShareValue(ctx sdk.Context, avsAddr, stakerId, operatorAddr string) (sdkmath.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixAVSOperatorStakerShareState)
	var ret restakingtype.ValueField
	key := restakingtype.GetJoinedStoreKey(avsAddr, stakerId, operatorAddr)
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.Int{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSOperatorStakerShareValue: key is %s", key))
	} else {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}
