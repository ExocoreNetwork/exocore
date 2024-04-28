package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
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

// UpdateOptedInAssetsState will update the USD share state related to asset, operator and AVS when
// the asset amount changes caused by delegation, undelegation, slashStaker and slashOperator.
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
	assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
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
				Amount: opAmount,
				Value:  opUSDValue,
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
	// avsAddr should be an evm contract address or a chain id.
	if !common.IsHexAddress(avsAddr) {
		if avsAddr != ctx.ChainID() { // TODO: other chain ids besides this chain's.
			return types.ErrInvalidAvsAddr
		}
	}
	// check optedIn info
	if k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrAlreadyOptedIn
	}
	// get the assets supported by the AVS
	// TODO: for x/dogfood, read the value from the params.
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avsAddr)
	if err != nil {
		return err
	}

	// get the Assets opted in the operator
	operatorAssets, err := k.assetsKeeper.GetOperatorAssetInfos(ctx, operatorAddress, avsSupportedAssets)
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
		assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
		if err != nil {
			return err
		}
		assetInfoRecord[assetID] = &AssetPriceAndDecimal{
			Price:        price,
			PriceDecimal: decimal,
			Decimal:      assetInfo.AssetBasicInfo.Decimals,
		}
		assetUSDValue := CalculateShare(operatorAssetState.TotalAmount, price, assetInfo.AssetBasicInfo.Decimals, decimal)
		operatorUSDValue := CalculateShare(operatorAssetState.OperatorAmount, price, assetInfo.AssetBasicInfo.Decimals, decimal)
		operatorOwnAssetUSDValue = operatorOwnAssetUSDValue.Add(operatorUSDValue)

		// UpdateStateForAsset
		changeState := types.OptedInAssetStateChange{
			Amount: operatorAssetState.TotalAmount,
			Value:  assetUSDValue,
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
	err = k.SetOptedInfo(ctx, operatorAddress.String(), avsAddr, optedInfo)
	if err != nil {
		return err
	}
	return nil
}

// OptOut call this function to opt out of AVS
func (k *Keeper) OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) error {
	if !k.IsOperator(ctx, operatorAddress) {
		return delegationtypes.ErrOperatorNotExist
	}
	// check optedIn info
	if !k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrNotOptedIn
	}
	if !common.IsHexAddress(avsAddr) {
		if avsAddr == ctx.ChainID() {
			found, _ := k.getOperatorConsKeyForChainID(ctx, operatorAddress, avsAddr)
			if found {
				// if the key exists, it should be in the process of being removed.
				// TODO: if slashing is moved to a snapshot approach, opt out should only be
				// performed if the key doesn't exist.
				if !k.IsOperatorRemovingKeyFromChainID(ctx, operatorAddress, avsAddr) {
					return types.ErrOperatorNotRemovingKey
				}
			}
		} else {
			return types.ErrInvalidAvsAddr
		}
	}

	// get the assets supported by the AVS
	avsSupportedAssets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avsAddr)
	if err != nil {
		return err
	}
	// get the Assets opted in the operator
	operatorAssets, err := k.assetsKeeper.GetOperatorAssetInfos(ctx, operatorAddress, avsSupportedAssets)
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
	handleFunc := func(info *types.OptedInfo) {
		// #nosec G701
		info.OptedOutHeight = uint64(ctx.BlockHeight())
	}
	err = k.HandleOptedInfo(ctx, operatorAddress.String(), avsAddr, handleFunc)
	if err != nil {
		return err
	}
	return nil
}
