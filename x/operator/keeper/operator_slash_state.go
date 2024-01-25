package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	operatortypes "github.com/exocore/x/operator/types"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) UpdateOperatorSlashInfo(ctx sdk.Context, operatorAddr, avsAddr, slashId string, slashInfo operatortypes.OperatorSlashInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)

	//check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return restakingtype.OperatorAddrIsNotAccAddr
	}
	slashInfoKey := restakingtype.GetJoinedStoreKey(operatorAddr, avsAddr, slashId)
	if store.Has(slashInfoKey) {
		return errorsmod.Wrap(operatortypes.ErrSlashInfoExist, fmt.Sprintf("slashInfoKey:%s", slashInfoKey))
	}
	// check the validation of slash info
	if slashInfo.SlashContract == "" {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err slashContract:%s", slashInfo.SlashContract))
	}
	if slashInfo.OccurredHeight > slashInfo.SlashHeight {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err SlashHeight:%v,OccurredHeight:%v", slashInfo.SlashHeight, slashInfo.OccurredHeight))
	}

	if slashInfo.SlashProportion.IsNil() || slashInfo.SlashProportion.IsNegative() || slashInfo.SlashProportion.GT(sdkmath.LegacyNewDec(1)) {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err SlashProportion:%v", slashInfo.SlashProportion))
	}

	//save single operator delegation state
	bz := k.cdc.MustMarshal(&slashInfo)
	store.Set(slashInfoKey, bz)
	return nil
}

func (k Keeper) GetOperatorSlashInfo(ctx sdk.Context, avsAddr, operatorAddr, slashId string) (changeState *operatortypes.OperatorSlashInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)
	slashInfoKey := restakingtype.GetJoinedStoreKey(operatorAddr, avsAddr, slashId)
	isExit := store.Has(slashInfoKey)
	operatorSlashInfo := operatortypes.OperatorSlashInfo{}
	if isExit {
		value := store.Get(slashInfoKey)
		k.cdc.MustUnmarshal(value, &operatorSlashInfo)
	} else {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorSlashInfo: key is %s", slashInfoKey))
	}
	return &operatorSlashInfo, nil
}

func (k Keeper) UpdateSlashAssetsState(ctx sdk.Context, assetId, stakerOrOperator string, completeHeight uint64, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixSlashAssetsState)
	var key []byte
	if stakerOrOperator == "" || assetId == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, fmt.Sprintf("assetId:%s,stakerOrOperator:%s", assetId, stakerOrOperator))
	}

	key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetId, stakerOrOperator)
	slashAmount := restakingtype.ValueField{Amount: sdkmath.NewInt(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &slashAmount)
	}
	err := restakingtype.UpdateAssetValue(&slashAmount.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&slashAmount)
	store.Set(key, bz)

	key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetId)
	totalSlashAmount := restakingtype.ValueField{Amount: sdkmath.NewInt(0)}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &totalSlashAmount)
	}
	err = restakingtype.UpdateAssetValue(&totalSlashAmount.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz = k.cdc.MustMarshal(&slashAmount)
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetSlashAssetsState(ctx sdk.Context, assetId, stakerOrOperator string, completeHeight uint64) (sdkmath.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixSlashAssetsState)
	var key []byte
	if stakerOrOperator == "" {
		key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetId)
	} else {
		key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetId, stakerOrOperator)
	}
	var ret restakingtype.ValueField
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.Int{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetSlashAssetsState: key is %s", key))
	} else {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &ret)
	}
	return ret.Amount, nil
}
