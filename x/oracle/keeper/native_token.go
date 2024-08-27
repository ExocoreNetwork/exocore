package keeper

import (
	"errors"
	"strings"

	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// deposit: update staker's totalDeposit
// withdoraw: update staker's totalDeposit
// delegate: update operator's price, operator's totalAmount, operator's totalShare, staker's share
// undelegate: update operator's price, operator's totalAmount, operator's totalShare, staker's share
// msg(refund or slash on beaconChain): update staker's price, operator's price

var stakerList types.StakerList

// TODO, NOTE: price changes will effect reward/slash calculation, every time one staker's price changed, it's reward/slash amount(LST) should be cleaned or recalculated immediately
// TODO: validatorIndex
// amount: represents for originalToken
func (k Keeper) UpdateNativeTokenByDepositOrWithdraw(ctx sdk.Context, assetID, stakerAddr string, amount sdkmath.Int) sdkmath.Int {
	// TODO: just convert the number for assets module, and don't store state in oracleModule, can use cache only here
	// TODO: we havn't included validatorIndex here, need the bridge info
	store := ctx.KVStore(k.storeKey)
	key := types.NativeTokenStakerKey(assetID, stakerAddr)
	stakerInfo := &types.StakerInfo{}
	if value := store.Get(key); value == nil {
		// create a new item for this staker
		stakerInfo = types.NewStakerInfo(stakerAddr)
	} else {
		k.cdc.MustUnmarshal(value, stakerInfo)
	}
	latestIndex := len(stakerInfo.PriceList) - 1
	// calculate amount of virtual LST from nativetoken with price
	amountInt := convertAmountOriginalIntToAmountFloat(amount, stakerInfo.PriceList[latestIndex].Price).RoundInt()
	stakerInfo.TotalDeposit = stakerInfo.TotalDeposit.Add(amountInt)
	// update totalDeposit of staker, and price won't change on either deposit or withdraw
	bz := k.cdc.MustMarshal(stakerInfo)
	store.Set(key, bz)
	return amountInt
}

// UpdateNativeTokenByDelegation update operator's price, operator's totalAmount, operator's totalShare, staker's share bsed on either delegation or undelegation
// this amount passed in from delegation hooks represent originalToken(not virtualLST)
func (k Keeper) UpdateNativeTokenByDelegation(ctx sdk.Context, assetID, operatorAddr, stakerAddr string, amountOriginal sdkmath.Int) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)
	keyOperator := types.NativeTokenOperatorKey(assetID, operatorAddr)
	operatorInfo := &types.OperatorInfo{}
	value := store.Get(keyOperator)
	if value == nil {
		operatorInfo = types.NewOperatorInfo(operatorAddr)
	} else {
		k.cdc.MustUnmarshal(value, operatorInfo)
	}
	stakerInfo := &types.StakerInfo{}
	keyStaker := types.NativeTokenStakerKey(assetID, stakerAddr)
	if value = store.Get(keyStaker); value == nil {
		panic("staker must exist before delegation")
	}
	k.cdc.MustUnmarshal(value, stakerInfo)

	operatorAmountFloat, operatorAmountOriginalFloat := getOperatorAmountFloat(operatorInfo)
	amountFloat, amountOriginalFloat := parseStakerAmountOriginalInt(amountOriginal, stakerInfo)

	operatorAmountOriginalFloat = operatorAmountOriginalFloat.Add(amountOriginalFloat)
	operatorAmountFloat = operatorAmountFloat.Add(amountFloat)

	// update operator's price for native token base on new delegation
	operatorInfo.PriceList = append(operatorInfo.PriceList, &types.PriceInfo{
		Price: operatorAmountOriginalFloat.Quo(operatorAmountFloat),
		Block: uint64(ctx.BlockHeight()),
	})

	// update operator's total amount for native token, for this 'amount' we don't disginguish different tokens from different stakers. That difference reflects in 'operator price'
	operatorInfo.TotalAmount = operatorAmountFloat.RoundInt()
	bz := k.cdc.MustMarshal(operatorInfo)
	store.Set(keyOperator, bz)
	amountInt := amountFloat.RoundInt()
	// update staker's related operator list
	keyDelegation := types.NativeTokenStakerDelegationKey(assetID, stakerAddr)
	stakerDelegation := &types.StakerDelegationInfo{}
	if value = store.Get(keyDelegation); value == nil {
		stakerDelegation.Delegations = []*types.DelegationInfo{
			{
				OperatorAddr: operatorAddr,
				Amount:       amountInt,
			},
		}
	} else {
		k.cdc.MustUnmarshal(value, stakerDelegation)
		for idx, delegationInfo := range stakerDelegation.Delegations {
			if delegationInfo.OperatorAddr == operatorAddr {
				if delegationInfo.Amount = delegationInfo.Amount.Add(amountInt); !delegationInfo.Amount.IsPositive() {
					stakerDelegation.Delegations = append(stakerDelegation.Delegations[0:idx], stakerDelegation.Delegations[idx:]...)
				}
				value = k.cdc.MustMarshal(stakerDelegation)
				store.Set(keyDelegation, value)
				return amountInt
			}
		}
		stakerDelegation.Delegations = append(stakerDelegation.Delegations, &types.DelegationInfo{
			OperatorAddr: operatorAddr,
			Amount:       amountInt,
		})
	}
	// update staker delegation infos for related operators
	value = k.cdc.MustMarshal(stakerDelegation)
	store.Set(keyDelegation, value)

	return amountInt
}

