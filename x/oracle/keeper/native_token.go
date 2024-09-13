package keeper

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// deposit: update staker's totalDeposit
// withdoraw: update staker's totalDeposit
// delegate: update operator's price, operator's totalAmount, operator's totalShare, staker's share
// undelegate: update operator's price, operator's totalAmount, operator's totalShare, staker's share
// msg(refund or slash on beaconChain): update staker's price, operator's price

var stakerList types.StakerList

// GetStakerInfo returns details about staker for native-restaking under asset of assetID
func (k Keeper) GetStakerInfo(ctx sdk.Context, assetID, stakerAddr string) types.StakerInfo {
	store := ctx.KVStore(k.storeKey)
	stakerInfo := types.StakerInfo{}
	value := store.Get(types.NativeTokenStakerKey(assetID, stakerAddr))
	if value == nil {
		return stakerInfo
	}
	k.cdc.MustUnmarshal(value, &stakerInfo)
	return stakerInfo
}

// TODO, NOTE: price changes will effect reward/slash calculation, every time one staker's price changed, it's reward/slash amount(LST) should be cleaned or recalculated immediately
// TODO: validatorIndex
// amount: represents for originalToken
func (k Keeper) UpdateNativeTokenByDepositOrWithdraw(ctx sdk.Context, assetID, stakerAddr string, amount sdkmath.Int, validatorIndex uint64) sdkmath.Int {
	// TODO: just convert the number for assets module, and don't store state in oracleModule, can use cache only here
	// TODO: we havn't included validatorIndex here, need the bridge info
	store := ctx.KVStore(k.storeKey)
	key := types.NativeTokenStakerKey(assetID, stakerAddr)
	stakerInfo := &types.StakerInfo{}
	if value := store.Get(key); value == nil {
		// create a new item for this staker
		stakerInfo = types.NewStakerInfo(stakerAddr, validatorIndex)
	} else {
		k.cdc.MustUnmarshal(value, stakerInfo)
	}

	latestIndex := len(stakerInfo.BalanceList) - 1
	newBalance := *(stakerInfo.BalanceList[latestIndex])
	newBalance.Index++
	newBalance.Block = uint64(ctx.BlockHeight())
	if amount.IsPositive() {
		newBalance.Change = types.BalanceInfo_ACTION_DEPOSIT
		// deopsit add a new validator into staker's validatorList
		stakerInfo.ValidatorIndexs = append(stakerInfo.ValidatorIndexs, validatorIndex)
	} else {
		// TODO: check if this validator has withdraw all its asset and then we can move it out from the staker's validatorList
		// currently when withdraw happened we assume this validator has left the staker's validatorList (deposit/withdraw all of that validator's staking ETH(<=32))
		newBalance.Change = types.BalanceInfo_ACTION_WITHDRAW
		for i, vIdx := range stakerInfo.ValidatorIndexs {
			if vIdx == validatorIndex {
				stakerInfo.ValidatorIndexs = append(stakerInfo.ValidatorIndexs[:i], stakerInfo.ValidatorIndexs[i+1:]...)
				break
			}
		}
	}

	newBalance.Balance += amount.Int64()

	keyStakerList := types.NativeTokenStakerListKey(assetID)
	valueStakerList := store.Get(keyStakerList)
	stakerList := &types.StakerList{}
	if valueStakerList != nil {
		k.cdc.MustUnmarshal(valueStakerList, stakerList)
	}
	exists := false
	for idx, stakerExists := range stakerList.StakerAddrs {
		if stakerExists == stakerAddr {
			if newBalance.Balance <= 0 {
				stakerList.StakerAddrs = append(stakerList.StakerAddrs[:idx], stakerList.StakerAddrs[idx+1:]...)
				valueStakerList = k.cdc.MustMarshal(stakerList)
				store.Set(keyStakerList, valueStakerList)
			}
			exists = true
			stakerInfo.StakerIndex = int64(idx)
			break
		}
	}

	if !exists {
		if newBalance.Balance <= 0 {
			// this should not happened, if a staker execute the 'withdraw' action, he must have already been in the stakerList
			return amount
		}
		stakerList.StakerAddrs = append(stakerList.StakerAddrs, stakerAddr)
		stakerInfo.StakerIndex = int64(len(stakerList.StakerAddrs) - 1)
		valueStakerList = k.cdc.MustMarshal(stakerList)
		store.Set(keyStakerList, valueStakerList)
	}

	if newBalance.Balance <= 0 {
		store.Delete(key)
	} else {
		bz := k.cdc.MustMarshal(stakerInfo)
		store.Set(key, bz)
	}
	return amount
}

