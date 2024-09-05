package keeper

import (
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationkeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	"github.com/ethereum/go-ethereum/common/hexutil"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerID string) (assetsInfo map[string]*assetstype.StakerAssetInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakerAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, []byte(stakerID))
	defer iterator.Close()

	ret := make(map[string]*assetstype.StakerAssetInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo assetstype.StakerAssetInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		keyList, err := assetstype.ParseJoinedStoreKey(iterator.Key(), 2)
		if err != nil {
			return nil, err
		}
		assetID := keyList[1]
		ret[assetID] = &stateInfo
	}
	// add exo-native-token info
	info, err := k.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetstype.NativeAssetID)
	if err != nil {
		return nil, err
	}
	ret[assetstype.NativeAssetID] = info
	return ret, nil
}

func (k Keeper) GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerID string, assetID string) (info *assetstype.StakerAssetInfo, err error) {
	if assetID == assetstype.NativeAssetID {
		stakerAddrStr, _, err := assetstype.ParseID(stakerID)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to parse stakerID")
		}
		stakerAccDecode, err := hexutil.Decode(stakerAddrStr)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to decode staker address")
		}
		stakerAcc := sdk.AccAddress(stakerAccDecode)
		balance := k.bk.GetBalance(ctx, stakerAcc, assetstype.NativeAssetDenom)
		info := &assetstype.StakerAssetInfo{
			TotalDepositAmount:  balance.Amount,
			WithdrawableAmount:  balance.Amount,
			WaitUnbondingAmount: math.NewInt(0),
		}

		delegationInfoRecords, err := k.dk.GetDelegationInfo(ctx, stakerID, assetID)
		if err != nil {
			return nil, errorsmod.Wrap(err, "failed to GetDelegationInfo")
		}
		for operator, record := range delegationInfoRecords.DelegationInfos {
			operatorAssetInfo, err := k.GetOperatorSpecifiedAssetInfo(ctx, sdk.MustAccAddressFromBech32(operator), assetID)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to GetOperatorSpecifiedAssetInfo")
			}
			undelegatableTokens, err := delegationkeeper.TokensFromShares(record.UndelegatableShare, operatorAssetInfo.TotalShare, operatorAssetInfo.TotalAmount)
			if err != nil {
				return nil, errorsmod.Wrap(err, "failed to get shares from token")
			}
			info.TotalDepositAmount = info.TotalDepositAmount.Add(undelegatableTokens).Add(record.WaitUndelegationAmount)
			info.WaitUnbondingAmount = info.WaitUnbondingAmount.Add(record.WaitUndelegationAmount)
		}
		return info, nil
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakerAssetInfos)
	key := assetstype.GetJoinedStoreKey(stakerID, assetID)
	value := store.Get(key)
	if value == nil {
		return nil, errorsmod.Wrap(assetstype.ErrNoStakerAssetKey, fmt.Sprintf("the key is:%s", key))
	}

	ret := assetstype.StakerAssetInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// UpdateStakerAssetState is used to update the staker asset state
// The input `changeAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is deposit or withdraw related to the specified staker.
func (k Keeper) UpdateStakerAssetState(ctx sdk.Context, stakerID string, assetID string, changeAmount assetstype.DeltaStakerSingleAsset) (err error) {
	// get the latest state,use the default initial state if the state hasn't been stored
	store := prefix.NewStore(ctx.KVStore(k.storeKey), assetstype.KeyPrefixReStakerAssetInfos)
	key := assetstype.GetJoinedStoreKey(stakerID, assetID)
	assetState := assetstype.StakerAssetInfo{
		TotalDepositAmount:  math.NewInt(0),
		WithdrawableAmount:  math.NewInt(0),
		WaitUnbondingAmount: math.NewInt(0),
	}
	value := store.Get(key)
	if value != nil {
		k.cdc.MustUnmarshal(value, &assetState)
	}
	// update all states of the specified restaker asset
	err = assetstype.UpdateAssetValue(&assetState.TotalDepositAmount, &changeAmount.TotalDepositAmount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState TotalDepositAmount error")
	}
	err = assetstype.UpdateAssetValue(&assetState.WithdrawableAmount, &changeAmount.WithdrawableAmount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState CanWithdrawAmountOrWantChangeValue error")
	}
	err = assetstype.UpdateAssetValue(&assetState.WaitUnbondingAmount, &changeAmount.WaitUnbondingAmount)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState WaitUndelegationAmountOrWantChangeValue error")
	}

	// store the updated state
	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)

	return nil
}
