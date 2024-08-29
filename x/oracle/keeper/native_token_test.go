package keeper_test

import (
	"encoding/binary"
	"strings"

	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imroc/biu"
)

// workflow:
// 1. Deposit. into staker_A
//  1. stakerInfo {totalDeposit, price} - new
//  2. stakerList - new
//
// 2. Delegate. into operator_A
//  1. stakerDelegation_AA {amount, operator} - new
//  2. operatorInfo_A {totalAmount, price} - new
//
// 3. Msg. minus staker_A's amountOriginal
//  1. stakerInfo_A {price-change} - update
//  2. operatorInfo_A {price-change} - update
//
// 4. Deposit more. into staker_A
//  1. stakerInfo {totalDeposit-change} -update
//
// 5. Msg. add staker_A's amountOriginal
//  1. stakerInfo_A {price-change} - update
//  2. operatorInfo_A {price-change} - update
//
// 6. delegate into operator_A
//  1. stakerDelegation_AA {amount-change} - update
//  2. operatorInfo_A {totalAmount-change} - update
//
// 7. Undelegate from operator_A
//  1. stakerDelegation_AA {amount-chagne} - update
//  2. operatorInfo_A {price-change, totalAmount-change} - update
//
// 8. UndelegateAll from operator_A
//  1. stakerDelegation_AA item-removed
//  2. operatorInfo_A totalAmount->0-> operatorInfo removed
//
// 9. withdrawAll from staker_A
//  1. stakerInfo removed
//  2. stakerList removed

