package types

import (
	"fmt"
	"strings"

	"github.com/ExocoreNetwork/exocore/utils"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	CrossChainActionLength       = 1
	CrossChainOpAmountLength     = 32
	GeneralAssetsAddrLength      = 32
	GeneralClientChainAddrLength = 32
	ClientChainLzIDIndexInTopics = 0
	ExoCoreOperatorAddrLength    = 42

	// MaxDecimal is set to prevent the overflow
	// during the calculation of share and usd value.
	MaxDecimal                  = 18
	MaxChainTokenNameLength     = 50
	MaxChainTokenMetaInfoLength = 200

	MinClientChainAddrLength = 20
)

const (
	Deposit CrossChainOpType = iota
	WithdrawPrincipal
	WithDrawReward
	DelegateTo
	UndelegateFrom
	Slash
)

type GeneralAssetsAddr [32]byte

type GeneralClientChainAddr [32]byte

type CrossChainOpType uint8

type WithdrawerAddress [32]byte

// DeltaStakerSingleAsset This is a struct to describe the desired change that matches with
// the StakerAssetInfo
type DeltaStakerSingleAsset StakerAssetInfo

// DeltaOperatorSingleAsset This is a struct to describe the desired change that matches
// with the OperatorAssetInfo
type DeltaOperatorSingleAsset OperatorAssetInfo

type CreateQueryContext func(height int64, prove bool) (sdk.Context, error)

// GetStakeIDAndAssetID stakerID = stakerAddress+'_'+clientChainLzID,assetID =
// assetAddress+'_'+clientChainLzID
func GetStakeIDAndAssetID(
	clientChainLzID uint64,
	stakerAddress []byte,
	assetsAddress []byte,
) (stakeID string, assetID string) {
	clientChainLzIDStr := hexutil.EncodeUint64(clientChainLzID)
	if stakerAddress != nil {
		stakeID = strings.Join([]string{hexutil.Encode(stakerAddress), clientChainLzIDStr}, utils.DelimiterForID)
	}

	if assetsAddress != nil {
		assetID = strings.Join([]string{hexutil.Encode(assetsAddress), clientChainLzIDStr}, utils.DelimiterForID)
	}
	return
}

// GetStakeIDAndAssetIDFromStr stakerID = stakerAddress+'_'+clientChainLzID,assetID =
// assetAddress+'_'+clientChainLzID
func GetStakeIDAndAssetIDFromStr(
	clientChainLzID uint64,
	stakerAddress string,
	assetsAddress string,
) (stakeID string, assetID string) {
	// hexutil always returns lowercase values
	clientChainLzIDStr := hexutil.EncodeUint64(clientChainLzID)
	if stakerAddress != "" {
		stakeID = strings.Join(
			[]string{strings.ToLower(stakerAddress), clientChainLzIDStr},
			utils.DelimiterForID,
		)
	}

	if assetsAddress != "" {
		assetID = strings.Join(
			[]string{strings.ToLower(assetsAddress), clientChainLzIDStr},
			utils.DelimiterForID,
		)
	}
	return
}

// UpdateAssetValue It's used to update asset state,negative or positive `changeValue`
// represents a decrease or increase in the asset state
// newValue = valueToUpdate + changeVale
func UpdateAssetValue(valueToUpdate *math.Int, changeValue *math.Int) error {
	if valueToUpdate == nil || changeValue == nil {
		return errorsmod.Wrap(
			ErrInputPointerIsNil,
			fmt.Sprintf("valueToUpdate:%v,changeValue:%v", valueToUpdate, changeValue),
		)
	}

	if !changeValue.IsNil() {
		if changeValue.IsNegative() {
			if valueToUpdate.LT(changeValue.Neg()) {
				return errorsmod.Wrap(
					ErrSubAmountIsMoreThanOrigin,
					fmt.Sprintf(
						"valueToUpdate:%s,changeValue:%s",
						*valueToUpdate,
						*changeValue,
					),
				)
			}
		}
		if !changeValue.IsZero() {
			*valueToUpdate = valueToUpdate.Add(*changeValue)
		}
	}
	return nil
}

// UpdateAssetDecValue It's used to update asset state,negative or positive `changeValue`
// represents a decrease or increase in the asset state
// newValue = valueToUpdate + changeVale
func UpdateAssetDecValue(valueToUpdate *math.LegacyDec, changeValue *math.LegacyDec) error {
	if valueToUpdate == nil || changeValue == nil {
		return errorsmod.Wrap(
			ErrInputPointerIsNil,
			fmt.Sprintf("valueToUpdate:%v,changeValue:%v", valueToUpdate, changeValue),
		)
	}

	if !changeValue.IsNil() {
		if changeValue.IsNegative() {
			if valueToUpdate.LT(changeValue.Neg()) {
				return errorsmod.Wrap(
					ErrSubAmountIsMoreThanOrigin,
					fmt.Sprintf(
						"valueToUpdate:%s,changeValue:%s",
						*valueToUpdate,
						*changeValue,
					),
				)
			}
		}
		if !changeValue.IsZero() {
			*valueToUpdate = valueToUpdate.Add(*changeValue)
		}
	}
	return nil
}
