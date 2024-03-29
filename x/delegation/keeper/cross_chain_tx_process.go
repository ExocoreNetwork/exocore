package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type DelegationOrUndelegationParams struct {
	ClientChainLzID uint64
	Action          types.CrossChainOpType
	AssetsAddress   []byte
	OperatorAddress sdk.AccAddress
	StakerAddress   []byte
	OpAmount        sdkmath.Int
	LzNonce         uint64
	TxHash          common.Hash
	// todo: The operator approved signature might be needed here in future
}

// The event hook process has been deprecated, now we use precompile contract to trigger the calls.
// solidity encode: bytes memory actionArgs = abi.encodePacked(token, operator, msg.sender, amount);
// _sendInterchainMsg(Action.DEPOSIT, actionArgs);
/*func (k Keeper) getParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*DelegationOrUndelegationParams, error) {
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
	clientChainInfo, err := k.restakingStateKeeper.GetClientChainInfoByIndex(ctx, clientChainLzID)
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

// DelegateTo : It doesn't need to check the active status of the operator in middlewares when delegating assets to the operator. This is because it adds assets to the operator's amount. But it needs to check if operator has been slashed or frozen.
func (k Keeper) DelegateTo(ctx sdk.Context, params *DelegationOrUndelegationParams) error {
	// check if the delegatedTo address is an operator
	if !k.IsOperator(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorNotExist
	}

	// check if the operator has been slashed or frozen
	if k.slashKeeper.IsOperatorFrozen(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorIsFrozen
	}

	// todo: The operator approved signature might be needed here in future

	// update the related states
	if params.OpAmount.IsNegative() {
		return delegationtype.ErrOpAmountIsNegative
	}

	stakerID, assetID := types.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)

	info, err := k.restakingStateKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
	if err != nil {
		return err
	}

	if info.CanWithdrawAmountOrWantChangeValue.LT(params.OpAmount) {
		return delegationtype.ErrDelegationAmountTooBig
	}

	err = k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, types.StakerSingleAssetOrChangeInfo{
		CanWithdrawAmountOrWantChangeValue: params.OpAmount.Neg(),
	})
	if err != nil {
		return err
	}

	err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, params.OperatorAddress, assetID, types.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue: params.OpAmount,
	})
	if err != nil {
		return err
	}

	delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
	delegatorAndAmount[params.OperatorAddress.String()] = &delegationtype.DelegationAmounts{
		CanUndelegationAmount: params.OpAmount,
	}
	err = k.UpdateDelegationState(ctx, stakerID, assetID, delegatorAndAmount)
	if err != nil {
		return err
	}
	err = k.UpdateStakerDelegationTotalAmount(ctx, stakerID, assetID, params.OpAmount)
	if err != nil {
		return err
	}
	return nil
}

// UndelegateFrom The undelegation needs to consider whether the operator's opted-in assets can exit from the AVS.
// Because only after the operator has served the AVS can the staking asset be undelegated.
// So we use two steps to handle the undelegation. Fist,record the undelegation request and the corresponding exit time which needs to be obtained from the operator opt-in module. Then,we handle the record when the exit time has expired.
func (k Keeper) UndelegateFrom(ctx sdk.Context, params *DelegationOrUndelegationParams) error {
	// check if the UndelegatedFrom address is an operator
	if !k.IsOperator(ctx, params.OperatorAddress) {
		return delegationtype.ErrOperatorNotExist
	}
	if params.OpAmount.IsNegative() {
		return delegationtype.ErrOpAmountIsNegative
	}
	// get staker delegation state, then check the validation of Undelegation amount
	stakerID, assetID := types.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)
	delegationState, err := k.GetSingleDelegationInfo(ctx, stakerID, assetID, params.OperatorAddress.String())
	if err != nil {
		return err
	}
	if params.OpAmount.GT(delegationState.CanUndelegationAmount) {
		return errorsmod.Wrap(delegationtype.ErrUndelegationAmountTooBig, fmt.Sprintf("UndelegationAmount:%s,CanUndelegationAmount:%s", params.OpAmount, delegationState.CanUndelegationAmount))
	}

	r := &delegationtype.UndelegationRecord{
		StakerID:              stakerID,
		AssetID:               assetID,
		OperatorAddr:          params.OperatorAddress.String(),
		TxHash:                params.TxHash.String(),
		IsPending:             true,
		LzTxNonce:             params.LzNonce,
		BlockNumber:           uint64(ctx.BlockHeight()),
		Amount:                params.OpAmount,
		ActualCompletedAmount: sdkmath.NewInt(0),
	}
	r.CompleteBlockNumber = k.operatorOptedInKeeper.GetOperatorCanUndelegateHeight(ctx, assetID, params.OperatorAddress, r.BlockNumber)
	err = k.SetUndelegationRecords(ctx, []*delegationtype.UndelegationRecord{r})
	if err != nil {
		return err
	}

	delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
	delegatorAndAmount[params.OperatorAddress.String()] = &delegationtype.DelegationAmounts{
		CanUndelegationAmount:  params.OpAmount.Neg(),
		WaitUndelegationAmount: params.OpAmount,
	}
	err = k.UpdateDelegationState(ctx, stakerID, assetID, delegatorAndAmount)
	if err != nil {
		return err
	}

	err = k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, types.StakerSingleAssetOrChangeInfo{
		WaitUndelegationAmountOrWantChangeValue: params.OpAmount,
	})
	if err != nil {
		return err
	}
	err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, params.OperatorAddress, assetID, types.OperatorSingleAssetOrChangeInfo{
		WaitUndelegationAmountOrWantChangeValue: params.OpAmount,
	})
	if err != nil {
		return err
	}
	return nil
}

/*func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	needLogs, err := k.depositKeeper.FilterCrossChainEventLogs(ctx, msg, receipt)
	if err != nil {
		return err
	}

	if len(needLogs) == 0 {
		log.Println("the hook message doesn't have any event needed to handle")
		return nil
	}

	for _, log := range needLogs {
		delegationParams, err := k.getParamsFromEventLog(ctx, log)
		if err != nil {
			return err
		}
		if delegationParams != nil {
			if delegationParams.Action == types.DelegateTo {
				err = k.DelegateTo(ctx, delegationParams)
			} else if delegationParams.Action == types.UndelegateFrom {
				err = k.UndelegateFrom(ctx, delegationParams)
			}
			if err != nil {
				// todo: need to test if the changed storage state will be reverted when there is an error occurred

				return err
			}
		}
	}
	return nil
}
*/