func (ks *KeeperSuite) TestNativeTokenLifeCycleOneStaker() {
	operator := ks.Operators[0]
	operatorStr := operator.String()
	stakerStr := common.Address(operator.Bytes()).String()
	assetID := assetstypes.NativeETHAssetID
	// 1. deposit amount 100
	amount100 := sdkmath.NewIntFromUint64(100)
	ks.k.UpdateNativeTokenByDepositOrWithdraw(ks.ctx, assetID, stakerStr, amount100)
	// - 1.1 check stakerInfo
	stakerInfo := ks.k.GetStakerInfo(ks.ctx, assetID, stakerStr)
	ks.Equal(stakerInfo.TotalDeposit, amount100)
	// - 1.2 check stakerList
	stakerList := ks.k.GetStakerList(ks.ctx, assetID)
	ks.Equal(stakerList.StakerAddrs[0], stakerStr)
	// 2. delegateTo operator with amount 80
	amount80 := sdkmath.NewIntFromUint64(80)
	ks.k.UpdateNativeTokenByDelegation(ks.ctx, assetID, operatorStr, stakerStr, amount80)
	// - 2.1 check stakerDelegatioin
	stakerDelegation := ks.k.GetStakerDelegations(ks.ctx, assetID, stakerStr)
	ks.Equal(len(stakerDelegation.Delegations), 1)
	ks.Equal(stakerDelegation.Delegations[0].OperatorAddr, operatorStr)
	ks.Equal(stakerDelegation.Delegations[0].Amount, amount80)
	// - 2.2 check operatorInfo
	operatorInfo := ks.k.GetOperatorInfo(ks.ctx, assetID, operatorStr)
	ks.Equal(operatorInfo, types.OperatorInfo{
		OperatorAddr: operatorStr,
		TotalAmount:  amount80,
		PriceList: []*types.PriceInfo{
			{
				Price:   sdkmath.LegacyNewDec(1),
				Block:   2,
				RoundID: 0,
			},
		},
	})
	// 3. Msg. minus staker's amountOriginal
	stakerChanges := [][]int{
		{0, -50},
	}
	rawData := convertBalanceChangeToBytes(stakerChanges)
	ks.k.UpdateNativeTokenByBalanceChange(ks.ctx, assetID, rawData, 9)
	// - 3.1 check stakerInfo
	stakerInfo = ks.k.GetStakerInfo(ks.ctx, assetID, stakerStr)
	ks.Equal(stakerInfo.PriceList[len(stakerInfo.PriceList)-1].Price, sdkmath.LegacyNewDecWithPrec(5, 1))
	// - 3.2 check operatorInfo
	operatorInfo = ks.k.GetOperatorInfo(ks.ctx, assetID, operatorStr)
	ks.Equal(operatorInfo.PriceList[len(operatorInfo.PriceList)-1].Price, sdkmath.LegacyNewDecWithPrec(5, 1))
	ks.Equal(operatorInfo.PriceList[len(operatorInfo.PriceList)-1].RoundID, uint64(9))

	// 4. deposit more. 100
	ks.k.UpdateNativeTokenByDepositOrWithdraw(ks.ctx, assetID, stakerStr, amount100)
	// - 4.1 check stakerInfo
	stakerInfo = ks.k.GetStakerInfo(ks.ctx, assetID, stakerStr)
	amount300 := sdkmath.NewInt(300)
	ks.Equal(stakerInfo.TotalDeposit, amount300)
	// 5. Msg. add staker's amountOriginal
	stakerChanges = [][]int{
		{0, 30},
	}
	rawData = convertBalanceChangeToBytes(stakerChanges)
	ks.k.UpdateNativeTokenByBalanceChange(ks.ctx, assetID, rawData, 11)
	// - 5.1 check stakerInfo
	stakerInfo = ks.k.GetStakerInfo(ks.ctx, assetID, stakerStr)
	ks.Equal(types.PriceInfo{
		Price:   sdkmath.LegacyNewDecWithPrec(6, 1),
		Block:   2,
		RoundID: 11,
	}, *stakerInfo.PriceList[2])
	ks.Equal(amount300, stakerInfo.TotalDeposit)
	// - 5.2 check operatorInfo
	operatorInfo = ks.k.GetOperatorInfo(ks.ctx, assetID, operatorStr)
	ks.Equal(sdkmath.LegacyNewDecWithPrec(6, 1), operatorInfo.PriceList[len(operatorInfo.PriceList)-1].Price)

	// 6. delegate more. 60->100
	amount60 := sdkmath.NewInt(60)
	ks.k.UpdateNativeTokenByDelegation(ks.ctx, assetID, operatorStr, stakerStr, amount60)
	// - 6.1 check delegation-record
	stakerDelegation = ks.k.GetStakerDelegations(ks.ctx, assetID, stakerStr)
	amount180 := sdkmath.NewInt(180)
	ks.Equal(amount180, stakerDelegation.Delegations[0].Amount)
	// - 6.2 check operatorInfo
	operatorInfo = ks.k.GetOperatorInfo(ks.ctx, assetID, operatorStr)
	ks.Equal(amount180, operatorInfo.TotalAmount)

	// 7. undelegate. 72->120
	amount72N := sdkmath.NewInt(-72)
	ks.k.UpdateNativeTokenByDelegation(ks.ctx, assetID, operatorStr, stakerStr, amount72N)
	// - 7.1 check delegation-record
	stakerDelegation = ks.k.GetStakerDelegations(ks.ctx, assetID, stakerStr)
	ks.Equal(amount60, stakerDelegation.Delegations[0].Amount)
	// - 7.2 check operatorInfo
	operatorInfo = ks.k.GetOperatorInfo(ks.ctx, assetID, operatorStr)
	ks.Equal(amount60, operatorInfo.TotalAmount)

	// 8. undelegate all
	amount36N := sdkmath.NewInt(-36)
	ks.k.UpdateNativeTokenByDelegation(ks.ctx, assetID, operatorStr, stakerStr, amount36N)
	// - 8.1 check delegation-record
	stakerDelegation = ks.k.GetStakerDelegations(ks.ctx, assetID, stakerStr)
	ks.Equal(0, len(stakerDelegation.Delegations))
	// - 8.2 check operatorInfo
	operatorInfo = ks.k.GetOperatorInfo(ks.ctx, assetID, operatorStr)
	ks.Equal(types.OperatorInfo{}, operatorInfo)

	// 9. withdraw all
	amount180N := sdkmath.NewInt(-180)
	ks.k.UpdateNativeTokenByDepositOrWithdraw(ks.ctx, assetID, stakerStr, amount180N)
	// - 9.1 check stakerInfo
	stakerInfo = ks.k.GetStakerInfo(ks.ctx, assetID, stakerStr)
	ks.Equal(types.StakerInfo{}, stakerInfo)
	// - 9.2 check stakerList
	stakerList = ks.k.GetStakerList(ks.ctx, assetID)
	ks.Equal(0, len(stakerList.StakerAddrs))
}

