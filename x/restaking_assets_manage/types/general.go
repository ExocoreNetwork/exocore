package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	CrossChainActionLength       = 1
	CrossChainOpAmountLength     = 32
	GeneralAssetsAddrLength      = 32
	GeneralClientChainAddrLength = 32

	ClientChainLzIdIndexInTopics = 0
	LzNonceIndexInTopics         = 2

	ExoCoreOperatorAddrLength = 44
)

type GeneralAssetsAddr [32]byte

type GeneralClientChainAddr [32]byte

type CrossChainOpType uint8

type WithdrawerAddress [32]byte

const (
	Deposit CrossChainOpType = iota
	WithdrawPrinciple
	WithDrawReward
	DelegateTo
	UndelegateFrom
	Slash
)

// GetStakeIDAndAssetId stakerId = stakerAddress+'_'+clientChainLzId,assetId = assetAddress+'_'+clientChainLzId
func GetStakeIDAndAssetId(clientChainLzId uint64, stakerAddress []byte, assetsAddress []byte) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(clientChainLzId)
	if stakerAddress != nil {
		stakeId = strings.Join([]string{hexutil.Encode(stakerAddress), clientChainLzIdStr}, "_")
	}

	if assetsAddress != nil {
		assetId = strings.Join([]string{hexutil.Encode(assetsAddress), clientChainLzIdStr}, "_")
	}
	return
}

// GetStakeIDAndAssetIdFromStr stakerId = stakerAddress+'_'+clientChainLzId,assetId = assetAddress+'_'+clientChainLzId
func GetStakeIDAndAssetIdFromStr(clientChainLzId uint64, stakerAddress string, assetsAddress string) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(clientChainLzId)
	if stakerAddress != "" {
		stakeId = strings.Join([]string{strings.ToLower(stakerAddress), clientChainLzIdStr}, "_")
	}

	if assetsAddress != "" {
		assetId = strings.Join([]string{strings.ToLower(assetsAddress), clientChainLzIdStr}, "_")
	}
	return
}

// UpdateAssetValue It's used to update asset state,negative or positive `changeValue` represents a decrease or increase in the asset state
// newValue = valueToUpdate + changeVale
func UpdateAssetValue(valueToUpdate *math.Int, changeValue *math.Int) error {
	if valueToUpdate == nil || changeValue == nil {
		return errorsmod.Wrap(ErrInputPointerIsNil, fmt.Sprintf("valueToUpdate:%v,changeValue:%v", valueToUpdate, changeValue))
	}

	if !changeValue.IsNil() {
		if changeValue.IsNegative() {
			if valueToUpdate.LT(changeValue.Neg()) {
				return errorsmod.Wrap(ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("valueToUpdate:%s,changeValue:%s", *valueToUpdate, *changeValue))
			}
		}
		if !changeValue.IsZero() {
			*valueToUpdate = valueToUpdate.Add(*changeValue)
		}
	}
	return nil
}

// UpdateAssetDecValue It's used to update asset state,negative or positive `changeValue` represents a decrease or increase in the asset state
// newValue = valueToUpdate + changeVale
func UpdateAssetDecValue(valueToUpdate *math.LegacyDec, changeValue *math.LegacyDec) error {
	if valueToUpdate == nil || changeValue == nil {
		return errorsmod.Wrap(ErrInputPointerIsNil, fmt.Sprintf("valueToUpdate:%v,changeValue:%v", valueToUpdate, changeValue))
	}

	if !changeValue.IsNil() {
		if changeValue.IsNegative() {
			if valueToUpdate.LT(changeValue.Neg()) {
				return errorsmod.Wrap(ErrSubAmountIsMoreThanOrigin, fmt.Sprintf("valueToUpdate:%s,changeValue:%s", *valueToUpdate, *changeValue))
			}
		}
		if !changeValue.IsZero() {
			*valueToUpdate = valueToUpdate.Add(*changeValue)
		}
	}
	return nil
}
