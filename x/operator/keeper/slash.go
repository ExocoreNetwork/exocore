package keeper

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetSlashIDForDogfood It use infractionType+'/'+'infractionHeight' as the slashID, because the slash event occurs in dogfood doesn't have a TxID. It isn't submitted through an external transaction.
func GetSlashIDForDogfood(infraction stakingtypes.Infraction, infractionHeight int64) string {
	// #nosec G701
	return string(assetstype.GetJoinedStoreKey(hexutil.EncodeUint64(uint64(infraction)), hexutil.EncodeUint64(uint64(infractionHeight))))
}

// GetAssetsAmountToSlash It will slash the assets that are opting into AVS first, and if there isn't enough to slash, then it will slash the assets that have requested to undelegate but still locked.
func (k *Keeper) GetAssetsAmountToSlash(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string, slashEventHeight int64, slashProportion sdkmath.LegacyDec) (*SlashAssets, error) {
	ret := &SlashAssets{
		slashStakerInfo:   make(map[string]map[string]*slashAmounts, 0),
		slashOperatorInfo: make(map[string]*slashAmounts, 0),
	}

	// get the state when the slash occurred
	historicalStateCtx, err := assetstype.ContextForHistoricalState(ctx, slashEventHeight)
	if err != nil {
		return nil, err
	}
	// get assetsInfo supported by AVS
	assetsFilter, err := k.avsKeeper.GetAvsSupportedAssets(historicalStateCtx, avsAddr)
	if err != nil {
		return nil, err
	}
	historyStakerAssets, err := k.delegationKeeper.DelegationStateByOperatorAssets(historicalStateCtx, operatorAddress.String(), assetsFilter)
	if err != nil {
		return nil, err
	}

	// get the Assets opted in the operator
	historyOperatorAssetsState, err := k.assetsKeeper.GetOperatorAssetInfos(historicalStateCtx, operatorAddress, assetsFilter)
	if err != nil {
		return nil, err
	}

	// calculate the actual slash amount according to the history and current state
	currentStakerAssets, err := k.delegationKeeper.DelegationStateByOperatorAssets(ctx, operatorAddress.String(), assetsFilter)
	if err != nil {
		return nil, err
	}
	// get the Assets opted in the operator
	currentOperatorAssetsState, err := k.assetsKeeper.GetOperatorAssetInfos(ctx, operatorAddress, assetsFilter)
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
			shouldSlashAmount := slashProportion.MulInt(historyAssetState.OperatorAmount).TruncateInt()
			if curAssetState.OperatorAmount.LT(shouldSlashAmount) {
				ret.slashOperatorInfo[assetID].AmountFromOptedIn = curAssetState.OperatorAmount
				remainShouldSlash := shouldSlashAmount.Sub(curAssetState.OperatorAmount)
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
			err = k.assetsKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, assetstype.StakerSingleAssetChangeInfo{
				TotalDepositAmount: slashSumValue.Neg(),
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
			err = k.assetsKeeper.UpdateOperatorAssetState(ctx, operatorAddress, assetID, assetstype.OperatorSingleAssetChangeInfo{
				TotalAmount: slashInfo.AmountFromOptedIn.Neg(),
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
		err := k.assetsKeeper.UpdateOperatorAssetState(ctx, operatorAddress, assetID, assetstype.OperatorSingleAssetChangeInfo{
			TotalAmount:                        slashSumValue.Neg(),
			OperatorAmount:                     slashInfo.AmountFromOptedIn.Neg(),
			OperatorUnbondableAmountAfterSlash: slashInfo.AmountFromUnbonding.Neg(),
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
func (k *Keeper) Slash(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr, slashContract, slashID string, slashEventHeight int64, slashProportion sdkmath.LegacyDec) error {
	height := ctx.BlockHeight()
	if slashEventHeight > height {
		return errorsmod.Wrap(types.ErrSlashOccurredHeight, fmt.Sprintf("slashEventHeight:%d,curHeight:%d", slashEventHeight, height))
	}

	// get the state when the slash occurred
	// get the opted-in info
	historicalSateCtx, err := assetstype.ContextForHistoricalState(ctx, slashEventHeight)
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
		EventHeight:     slashEventHeight,
		SlashProportion: slashProportion,
		ProcessedHeight: height + types.SlashVetoDuration,
	}
	err = k.UpdateOperatorSlashInfo(ctx, operatorAddress.String(), avsAddr, slashID, slashInfo)
	if err != nil {
		return err
	}

	// get the assets and amounts that should be slashed
	assetsSlashInfo, err := k.GetAssetsAmountToSlash(ctx, operatorAddress, avsAddr, slashEventHeight, slashProportion)
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

// SlashWithInfractionReason is an expected slash interface for the dogfood module.
func (k Keeper) SlashWithInfractionReason(
	ctx sdk.Context, addr sdk.AccAddress, infractionHeight, _ int64,
	slashFactor sdk.Dec, infraction stakingtypes.Infraction,
) sdkmath.Int {
	chainID := ctx.ChainID()
	avsAddr, err := k.avsKeeper.GetAvsAddrByChainID(ctx, chainID)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), chainID)
		return sdkmath.NewInt(0)
	}
	slashContract, err := k.avsKeeper.GetAvsSlashContract(ctx, avsAddr)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), avsAddr)
		return sdkmath.NewInt(0)
	}
	slashID := GetSlashIDForDogfood(infraction, infractionHeight)
	err = k.Slash(ctx, addr, avsAddr, slashContract, slashID, infractionHeight, slashFactor)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), avsAddr)
		return sdkmath.NewInt(0)
	}
	// todo: The returned value should be the amount of burned Exo if we considering a slash from the reward
	// Now it doesn't slash from the reward, so just return 0
	return sdkmath.NewInt(0)
}

// IsOperatorJailedForChainID add for dogfood
func (k Keeper) IsOperatorJailedForChainID(ctx sdk.Context, consAddr sdk.ConsAddress, chainID string) bool {
	found, operatorAddr := k.GetOperatorAddressForChainIDAndConsAddr(ctx, chainID, consAddr)
	if !found {
		k.Logger(ctx).Info("couldn't find operator by consensus address and chainID", consAddr, chainID)
		return false
	}

	avsAddr, err := k.avsKeeper.GetAvsAddrByChainID(ctx, chainID)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), chainID)
		return false
	}

	optInfo, err := k.GetOptedInfo(ctx, operatorAddr.String(), avsAddr)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), operatorAddr, avsAddr)
		return false
	}
	return optInfo.Jailed
}

func (k *Keeper) SetJailedState(ctx sdk.Context, consAddr sdk.ConsAddress, chainID string, jailed bool) {
	found, operatorAddr := k.GetOperatorAddressForChainIDAndConsAddr(ctx, chainID, consAddr)
	if !found {
		k.Logger(ctx).Info("couldn't find operator by consensus address and chainID", consAddr, chainID)
		return
	}

	avsAddr, err := k.avsKeeper.GetAvsAddrByChainID(ctx, chainID)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), chainID)
		return
	}

	handleFunc := func(info *types.OptedInfo) {
		info.Jailed = jailed
	}
	err = k.HandleOptedInfo(ctx, operatorAddr.String(), avsAddr, handleFunc)
	if err != nil {
		k.Logger(ctx).Error(err.Error(), chainID)
	}
}

// Jail an operator
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress, chainID string) {
	k.SetJailedState(ctx, consAddr, chainID, true)
}

// Unjail an operator
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress, chainID string) {
	k.SetJailedState(ctx, consAddr, chainID, true)
}