func convertBalanceChangeToBytes(stakerChanges [][]int) []byte {
	if len(stakerChanges) == 0 {
		return nil
	}
	str := ""
	index := 0
	changeBytesList := make([][]byte, 0, len(stakerChanges))
	bitsList := make([]int, 0, len(stakerChanges))
	for _, stakerChange := range stakerChanges {
		str += strings.Repeat("0", stakerChange[0]-index) + "1"
		index = stakerChange[0] + 1

		// change amount -> bytes
		change := stakerChange[1]
		var changeBytes []byte
		symbol := 1
		if change < 0 {
			symbol = -1
			change *= -1
		}
		change--
		bits := 0
		if change == 0 {
			bits = 1
			changeBytes = []byte{byte(0)}
		} else {
			tmpChange := change
			for tmpChange > 0 {
				bits++
				tmpChange /= 2
			}
			if change < 256 {
				// 1 byte
				changeBytes = []byte{byte(change)}
				changeBytes[0] <<= (8 - bits)
			} else {
				// 2 byte
				changeBytes = make([]byte, 2)
				binary.BigEndian.PutUint16(changeBytes, uint16(change))
				moveLength := 16 - bits
				changeBytes[0] <<= moveLength
				tmp := changeBytes[1] >> (8 - moveLength)
				changeBytes[0] |= tmp
				changeBytes[1] <<= moveLength
			}
		}

		// use lower 4 bits to represent the length of valid change value in bits format
		bitsLengthBytes := []byte{byte(bits)}
		bitsLengthBytes[0] <<= 4
		if symbol < 0 {
			bitsLengthBytes[0] |= 8
		}

		tmp := changeBytes[0] >> 5
		bitsLengthBytes[0] |= tmp
		if bits <= 3 {
			changeBytes = nil
		} else {
			changeBytes[0] <<= 3
		}

		if len(changeBytes) == 2 {
			tmp = changeBytes[1] >> 5
			changeBytes[0] |= tmp
			if bits <= 11 {
				changeBytes = changeBytes[:1]
			} else {
				changeBytes[1] <<= 3
			}
		}
		bitsLengthBytes = append(bitsLengthBytes, changeBytes...)
		changeBytesList = append(changeBytesList, bitsLengthBytes)
		bitsList = append(bitsList, bits)
	}

	l := len(bitsList)
	changeResult := changeBytesList[l-1]
	bitsList[len(bitsList)-1] = bitsList[len(bitsList)-1] + 5
	for i := l - 2; i >= 0; i-- {
		prev := changeBytesList[i]

		byteLength := 8 * len(prev)
		bitsLength := bitsList[i] + 5
		// delta must <8
		delta := byteLength - bitsLength
		if delta == 0 {
			changeResult = append(prev, changeResult...)
			bitsList[i] = bitsLength + bitsList[i+1]
		} else {
			// delta : (0,8)
			tmp := changeResult[0] >> (8 - delta)
			prev[len(prev)-1] |= tmp
			if len(changeResult) > 1 {
				for j := 1; j < len(changeResult); j++ {
					changeResult[j-1] <<= delta
					tmp := changeResult[j] >> (8 - delta)
					changeResult[j-1] |= tmp
				}
			}
			changeResult[len(changeResult)-1] <<= delta
			left := bitsList[i+1] % 8
			if bitsList[i+1] > 0 && left == 0 {
				left = 8
			}
			if left <= delta {
				changeResult = changeResult[:len(changeResult)-1]
			}
			changeResult = append(prev, changeResult...)
			bitsList[i] = bitsLength + bitsList[i+1]
		}
	}
	str += strings.Repeat("0", 256-index)
	bytesIndex := biu.BinaryStringToBytes(str)

	result := append(bytesIndex, changeResult...)
	return result
}

