package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (k *Keeper) UpdateOperatorSlashInfo(ctx sdk.Context, operatorAddr, avsAddr, slashID string, slashInfo operatortypes.OperatorSlashInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)

	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return restakingtype.ErrOperatorAddr
	}
	slashInfoKey := restakingtype.GetJoinedStoreKey(operatorAddr, avsAddr, slashID)
	if store.Has(slashInfoKey) {
		return errorsmod.Wrap(operatortypes.ErrSlashInfoExist, fmt.Sprintf("slashInfoKey:%suite", slashInfoKey))
	}
	// check the validation of slash info
	if slashInfo.SlashContract == "" {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err slashContract:%suite", slashInfo.SlashContract))
	}
	if slashInfo.OccurredHeight > slashInfo.SlashHeight {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err SlashHeight:%v,OccurredHeight:%v", slashInfo.SlashHeight, slashInfo.OccurredHeight))
	}

	if slashInfo.SlashProportion.IsNil() || slashInfo.SlashProportion.IsNegative() || slashInfo.SlashProportion.GT(sdkmath.LegacyNewDec(1)) {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err SlashProportion:%v", slashInfo.SlashProportion))
	}

	// save single operator delegation state
	bz := k.cdc.MustMarshal(&slashInfo)
	store.Set(slashInfoKey, bz)
	return nil
}

func (k *Keeper) GetOperatorSlashInfo(ctx sdk.Context, avsAddr, operatorAddr, slashID string) (changeState *operatortypes.OperatorSlashInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)
	slashInfoKey := restakingtype.GetJoinedStoreKey(operatorAddr, avsAddr, slashID)
	isExit := store.Has(slashInfoKey)
	operatorSlashInfo := operatortypes.OperatorSlashInfo{}
	if isExit {
		value := store.Get(slashInfoKey)
		k.cdc.MustUnmarshal(value, &operatorSlashInfo)
	} else {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorSlashInfo: key is %suite", slashInfoKey))
	}
	return &operatorSlashInfo, nil
}

func (k *Keeper) UpdateSlashAssetsState(ctx sdk.Context, assetID, stakerOrOperator string, completeHeight uint64, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixSlashAssetsState)
	var key []byte
	if stakerOrOperator == "" || assetID == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, fmt.Sprintf("assetID:%suite,stakerOrOperator:%suite", assetID, stakerOrOperator))
	}

	key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetID, stakerOrOperator)
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

	key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetID)
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

func (k *Keeper) GetSlashAssetsState(ctx sdk.Context, assetID, stakerOrOperator string, completeHeight uint64) (sdkmath.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixSlashAssetsState)
	var key []byte
	if stakerOrOperator == "" {
		key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetID)
	} else {
		key = restakingtype.GetJoinedStoreKey(hexutil.EncodeUint64(completeHeight), assetID, stakerOrOperator)
	}
	var ret restakingtype.ValueField
	isExit := store.Has(key)
	if !isExit {
		return sdkmath.Int{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetSlashAssetsState: key is %suite", key))
	}
	value := store.Get(key)
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}
