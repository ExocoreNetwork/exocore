package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	types2 "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AssetPriceAndDecimal struct {
	Price        sdkmath.Int
	PriceDecimal uint8
	Decimal      uint32
}

type slashAmounts struct {
	AmountFromUnbonding sdkmath.Int
	AmountFromOptedIn   sdkmath.Int
}
type SlashAssets struct {
	slashStakerInfo   map[string]map[string]*slashAmounts
	slashOperatorInfo map[string]*slashAmounts
}

func (k *Keeper) UpdateOptedInAssetsState(ctx sdk.Context, stakerID, assetID, operatorAddr string, opAmount sdkmath.Int) error {
	// get the AVS opted-in by the operator
	avsList, err := k.GetOptedInAVSForOperator(ctx, operatorAddr)
	if err != nil {
		return err
	}
	// get price and priceDecimal from oracle
	price, decimal, err := k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetID)
	if err != nil {
		return err
	}

	// get the decimal of asset
	assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetID)
	if err != nil {
		return err
	}
	opUSDValue := CalculateShare(opAmount, price, assetInfo.AssetBasicInfo.Decimals, decimal)
	for _, avs := range avsList {
		// get the assets supported by the AVS
		avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avs)
		if err != nil {
			return err
		}

		if _, ok := avsSupportedAssets[assetID]; ok {
			// UpdateStakerShare
			err = k.UpdateStakerShare(ctx, avs, stakerID, operatorAddr, opUSDValue)
			if err != nil {
				return err
			}

			// UpdateStateForAsset
			changeState := types.OptedInAssetStateChange{
				ChangeForAmount: opAmount,
				ChangeForValue:  opUSDValue,
			}
			err = k.UpdateStateForAsset(ctx, assetID, avs, operatorAddr, changeState)
			if err != nil {
				return err
			}

			// UpdateOperatorShare
			err = k.UpdateOperatorShare(ctx, avs, operatorAddr, opUSDValue)
			if err != nil {
				return err
			}

			// UpdateAVSShare
			err = k.UpdateAVSShare(ctx, avs, opUSDValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// OptIn call this function to opt in AVS
func (k *Keeper) OptIn(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) error {
	// check optedIn info
	if k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrAlreadyOptedIn
	}
	// get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avsAddr)
	if err != nil {
		return err
	}

	// get the Assets opted in the operator
	operatorAssets, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress, avsSupportedAssets)
	if err != nil {
		return err
	}

	totalAssetUSDValue := sdkmath.LegacyNewDec(0)
	operatorOwnAssetUSDValue := sdkmath.LegacyNewDec(0)
	assetFilter := make(map[string]interface{})
	assetInfoRecord := make(map[string]*AssetPriceAndDecimal)

	for assetID, operatorAssetState := range operatorAssets {
		// get price and priceDecimal from oracle
		price, decimal, err := k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetID)
		if err != nil {
			return err
		}

		// get the decimal of asset
		assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetID)
		if err != nil {
			return err
		}
		assetInfoRecord[assetID] = &AssetPriceAndDecimal{
			Price:        price,
			PriceDecimal: decimal,
			Decimal:      assetInfo.AssetBasicInfo.Decimals,
		}
		assetUSDValue := CalculateShare(operatorAssetState.TotalAmount, price, assetInfo.AssetBasicInfo.Decimals, decimal)
		operatorUSDValue := CalculateShare(operatorAssetState.OperatorOwnAmount, price, assetInfo.AssetBasicInfo.Decimals, decimal)
		operatorOwnAssetUSDValue = operatorOwnAssetUSDValue.Add(operatorUSDValue)

		// UpdateStateForAsset
		changeState := types.OptedInAssetStateChange{
			ChangeForAmount: operatorAssetState.TotalAmount,
			ChangeForValue:  assetUSDValue,
		}
		err = k.UpdateStateForAsset(ctx, assetID, avsAddr, operatorAddress.String(), changeState)
		if err != nil {
			return err
		}
		totalAssetUSDValue = totalAssetUSDValue.Add(assetUSDValue)
		assetFilter[assetID] = nil
	}

	// update the share value of operator itself, the input stakerID should be empty
	err = k.UpdateStakerShare(ctx, avsAddr, "", operatorAddress.String(), operatorOwnAssetUSDValue)
	if err != nil {
		return err
	}

	// UpdateAVSShare
	err = k.UpdateAVSShare(ctx, avsAddr, totalAssetUSDValue)
	if err != nil {
		return err
	}
	// UpdateOperatorShare
	err = k.UpdateOperatorShare(ctx, avsAddr, operatorAddress.String(), totalAssetUSDValue)
	if err != nil {
		return err
	}

	// UpdateStakerShare
	relatedAssetsState, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetFilter)
	if err != nil {
		return err
	}

	for stakerID, assetState := range relatedAssetsState {
		stakerAssetsUSDValue := sdkmath.LegacyNewDec(0)
		for assetID, amount := range assetState {
			singleAssetUSDValue := CalculateShare(amount.UndelegatableAmount, assetInfoRecord[assetID].Price, assetInfoRecord[assetID].Decimal, assetInfoRecord[assetID].PriceDecimal)
			stakerAssetsUSDValue = stakerAssetsUSDValue.Add(singleAssetUSDValue)
		}

		err = k.UpdateStakerShare(ctx, avsAddr, stakerID, operatorAddress.String(), stakerAssetsUSDValue)
		if err != nil {
			return err
		}
	}

	// update opted-in info
	slashContract, err := k.avsKeeper.GetAvsSlashContract(ctx, avsAddr)
	if err != nil {
		return err
	}
	optedInfo := &types.OptedInfo{
		SlashContract: slashContract,
		// #nosec G701
		OptedInHeight:  uint64(ctx.BlockHeight()),
		OptedOutHeight: types.DefaultOptedOutHeight,
	}
	err = k.UpdateOptedInfo(ctx, operatorAddress.String(), avsAddr, optedInfo)
	if err != nil {
		return err
	}
	return nil
}

