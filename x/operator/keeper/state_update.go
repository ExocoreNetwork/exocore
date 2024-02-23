package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	types2 "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
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

func (k *Keeper) UpdateOptedInAssetsState(ctx sdk.Context, stakerId, assetId, operatorAddr string, opAmount sdkmath.Int) error {
	//get the AVS opted-in by the operator
	avsList, err := k.GetOptedInAVSForOperator(ctx, operatorAddr)
	if err != nil {
		return err
	}
	//get price and priceDecimal from oracle
	price, decimal, err := k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetId)
	if err != nil {
		return err
	}

	//get the decimal of asset
	assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetId)
	if err != nil {
		return err
	}
	opUSDValue := CalculateShare(opAmount, price, assetInfo.AssetBasicInfo.Decimals, decimal)
	for _, avs := range avsList {
		//get the assets supported by the AVS
		avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avs)
		if err != nil {
			return err
		}

		if _, ok := avsSupportedAssets[assetId]; ok {
			//UpdateStakerShare
			err = k.UpdateStakerShare(ctx, avs, stakerId, operatorAddr, opUSDValue)
			if err != nil {
				return err
			}

			//UpdateStateForAsset
			changeState := types.AssetOptedInState{
				Amount: opAmount,
				Value:  opUSDValue,
			}
			err = k.UpdateStateForAsset(ctx, assetId, avs, operatorAddr, changeState)
			if err != nil {
				return err
			}

			//UpdateOperatorShare
			err = k.UpdateOperatorShare(ctx, avs, operatorAddr, opUSDValue)
			if err != nil {
				return err
			}

			//UpdateAVSShare
			err = k.UpdateAVSShare(ctx, avs, opUSDValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// OptIn call this function to opt in AVS
func (k *Keeper) OptIn(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string) error {
	//check optedIn info
	if k.IsOptedIn(ctx, operatorAddress.String(), AVSAddr) {
		return types.ErrAlreadyOptedIn
	}
	//get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, AVSAddr)
	if err != nil {
		return err
	}

	//get the Assets opted in the operator
	operatorAssets, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress, avsSupportedAssets)
	if err != nil {
		return err
	}

	totalAssetUSDValue := sdkmath.LegacyNewDec(0)
	operatorOwnAssetUSDValue := sdkmath.LegacyNewDec(0)
	assetFilter := make(map[string]interface{})
	assetInfoRecord := make(map[string]*AssetPriceAndDecimal)

	for assetId, operatorAssetState := range operatorAssets {
		//get price and priceDecimal from oracle
		price, decimal, err := k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetId)
		if err != nil {
			return err
		}

		//get the decimal of asset
		assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetId)
		if err != nil {
			return err
		}
		assetInfoRecord[assetId] = &AssetPriceAndDecimal{
			Price:        price,
			PriceDecimal: decimal,
			Decimal:      assetInfo.AssetBasicInfo.Decimals,
		}
		assetUSDValue := CalculateShare(operatorAssetState.TotalAmountOrWantChangeValue, price, assetInfo.AssetBasicInfo.Decimals, decimal)
		operatorUSDValue := CalculateShare(operatorAssetState.OperatorOwnAmountOrWantChangeValue, price, assetInfo.AssetBasicInfo.Decimals, decimal)
		operatorOwnAssetUSDValue = operatorOwnAssetUSDValue.Add(operatorUSDValue)

		//UpdateStateForAsset
		changeState := types.AssetOptedInState{
			Amount: operatorAssetState.TotalAmountOrWantChangeValue,
			Value:  assetUSDValue,
		}
		err = k.UpdateStateForAsset(ctx, assetId, AVSAddr, operatorAddress.String(), changeState)
		if err != nil {
			return err
		}
		totalAssetUSDValue = totalAssetUSDValue.Add(assetUSDValue)
		assetFilter[assetId] = nil
	}

	//update the share value of operator itself, the input stakerId should be empty
	err = k.UpdateStakerShare(ctx, AVSAddr, "", operatorAddress.String(), operatorOwnAssetUSDValue)
	if err != nil {
		return err
	}

	//UpdateAVSShare
	err = k.UpdateAVSShare(ctx, AVSAddr, totalAssetUSDValue)
	if err != nil {
		return err
	}
	//UpdateOperatorShare
	err = k.UpdateOperatorShare(ctx, AVSAddr, operatorAddress.String(), totalAssetUSDValue)
	if err != nil {
		return err
	}

	//UpdateStakerShare
	relatedAssetsState, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetFilter)
	if err != nil {
		return err
	}

	for stakerId, assetState := range relatedAssetsState {
		stakerAssetsUSDValue := sdkmath.LegacyNewDec(0)
		for assetId, amount := range assetState {
			singleAssetUSDValue := CalculateShare(amount.CanUndelegationAmount, assetInfoRecord[assetId].Price, assetInfoRecord[assetId].Decimal, assetInfoRecord[assetId].PriceDecimal)
			stakerAssetsUSDValue = stakerAssetsUSDValue.Add(singleAssetUSDValue)
		}

		err = k.UpdateStakerShare(ctx, AVSAddr, stakerId, operatorAddress.String(), stakerAssetsUSDValue)
		if err != nil {
			return err
		}
	}

	//update opted-in info
	slashContract, err := k.avsKeeper.GetAvsSlashContract(ctx, AVSAddr)
	if err != nil {
		return err
	}
	optedInfo := &types.OptedInfo{
		SlashContract:  slashContract,
		OptedInHeight:  uint64(ctx.BlockHeight()),
		OptedOutHeight: types.DefaultOptedOutHeight,
	}
	err = k.UpdateOptedInfo(ctx, operatorAddress.String(), AVSAddr, optedInfo)
	if err != nil {
		return err
	}
	return nil
}

