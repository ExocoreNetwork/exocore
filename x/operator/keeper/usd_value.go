package keeper

import (
	"errors"
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationkeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	oracletype "github.com/ExocoreNetwork/exocore/x/oracle/types"

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
	var key []byte
	if operatorAddr == "" {
		return errorsmod.Wrap(operatortypes.ErrParameterInvalid, "UpdateOperatorUSDValue the operatorAddr is empty")
	}
	key = assetstype.GetJoinedStoreKey(avsAddr, operatorAddr)
	store.Delete(key)

	return nil
}

func (k *Keeper) DeleteAllOperatorsUSDValueForAVS(ctx sdk.Context, avsAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
	iterator := sdk.KVStorePrefixIterator(store, operatortypes.IterateOperatorsForAVSPrefix(avsAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
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

	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForAVS)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForAVS)
	key := []byte(avsAddr)
	setValue := operatortypes.DecValueField{Amount: amount}
	bz := k.cdc.MustMarshal(&setValue)
	store.Set(key, bz)
	return nil
}

func (k *Keeper) DeleteAVSUSDValue(ctx sdk.Context, avsAddr string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForAVS)
	key := []byte(avsAddr)
	store.Delete(key)
	return nil
}

// GetAVSUSDValue is a function to retrieve the USD share of specified Avs,
// The key and value to retrieve is:
// AVSAddr -> types.DecValueField（the total USD share of specified Avs）
func (k *Keeper) GetAVSUSDValue(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error) {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		operatortypes.KeyPrefixUSDValueForAVS,
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
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
	ctx sdk.Context, operators []sdk.AccAddress, chainIDWithoutRevision string,
) ([]int64, error) {
	isAvs, avsAddrString := k.avsKeeper.IsAVSByChainID(ctx, chainIDWithoutRevision)
	if !isAvs {
		return nil, errorsmod.Wrap(operatortypes.ErrUnknownChainID, fmt.Sprintf("GetVotePowerForChainID: chainIDWithoutRevision is %s", chainIDWithoutRevision))
	}
	ret := make([]int64, 0)
	for _, operator := range operators {
		// this already filters by the required assetIDs
		optedUSDValues, err := k.GetOperatorOptedUSDValue(ctx, avsAddrString, operator.String())
		if err != nil {
			return nil, err
		}
		// truncate the USD value to int64, so if the usd value is smaller than 1U,
		// the returned value is 0.
		ret = append(ret, optedUSDValues.ActiveUSDValue.TruncateInt64())
	}
	return ret, nil
}

func (k *Keeper) GetOperatorAssetValue(ctx sdk.Context, operator sdk.AccAddress, chainIDWithoutRevision string) (int64, error) {
	isAvs, avsAddr := k.avsKeeper.IsAVSByChainID(ctx, chainIDWithoutRevision)
	if !isAvs {
		return 0, errorsmod.Wrap(operatortypes.ErrUnknownChainID, fmt.Sprintf("GetOperatorAssetValue: chainIDWithoutRevision is %s", chainIDWithoutRevision))
	}
	optedUSDValues, err := k.GetOperatorOptedUSDValue(ctx, operator.String(), avsAddr)
	if err != nil {
		return 0, err
	}
	// truncate the USD value to int64
	return optedUSDValues.ActiveUSDValue.TruncateInt64(), nil
}

func (k *Keeper) SetAllOperatorUSDValues(ctx sdk.Context, usdValues []operatortypes.OperatorUSDValue) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
	for i := range usdValues {
		usdValue := usdValues[i]
		bz := k.cdc.MustMarshal(&usdValue.OptedUSDValue)
		store.Set([]byte(usdValue.Key), bz)
	}
	return nil
}

func (k *Keeper) GetAllOperatorUSDValues(ctx sdk.Context) ([]operatortypes.OperatorUSDValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForOperator)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	ret := make([]operatortypes.OperatorUSDValue, 0)
	for ; iterator.Valid(); iterator.Next() {
		var usdValues operatortypes.OperatorOptedUSDValue
		k.cdc.MustUnmarshal(iterator.Value(), &usdValues)
		ret = append(ret, operatortypes.OperatorUSDValue{
			Key:           string(iterator.Key()),
			OptedUSDValue: usdValues,
		})
	}
	return ret, nil
}

func (k *Keeper) SetAllAVSUSDValues(ctx sdk.Context, usdValues []operatortypes.AVSUSDValue) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForAVS)
	for i := range usdValues {
		usdValue := usdValues[i]
		bz := k.cdc.MustMarshal(&usdValue.Value)
		store.Set([]byte(usdValue.AVSAddr), bz)
	}
	return nil
}

