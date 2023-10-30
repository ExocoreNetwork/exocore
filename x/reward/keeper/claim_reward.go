package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/evmos/v14/rpc/namespaces/ethereum/eth/filters"
	"github.com/exocore/x/restaking_assets_manage/types"
	rtypes "github.com/exocore/x/reward/types"
)

type RewardParams struct {
	clientChainLzId       uint64
	action                types.CrossChainOpType
	assetsAddress         types.GeneralAssetsAddr
	withdrawRewardAddress types.GeneralClientChainAddr
	opAmount              sdkmath.Int
}

func getRewardParamsFromEventLog(log *ethtypes.Log) (*RewardParams, error) {
	// check if action is to get reward
	var action types.CrossChainOpType
	var err error
	readStart := 0
	readEnd := types.CrossChainActionLength
	r := bytes.NewReader(log.Data[readStart:readEnd])
	err = binary.Read(r, binary.BigEndian, &action)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read action")
	}
	if action != types.WithDrawReward {
		// not handle the actions that isn't deposit
		return nil, nil
	}

	//decode the action parameters
	readStart = readEnd
	readEnd += types.GeneralAssetsAddrLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	var assetsAddress types.GeneralAssetsAddr
	err = binary.Read(r, binary.BigEndian, assetsAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += types.GeneralClientChainAddrLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	var rewardAddr types.GeneralClientChainAddr
	err = binary.Read(r, binary.BigEndian, rewardAddr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += types.CrossChainOpAmountLength
	amount := sdkmath.NewIntFromBigInt(big.NewInt(0).SetBytes(log.Data[readStart:readEnd]))

	var clientChainLzId uint64
	r = bytes.NewReader(log.Topics[types.ClientChainLzIdIndexInTopics][:])
	err = binary.Read(r, binary.BigEndian, &clientChainLzId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read clientChainLzId from topic")
	}

	return &RewardParams{
		clientChainLzId:       clientChainLzId,
		action:                action,
		assetsAddress:         assetsAddress,
		withdrawRewardAddress: rewardAddr,
		opAmount:              amount,
	}, nil
}

func getStakeIDAndAssetId(params *RewardParams) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(params.clientChainLzId)
	stakeId = strings.Join([]string{hexutil.Encode(params.withdrawRewardAddress[:]), clientChainLzIdStr}, "_")
	assetId = strings.Join([]string{hexutil.Encode(params.assetsAddress[:]), clientChainLzIdStr}, "_")
	return
}

func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	//TODO check if contract address is valid layerZero relayer address
	//check if log address and topicId is valid
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}
	//filter needed logs
	addresses := []common.Address{common.HexToAddress(params.ExoCoreLzAppAddress)}
	topics := [][]common.Hash{
		{common.HexToHash(params.ExoCoreLzAppEventTopic)},
	}
	needLogs := filters.FilterLogs(receipt.Logs, nil, nil, addresses, topics)
	if len(needLogs) == 0 {
		log.Println("the hook message doesn't have any event needed to handle")
		return nil
	}

	for _, log := range needLogs {
		rewardParams, err := getRewardParamsFromEventLog(log)
		if err != nil {
			return err
		}
		if rewardParams != nil {
			err = k.RewardForWithdraw(ctx, rewardParams)
			if err != nil {
				// todo: need to test if the changed storage state will be reverted if there is an error occurred
				return err
			}
		}
	}
	return nil
}

func (k Keeper) RewardForWithdraw(ctx sdk.Context, event *RewardParams) error {
	//check event parameter then execute RewardForWithdraw operation
	if event.opAmount.IsNegative() {
		return errorsmod.Wrap(rtypes.ErrRewardAmountIsNegative, fmt.Sprintf("the amount is:%s", event.opAmount))
	}
	stakeId, assetId := getStakeIDAndAssetId(event)
	//check is asset exist
	if !k.retakingStateKeeper.StakingAssetIsExist(ctx, assetId) {
		return errorsmod.Wrap(rtypes.ErrRewardAssetNotExist, fmt.Sprintf("the assetId is:%s", assetId))
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