// OptOut call this function to opt out of AVS
func (k *Keeper) OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) error {
	// check optedIn info
	if !k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrNotOptedIn
	}

	// get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avsAddr)
	if err != nil {
		return err
	}
	// get the Assets opted in the operator
	operatorAssets, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress, avsSupportedAssets)
	if err != nil {
		return err
	}

	assetFilter := make(map[string]interface{})

	for assetID := range operatorAssets {
		err = k.DeleteAssetState(ctx, assetID, avsAddr, operatorAddress.String())
		if err != nil {
			return err
		}
		assetFilter[assetID] = nil
	}

	avsOperatorTotalValue, err := k.GetOperatorShare(ctx, avsAddr, operatorAddress.String())
	if err != nil {
		return err
	}
	if avsOperatorTotalValue.IsNegative() {
		return errorsmod.Wrap(types.ErrTheValueIsNegative, fmt.Sprintf("OptOut,avsOperatorTotalValue:%suite", avsOperatorTotalValue))
	}

	// delete the share value of operator itself, the input stakerID should be empty
	err = k.DeleteStakerShare(ctx, avsAddr, "", operatorAddress.String())
	if err != nil {
		return err
	}

	// UpdateAVSShare
	err = k.UpdateAVSShare(ctx, avsAddr, avsOperatorTotalValue.Neg())
	if err != nil {
		return err
	}
	// DeleteOperatorShare
	err = k.DeleteOperatorShare(ctx, avsAddr, operatorAddress.String())
	if err != nil {
		return err
	}

	// DeleteStakerShare
	relatedAssetsState, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetFilter)
	if err != nil {
		return err
	}
	for stakerID := range relatedAssetsState {
		err = k.DeleteStakerShare(ctx, avsAddr, stakerID, operatorAddress.String())
		if err != nil {
			return err
		}
	}

	// set opted-out height
	optedInfo, err := k.GetOptedInfo(ctx, operatorAddress.String(), avsAddr)
	if err != nil {
		return err
	}
	// #nosec G701
	optedInfo.OptedOutHeight = uint64(ctx.BlockHeight())
	err = k.UpdateOptedInfo(ctx, operatorAddress.String(), avsAddr, optedInfo)
	if err != nil {
		return err
	}
	return nil
}

