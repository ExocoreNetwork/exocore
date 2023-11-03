package keeper

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
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
	clientChainLzId uint64
	action          types.CrossChainOpType
	assetsAddress   []byte
	operatorAddress sdk.AccAddress
	stakerAddress   []byte
	opAmount        sdkmath.Int
	lzNonce         uint64
	txHash          common.Hash
	//todo: The operator approved signature might be needed here in future
}

// solidity encode: bytes memory actionArgs = abi.encodePacked(token, operator, msg.sender, amount);
// _sendInterchainMsg(Action.DEPOSIT, actionArgs);
func (k Keeper) getParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*DelegationOrUnDelegationParams, error) {
	// check if action is deposit
	var action types.CrossChainOpType
	var err error
	readStart := uint32(0)
	readEnd := uint32(types.CrossChainActionLength)
	r := bytes.NewReader(log.Data[readStart:readEnd])
	err = binary.Read(r, binary.BigEndian, &action)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read action")
	}
	if action != types.DelegationTo && action != types.UnDelegationFrom {
		// not handle the actions that isn't deposit
		return nil, nil
	}

	var clientChainLzId uint64
	r = bytes.NewReader(log.Topics[types.ClientChainLzIdIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &clientChainLzId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read clientChainLzId from topic")
	}
	clientChainInfo, err := k.retakingStateKeeper.GetClientChainInfoByIndex(ctx, clientChainLzId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when get client chain info")
	}

	var lzNonce uint64
	r = bytes.NewReader(log.Topics[types.LzNonceIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &lzNonce)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read lzNonce from topic")
	}

	//decode the action parameters
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
		clientChainLzId: clientChainLzId,
		action:          action,
		assetsAddress:   assetsAddress,
		stakerAddress:   stakerAddress,
		operatorAddress: opAccAddr,
		opAmount:        amount,
		lzNonce:         lzNonce,
		txHash:          log.TxHash,
	}, nil
}

// DelegateTo : It doesn't need to check the active status of the operator in middlewares when delegating assets to the operator. This is because it adds assets to the operator's amount. But it needs to check if operator has been slashed or frozen.
func (k Keeper) DelegateTo(ctx sdk.Context, params *DelegationOrUnDelegationParams) error {
	// check if the delegatedTo address is an operator
	if !k.IsOperator(ctx, params.operatorAddress) {
		return types2.ErrOperatorNotExist
	}

	// check if the operator has been slashed or frozen
	if k.slashKeeper.IsOperatorFrozen(params.operatorAddress) {
		return types2.ErrOperatorIsFrozen
	}

	//todo: The operator approved signature might be needed here in future

	//update the related states
	if params.opAmount.IsNegative() {
		return types2.ErrOpAmountIsNegative
	}

	stakerId, assetId := types.GetStakeIDAndAssetId(params.clientChainLzId, params.stakerAddress, params.assetsAddress)
	err := k.retakingStateKeeper.UpdateStakerAssetState(ctx, stakerId, assetId, types.StakerSingleAssetOrChangeInfo{
		CanWithdrawAmountOrWantChangeValue: params.opAmount.Neg(),
	})
	if err != nil {
		return err
	}

	err = k.retakingStateKeeper.UpdateOperatorAssetState(ctx, params.operatorAddress, assetId, types.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue: params.opAmount,
	})
	if err != nil {
		return err
	}

	delegatorAndAmount := make(map[string]sdkmath.Int)
	delegatorAndAmount[params.operatorAddress.String()] = params.opAmount
	err = k.UpdateDelegationState(ctx, stakerId, assetId, delegatorAndAmount)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) UnDelegateFrom(ctx sdk.Context, params *DelegationOrUnDelegationParams) error {

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
			if delegationParams.action == types.DelegationTo {
				err = k.DelegateTo(ctx, delegationParams)
			} else if delegationParams.action == types.UnDelegationFrom {
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