func (k *Keeper) GetAllAVSUSDValues(ctx sdk.Context) ([]operatortypes.AVSUSDValue, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixUSDValueForAVS)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	ret := make([]operatortypes.AVSUSDValue, 0)
	for ; iterator.Valid(); iterator.Next() {
		var usdValue operatortypes.DecValueField
		k.cdc.MustUnmarshal(iterator.Value(), &usdValue)
		ret = append(ret, operatortypes.AVSUSDValue{
			AVSAddr: string(iterator.Key()),
			Value:   usdValue,
		})
	}
	return ret, nil
}

// CalculateUSDValueForOperator calculates the total and self usd value for the
// operator according to the input assets filter and prices.
// This function will be used in slashing calculations and voting power updates per epoch.
// The inputs/outputs and calculation logic for these two cases are different,
// so an `isForSlash` flag is used to distinguish between them.
// When it's called by the voting power update, the needed outputs are the current total
// staking amount and the self-staking amount of the operator. The current total
// staking amount excludes the pending unbonding amount, so it's used to calculate the voting power.
// The self-staking amount is also needed to check if the operator's self-staking is sufficient.
// At the same time, the prices of all assets have been retrieved in the caller's function, so they
// are inputted as a parameter.
// When it's called by the slash execution, the needed output is the sum of the current total amount and
// the pending unbonding amount, because the undelegation also needs to be slashed. And the prices of
// all assets haven't been prepared by the caller, so the prices should be retrieved in this function.
func (k *Keeper) CalculateUSDValueForOperator(
	ctx sdk.Context,
	isForSlash bool,
	operator string,
	assetsFilter map[string]interface{},
	decimals map[string]uint32,
	prices map[string]oracletype.Price,
) (operatortypes.OperatorStakingInfo, error) {
	var err error
	ret := operatortypes.OperatorStakingInfo{
		Staking:                 sdkmath.LegacyNewDec(0),
		SelfStaking:             sdkmath.LegacyNewDec(0),
		StakingAndWaitUnbonding: sdkmath.LegacyNewDec(0),
	}
	// iterate all assets owned by the operator to calculate its voting power
	opFuncToIterateAssets := func(assetID string, state *assetstype.OperatorAssetInfo) error {
		//		var price operatortypes.Price
		var price oracletype.Price
		var decimal uint32
		if isForSlash {
			// when calculated the USD value for slashing, the input prices map is null
			// so the price needs to be retrieved here
			price, err = k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetID)
			if err != nil {
				// TODO: when assetID is not registered in oracle module, this error will finally lead to panic
				if !errors.Is(err, oracletype.ErrGetPriceRoundNotFound) {
					return err
				}
				// TODO: for now, we ignore the error when the price round is not found and set the price to 1 to avoid panic
			}
			assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
			if err != nil {
				return err
			}
			decimal = assetInfo.AssetBasicInfo.Decimals
			ret.StakingAndWaitUnbonding = ret.StakingAndWaitUnbonding.Add(CalculateUSDValue(state.TotalAmount.Add(state.PendingUndelegationAmount), price.Value, decimal, price.Decimal))
		} else {
			if prices == nil {
				return errorsmod.Wrap(operatortypes.ErrValueIsNilOrZero, "CalculateUSDValueForOperator prices map is nil")
			}
			price, ok := prices[assetID]
			if !ok {
				return errorsmod.Wrap(operatortypes.ErrKeyNotExistInMap, "CalculateUSDValueForOperator map: prices, key: assetID")
			}
			decimal, ok := decimals[assetID]
			if !ok {
				return errorsmod.Wrap(operatortypes.ErrKeyNotExistInMap, "CalculateUSDValueForOperator map: decimals, key: assetID")
			}
			ret.Staking = ret.Staking.Add(CalculateUSDValue(state.TotalAmount, price.Value, decimal, price.Decimal))
			// calculate the token amount from the share for the operator
			selfAmount, err := delegationkeeper.TokensFromShares(state.OperatorShare, state.TotalShare, state.TotalAmount)
			if err != nil {
				return err
			}
			ret.SelfStaking = ret.SelfStaking.Add(CalculateUSDValue(selfAmount, price.Value, decimal, price.Decimal))
		}
		return nil
	}
	err = k.assetsKeeper.IterateAssetsForOperator(ctx, false, operator, assetsFilter, opFuncToIterateAssets)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (k Keeper) GetOrCalculateOperatorUSDValues(
	ctx sdk.Context,
	operator sdk.AccAddress,
	chainIDWithoutRevision string,
) (optedUSDValues operatortypes.OperatorOptedUSDValue, err error) {
	isAvs, avsAddrString := k.avsKeeper.IsAVSByChainID(ctx, chainIDWithoutRevision)
	if !isAvs {
		return operatortypes.OperatorOptedUSDValue{}, errorsmod.Wrap(operatortypes.ErrUnknownChainID, fmt.Sprintf("GetOrCalculateOperatorUSDValues: chainIDWithoutRevision is %s", chainIDWithoutRevision))
	}
	// the usd values will be deleted if the operator opts out, so recalculate the
	// voting power to set the tokens and shares for this case.
	if !k.IsOptedIn(ctx, operator.String(), avsAddrString) {
		// get assets supported by the AVS
		assets, err := k.avsKeeper.GetAVSSupportedAssets(ctx, avsAddrString)
		if err != nil {
			return operatortypes.OperatorOptedUSDValue{}, err
		}
		if assets == nil {
			return operatortypes.OperatorOptedUSDValue{}, err
		}
		// get the prices and decimals of assets
		decimals, err := k.assetsKeeper.GetAssetsDecimal(ctx, assets)
		if err != nil {
			return operatortypes.OperatorOptedUSDValue{}, err
		}
		prices, err := k.oracleKeeper.GetMultipleAssetsPrices(ctx, assets)
		if err != nil {
			return operatortypes.OperatorOptedUSDValue{}, err
		}
		stakingInfo, err := k.CalculateUSDValueForOperator(ctx, false, operator.String(), assets, decimals, prices)
		if err != nil {
			return operatortypes.OperatorOptedUSDValue{}, err
		}
		optedUSDValues.SelfUSDValue = stakingInfo.SelfStaking
		optedUSDValues.TotalUSDValue = stakingInfo.Staking
	} else {
		optedUSDValues, err = k.GetOperatorOptedUSDValue(ctx, avsAddrString, operator.String())
		if err != nil {
			return operatortypes.OperatorOptedUSDValue{}, err
		}
	}
	return optedUSDValues, nil
}