func (k Keeper) GetNativeTokenPriceUSDForOperator(ctx sdk.Context, assetID string) (types.Price, error) {
	parsedAssetID := strings.Split(assetID, "_")
	if len(parsedAssetID) != 3 {
		return types.Price{}, types.ErrGetPriceAssetNotFound
	}
	assetID = strings.Join([]string{parsedAssetID[0], parsedAssetID[1]}, "_")
	operatorAddr := parsedAssetID[2]

	store := ctx.KVStore(k.storeKey)
	key := types.NativeTokenOperatorKey(assetID, operatorAddr)
	if value := store.Get(key); value == nil {
		return types.Price{}, types.ErrGetPriceAssetNotFound
	} else {
		operatorInfo := &types.OperatorInfo{}
		k.cdc.MustUnmarshal(value, operatorInfo)
		baseTokenUSDPrice, err := k.GetSpecifiedAssetsPrice(ctx, assetstypes.GetBaseTokenForNativeToken(assetID))
		if err != nil {
			return types.Price{}, types.ErrGetPriceAssetNotFound
		}
		operatorPriceFloat := getLatestOperatorPriceFloat(operatorInfo)
		baseTokenUSDPrice.Value = (baseTokenUSDPrice.Value.ToLegacyDec().Mul(operatorPriceFloat)).RoundInt()
		return baseTokenUSDPrice, nil
	}
}

func (k Keeper) GetStakerList(ctx sdk.Context, assetID string) types.StakerList {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.NativeTokenStakerListKey(assetID))
	if value == nil {
		return types.StakerList{}
	}
	stakerList := &types.StakerList{}
	k.cdc.MustUnmarshal(value, stakerList)
	return *stakerList
}

func (k Keeper) UpdateNativeTokenByBalanceChange(ctx sdk.Context, assetID string, rawData []byte, roundID uint64) error {
	sl := k.getStakerList(ctx, assetID)
	if len(sl.StakerAddrs) == 0 {
		return errors.New("staker list is empty")
	}
	stakerChanges, err := parseBalanceChange(rawData, sl)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	for stakerAddr, change := range stakerChanges {
		key := types.NativeTokenStakerKey(assetID, stakerAddr)
		value := store.Get(key)
		if value == nil {
			return errors.New("stakerInfo does not exist")
		}
		stakerInfo := &types.StakerInfo{}
		k.cdc.MustUnmarshal(value, stakerInfo)
		changeOriginalFloat := sdkmath.LegacyNewDec(int64(change))
		totalAmountFloat, totalAmountOriginalFloat := parseStakerAmountInt(stakerInfo.TotalDeposit, stakerInfo)
		totalAmountOriginalFloat = totalAmountOriginalFloat.Add(changeOriginalFloat)
		prevStakerPrice := getLatestStakerPriceFloat(stakerInfo)
		// update staker price based on beacon chain effective balance change
		stakerPrice := totalAmountOriginalFloat.Quo(totalAmountFloat)
		stakerInfo.PriceList = append(stakerInfo.PriceList, &types.PriceInfo{
			Price:   stakerPrice,
			Block:   uint64(ctx.BlockHeight()),
			RoundID: roundID,
		})
		bz := k.cdc.MustMarshal(stakerInfo)
		store.Set(key, bz)
		// update related operator's price
		keyStakerDelegations := types.NativeTokenStakerDelegationKey(assetID, stakerAddr)
		value = store.Get(keyStakerDelegations)
		if value != nil {
			delegationInfo := &types.StakerDelegationInfo{}
			k.cdc.MustUnmarshal(value, delegationInfo)
			for _, delegation := range delegationInfo.Delegations {
				keyOperator := types.NativeTokenOperatorKey(assetID, delegation.OperatorAddr)
				value = store.Get(keyOperator)
				if value == nil {
					panic("staker delegation related to operator not exists")
				}
				operatorInfo := &types.OperatorInfo{}
				k.cdc.MustUnmarshal(value, operatorInfo)
				AmountFloat, prevAmountOriginalFloat := getOperatorAmountFloat(operatorInfo)
				delta := delegation.Amount.ToLegacyDec().Mul(stakerPrice.Sub(prevStakerPrice))
				operatorInfo.PriceList = append(operatorInfo.PriceList, &types.PriceInfo{
					Price:   prevAmountOriginalFloat.Add(delta).Quo(AmountFloat),
					Block:   uint64(ctx.BlockHeight()),
					RoundID: roundID,
				})
				bz := k.cdc.MustMarshal(operatorInfo)
				store.Set(keyOperator, bz)
			}

		}

	}
	return nil
}

