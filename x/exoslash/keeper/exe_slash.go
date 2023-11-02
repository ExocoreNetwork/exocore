package keeper

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/evmos/v14/rpc/namespaces/ethereum/eth/filters"
	rtypes "github.com/exocore/x/exoslash/types"
	"github.com/exocore/x/restaking_assets_manage/types"
	"log"
	"math/big"
	"strings"
)

type SlashParams struct {
	clientChainLzId uint64
	action          types.CrossChainOpType
	assetsAddress   []byte
	operatorAddress sdk.AccAddress
	stakerAddress   []byte
	opAmount        sdkmath.Int
}

func (k Keeper) getParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*SlashParams, error) {
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

	return &SlashParams{
		clientChainLzId: clientChainLzId,
		action:          action,
		assetsAddress:   assetsAddress,
		stakerAddress:   stakerAddress,
		operatorAddress: opAccAddr,
		opAmount:        amount,
	}, nil
}

func getStakeIDAndAssetId(params *SlashParams) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(params.clientChainLzId)
	stakeId = strings.Join([]string{hexutil.Encode(params.stakerAddress[:]), clientChainLzIdStr}, "_")
	assetId = strings.Join([]string{hexutil.Encode(params.assetsAddress[:]), clientChainLzIdStr}, "_")
	return
}

func (k Keeper) FilterCrossChainEventLogs(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) ([]*ethtypes.Log, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	//filter needed logs
	addresses := []common.Address{common.HexToAddress(params.ExoCoreLzAppAddress)}
	topics := [][]common.Hash{
		{common.HexToHash(params.ExoCoreLzAppEventTopic)},
	}
	needLogs := filters.FilterLogs(receipt.Logs, nil, nil, addresses, topics)
	return needLogs, nil
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
		slashParams, err := k.getParamsFromEventLog(ctx, log)
		if err != nil {
			return err
		}
		if slashParams != nil {
			err = k.Slash(ctx, slashParams)
			if err != nil {
				// todo: need to test if the changed storage state will be reverted if there is an error occurred
				return err
			}
		}
	}
	return nil
}

func (k Keeper) Slash(ctx sdk.Context, event *SlashParams) error {
	//check event parameter then execute slash operation
	if event.opAmount.IsNegative() {
		return errorsmod.Wrap(rtypes.ErrSlashAmountIsNegative, fmt.Sprintf("the amount is:%s", event.opAmount))
	}
	stakeId, assetId := getStakeIDAndAssetId(event)
	//check is asset exist
	if !k.retakingStateKeeper.StakingAssetIsExist(ctx, assetId) {
		return errorsmod.Wrap(rtypes.ErrSlashAssetNotExist, fmt.Sprintf("the assetId is:%s", assetId))
	}

	//TODO
	changeAmount := types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: event.opAmount,
		CanWithdrawAmountOrWantChangeValue:  event.opAmount,
	}
	err := k.retakingStateKeeper.UpdateStakerAssetState(ctx, stakeId, assetId, changeAmount)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) FreezeOperator(ctx sdk.Context, event *SlashParams) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ResetFrozenStatus(ctx sdk.Context, event *SlashParams) error {
	//TODO implement me
	panic("implement me")
}
func (k Keeper) IsOperatorFrozen(ctx sdk.Context, event *SlashParams) error {
	//TODO implement me
	panic("implement me")
}