// OptOut call this function to opt out of AVS
func (k *Keeper) OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string) error {
	//check optedIn info
	if !k.IsOptedIn(ctx, operatorAddress.String(), AVSAddr) {
		return types.ErrNotOptedIn
	}

	//get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, AVSAddr)
	if err != nil {
		return err
	}
	//get the Assets opted in the operator
	operatorAssets, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress, avsSupportedAssets)
	if err != nil {
		return err
	}

	assetFilter := make(map[string]interface{})

	for assetId := range operatorAssets {
		err = k.DeleteAssetState(ctx, assetId, AVSAddr, operatorAddress.String())
		if err != nil {
			return err
		}
		assetFilter[assetId] = nil
	}

	avsOperatorTotalValue, err := k.GetOperatorShare(ctx, AVSAddr, operatorAddress.String())
	if err != nil {
		return err
	}
	if avsOperatorTotalValue.IsNegative() {
		return errorsmod.Wrap(types.ErrTheValueIsNegative, fmt.Sprintf("OptOut,avsOperatorTotalValue:%suite", avsOperatorTotalValue))
	}

	//delete the share value of operator itself, the input stakerId should be empty
	err = k.DeleteStakerShare(ctx, AVSAddr, "", operatorAddress.String())
	if err != nil {
		return err
	}

	//UpdateAVSShare
	err = k.UpdateAVSShare(ctx, AVSAddr, avsOperatorTotalValue.Neg())
	if err != nil {
		return err
	}
	//DeleteOperatorShare
	err = k.DeleteOperatorShare(ctx, AVSAddr, operatorAddress.String())
	if err != nil {
		return err
	}

	//DeleteStakerShare
	relatedAssetsState, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetFilter)
	if err != nil {
		return err
	}
	for stakerId := range relatedAssetsState {
		err = k.DeleteStakerShare(ctx, AVSAddr, stakerId, operatorAddress.String())
		if err != nil {
			return err
		}
	}

	//set opted-out height
	optedInfo, err := k.GetOptedInfo(ctx, operatorAddress.String(), AVSAddr)
	if err != nil {
		return err
	}
	optedInfo.OptedOutHeight = uint64(ctx.BlockHeight())
	err = k.UpdateOptedInfo(ctx, operatorAddress.String(), AVSAddr, optedInfo)
	if err != nil {
		return err
	}
	return nil
}

