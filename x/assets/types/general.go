package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	CrossChainActionLength       = 1
	CrossChainOpAmountLength     = 32
	GeneralAssetsAddrLength      = 32
	GeneralClientChainAddrLength = 32

	ClientChainLzIDIndexInTopics = 0
	LzNonceIndexInTopics         = 2

	ExoCoreOperatorAddrLength = 42
)

const (
	Deposit CrossChainOpType = iota
	WithdrawPrinciple
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

// GetStakeIDAndAssetID stakerID = stakerAddress+'_'+clientChainLzID,assetID =
// assetAddress+'_'+clientChainLzID
func GetStakeIDAndAssetID(
	clientChainLzID uint64,
	stakerAddress []byte,
	assetsAddress []byte,
) (stakeID string, assetID string) {
	clientChainLzIDStr := hexutil.EncodeUint64(clientChainLzID)
	if stakerAddress != nil {
		stakeID = strings.Join([]string{hexutil.Encode(stakerAddress), clientChainLzIDStr}, "_")
	}

	if assetsAddress != nil {
		assetID = strings.Join([]string{hexutil.Encode(assetsAddress), clientChainLzIDStr}, "_")
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
			"_",
		)
	}

	if assetsAddress != "" {
		assetID = strings.Join(
			[]string{strings.ToLower(assetsAddress), clientChainLzIDStr},
			"_",
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

func ContextForHistoricalState(ctx sdk.Context, height int64) (sdk.Context, error) {
	if height < 0 {
		return sdk.Context{}, errorsmod.Wrap(
			sdkerrors.ErrInvalidHeight,
			fmt.Sprintf("height:%v", height),
		)
	}
	cms := ctx.MultiStore()
	lastBlockHeight := cms.LatestVersion()
	if lastBlockHeight == 0 {
		return sdk.Context{}, errorsmod.Wrap(
			sdkerrors.ErrInvalidHeight,
			"app is not ready; please wait for first block",
		)
	}
	if height > lastBlockHeight {
		return sdk.Context{},
			errorsmod.Wrap(
				sdkerrors.ErrInvalidHeight,
				"cannot query with height in the future; please provide a valid height",
			)
	}

	// when the caller did not provide a query height, manually inject the latest
	if height == 0 {
		height = lastBlockHeight
	}
	cacheMS, err := cms.CacheMultiStoreWithVersion(height)
	if err != nil {
		return sdk.Context{},
			errorsmod.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"failed to load state at height %d; %s (latest height: %d)",
				height,
				err,
				lastBlockHeight,
			)
	}

	// branch the commit-multistore for safety
	historicalStateCtx := sdk.NewContext(cacheMS, ctx.BlockHeader(), true, ctx.Logger()).
		WithMinGasPrices(ctx.MinGasPrices()).
		WithBlockHeight(height)
	if height != lastBlockHeight {
		rms, ok := cms.(*rootmulti.Store)
		if ok {
			cInfo, err := rms.GetCommitInfo(height)
			if cInfo != nil && err == nil {
				historicalStateCtx = historicalStateCtx.WithBlockTime(cInfo.Timestamp)
			}
		}
	}
	return historicalStateCtx, nil
}