// GetstakerList return stakerList for native-restaking asset of assetID
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

// UpdateNativeTokenByBalanceChange updates balance info for staker under native-restaking asset of assetID when its balance changed by slash/refund on the source chain (beacon chain for eth)
func (k Keeper) UpdateNativeTokenByBalanceChange(ctx sdk.Context, assetID string, rawData []byte, roundID uint64) error {
	if len(rawData) < 32 {
		return errors.New("length of indicate maps for stakers shoule be exactly 32 bytes")
	}
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
		//		changeOriginalFloat := sdkmath.LegacyNewDec(int64(change))
		//		changeFloat := sdkmath.LegacyNewDec(int64(change))
		length := len(stakerInfo.BalanceList)
		balance := stakerInfo.BalanceList[length-1]
		newBalance := *balance
		newBalance.Block = uint64(ctx.BlockHeight())
		if newBalance.RoundID == roundID {
			newBalance.Index++
		} else {
			newBalance.RoundID = roundID
		}
		newBalance.Change = types.BalanceInfo_ACTION_SLASH_REFUND
		newBalance.Balance += int64(change)
		stakerInfo.Append(&newBalance)
		bz := k.cdc.MustMarshal(stakerInfo)
		store.Set(key, bz)
		// TODO: call assetsmodule. func(k Keeper) UpdateNativeRestakingBalance(ctx sdk.Context, stakerID, assetID string, amount sdkmath.Int) error
	}
	return nil
}

// getStakerList returns all Stakers for native-restaking of assetID, this is used for cache
func (k Keeper) getStakerList(ctx sdk.Context, assetID string) types.StakerList {
	if len(stakerList.StakerAddrs) == 0 {
		stakerList = k.GetStakerList(ctx, assetID)
	}
	return stakerList
}

// parseBalanceChange parses rawData to details of amount change for all stakers relative to native restaking
func parseBalanceChange(rawData []byte, sl types.StakerList) (map[string]int, error) {
	indexs := rawData[:32]
	changes := rawData[32:]
	index := -1
	byteIndex := 0
	bitOffset := 0
	lengthBits := 5
	stakerChanges := make(map[string]int)
	for _, b := range indexs {
		for i := 7; i >= 0; i-- {
			index++
			if (b>>i)&1 == 1 {
				lenValue := changes[byteIndex] << bitOffset
				bitsLeft := 8 - bitOffset
				lenValue >>= (8 - lengthBits)
				if bitsLeft < lengthBits {
					byteIndex++
					lenValue |= changes[byteIndex] >> (8 - lengthBits + bitsLeft)
					bitOffset = lengthBits - bitsLeft
				} else {
					if bitOffset += lengthBits; bitOffset == 8 {
						bitOffset = 0
					}
					if bitsLeft == lengthBits {
						byteIndex++
					}
				}

				symbol := lenValue & 1
				lenValue >>= 1
				if lenValue <= 0 {
					return stakerChanges, errors.New("length of change value must be at least 1 bit")
				}

				bitsExtracted := 0
				stakerChange := 0
				for bitsExtracted < int(lenValue) {
					bitsLeft := 8 - bitOffset
					byteValue := changes[byteIndex] << bitOffset
					if (int(lenValue) - bitsExtracted) < bitsLeft {
						bitsLeft = int(lenValue) - bitsExtracted
						bitOffset += bitsLeft
					} else {
						byteIndex++
						bitOffset = 0
					}
					byteValue >>= (8 - bitsLeft)
					stakerChange = (stakerChange << bitsLeft) | int(byteValue)
					bitsExtracted += bitsLeft
				}
				stakerChange++
				if symbol == 1 {
					stakerChange *= -1
				}
				stakerChanges[sl.StakerAddrs[index]] = stakerChange
			}
		}
	}
	return stakerChanges, nil
}