func (k Keeper) getStakerList(ctx sdk.Context, assetID string) types.StakerList {
	if len(stakerList.StakerAddrs) == 0 {
		stakerList = k.GetStakerList(ctx, assetID)
	}
	return stakerList
}

func parseBalanceChange(rawData []byte, sl types.StakerList) (map[string]int, error) {
	indexs := rawData[:32]
	changes := rawData[32:]
	//	lenChanges := len(changes)
	index := -1
	byteIndex := -1
	bitOffset := 5
	stakerChanges := make(map[string]int)
	for _, b := range indexs {
		for i := 7; i >= 0; i-- {
			// staker's index start from 1
			index++
			if (b>>i)&1 == 1 {
				// effect balance  f stakerAddr[index] has changed
				lenValue := int(changes[byteIndex] >> 4)
				if lenValue <= 0 {
					return stakerChanges, errors.New("length of change value must be at least 1 bit")
				}
				symbol := (changes[byteIndex] >> 3) & 1
				bitsExtracted := 0
				stakerChange := 0
				for j := 0; j < lenValue; j++ {
					byteIndex++
					byteValue := changes[byteIndex] << bitOffset
					// byteValue <<= bitOffset
					bitsLeft := 8 - bitOffset
					if bitsExtracted+bitsLeft > lenValue {
						bitsLeft = lenValue - bitsExtracted
						bitOffset = bitsLeft
					} else {
						bitOffset = 0
					}
					byteValue = (byteValue >> (8 - bitsLeft)) & ((1 << bitsLeft) - 1)
					stakerChange = (stakerChange << bitsLeft) | int(byteValue)
				}
				if symbol == 1 {
					stakerChange *= -1
				}
				stakerChanges[sl.StakerAddrs[index]] = stakerChange
			}
		}
	}
	return stakerChanges, nil
}

func getLatestOperatorPriceFloat(operatorInfo *types.OperatorInfo) sdkmath.LegacyDec {
	latestIndex := len(operatorInfo.PriceList) - 1
	return operatorInfo.PriceList[latestIndex].Price
}

func getLatestStakerPriceFloat(stakerInfo *types.StakerInfo) sdkmath.LegacyDec {
	latestIndex := len(stakerInfo.PriceList) - 1
	return stakerInfo.PriceList[latestIndex].Price
}

func convertAmountOriginalIntToAmountFloat(amount sdkmath.Int, price sdkmath.LegacyDec) sdkmath.LegacyDec {
	amountFloat := amount.ToLegacyDec()
	return amountFloat.Quo(price)
}

func getOperatorAmountFloat(operatorInfo *types.OperatorInfo) (amountFloat, amountOriginalFloat sdkmath.LegacyDec) {
	latestIndexOperator := len(operatorInfo.PriceList) - 1
	price := operatorInfo.PriceList[latestIndexOperator].Price
	amountFloat = operatorInfo.TotalAmount.ToLegacyDec()
	amountOriginalFloat = amountFloat.Mul(price)
	return
}

func parseStakerAmountInt(amount sdkmath.Int, stakerInfo *types.StakerInfo) (amountFloat, amountOriginalFloat sdkmath.LegacyDec) {
	latestIndex := len(stakerInfo.PriceList) - 1
	price := stakerInfo.PriceList[latestIndex].Price
	amountFloat = amount.ToLegacyDec()
	amountOriginalFloat = amountFloat.Mul(price)
	return
}

func parseStakerAmountOriginalInt(amountOriginalInt sdkmath.Int, stakerInfo *types.StakerInfo) (amountFloat, amountOriginalFloat sdkmath.LegacyDec) {
	latestIndex := len(stakerInfo.PriceList) - 1
	price := stakerInfo.PriceList[latestIndex].Price
	amountOriginalFloat = amountOriginalInt.ToLegacyDec()
	amountFloat = amountOriginalFloat.Quo(price)
	return
}
