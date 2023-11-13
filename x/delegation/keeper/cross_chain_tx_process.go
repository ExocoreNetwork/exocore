package keeper

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	types2 "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/restaking_assets_manage/types"
	"log"
	"math/big"
)

type DelegationOrUnDelegationParams struct {
	ClientChainLzId uint64
	Action          types.CrossChainOpType
	AssetsAddress   []byte
	OperatorAddress sdk.AccAddress
	StakerAddress   []byte
	OpAmount        sdkmath.Int
	LzNonce         uint64
	TxHash          common.Hash
	//todo: The operator approved signature might be needed here in future
}

// solidity encode: bytes memory actionArgs = abi.encodePacked(token, operator, msg.sender, amount);
// _sendInterchainMsg(Action.DEPOSIT, actionArgs);
func (k Keeper) getParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*DelegationOrUnDelegationParams, error) {
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
	if action != types.DelegationTo && action != types.UnDelegationFrom {
		// not handle the actions that isn't deposit
		return nil, nil
	}

	var clientChainLzId uint64
	r = bytes.NewReader(log.Topics[types.ClientChainLzIdIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &clientChainLzId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read ClientChainLzId from topic")
	}
	clientChainInfo, err := k.retakingStateKeeper.GetClientChainInfoByIndex(ctx, clientChainLzId)
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

	return &DelegationOrUnDelegationParams{
		ClientChainLzId: clientChainLzId,
		Action:          action,
		AssetsAddress:   assetsAddress,
		StakerAddress:   stakerAddress,
		OperatorAddress: opAccAddr,
		OpAmount:        amount,
		LzNonce:         lzNonce,
		TxHash:          log.TxHash,
	}, nil
}

// DelegateTo : It doesn't need to check the active status of the operator in middlewares when delegating assets to the operator. This is because it adds assets to the operator's amount. But it needs to check if operator has been slashed or frozen.
func (k Keeper) DelegateTo(ctx sdk.Context, params *DelegationOrUnDelegationParams) error {
	// check if the delegatedTo address is an operator
	if !k.IsOperator(ctx, params.OperatorAddress) {
		return types2.ErrOperatorNotExist
	}

	// check if the operator has been slashed or frozen
	if k.slashKeeper.IsOperatorFrozen(ctx, params.OperatorAddress) {
		return types2.ErrOperatorIsFrozen
	}

	//todo: The operator approved signature might be needed here in future

	//update the related states
	if params.OpAmount.IsNegative() {
		return types2.ErrOpAmountIsNegative
	}

	stakerId, assetId := types.GetStakeIDAndAssetId(params.ClientChainLzId, params.StakerAddress, params.AssetsAddress)
	err := k.retakingStateKeeper.UpdateStakerAssetState(ctx, stakerId, assetId, types.StakerSingleAssetOrChangeInfo{
		CanWithdrawAmountOrWantChangeValue: params.OpAmount.Neg(),
	})
	if err != nil {
		return err
	}

	err = k.retakingStateKeeper.UpdateOperatorAssetState(ctx, params.OperatorAddress, assetId, types.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue: params.OpAmount,
	})
	if err != nil {
		return err
	}

	delegatorAndAmount := make(map[string]*types2.DelegationAmounts)
	delegatorAndAmount[params.OperatorAddress.String()] = &types2.DelegationAmounts{
		CanUnDelegationAmount: params.OpAmount,
	}
	err = k.UpdateDelegationState(ctx, stakerId, assetId, delegatorAndAmount)
	if err != nil {
		return err
	}
	err = k.UpdateStakerDelegationTotalAmount(ctx, stakerId, assetId, params.OpAmount)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) UnDelegateFrom(ctx sdk.Context, params *DelegationOrUnDelegationParams) error {
	// check if the unDelegatedFrom address is an operator
	if !k.IsOperator(ctx, params.OperatorAddress) {
		return types2.ErrOperatorNotExist
	}
	if params.OpAmount.IsNegative() {
		return types2.ErrOpAmountIsNegative
	}
	// get staker delegation state, then check the validation of unDelegation amount
	stakerId, assetId := types.GetStakeIDAndAssetId(params.ClientChainLzId, params.StakerAddress, params.AssetsAddress)
	delegationState, err := k.GetSingleDelegationInfo(ctx, stakerId, assetId, params.OperatorAddress.String())
	if err != nil {
		return err
	}
	if params.OpAmount.GT(delegationState.CanUnDelegationAmount) {
		return errorsmod.Wrap(types2.ErrUnDelegationAmountTooBig, fmt.Sprintf("unDelegationAmount:%s,CanUnDelegationAmount:%s", params.OpAmount, delegationState.CanUnDelegationAmount))
	}

	//record unDelegation event
	r := &types2.UnDelegationRecord{
		StakerId:              stakerId,
		AssetId:               assetId,
		OperatorAddr:          params.OperatorAddress.String(),
		TxHash:                params.TxHash.String(),
		IsPending:             true,
		LzTxNonce:             params.LzNonce,
		BlockNumber:           uint64(ctx.BlockHeight()),
		Amount:                params.OpAmount,
		ActualCompletedAmount: sdkmath.NewInt(0),
	}
	r.CompleteBlockNumber = k.operatorOptedInKeeper.GetOperatorCanUnDelegateHeight(ctx, assetId, params.OperatorAddress, r.BlockNumber)
	err = k.SetUnDelegationStates(ctx, []*types2.UnDelegationRecord{r})
	if err != nil {
		return err
	}

	//update delegation state
	delegatorAndAmount := make(map[string]*types2.DelegationAmounts)
	delegatorAndAmount[params.OperatorAddress.String()] = &types2.DelegationAmounts{
		CanUnDelegationAmount:  params.OpAmount.Neg(),
		WaitUnDelegationAmount: params.OpAmount,
	}
	err = k.UpdateDelegationState(ctx, stakerId, assetId, delegatorAndAmount)
	if err != nil {
		return err
	}

	//update staker and operator assets state
	err = k.retakingStateKeeper.UpdateStakerAssetState(ctx, stakerId, assetId, types.StakerSingleAssetOrChangeInfo{
		WaitUnDelegationAmountOrWantChangeValue: params.OpAmount,
	})
	if err != nil {
		return err
	}
	err = k.retakingStateKeeper.UpdateOperatorAssetState(ctx, params.OperatorAddress, assetId, types.OperatorSingleAssetOrChangeInfo{
		WaitUnDelegationAmountOrWantChangeValue: params.OpAmount,
	})
	if err != nil {
		return err
	}
	return nil
}
func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
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
			if delegationParams.Action == types.DelegationTo {
				err = k.DelegateTo(ctx, delegationParams)
			} else if delegationParams.Action == types.UnDelegationFrom {
				err = k.UnDelegateFrom(ctx, delegationParams)
			}
			if err != nil {
				// todo: need to test if the changed storage state will be reverted when there is an error occurred

				return err
			}
		}
	}
	return nil
}