// GetAssetsAmountToSlash It will slash the assets that are opting into AVS first, and if there isn't enough to slash, then it will slash the assets that have requested to undelegate but still locked.
func (k *Keeper) GetAssetsAmountToSlash(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string, occurredSateHeight int64, slashProportion sdkmath.LegacyDec) (*SlashAssets, error) {
	ret := &SlashAssets{
		slashStakerInfo:   make(map[string]map[string]*slashAmounts, 0),
		slashOperatorInfo: make(map[string]*slashAmounts, 0),
	}

	// get the state when the slash occurred
	historicalSateCtx, err := types2.ContextForHistoricalState(ctx, occurredSateHeight)
	if err != nil {
		return nil, err
	}
	// get assetsInfo supported by AVS
	assetsFilter, err := k.avsKeeper.GetAvsSupportedAssets(historicalSateCtx, avsAddr)
	if err != nil {
		return nil, err
	}
	historyStakerAssets, err := k.delegationKeeper.DelegationStateByOperatorAssets(historicalSateCtx, operatorAddress.String(), assetsFilter)
	if err != nil {
		return nil, err
	}

	// get the Assets opted in the operator
	historyOperatorAssetsState, err := k.restakingStateKeeper.GetOperatorAssetInfos(historicalSateCtx, operatorAddress, assetsFilter)
	if err != nil {
		return nil, err
	}

	// calculate the actual slash amount according to the history and current state
	currentStakerAssets, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetsFilter)
	if err != nil {
		return nil, err
	}
	// get the Assets opted in the operator
	currentOperatorAssetsState, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress, assetsFilter)
	if err != nil {
		return nil, err
	}

	// calculate the actual slash amount for staker
	for stakerID, assetsState := range currentStakerAssets {
		if historyAssetState, ok := historyStakerAssets[stakerID]; ok {
			for assetID, curState := range assetsState {
				if historyState, isExist := historyAssetState[assetID]; isExist {
					if _, exist := ret.slashStakerInfo[stakerID]; !exist {
						ret.slashStakerInfo[stakerID] = make(map[string]*slashAmounts, 0)
					}
					shouldSlashAmount := slashProportion.MulInt(historyState.UndelegatableAmount).TruncateInt()
					if curState.UndelegatableAmount.LT(shouldSlashAmount) {
						ret.slashStakerInfo[stakerID][assetID].AmountFromOptedIn = curState.UndelegatableAmount
						remainShouldSlash := shouldSlashAmount.Sub(curState.UndelegatableAmount)
						if curState.UndelegatableAfterSlash.LT(remainShouldSlash) {
							ret.slashStakerInfo[stakerID][assetID].AmountFromUnbonding = curState.UndelegatableAfterSlash
						} else {
							ret.slashStakerInfo[stakerID][assetID].AmountFromUnbonding = remainShouldSlash
						}
					} else {
						ret.slashStakerInfo[stakerID][assetID].AmountFromOptedIn = shouldSlashAmount
					}
				}
			}
		}
	}

	// calculate the actual slash amount for operator
	for assetID, curAssetState := range currentOperatorAssetsState {
		if historyAssetState, ok := historyOperatorAssetsState[assetID]; ok {
			shouldSlashAmount := slashProportion.MulInt(historyAssetState.OperatorOwnAmount).TruncateInt()
			if curAssetState.OperatorOwnAmount.LT(shouldSlashAmount) {
				ret.slashOperatorInfo[assetID].AmountFromOptedIn = curAssetState.OperatorOwnAmount
				remainShouldSlash := shouldSlashAmount.Sub(curAssetState.OperatorOwnAmount)
				if curAssetState.OperatorUnbondableAmountAfterSlash.LT(remainShouldSlash) {
					ret.slashOperatorInfo[assetID].AmountFromUnbonding = curAssetState.OperatorUnbondableAmountAfterSlash
				} else {
					ret.slashOperatorInfo[assetID].AmountFromUnbonding = remainShouldSlash
				}
			} else {
				ret.slashOperatorInfo[assetID].AmountFromOptedIn = shouldSlashAmount
			}
		}
	}
	return ret, nil
}

