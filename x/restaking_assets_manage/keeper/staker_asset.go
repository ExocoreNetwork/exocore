package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]*restakingtype.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	iterator := sdk.KVStorePrefixIterator(store, []byte(stakerId))
	defer iterator.Close()

	ret := make(map[string]*restakingtype.StakerSingleAssetOrChangeInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var stateInfo restakingtype.StakerSingleAssetOrChangeInfo
		k.cdc.MustUnmarshal(iterator.Value(), &stateInfo)
		keyList, err := restakingtype.ParseJoinedStoreKey(iterator.Key(), 2)
		if err != nil {
			return nil, err
		}
		assetId := keyList[1]
		ret[assetId] = &stateInfo
	}
	return ret, nil
}

func (k Keeper) GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerId string, assetId string) (info *restakingtype.StakerSingleAssetOrChangeInfo, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	key := restakingtype.GetJoinedStoreKey(stakerId, assetId)
	ifExist := store.Has(key)
	if !ifExist {
		return nil, errorsmod.Wrap(restakingtype.ErrNoStakerAssetKey, fmt.Sprintf("the key is:%s", key))
	}

	value := store.Get(key)

	ret := restakingtype.StakerSingleAssetOrChangeInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// UpdateStakerAssetState It's used to update the staker asset state
// The input `changeAmount` represents the values that you want to add or decrease,using positive or negative values for increasing and decreasing,respectively. The function will calculate and update new state after a successful check.
// The function will be called when there is deposit or withdraw related to the specified staker.
func (k Keeper) UpdateStakerAssetState(ctx sdk.Context, stakerId string, assetId string, changeAmount restakingtype.StakerSingleAssetOrChangeInfo) (err error) {
	//get the latest state,use the default initial state if the state hasn't been stored
	store := prefix.NewStore(ctx.KVStore(k.storeKey), restakingtype.KeyPrefixReStakerAssetInfos)
	key := restakingtype.GetJoinedStoreKey(stakerId, assetId)
	assetState := restakingtype.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     math.NewInt(0),
		CanWithdrawAmountOrWantChangeValue:      math.NewInt(0),
		WaitUndelegationAmountOrWantChangeValue: math.NewInt(0),
	}
	if store.Has(key) {
		value := store.Get(key)
		k.cdc.MustUnmarshal(value, &assetState)
	}

	// update all states of the specified restaker asset
	err = restakingtype.UpdateAssetValue(&assetState.TotalDepositAmountOrWantChangeValue, &changeAmount.TotalDepositAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState TotalDepositAmountOrWantChangeValue error")
	}
	err = restakingtype.UpdateAssetValue(&assetState.CanWithdrawAmountOrWantChangeValue, &changeAmount.CanWithdrawAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState CanWithdrawAmountOrWantChangeValue error")
	}
	err = restakingtype.UpdateAssetValue(&assetState.WaitUndelegationAmountOrWantChangeValue, &changeAmount.WaitUndelegationAmountOrWantChangeValue)
	if err != nil {
		return errorsmod.Wrap(err, "UpdateStakerAssetState WaitUndelegationAmountOrWantChangeValue error")
	}

	// store the updated state
	bz := k.cdc.MustMarshal(&assetState)
	store.Set(key, bz)

	return nil
}
