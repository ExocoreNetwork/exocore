package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/operator/types"
)

func (k Keeper) UpdateOptedInAssetsState(ctx sdk.Context, stakerId, assetId, operatorAddr string, opAmount sdkmath.Int) error {
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

	//opUsdValue = (opAmount*price*10^UsdValueDefaultDecimal)/(10^(asset.decimal+priceDecimal))
	value := opAmount.Mul(price).Mul(sdkmath.NewIntWithDecimal(1, int(types.UsdValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetInfo.AssetBasicInfo.Decimals)+int(decimal)))
	opUsdValue := sdkmath.LegacyNewDecFromBigIntWithPrec(value.BigInt(), int64(types.UsdValueDefaultDecimal))

	for _, avs := range avsList {
		//get the assets supported by the AVS
		avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avs)
		if err != nil {
			return err
		}

		if _, ok := avsSupportedAssets[assetId]; ok {
			//UpdateAVSOperatorStakerShareValue
			err = k.UpdateAVSOperatorStakerShareValue(ctx, avs, stakerId, operatorAddr, opUsdValue)
			if err != nil {
				return err
			}

			//UpdateOperatorAVSAssetsState
			changeState := types.AssetOptedInState{
				Amount: opAmount,
				Value:  opUsdValue,
			}
			err = k.UpdateOperatorAVSAssetsState(ctx, assetId, avs, operatorAddr, changeState)
			if err != nil {
				return err
			}

			//UpdateAVSOperatorTotalValue
			err = k.UpdateAVSOperatorTotalValue(ctx, avs, operatorAddr, opUsdValue)
			if err != nil {
				return err
			}

			//UpdateAVSTotalValue
			err = k.UpdateAVSTotalValue(ctx, avs, opUsdValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type AssetPriceAndDecimal struct {
	Price        sdkmath.Int
	PriceDecimal uint8
	Decimal      uint32
}

// OptIn call this function to opt in AVS
func (k Keeper) OptIn(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string) error {
	//check optedIn info
	if k.IsOptedIn(ctx, operatorAddress.String(), AVSAddr) {
		return types.ErrAlreadyOptedIn
	}

	//get the Assets opted in the operator
	operatorAssetsState, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress)
	if err != nil {
		return err
	}
	//get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, AVSAddr)
	if err != nil {
		return err
	}
	totalAssetUsdValue := sdkmath.LegacyNewDec(0)
	operatorOwnAssetUsdValue := sdkmath.LegacyNewDec(0)
	assetFilter := make(map[string]interface{})
	assetInfoRecord := make(map[string]*AssetPriceAndDecimal)

	for assetId := range avsSupportedAssets {
		operatorAssetState, ok := operatorAssetsState[assetId]
		if !ok {
			continue
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
		assetInfoRecord[assetId] = &AssetPriceAndDecimal{
			Price:        price,
			PriceDecimal: decimal,
			Decimal:      assetInfo.AssetBasicInfo.Decimals,
		}

		//assetValue = (amount*price*10^UsdValueDefaultDecimal)/(10^(asset.decimal+priceDecimal))
		assetValue := operatorAssetState.TotalAmountOrWantChangeValue.Mul(price).Mul(sdkmath.NewIntWithDecimal(1, int(types.UsdValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetInfo.AssetBasicInfo.Decimals)+int(decimal)))
		assetUsdValue := sdkmath.LegacyNewDecFromBigIntWithPrec(assetValue.BigInt(), int64(types.UsdValueDefaultDecimal))

		operatorOwnAssetValue := operatorAssetState.OperatorOwnAmountOrWantChangeValue.Mul(price).Mul(sdkmath.NewIntWithDecimal(1, int(types.UsdValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetInfo.AssetBasicInfo.Decimals)+int(decimal)))
		operatorOwnAssetUsdValue = operatorOwnAssetUsdValue.Add(sdkmath.LegacyNewDecFromBigIntWithPrec(operatorOwnAssetValue.BigInt(), int64(types.UsdValueDefaultDecimal)))

		//UpdateOperatorAVSAssetsState
		changeState := types.AssetOptedInState{
			Amount: operatorAssetState.TotalAmountOrWantChangeValue,
			Value:  assetUsdValue,
		}
		err = k.UpdateOperatorAVSAssetsState(ctx, assetId, AVSAddr, operatorAddress.String(), changeState)
		if err != nil {
			return err
		}
		totalAssetUsdValue = totalAssetUsdValue.Add(assetUsdValue)
		assetFilter[assetId] = nil
	}

	//update the share value of operator itself, the input stakerId should be empty
	err = k.UpdateAVSOperatorStakerShareValue(ctx, AVSAddr, "", operatorAddress.String(), operatorOwnAssetUsdValue)
	if err != nil {
		return err
	}

	//UpdateAVSTotalValue
	err = k.UpdateAVSTotalValue(ctx, AVSAddr, totalAssetUsdValue)
	if err != nil {
		return err
	}
	//UpdateAVSOperatorTotalValue
	err = k.UpdateAVSOperatorTotalValue(ctx, AVSAddr, operatorAddress.String(), totalAssetUsdValue)
	if err != nil {
		return err
	}

	//UpdateAVSOperatorStakerShareValue
	relatedAssetsState, err := k.delegationKeeper.GetDelegationStateByOperatorAndAssetList(ctx, operatorAddress.String(), assetFilter)
	if err != nil {
		return err
	}

	for stakerId, assetState := range relatedAssetsState {
		stakerAssetsUsdValue := sdkmath.LegacyNewDec(0)
		for assetId, amount := range assetState {
			singleAssetValue := amount.Mul(assetInfoRecord[assetId].Price).Mul(sdkmath.NewIntWithDecimal(1, int(types.UsdValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetInfoRecord[assetId].Decimal)+int(assetInfoRecord[assetId].PriceDecimal)))
			singleAssetUsdValue := sdkmath.LegacyNewDecFromBigIntWithPrec(singleAssetValue.BigInt(), int64(types.UsdValueDefaultDecimal))
			stakerAssetsUsdValue = stakerAssetsUsdValue.Add(singleAssetUsdValue)
		}

		err = k.UpdateAVSOperatorStakerShareValue(ctx, AVSAddr, stakerId, operatorAddress.String(), stakerAssetsUsdValue)
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
func (k Keeper) OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string) error {
	//check optedIn info
	if !k.IsOptedIn(ctx, operatorAddress.String(), AVSAddr) {
		return types.ErrNotOptedIn
	}

	//get the Assets opted in the operator
	operatorAssetsState, err := k.restakingStateKeeper.GetOperatorAssetInfos(ctx, operatorAddress)
	if err != nil {
		return err
	}
	//get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, AVSAddr)
	if err != nil {
		return err
	}
	assetFilter := make(map[string]interface{})

	for assetId := range avsSupportedAssets {
		_, ok := operatorAssetsState[assetId]
		if !ok {
			continue
		}
		err = k.DeleteOperatorAVSAssetsState(ctx, assetId, AVSAddr, operatorAddress.String())
		if err != nil {
			return err
		}
		assetFilter[assetId] = nil
	}

	avsOperatorTotalValue, err := k.GetAVSOperatorTotalValue(ctx, AVSAddr, operatorAddress.String())
	if err != nil {
		return err
	}
	if avsOperatorTotalValue.IsNegative() {
		return errorsmod.Wrap(types.ErrTheValueIsNegative, fmt.Sprintf("OptOut,avsOperatorTotalValue:%s", avsOperatorTotalValue))
	}

	//delete the share value of operator itself, the input stakerId should be empty
	err = k.DeleteAVSOperatorStakerShareValue(ctx, AVSAddr, "", operatorAddress.String())
	if err != nil {
		return err
	}

	//UpdateAVSTotalValue
	err = k.UpdateAVSTotalValue(ctx, AVSAddr, avsOperatorTotalValue.Neg())
	if err != nil {
		return err
	}
	//DeleteAVSOperatorTotalValue
	err = k.DeleteAVSOperatorTotalValue(ctx, AVSAddr, operatorAddress.String())
	if err != nil {
		return err
	}

	//DeleteAVSOperatorStakerShareValue
	relatedAssetsState, err := k.delegationKeeper.GetDelegationStateByOperatorAndAssetList(ctx, operatorAddress.String(), assetFilter)
	if err != nil {
		return err
	}
	for stakerId := range relatedAssetsState {
		err = k.DeleteAVSOperatorStakerShareValue(ctx, AVSAddr, stakerId, operatorAddress.String())
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

type SlashAssetsAndAmount struct {
	slashFromStakerUnbonding map[string]map[string]sdkmath.Int
	slashFromStakerOptedIn   map[string]map[string]sdkmath.Int

	slashFromOperatorUnbonding map[string]sdkmath.Int
	slashFromOperatorOptedIn   map[string]sdkmath.Int
}

// GetAssetsAndAmountToSlash It will slash the assets that have been undelegated but still locked first, and if there isn't enough to slash, then it will slash the assets that are opting into AVS.
func (k Keeper) GetAssetsAndAmountToSlash(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string, occurredHeight int64, slashProportion sdkmath.LegacyDec) (*SlashAssetsAndAmount, error) {
	ret := SlashAssetsAndAmount{
		slashFromStakerUnbonding:   make(map[string]map[string]sdkmath.Int, 0),
		slashFromStakerOptedIn:     make(map[string]map[string]sdkmath.Int, 0),
		slashFromOperatorUnbonding: make(map[string]sdkmath.Int, 0),
		slashFromOperatorOptedIn:   make(map[string]sdkmath.Int, 0),
	}

}

func (k Keeper) Slash(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr, slashContract, slashId string, occurredHeight int64, slashProportion sdkmath.LegacyDec) error {
	height := ctx.BlockHeight()
	if occurredHeight > height {
		return errorsmod.Wrap(types.ErrSlashOccurredHeight, fmt.Sprintf("occurredHeight:%d,curHeight:%d", occurredHeight, height))
	}

	//get the state when the slash occurred
	//get the opted-in info
	ctx = ctx.WithBlockHeight(occurredHeight)
	optedInfo, err := k.GetOptedInfo(ctx, operatorAddress.String(), AVSAddr)
	if err != nil {
		return err
	}
	// reset context height
	ctx = ctx.WithBlockHeight(height)

	if optedInfo.SlashContract != slashContract {
		return errorsmod.Wrap(types.ErrSlashContractNotMatch, fmt.Sprintf("input slashContract:%s, opted-in slash contract:%s", slashContract, optedInfo.SlashContract))
	}

	//todo: recording the slash event might be moved to the slash module
	slashInfo := types.OperatorSlashInfo{
		SlashContract:   slashContract,
		SlashHeight:     height,
		OccurredHeight:  occurredHeight,
		SlashProportion: slashProportion,
		ExecuteHeight:   height + types.SlashVetoDuration,
	}
	err = k.UpdateOperatorSlashInfo(ctx, operatorAddress.String(), AVSAddr, slashId, slashInfo)
	if err != nil {
		return err
	}

	// get the assets and amounts that should be slashed
	assetsSlashInfo, err := k.GetAssetsAndAmountToSlash(ctx, operatorAddress, AVSAddr, occurredHeight, slashProportion)
	if err != nil {
		return err
	}

	// update delegation assets state

	// update staker assets state

	// update opted-in state

}