func (k *Keeper) CalculateUSDValueForStaker(ctx sdk.Context, stakerID, avsAddr string, operator sdk.AccAddress) (sdkmath.LegacyDec, error) {
	if !k.IsActive(ctx, operator, avsAddr) {
		return sdkmath.LegacyNewDec(0), nil
	}
	optedUSDValues, err := k.GetOperatorOptedUSDValue(ctx, avsAddr, operator.String())
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	if optedUSDValues.ActiveUSDValue.IsZero() {
		return sdkmath.LegacyNewDec(0), err
	}

	// calculate the active voting power for staker
	assets, err := k.avsKeeper.GetAVSSupportedAssets(ctx, avsAddr)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	if assets == nil {
		return sdkmath.LegacyNewDec(0), nil
	}
	prices, err := k.oracleKeeper.GetMultipleAssetsPrices(ctx, assets)
	// we don't ignore the error regarding the price round not found here, because it's used to
	// distribute the reward.
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	if prices == nil {
		return sdkmath.LegacyDec{}, errorsmod.Wrap(operatortypes.ErrValueIsNilOrZero, "CalculateUSDValueForStaker prices map is nil")
	}
	totalUSDValue := sdkmath.LegacyNewDec(0)
	opFunc := func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) (bool, error) {
		// Return true to stop iteration, false to continue iterating
		if keys.OperatorAddr == operator.String() {
			if _, ok := assets[keys.AssetID]; ok {
				price, ok := prices[keys.AssetID]
				if !ok {
					return true, errorsmod.Wrapf(operatortypes.ErrKeyNotExistInMap, "CalculateUSDValueForStaker Price not found for assetID: %s", keys.AssetID)
				}
				operatorAsset, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, operator, keys.AssetID)
				if err != nil {
					return true, err
				}
				amount, err := delegationkeeper.TokensFromShares(amounts.UndelegatableShare, operatorAsset.TotalShare, operatorAsset.TotalAmount)
				if err != nil {
					return true, err
				}
				assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, keys.AssetID)
				if err != nil {
					return true, err
				}
				usdValue := CalculateUSDValue(amount, price.Value, assetInfo.AssetBasicInfo.Decimals, price.Decimal)
				totalUSDValue = totalUSDValue.Add(usdValue)
			}
		}
		return false, nil
	}
	err = k.delegationKeeper.IterateDelegationsForStaker(ctx, stakerID, opFunc)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	return totalUSDValue, nil
}