// GetAssetsAmountToSlash It will slash the assets that are opting into AVS first, and if there isn't enough to slash, then it will slash the assets that have requested to undelegate but still locked.
func (k *Keeper) GetAssetsAmountToSlash(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string, occurredSateHeight int64, slashProportion sdkmath.LegacyDec) (*SlashAssets, error) {
	ret := &SlashAssets{
		slashStakerInfo:   make(map[string]map[string]*slashAmounts, 0),
		slashOperatorInfo: make(map[string]*slashAmounts, 0),
	}

	//get the state when the slash occurred
	historicalSateCtx, err := types2.ContextForHistoricalState(ctx, occurredSateHeight)
	if err != nil {
		return nil, err
	}
	//get assetsInfo supported by AVS
	assetsFilter, err := k.avsKeeper.GetAvsSupportedAssets(historicalSateCtx, AVSAddr)
	if err != nil {
		return nil, err
	}
	historyStakerAssets, err := k.delegationKeeper.DelegationStateByOperatorAssets(historicalSateCtx, operatorAddress.String(), assetsFilter)
	if err != nil {
		return nil, err
	}

	//get the Assets opted in the operator
	historyOperatorAssetsState, err := k.restakingStateKeeper.GetOperatorAssetInfos(historicalSateCtx, operatorAddress, assetsFilter)
	if err != nil {
		return nil, err
	}

	//calculate the actual slash amount according to the history and current state
	currentStakerAssets, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetsFilter)
	if err != nil {
		return nil, err
	}
	//get the Assets opted in the operator
	currentOperatorAssetsState, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress, assetsFilter)
	if err != nil {
		return nil, err
	}

	//calculate the actual slash amount for staker
	for stakerId, assetsState := range currentStakerAssets {
		if historyAssetState, ok := historyStakerAssets[stakerId]; ok {
			for assetId, curState := range assetsState {
				if historyState, isExist := historyAssetState[assetId]; isExist {
					if _, exist := ret.slashStakerInfo[stakerId]; !exist {
						ret.slashStakerInfo[stakerId] = make(map[string]*slashAmounts, 0)
					}
					shouldSlashAmount := slashProportion.MulInt(historyState.CanUndelegationAmount).TruncateInt()
					if curState.CanUndelegationAmount.LT(shouldSlashAmount) {
						ret.slashStakerInfo[stakerId][assetId].AmountFromOptedIn = curState.CanUndelegationAmount
						remainShouldSlash := shouldSlashAmount.Sub(curState.CanUndelegationAmount)
						if curState.UndelegatableAmountAfterSlash.LT(remainShouldSlash) {
							ret.slashStakerInfo[stakerId][assetId].AmountFromUnbonding = curState.UndelegatableAmountAfterSlash
						} else {
							ret.slashStakerInfo[stakerId][assetId].AmountFromUnbonding = remainShouldSlash
						}
					} else {
						ret.slashStakerInfo[stakerId][assetId].AmountFromOptedIn = shouldSlashAmount
					}
				}
			}
		}
	}

	//calculate the actual slash amount for operator
	for assetId, curAssetState := range currentOperatorAssetsState {
		if historyAssetState, ok := historyOperatorAssetsState[assetId]; ok {
			shouldSlashAmount := slashProportion.MulInt(historyAssetState.OperatorOwnAmountOrWantChangeValue).TruncateInt()
			if curAssetState.OperatorOwnAmountOrWantChangeValue.LT(shouldSlashAmount) {
				ret.slashOperatorInfo[assetId].AmountFromOptedIn = curAssetState.OperatorOwnAmountOrWantChangeValue
				remainShouldSlash := shouldSlashAmount.Sub(curAssetState.OperatorOwnAmountOrWantChangeValue)
				if curAssetState.OperatorUnbondableAmountAfterSlash.LT(remainShouldSlash) {
					ret.slashOperatorInfo[assetId].AmountFromUnbonding = curAssetState.OperatorUnbondableAmountAfterSlash
				} else {
					ret.slashOperatorInfo[assetId].AmountFromUnbonding = remainShouldSlash
				}
			} else {
				ret.slashOperatorInfo[assetId].AmountFromOptedIn = shouldSlashAmount
			}
		}
	}
	return ret, nil
}