func (k *Keeper) SlashStaker(ctx sdk.Context, operatorAddress sdk.AccAddress, slashStakerInfo map[string]map[string]*slashAmounts, processedHeight uint64) error {
	for stakerID, slashAssets := range slashStakerInfo {
		for assetID, slashInfo := range slashAssets {
			// handle the state that needs to be updated when slashing both opted-in and unbonding assets
			// update delegation state
			delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
			delegatorAndAmount[operatorAddress.String()] = &delegationtype.DelegationAmounts{
				UndelegatableAmount:     slashInfo.AmountFromOptedIn.Neg(),
				UndelegatableAfterSlash: slashInfo.AmountFromUnbonding.Neg(),
			}
			err := k.delegationKeeper.UpdateDelegationState(ctx, stakerID, assetID, delegatorAndAmount)
			if err != nil {
				return err
			}
			err = k.delegationKeeper.UpdateStakerDelegationTotalAmount(ctx, stakerID, assetID, slashInfo.AmountFromOptedIn.Neg())
			if err != nil {
				return err
			}

			slashSumValue := slashInfo.AmountFromUnbonding.Add(slashInfo.AmountFromOptedIn)
			// update staker and operator assets state
			err = k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, types2.StakerSingleAssetChangeInfo{
				ChangeForTotalDeposit: slashSumValue.Neg(),
			})
			if err != nil {
				return err
			}

			// Record the slash information for scheduled tasks and send it to the client chain once the veto duration expires.
			err = k.UpdateSlashAssetsState(ctx, assetID, stakerID, processedHeight, slashSumValue)
			if err != nil {
				return err
			}

			// handle the state that needs to be updated when slashing opted-in assets
			err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAddress, assetID, types2.OperatorSingleAssetChangeInfo{
				ChangeForTotalAmount: slashInfo.AmountFromOptedIn.Neg(),
			})
			if err != nil {
				return err
			}
			// decrease the related share value
			err = k.UpdateOptedInAssetsState(ctx, stakerID, assetID, operatorAddress.String(), slashInfo.AmountFromOptedIn.Neg())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *Keeper) SlashOperator(ctx sdk.Context, operatorAddress sdk.AccAddress, slashOperatorInfo map[string]*slashAmounts, processedHeight uint64) error {
	for assetID, slashInfo := range slashOperatorInfo {
		slashSumValue := slashInfo.AmountFromUnbonding.Add(slashInfo.AmountFromOptedIn)
		// handle the state that needs to be updated when slashing both opted-in and unbonding assets
		err := k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAddress, assetID, types2.OperatorSingleAssetChangeInfo{
			ChangeForTotalAmount:          slashSumValue.Neg(),
			ChangeForOperatorOwn:          slashInfo.AmountFromOptedIn.Neg(),
			ChangeForUnbondableAfterSlash: slashInfo.AmountFromUnbonding.Neg(),
		})
		if err != nil {
			return err
		}
		// Record the slash information for scheduled tasks and send it to the client chain once the veto duration expires.
		err = k.UpdateSlashAssetsState(ctx, assetID, operatorAddress.String(), processedHeight, slashSumValue)
		if err != nil {
			return err
		}

		// handle the state that needs to be updated when slashing opted-in assets
		// decrease the related share value
		err = k.UpdateOptedInAssetsState(ctx, "", assetID, operatorAddress.String(), slashInfo.AmountFromOptedIn.Neg())
		if err != nil {
			return err
		}
	}
	return nil
}

// Slash The occurredSateHeight should be the height that has the latest stable state.
func (k *Keeper) Slash(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr, slashContract, slashID string, occurredSateHeight int64, slashProportion sdkmath.LegacyDec) error {
	height := ctx.BlockHeight()
	if occurredSateHeight > height {
		return errorsmod.Wrap(types.ErrSlashOccurredHeight, fmt.Sprintf("occurredSateHeight:%d,curHeight:%d", occurredSateHeight, height))
	}

	// get the state when the slash occurred
	// get the opted-in info
	historicalSateCtx, err := types2.ContextForHistoricalState(ctx, occurredSateHeight)
	if err != nil {
		return err
	}
	if !k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrNotOptedIn
	}
	optedInfo, err := k.GetOptedInfo(historicalSateCtx, operatorAddress.String(), avsAddr)
	if err != nil {
		return err
	}
	if optedInfo.SlashContract != slashContract {
		return errorsmod.Wrap(types.ErrSlashContractNotMatch, fmt.Sprintf("input slashContract:%suite, opted-in slash contract:%suite", slashContract, optedInfo.SlashContract))
	}

	// todo: recording the slash event might be moved to the slash module
	slashInfo := types.OperatorSlashInfo{
		SlashContract:   slashContract,
		SubmittedHeight: height,
		EventHeight:     occurredSateHeight,
		SlashProportion: slashProportion,
		ProcessedHeight: height + types.SlashVetoDuration,
	}
	err = k.UpdateOperatorSlashInfo(ctx, operatorAddress.String(), avsAddr, slashID, slashInfo)
	if err != nil {
		return err
	}

	// get the assets and amounts that should be slashed
	assetsSlashInfo, err := k.GetAssetsAmountToSlash(ctx, operatorAddress, avsAddr, occurredSateHeight, slashProportion)
	if err != nil {
		return err
	}
	// #nosec G701
	err = k.SlashStaker(ctx, operatorAddress, assetsSlashInfo.slashStakerInfo, uint64(slashInfo.ProcessedHeight))
	if err != nil {
		return err
	}
	// #nosec G701
	err = k.SlashOperator(ctx, operatorAddress, assetsSlashInfo.slashOperatorInfo, uint64(slashInfo.ProcessedHeight))
	if err != nil {
		return err
	}
	return nil
}
