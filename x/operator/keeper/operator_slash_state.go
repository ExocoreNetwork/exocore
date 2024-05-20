package keeper

import (
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// UpdateOperatorSlashInfo This is a function to store the slash info related to an operator
// The stored state is: operator + '/' + AVSAddr + '/' + slashId -> OperatorSlashInfo
// Now this function will be called by `slash` function implemented in 'state_update.go' when there is a slash event occurs.
func (k *Keeper) UpdateOperatorSlashInfo(ctx sdk.Context, operatorAddr, avsAddr, slashID string, slashInfo operatortypes.OperatorSlashInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)

	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return assetstype.ErrInvalidOperatorAddr
	}
	slashInfoKey := assetstype.GetJoinedStoreKey(operatorAddr, avsAddr, slashID)
	if store.Has(slashInfoKey) {
		return errorsmod.Wrap(operatortypes.ErrSlashInfoExist, fmt.Sprintf("slashInfoKey:%s", slashInfoKey))
	}
	// check the validation of slash info
	getSlashContract, err := k.avsKeeper.GetAVSSlashContract(ctx, avsAddr)
	if err != nil {
		return err
	}
	if slashInfo.SlashContract != getSlashContract {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err slashContract:%s, stored contract:%s", slashInfo.SlashContract, getSlashContract))
	}
	if slashInfo.EventHeight > slashInfo.SubmittedHeight {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err SubmittedHeight:%v,EventHeight:%v", slashInfo.SubmittedHeight, slashInfo.EventHeight))
	}

	if slashInfo.SlashProportion.IsNil() || slashInfo.SlashProportion.IsNegative() || slashInfo.SlashProportion.GT(sdkmath.LegacyNewDec(1)) {
		return errorsmod.Wrap(operatortypes.ErrSlashInfo, fmt.Sprintf("err SlashProportion:%v", slashInfo.SlashProportion))
	}

	// save single operator delegation state
	bz := k.cdc.MustMarshal(&slashInfo)
	store.Set(slashInfoKey, bz)
	return nil
}

// GetOperatorSlashInfo This is a function to retrieve the slash info related to an operator
// Now this function hasn't been called. In the future, it might be called by the grpc query.
// Additionally, it might be used when implementing the veto function
func (k *Keeper) GetOperatorSlashInfo(ctx sdk.Context, avsAddr, operatorAddr, slashID string) (changeState *operatortypes.OperatorSlashInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)
	slashInfoKey := assetstype.GetJoinedStoreKey(operatorAddr, avsAddr, slashID)
	value := store.Get(slashInfoKey)
	if value == nil {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorSlashInfo: key is %s", slashInfoKey))
	}
	operatorSlashInfo := operatortypes.OperatorSlashInfo{}
	k.cdc.MustUnmarshal(value, &operatorSlashInfo)
	return &operatorSlashInfo, nil
}

// AllOperatorSlashInfo return all slash information for the specified operator and AVS
func (k *Keeper) AllOperatorSlashInfo(ctx sdk.Context, avsAddr, operatorAddr string) (map[string]*operatortypes.OperatorSlashInfo, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorSlashInfo)
	prefix := assetstype.GetJoinedStoreKey(operatorAddr, avsAddr)

	ret := make(map[string]*operatortypes.OperatorSlashInfo, 0)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var slashInfo operatortypes.OperatorSlashInfo
		k.cdc.MustUnmarshal(iterator.Value(), &slashInfo)
		keys, err := assetstype.ParseJoinedKey(iterator.Key())
		if err != nil {
			return nil, err
		}
		ret[keys[2]] = &slashInfo
	}
	return ret, nil
}

// UpdateSlashAssetsState This is a function to update the assets amount that need to be slashed
// The stored state is:
// KeyPrefixSlashAssetsState key-value:
// processedSlashHeight + '/' + assetID -> SlashAmount
// processedSlashHeight + '/' + assetID + '/' + stakerID -> SlashAmount
// processedSlashHeight + '/' + assetID + '/' + operatorAddr -> SlashAmount
// The slashed assets info won't be sent to the client chain immediately after the slash event being processed, env if
// the asset amounts of related operator and staker have been decreased. This is because we need to wait a veto period.
// The state updated by this function will be sent to the client chain once the veto period has expired.
// This function will be called by `SlashStaker` and `SlashOperator` implemented in the 'state_update.go' file.
func (k *Keeper) UpdateSlashAssetsState(ctx sdk.Context, assetID, stakerOrOperator string, processedHeight uint64, opAmount sdkmath.Int) error {
	if opAmount.IsNil() || opAmount.IsZero() {
		return nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixSlashAssetsState)
	var key []byte
	if stakerOrOperator == "" || assetID == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, fmt.Sprintf("assetID:%s,stakerOrOperator:%s", assetID, stakerOrOperator))
	}

	key = assetstype.GetJoinedStoreKey(hexutil.EncodeUint64(processedHeight), assetID, stakerOrOperator)
	slashAmount := assetstype.ValueField{Amount: sdkmath.NewInt(0)}
	value := store.Get(key)
	if value != nil {
		k.cdc.MustUnmarshal(value, &slashAmount)
	}

	err := assetstype.UpdateAssetValue(&slashAmount.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz := k.cdc.MustMarshal(&slashAmount)
	store.Set(key, bz)

	key = assetstype.GetJoinedStoreKey(hexutil.EncodeUint64(processedHeight), assetID)
	totalSlashAmount := assetstype.ValueField{Amount: sdkmath.NewInt(0)}
	value = store.Get(key)
	if value != nil {
		k.cdc.MustUnmarshal(value, &totalSlashAmount)
	}

	err = assetstype.UpdateAssetValue(&totalSlashAmount.Amount, &opAmount)
	if err != nil {
		return err
	}
	bz = k.cdc.MustMarshal(&slashAmount)
	store.Set(key, bz)
	return nil
}

// GetSlashAssetsState This is a function to retrieve the assets awaiting transfer to the client chain for slashing.
// Now this function hasn't been called, it might be called by the grpc query in the future.
// Additionally, this function might be called in the schedule function `EndBlock` to send the slash info to client chain.
// todo: It's to be determined about how to send the slash info to client chain. If we send them in `EndBlock`, then the native code needs to call the gateway contract deployed in exocore. This seems a little bit odd.
func (k *Keeper) GetSlashAssetsState(ctx sdk.Context, assetID, stakerOrOperator string, processedHeight uint64) (sdkmath.Int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixSlashAssetsState)
	var key []byte
	if stakerOrOperator == "" {
		key = assetstype.GetJoinedStoreKey(hexutil.EncodeUint64(processedHeight), assetID)
	} else {
		key = assetstype.GetJoinedStoreKey(hexutil.EncodeUint64(processedHeight), assetID, stakerOrOperator)
	}
	value := store.Get(key)
	if value == nil {
		return sdkmath.Int{}, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetSlashAssetsState: key is %s", key))
	}
	var ret assetstype.ValueField
	k.cdc.MustUnmarshal(value, &ret)

	return ret.Amount, nil
}