func (k *Keeper) SlashStaker(ctx sdk.Context, operatorAddress sdk.AccAddress, slashStakerInfo map[string]map[string]*slashAmounts, executeHeight uint64) error {
	for stakerId, slashAssets := range slashStakerInfo {
		for assetId, slashInfo := range slashAssets {
			//handle the state that needs to be updated when slashing both opted-in and unbonding assets
			//update delegation state
			delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
			delegatorAndAmount[operatorAddress.String()] = &delegationtype.DelegationAmounts{
				CanUndelegationAmount:         slashInfo.AmountFromOptedIn.Neg(),
				UndelegatableAmountAfterSlash: slashInfo.AmountFromUnbonding.Neg(),
			}
			err := k.delegationKeeper.UpdateDelegationState(ctx, stakerId, assetId, delegatorAndAmount)
			if err != nil {
				return err
			}
			err = k.delegationKeeper.UpdateStakerDelegationTotalAmount(ctx, stakerId, assetId, slashInfo.AmountFromOptedIn.Neg())
			if err != nil {
				return err
			}

			slashSumValue := slashInfo.AmountFromUnbonding.Add(slashInfo.AmountFromOptedIn)
			//update staker and operator assets state
			err = k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakerId, assetId, types2.StakerSingleAssetOrChangeInfo{
				TotalDepositAmountOrWantChangeValue: slashSumValue.Neg(),
			})
			if err != nil {
				return err
			}

			//Record the slash information for scheduled tasks and send it to the client chain once the veto duration expires.
			err = k.UpdateSlashAssetsState(ctx, assetId, stakerId, executeHeight, slashSumValue)
			if err != nil {
				return err
			}

			//handle the state that needs to be updated when slashing opted-in assets
			err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAddress, assetId, types2.OperatorSingleAssetOrChangeInfo{
				TotalAmountOrWantChangeValue: slashInfo.AmountFromOptedIn.Neg(),
			})
			if err != nil {
				return err
			}
			//decrease the related share value
			err = k.UpdateOptedInAssetsState(ctx, stakerId, assetId, operatorAddress.String(), slashInfo.AmountFromOptedIn.Neg())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *Keeper) SlashOperator(ctx sdk.Context, operatorAddress sdk.AccAddress, slashOperatorInfo map[string]*slashAmounts, executeHeight uint64) error {
	for assetId, slashInfo := range slashOperatorInfo {
		slashSumValue := slashInfo.AmountFromUnbonding.Add(slashInfo.AmountFromOptedIn)
		//handle the state that needs to be updated when slashing both opted-in and unbonding assets
		err := k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAddress, assetId, types2.OperatorSingleAssetOrChangeInfo{
			TotalAmountOrWantChangeValue:       slashSumValue.Neg(),
			OperatorOwnAmountOrWantChangeValue: slashInfo.AmountFromOptedIn.Neg(),
			OperatorUnbondableAmountAfterSlash: slashInfo.AmountFromUnbonding.Neg(),
		})
		if err != nil {
			return err
		}
		//Record the slash information for scheduled tasks and send it to the client chain once the veto duration expires.
		err = k.UpdateSlashAssetsState(ctx, assetId, operatorAddress.String(), executeHeight, slashSumValue)
		if err != nil {
			return err
		}

		//handle the state that needs to be updated when slashing opted-in assets
		//decrease the related share value
		err = k.UpdateOptedInAssetsState(ctx, "", assetId, operatorAddress.String(), slashInfo.AmountFromOptedIn.Neg())
		if err != nil {
			return err
		}
	}
	return nil
}

// Slash The occurredSateHeight should be the height that has the latest stable state.
func (k *Keeper) Slash(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr, slashContract, slashId string, occurredSateHeight int64, slashProportion sdkmath.LegacyDec) error {
	height := ctx.BlockHeight()
	if occurredSateHeight > height {
		return errorsmod.Wrap(types.ErrSlashOccurredHeight, fmt.Sprintf("occurredSateHeight:%d,curHeight:%d", occurredSateHeight, height))
	}

	//get the state when the slash occurred
	//get the opted-in info
	historicalSateCtx, err := types2.ContextForHistoricalState(ctx, occurredSateHeight)
	if err != nil {
		return err
	}
	if !k.IsOptedIn(ctx, operatorAddress.String(), AVSAddr) {
		return types.ErrNotOptedIn
	}
	optedInfo, err := k.GetOptedInfo(historicalSateCtx, operatorAddress.String(), AVSAddr)
	if err != nil {
		return err
	}
	if optedInfo.SlashContract != slashContract {
		return errorsmod.Wrap(types.ErrSlashContractNotMatch, fmt.Sprintf("input slashContract:%suite, opted-in slash contract:%suite", slashContract, optedInfo.SlashContract))
	}

	//todo: recording the slash event might be moved to the slash module
	slashInfo := types.OperatorSlashInfo{
		SlashContract:   slashContract,
		SlashHeight:     height,
		OccurredHeight:  occurredSateHeight,
		SlashProportion: slashProportion,
		ExecuteHeight:   height + types.SlashVetoDuration,
	}
	err = k.UpdateOperatorSlashInfo(ctx, operatorAddress.String(), AVSAddr, slashId, slashInfo)
	if err != nil {
		return err
	}

	// get the assets and amounts that should be slashed
	assetsSlashInfo, err := k.GetAssetsAmountToSlash(ctx, operatorAddress, AVSAddr, occurredSateHeight, slashProportion)
	if err != nil {
		return err
	}

	err = k.SlashStaker(ctx, operatorAddress, assetsSlashInfo.slashStakerInfo, uint64(slashInfo.ExecuteHeight))
	if err != nil {
		return err
	}

	err = k.SlashOperator(ctx, operatorAddress, assetsSlashInfo.slashOperatorInfo, uint64(slashInfo.ExecuteHeight))
	if err != nil {
		return err
	}
	return nil
}
