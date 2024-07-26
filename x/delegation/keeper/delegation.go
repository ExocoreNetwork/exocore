package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	errorsmod "cosmossdk.io/errors"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegateTo : It doesn't need to check the active status of the operator in middlewares when
// delegating assets to the operator. This is because it adds assets to the operator's amount.
// But it needs to check if operator has been slashed or frozen.
func (k Keeper) DelegateTo(ctx sdk.Context, params *delegationtype.DelegationOrUndelegationParams) error {
	return k.delegateTo(ctx, params, true)
}

// delegateTo is the internal private version of DelegateTo. if the notGenesis parameter is
// false, the operator keeper and the delegation hooks are not called.
func (k *Keeper) delegateTo(
	ctx sdk.Context,
	params *delegationtype.DelegationOrUndelegationParams,
	notGenesis bool,
) error {
	if !params.OpAmount.IsPositive() {
		return delegationtype.ErrAmountIsNotPositive
	}
	// check if the delegatedTo address is an operator
	if !k.operatorKeeper.IsOperator(ctx, params.OperatorAddress) {
		return errorsmod.Wrap(delegationtype.ErrOperatorNotExist, fmt.Sprintf("input operatorAddr is:%s", params.OperatorAddress))
	}

	// check if the operator has been slashed or frozen
	// skip the check if not genesis (or chain restart)
	if notGenesis && k.slashKeeper.IsOperatorFrozen(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorIsFrozen
	}

	stakerID, assetID := assetstype.GetStakeIDAndAssetID(params.ClientChainID, params.StakerAddress, params.AssetsAddress)

	// check if the staker asset has been deposited and the canWithdraw amount is bigger than the delegation amount
	info, err := k.assetsKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
	if err != nil {
		return err
	}

	if info.WithdrawableAmount.LT(params.OpAmount) {
		return errorsmod.Wrap(delegationtype.ErrDelegationAmountTooBig, fmt.Sprintf("the opAmount is:%s the WithdrawableAmount amount is:%s", params.OpAmount, info.WithdrawableAmount))
	}

	// update staker asset state
	err = k.assetsKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, assetstype.DeltaStakerSingleAsset{
		WithdrawableAmount: params.OpAmount.Neg(),
	})
	if err != nil {
		return err
	}

	// calculate the share from the delegation amount
	share, err := k.CalculateShare(ctx, params.OperatorAddress, assetID, params.OpAmount)
	if err != nil {
		return err
	}

	deltaOperatorAsset := assetstype.DeltaOperatorSingleAsset{
		TotalAmount: params.OpAmount,
		TotalShare:  share,
	}
	// Check if the staker belongs to the delegated operator. Increase the operator's share if yes.
	operator, err := k.GetSelfDelegatedOperator(ctx, stakerID)
	if err != nil {
		return err
	}
	if operator == params.OperatorAddress.String() {
		deltaOperatorAsset.OperatorShare = share
	}

	err = k.assetsKeeper.UpdateOperatorAssetState(ctx, params.OperatorAddress, assetID, deltaOperatorAsset)
	if err != nil {
		return err
	}

	deltaAmount := &delegationtype.DeltaDelegationAmounts{
		UndelegatableShare: share,
	}
	_, err = k.UpdateDelegationState(ctx, stakerID, assetID, params.OperatorAddress.String(), deltaAmount)
	if err != nil {
		return err
	}
	err = k.AppendStakerForOperator(ctx, params.OperatorAddress.String(), assetID, stakerID)
	if err != nil {
		return err
	}

	if notGenesis {
		// call the hooks registered by the other modules
		k.Hooks().AfterDelegation(ctx, params.OperatorAddress)
	}
	return nil
}

// UndelegateFrom: The undelegation needs to consider whether the operator's opted-in assets can exit from the AVS.
// Because only after the operator has served the AVS can the staking asset be undelegated.
// So we use two steps to handle the undelegation. Fist,record the undelegation request and the corresponding exit time which needs to be obtained from the operator opt-in module. Then,we handle the record when the exit time has expired.
func (k *Keeper) UndelegateFrom(ctx sdk.Context, params *delegationtype.DelegationOrUndelegationParams) error {
	if !params.OpAmount.IsPositive() {
		return delegationtype.ErrAmountIsNotPositive
	}
	// check if the UndelegatedFrom address is an operator
	if !k.operatorKeeper.IsOperator(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorNotExist
	}
	// get staker delegation state, then check the validation of Undelegation amount
	stakerID, assetID := assetstype.GetStakeIDAndAssetID(params.ClientChainID, params.StakerAddress, params.AssetsAddress)

	// verify the undelegation amount
	share, err := k.ValidateUndelegationAmount(ctx, params.OperatorAddress, stakerID, assetID, params.OpAmount)
	if err != nil {
		return err
	}

	// remove share
	removeToken, err := k.RemoveShare(ctx, true, params.OperatorAddress, stakerID, assetID, share)
	if err != nil {
		return err
	}
	// record Undelegation event
	r := &delegationtype.UndelegationRecord{
		StakerID:              stakerID,
		AssetID:               assetID,
		OperatorAddr:          params.OperatorAddress.String(),
		TxHash:                params.TxHash.String(),
		IsPending:             true,
		LzTxNonce:             params.LzNonce,
		BlockNumber:           uint64(ctx.BlockHeight()),
		Amount:                removeToken,
		ActualCompletedAmount: removeToken,
	}
	r.CompleteBlockNumber = k.operatorKeeper.GetUnbondingExpirationBlockNumber(ctx, params.OperatorAddress, r.BlockNumber)
	err = k.SetUndelegationRecords(ctx, []*delegationtype.UndelegationRecord{r})
	if err != nil {
		return err
	}

	// call the hooks registered by the other modules
	return k.Hooks().AfterUndelegationStarted(ctx, params.OperatorAddress, delegationtype.GetUndelegationRecordKey(r.BlockNumber, r.LzTxNonce, r.TxHash, r.OperatorAddr))
}

