package keeper

import (
	"fmt"

	delegationkeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SlashInputInfo struct {
	SlashType        types.SlashType
	Operator         sdk.AccAddress
	AVSAddr          string
	SlashContract    string
	SlashID          string
	SlashEventHeight int64
	SlashProportion  sdkmath.LegacyDec
}

// GetSlashIDForDogfood It use infractionType+'/'+'infractionHeight' as the slashID, because the slash event occurs in dogfood doesn't have a TxID. It isn't submitted through an external transaction.
func GetSlashIDForDogfood(infraction stakingtypes.Infraction, infractionHeight int64) string {
	// #nosec G701
	return string(assetstype.GetJoinedStoreKey(hexutil.EncodeUint64(uint64(infraction)), hexutil.EncodeUint64(uint64(infractionHeight))))
}

func SlashFromUndelegation(undelegation *delegationtype.UndelegationRecord, totalSlashAmount, slashAmount sdkmath.Int) (sdkmath.Int, error) {
	if undelegation.ActualCompletedAmount.IsZero() {
		return totalSlashAmount, nil
	}
	// reduce the actual_completed_amount in the record
	if slashAmount.GTE(undelegation.ActualCompletedAmount) {
		slashAmount = undelegation.ActualCompletedAmount
		undelegation.ActualCompletedAmount = sdkmath.NewInt(0)
	} else {
		undelegation.ActualCompletedAmount = undelegation.ActualCompletedAmount.Sub(slashAmount)
	}
	// slashing from the operator isn't needed if the remainingSlashAmount isn't positive
	remainingSlashAmount := totalSlashAmount.Sub(slashAmount)
	return remainingSlashAmount, nil
}

func (k *Keeper) VerifySlashEvent(ctx sdk.Context, parameter *SlashInputInfo) (sdk.Context, error) {
	height := ctx.BlockHeight()
	if parameter.SlashEventHeight > height {
		return ctx, errorsmod.Wrap(types.ErrSlashOccurredHeight, fmt.Sprintf("slashEventHeight:%d,curHeight:%d", parameter.SlashEventHeight, height))
	}

	// get the state when the slash occurred
	// get the opted-in info
	// When the slash occurs, retrieves the end block height of the epoch
	// where the used voting power resides.
	heightForVotingPower, err := k.avsKeeper.GetHeightForVotingPower(ctx, parameter.AVSAddr, parameter.SlashEventHeight)
	if err != nil {
		return ctx, err
	}
	if k.historicalCtx == nil {
		return ctx, errorsmod.Wrap(types.ErrValueIsNilOrZero, "VerifySlashEvent the historicalCtx is nil")
	}
	historicalStateCtx, err := k.historicalCtx(heightForVotingPower, false)
	if err != nil {
		return ctx, err
	}
	if !k.IsOptedIn(ctx, parameter.Operator.String(), parameter.AVSAddr) {
		return ctx, types.ErrNotOptedIn
	}
	optedInfo, err := k.GetOptedInfo(historicalStateCtx, parameter.Operator.String(), parameter.AVSAddr)
	if err != nil {
		return ctx, err
	}
	if optedInfo.SlashContract != parameter.SlashContract {
		return ctx, errorsmod.Wrap(types.ErrSlashContractNotMatch, fmt.Sprintf("input slashContract:%s, opted-in slash contract:%suite", parameter.SlashContract, optedInfo.SlashContract))
	}

	return historicalStateCtx, nil
}

// NoInstantaneousSlash indicates that the slash event will be processed after a certain
// period of time, thus requiring a reduction in the share of the corresponding staker.
// It will slash the assets from the undelegation firstly, then slash the asset from the
// staker's share.
// Compared to the instant slash, the con is the handling isn't efficient, but the pro is
// there isn't any slash mistake for the new-coming delegations after the slash event.
func (k *Keeper) NoInstantaneousSlash(ctx, historicalStateCtx sdk.Context, parameter *SlashInputInfo) error {
	// get assetsInfo supported by AVS
	assetsFilter, err := k.avsKeeper.GetAVSSupportedAssets(historicalStateCtx, parameter.AVSAddr)
	if err != nil {
		return err
	}
	for assetID := range assetsFilter {
		historyOperatorAsset, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(historicalStateCtx, parameter.Operator, assetID)
		if err != nil {
			return err
		}
		stakerMap, err := k.delegationKeeper.GetStakersByOperator(historicalStateCtx, parameter.Operator.String(), assetID)
		if err != nil {
			return err
		}
		// slash the staker share according to the historical and current state
		for stakerID := range stakerMap.Stakers {
			delegationState, err := k.delegationKeeper.GetSingleDelegationInfo(historicalStateCtx, stakerID, assetID, parameter.Operator.String())
			if err != nil {
				return err
			}
			assetAmount, err := delegationkeeper.TokensFromShares(delegationState.UndelegatableShare, historyOperatorAsset.TotalShare, historyOperatorAsset.TotalAmount)
			if err != nil {
				return err
			}
			shouldSlashAmount := parameter.SlashProportion.MulInt(assetAmount).TruncateInt()
			// slash the asset from the undelegation firstly.
			undelegations, err := k.delegationKeeper.GetStakerUndelegationRecords(ctx, stakerID, assetID)
			if err != nil {
				return err
			}
			isSlashFromShare := true
			for _, undelegation := range undelegations {
				if undelegation.OperatorAddr == parameter.Operator.String() {
					remainingSlashAmount, err := SlashFromUndelegation(undelegation, shouldSlashAmount, shouldSlashAmount)
					if err != nil {
						return err
					}
					// update the undelegation state
					_, err = k.delegationKeeper.SetSingleUndelegationRecord(ctx, undelegation)
					if err != nil {
						return err
					}
					if remainingSlashAmount.IsZero() {
						// all amount has been slashed from the undelegaion,
						// so the share of staker shouldn't be slashed.
						isSlashFromShare = false
						break
					}
					// slash the remaining amount from the next undelegation
					shouldSlashAmount = remainingSlashAmount
				}
			}
			// slash the asset from the staker's asset share if there is still remaining slash amount after
			// slashing from the undelegation.
			if isSlashFromShare {
				slashShare, err := k.delegationKeeper.CalculateSlashShare(ctx, parameter.Operator, stakerID, assetID, shouldSlashAmount)
				if err != nil {
					return err
				}
				_, err = k.delegationKeeper.RemoveShare(ctx, false, parameter.Operator, stakerID, assetID, slashShare)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// InstantSlash represents the slash events that will be handled instantly after occurring
// so the slash will reduce the amount of related operator's assets pool.
// The implementation is similar to the cosmos-sdk, but the difference is the comos-sdk
// slashes a proportion from every undelegations, which may result in the actual total slash amount is greater
// than intended. But this function won't slash the remaining undelegation if the total slash amount
// is reached. So some undelegations might escape from the slash because of the other new-coming
// delegations.
// The new delegated assets after the slash event might be slashed incorrectly, which is same as cosmos-sdk.
// So this function only apply to the situation that the slash events  will be handled instantly after occurring.
// The pro is that the handling is more efficient compared to the no-instantaneous slash.
func (k *Keeper) InstantSlash(ctx, historicalStateCtx sdk.Context, parameter *SlashInputInfo) error {
	// todo: should we check if the slashEventHeight is too smaller than the current block height?
	// Then force the AVS using the `NoInstantaneousSlash` if the slash event is too old
	// This may depend on the level of authority we want to grant to AVS when
	// setting their own slash logic.

	slashAssets := make(map[string]sdkmath.Int, 0)
	// get assetsInfo supported by AVS
	assetsFilter, err := k.avsKeeper.GetAVSSupportedAssets(historicalStateCtx, parameter.AVSAddr)
	if err != nil {
		return err
	}

	// get the Assets opted in the operator
	historyOperatorAssets, err := k.assetsKeeper.GetOperatorAssetInfos(historicalStateCtx, parameter.Operator, assetsFilter)
	if err != nil {
		return err
	}

	// calculate the assets amount that should be slashed
	for assetID, state := range historyOperatorAssets {
		slashAmount := parameter.SlashProportion.MulInt(state.TotalAmount).TruncateInt()
		slashAssets[assetID] = slashAmount
	}

	// slash from the unbonding stakers
	if parameter.SlashEventHeight < ctx.BlockHeight() {
		// get the undelegations that are submitted after the slash.
		opFunc := func(undelegation *delegationtype.UndelegationRecord) error {
			totalSlashAmount, ok := slashAssets[undelegation.AssetID]
			if ok {
				slashAmount := parameter.SlashProportion.MulInt(undelegation.Amount).TruncateInt()
				remainingSlashAmount, err := SlashFromUndelegation(undelegation, totalSlashAmount, slashAmount)
				if err != nil {
					return err
				}
				// slashing from the operator isn't needed if the remainingSlashAmount isn't positive
				if !remainingSlashAmount.IsPositive() {
					delete(slashAssets, undelegation.AssetID)
				} else {
					slashAssets[undelegation.AssetID] = remainingSlashAmount
				}
			}
			return nil
		}
		// #nosec G701
		heightFilter := uint64(parameter.SlashEventHeight)
		err = k.delegationKeeper.IterateUndelegationsByOperator(ctx, parameter.Operator.String(), &heightFilter, true, opFunc)
		if err != nil {
			return err
		}
	}

	// slash the remaining from the assets pool of the operator
	for assetID, slashAmount := range slashAssets {
		operatorAsset, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, parameter.Operator, assetID)
		if err != nil {
			return err
		}
		isClearUselessShare := false
		if slashAmount.GTE(operatorAsset.TotalAmount) {
			slashAmount = operatorAsset.TotalAmount
			isClearUselessShare = true
		}

		changeAmount := assetstype.DeltaOperatorSingleAsset{
			TotalAmount: slashAmount.Neg(),
		}
		// all shares need to be cleared if the asset amount is slashed to zero,
		// otherwise there will be a problem in updating the shares when handling
		// the new delegations.
		if isClearUselessShare {
			// clear the share of other stakers
			stakerMap, err := k.delegationKeeper.GetStakersByOperator(ctx, parameter.Operator.String(), assetID)
			if err != nil {
				return err
			}
			err = k.delegationKeeper.SetStakerShareToZero(ctx, parameter.Operator.String(), assetID, stakerMap)
			if err != nil {
				return err
			}
			err = k.delegationKeeper.DeleteStakerMapForOperator(ctx, parameter.Operator.String(), assetID)
			if err != nil {
				return err
			}
			changeAmount.TotalShare = sdkmath.LegacyNewDec(0)
			changeAmount.OperatorAmount = sdkmath.NewInt(0)
			changeAmount.OperatorShare = sdkmath.LegacyNewDec(0)
		}

		err = k.assetsKeeper.UpdateOperatorAssetState(ctx, parameter.Operator, assetID, changeAmount)
		if err != nil {
			return err
		}
	}
	return nil
}

// Slash performs all slash events include instant slash and no-instantaneous slash.
func (k *Keeper) Slash(ctx sdk.Context, parameter *SlashInputInfo) error {
	historicalStateCtx, err := k.VerifySlashEvent(ctx, parameter)
	if err != nil {
		return err
	}
	switch parameter.SlashType {
	case types.SlashType_SLASH_TYPE_INSTANT_SLASH:
		err = k.InstantSlash(ctx, historicalStateCtx, parameter)
	case types.SlashType_SLASH_TYPE_NO_INSTANTANEOUS_SLASH:
		err = k.NoInstantaneousSlash(ctx, historicalStateCtx, parameter)
	default:
		return errorsmod.Wrap(types.ErrInvalidSlashType, fmt.Sprintf("the slash type is:%v", parameter.SlashType))
	}
	if err != nil {
		return err
	}

	// todo: recording the slash event might be moved to the slash module
	height := ctx.BlockHeight()
	slashInfo := types.OperatorSlashInfo{
		SlashContract:   parameter.SlashContract,
		SubmittedHeight: height,
		EventHeight:     parameter.SlashEventHeight,
		SlashProportion: parameter.SlashProportion,
		ProcessedHeight: height + types.SlashVetoDuration,
		SlashType:       parameter.SlashType,
	}
	err = k.UpdateOperatorSlashInfo(ctx, parameter.Operator.String(), parameter.AVSAddr, parameter.SlashID, slashInfo)
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
	// todo: disable the slash currently, waiting for the new slash implementation
	k.Logger(ctx).Info("slash occurs", addr, infractionHeight, slashFactor, infraction)
	/*	chainID := ctx.ChainID()
		avsAddr, err := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
		if err != nil {
			k.Logger(ctx).Error(err.Error(), chainID)
			return sdkmath.NewInt(0)
		}
		slashContract, err := k.avsKeeper.GetAVSSlashContract(ctx, avsAddr)
		if err != nil {
			k.Logger(ctx).Error(err.Error(), avsAddr)
			return sdkmath.NewInt(0)
		}
		slashID := GetSlashIDForDogfood(infraction, infractionHeight)
		slashParam := &SlashInputInfo{
			SlashType:        types.SlashType_SLASH_TYPE_INSTANT_SLASH,
			Operator:         addr,
			AVSAddr:          avsAddr,
			SlashContract:    slashContract,
			SlashID:          slashID,
			SlashEventHeight: infractionHeight,
			SlashProportion:  slashFactor,
		}
		err = k.Slash(ctx, slashParam)
		if err != nil {
			k.Logger(ctx).Error(err.Error(), avsAddr)
			return sdkmath.NewInt(0)
		}*/
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

	avsAddr, err := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
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

	avsAddr, err := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
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
