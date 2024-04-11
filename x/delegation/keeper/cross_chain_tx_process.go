package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The event hook process has been deprecated, now we use precompile contract to trigger the calls.
// solidity encode: bytes memory actionArgs = abi.encodePacked(token, operator, msg.sender, amount);
// _sendInterchainMsg(Action.DEPOSIT, actionArgs);
/*func (k *Keeper) getParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*DelegationOrUndelegationParams, error) {
	// check if Action is deposit
	var action types.CrossChainOpType
	var err error
	readStart := uint32(0)
	readEnd := uint32(types.CrossChainActionLength)
	r := bytes.NewReader(log.Data[readStart:readEnd])
	err = binary.Read(r, binary.BigEndian, &action)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read Action")
	}
	if action != types.DelegateTo && action != types.UndelegateFrom {
		// not handle the actions that isn't deposit
		return nil, nil
	}

	var clientChainLzID uint64
	r = bytes.NewReader(log.Topics[types.ClientChainLzIDIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &clientChainLzID)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read ClientChainLzID from topic")
	}
	clientChainInfo, err := k.assetsKeeper.GetClientChainInfoByIndex(ctx, clientChainLzID)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when get client chain info")
	}

	var lzNonce uint64
	r = bytes.NewReader(log.Topics[types.LzNonceIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &lzNonce)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read LzNonce from topic")
	}

	//decode the Action parameters
	readStart = readEnd
	readEnd += clientChainInfo.AddressLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	assetsAddress := make([]byte, clientChainInfo.AddressLength)
	err = binary.Read(r, binary.BigEndian, assetsAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += types.ExoCoreOperatorAddrLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	operatorAddress := [types.ExoCoreOperatorAddrLength]byte{}
	err = binary.Read(r, binary.BigEndian, operatorAddress[:])
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read operator address")
	}
	opAccAddr, err := sdk.AccAddressFromBech32(string(operatorAddress[:]))
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when parse acc address from Bech32")
	}

	readStart = readEnd
	readEnd += clientChainInfo.AddressLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	stakerAddress := make([]byte, clientChainInfo.AddressLength)
	err = binary.Read(r, binary.BigEndian, stakerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read staker address")
	}

	readStart = readEnd
	readEnd += types.CrossChainOpAmountLength
	amount := sdkmath.NewIntFromBigInt(big.NewInt(0).SetBytes(log.Data[readStart:readEnd]))

	return &DelegationOrUndelegationParams{
		ClientChainLzID: clientChainLzID,
		Action:          action,
		AssetsAddress:   assetsAddress,
		StakerAddress:   stakerAddress,
		OperatorAddress: opAccAddr,
		OpAmount:        amount,
		LzNonce:         lzNonce,
		TxHash:          log.TxHash,
	}, nil
}*/

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
	// check if the delegatedTo address is an operator
	if !k.operatorKeeper.IsOperator(ctx, params.OperatorAddress) {
		return errorsmod.Wrap(delegationtype.ErrOperatorNotExist, fmt.Sprintf("input operatorAddr is:%s", params.OperatorAddress))
	}

	// check if the operator has been slashed or frozen
	// skip the check if not genesis (or chain restart)
	if notGenesis && k.slashKeeper.IsOperatorFrozen(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorIsFrozen
	}

	// todo: The operator approved signature might be needed here in future

	// update the related states
	if params.OpAmount.IsNegative() {
		return delegationtype.ErrOpAmountIsNegative
	}

	stakerID, assetID := assetstype.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)

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

	err = k.assetsKeeper.UpdateOperatorAssetState(ctx, params.OperatorAddress, assetID, assetstype.DeltaOperatorSingleAsset{
		TotalAmount: params.OpAmount,
		TotalShare:  share,
	})
	if err != nil {
		return err
	}

	delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
	delegatorAndAmount[params.OperatorAddress.String()] = &delegationtype.DelegationAmounts{
		UndelegatableAmount: params.OpAmount,
		UndelegatableShare:  share,
	}
	err = k.UpdateDelegationState(ctx, stakerID, assetID, delegatorAndAmount)
	if err != nil {
		return err
	}

	if notGenesis {
		// call the hooks registered by the other modules
		k.Hooks().AfterDelegation(ctx, params.OperatorAddress)
	}
	return nil
}

// UndelegateFrom The undelegation needs to consider whether the operator's opted-in assets can exit from the AVS.
// Because only after the operator has served the AVS can the staking asset be undelegated.
// So we use two steps to handle the undelegation. Fist,record the undelegation request and the corresponding exit time which needs to be obtained from the operator opt-in module. Then,we handle the record when the exit time has expired.
func (k *Keeper) UndelegateFrom(ctx sdk.Context, params *delegationtype.DelegationOrUndelegationParams) error {
	// check if the UndelegatedFrom address is an operator
	if !k.operatorKeeper.IsOperator(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorNotExist
	}
	if params.OpAmount.IsNegative() {
		return delegationtype.ErrOpAmountIsNegative
	}
	// get staker delegation state, then check the validation of Undelegation amount
	stakerID, assetID := assetstype.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)

	// verify the undelegation amount
	share, err := k.ValidateUndeleagtionAmount(ctx, params.OperatorAddress, stakerID, assetID, params.OpAmount)
	if err != nil {
		return err
	}

	// remove share from operator
	removeToken, err := k.RemoveShareFromOperator(ctx, params.OperatorAddress, assetID, share)
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
		ActualCompletedAmount: sdkmath.NewInt(0),
	}
	r.CompleteBlockNumber = k.operatorKeeper.GetUnbondingExpirationBlockNumber(ctx, params.OperatorAddress, r.BlockNumber)
	err = k.SetUndelegationRecords(ctx, []*delegationtype.UndelegationRecord{r})
	if err != nil {
		return err
	}

	// update delegation state
	delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
	delegatorAndAmount[params.OperatorAddress.String()] = &delegationtype.DelegationAmounts{
		UndelegatableAmount:     removeToken.Neg(),
		WaitUndelegationAmount:  removeToken,
		UndelegatableAfterSlash: removeToken,
		UndelegatableShare:      share.Neg(),
	}
	err = k.UpdateDelegationState(ctx, stakerID, assetID, delegatorAndAmount)
	if err != nil {
		return err
	}

	// update staker and operator assets state
	// todo: TotalDepositAmount might be influenced by slash and precision loss,
	// consider removing it, it can be recalculated from the share for RPC query.
	err = k.assetsKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, assetstype.DeltaStakerSingleAsset{
		WaitUnbondingAmount: removeToken,
	})
	if err != nil {
		return err
	}

	// call the hooks registered by the other modules
	return k.Hooks().AfterUndelegationStarted(ctx, params.OperatorAddress, delegationtype.GetUndelegationRecordKey(r.LzTxNonce, r.TxHash, r.OperatorAddr))
}