// AssociateOperatorWithStaker marks that a staker is claiming to be associated with an operator.
// In other words, the staker's delegations will be marked as self-delegations for the operator.
// Each stakerID can associate, at most, to one operator. To change that operator, the staker must
// first call DissociateOperatorFromStaker.
// However, each operator may be associated with multiple stakers.
// This function is not available for end users to call directly, and such calls must be made
// via the ExocoreGateway. The gateway associates the `msg.sender` of the call, along with the
// chain id, with the operator. Since it associates `msg.sender`, we do not need to check that
// this request is signed by the staker.
// TODO: It may be possible that an address, which is an EOA staker on a client chain, is a
// contract on Exocore. When that happens, the contract may be able to call the Gateway to
// maliciously associate the staker with an operator. The probability of this, however, is
// infinitesimal (10^-60) so we will not do anything about it for now.
// The solution could be to require that such associations originate from the client chain.
func (k *Keeper) AssociateOperatorWithStaker(
	ctx sdk.Context,
	clientChainID uint64,
	operatorAddress sdk.AccAddress,
	stakerAddress common.Address,
) error {
	if !k.assetsKeeper.ClientChainExists(ctx, clientChainID) {
		return delegationtype.ErrClientChainNotExist
	}
	if !k.operatorKeeper.IsOperator(ctx, operatorAddress) {
		return delegationtype.ErrOperatorNotExist
	}

	stakerID, _ := assetstype.GetStakeIDAndAssetID(clientChainID, stakerAddress[:], nil)
	associatedOperator, err := k.GetSelfDelegatedOperator(ctx, stakerID)
	if err != nil {
		return err
	}
	if associatedOperator != "" {
		return delegationtype.ErrOperatorAlreadyAssociated
	}

	opFunc := func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) error {
		// increase the share of new marked operator
		if keys.OperatorAddr == operatorAddress.String() {
			err = k.assetsKeeper.UpdateOperatorAssetState(ctx, operatorAddress, keys.AssetID, assetstype.DeltaOperatorSingleAsset{
				OperatorShare: amounts.UndelegatableShare,
			})
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = k.IterateDelegationsForStaker(ctx, stakerID, opFunc)
	if err != nil {
		return err
	}

	// update the marking information
	err = k.SetSelfDelegatedOperator(ctx, stakerID, operatorAddress.String())
	if err != nil {
		return err
	}

	return nil
}

// DissociateOperatorFromStaker deletes the associative relationship between a staker
// (address + chain id combination) and an operator. See AssociateOperatorWithStaker for more
// information on the relationship and restrictions.
func (k *Keeper) DissociateOperatorFromStaker(
	ctx sdk.Context,
	clientChainID uint64,
	stakerAddress common.Address,
) error {
	stakerID, _ := assetstype.GetStakeIDAndAssetID(clientChainID, stakerAddress[:], nil)
	associatedOperator, err := k.GetSelfDelegatedOperator(ctx, stakerID)
	if err != nil {
		return err
	}
	if associatedOperator == "" {
		return delegationtype.ErrNoAssociatedOperatorByStaker
	}
	oldOperatorAccAddr, err := sdk.AccAddressFromBech32(associatedOperator)
	if err != nil {
		return delegationtype.OperatorAddrIsNotAccAddr
	}

	opFunc := func(keys *delegationtype.SingleDelegationInfoReq, amounts *delegationtype.DelegationAmounts) error {
		// decrease the share of old operator
		if keys.OperatorAddr == associatedOperator {
			err = k.assetsKeeper.UpdateOperatorAssetState(ctx, oldOperatorAccAddr, keys.AssetID, assetstype.DeltaOperatorSingleAsset{
				OperatorShare: amounts.UndelegatableShare.Neg(),
			})
		}
		if err != nil {
			return err
		}
		return nil
	}
	err = k.IterateDelegationsForStaker(ctx, stakerID, opFunc)
	if err != nil {
		return err
	}

	// delete the marking information
	err = k.DeleteSelfDelegatedOperator(ctx, stakerID)
	if err != nil {
		return err
	}

	return nil
}