// func convertBalanceChangeToBytes(stakerChanges [][]int) []byte {
// 	if len(stakerChanges) == 0 {
// 		return nil
// 	}
// 	str := ""
// 	index := 0
// 	changeBytesList := make([][]byte, 0, len(stakerChanges))
// 	bitsList := make([]int, 0, len(stakerChanges))
// 	for _, stakerChange := range stakerChanges {
// 		str += strings.Repeat("0", stakerChange[0]-index) + "1"
// 		index = stakerChange[0] + 1
//
// 		// change amount -> bytes
// 		change := stakerChange[1]
// 		var changeBytes []byte
// 		symbol := 1
// 		if change < 0 {
// 			symbol = -1
// 			change *= -1
// 			change--
// 		}
// 		bits := 0
// 		if change == 0 {
// 			bits = 1
// 			changeBytes = []byte{byte(0)}
// 		} else {
// 			for change > 0 {
// 				bits++
// 				change /= 2
// 			}
// 			if change < 256 {
// 				// 1 byte
// 				changeBytes = []byte{byte(change)}
// 				changeBytes[0] <<= (8 - bits)
// 			} else {
// 				// 2 byte
// 				changeBytes = make([]byte, 0, 2)
// 				binary.BigEndian.PutUint16(changeBytes, uint16(change))
// 				moveLength := 16 - bits
// 				changeBytes[0] <<= moveLength
// 				tmp := changeBytes[1] >> (8 - moveLength)
// 				changeBytes[0] |= tmp
// 				changeBytes[1] <<= moveLength
// 			}
// 		}
//
// 		// use lower 4 bits to represent the length of valid change value in bits format
// 		bitsLengthBytes := []byte{byte(bits)}
// 		bitsLengthBytes[0] <<= 4
// 		if symbol < 0 {
// 			bitsLengthBytes[0] |= 8
// 		}
//
// 		tmp := changeBytes[0] >> 5
// 		bitsLengthBytes[0] |= tmp
// 		if bits <= 3 {
// 			changeBytes = nil
// 		} else {
// 			changeBytes[0] <<= 3
// 		}
//
// 		if len(changeBytes) == 2 {
// 			tmp = changeBytes[1] >> 5
// 			changeBytes[0] |= tmp
// 			if bits <= 11 {
// 				changeBytes = changeBytes[:1]
// 			} else {
// 				changeBytes[1] <<= 3
// 			}
// 		}
// 		bitsLengthBytes = append(bitsLengthBytes, changeBytes...)
// 		changeBytesList = append(changeBytesList, bitsLengthBytes)
// 		bitsList = append(bitsList, bits)
// 	}
//
// 	l := len(bitsList)
// 	changeResult := changeBytesList[l-1]
// 	bitsList[len(bitsList)-1] = bitsList[len(bitsList)-1] + 5
// 	for i := l - 2; i >= 0; i-- {
// 		prev := changeBytesList[i]
//
// 		byteLength := 8 * len(prev)
// 		bitsLength := bitsList[i] + 5
// 		// delta must <8
// 		delta := byteLength - bitsLength
// 		if delta == 0 {
// 			changeResult = append(prev, changeResult...)
// 			bitsList[i] = bitsLength + bitsList[i+1]
// 		} else {
// 			// delta : (0,8)
// 			tmp := changeResult[0] >> (8 - delta)
// 			prev[len(prev)-1] |= tmp
// 			if len(changeResult) > 1 {
// 				for j := 1; j < len(changeResult); j++ {
// 					changeResult[j-1] <<= delta
// 					tmp := changeResult[j] >> (8 - delta)
// 					changeResult[j-1] |= tmp
// 				}
// 			}
// 			changeResult[len(changeResult)-1] <<= delta
// 			if bitsList[i+1]%8 <= delta {
// 				changeResult = changeResult[:len(changeResult)-1]
// 			}
// 			changeResult = append(prev, changeResult...)
// 			bitsList[i] = bitsLength + bitsList[i+1]
// 		}
// 	}
// 	str += strings.Repeat("0", 256-index)
// 	bytesIndex := biu.BinaryStringToBytes(str)
//
// 	result := append(bytesIndex, changeResult...)
// 	return result
// }
