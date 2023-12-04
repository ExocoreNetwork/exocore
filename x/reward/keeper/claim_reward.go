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
	ClientChainLzId       uint64
	Action                types.CrossChainOpType
	AssetsAddress         []byte
	WithdrawRewardAddress []byte
	OpAmount              sdkmath.Int
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
	var assetsAddress []byte
	err = binary.Read(r, binary.BigEndian, assetsAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += types.GeneralClientChainAddrLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	var rewardAddr []byte
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
		ClientChainLzId:       clientChainLzId,
		Action:                action,
		AssetsAddress:         assetsAddress,
		WithdrawRewardAddress: rewardAddr,
		OpAmount:              amount,
	}, nil
}

func getStakeIDAndAssetId(params *RewardParams) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(params.ClientChainLzId)
	stakeId = strings.Join([]string{hexutil.Encode(params.WithdrawRewardAddress[:]), clientChainLzIdStr}, "_")
	assetId = strings.Join([]string{hexutil.Encode(params.AssetsAddress[:]), clientChainLzIdStr}, "_")
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
	if event.OpAmount.IsNegative() {
		return errorsmod.Wrap(rtypes.ErrRewardAmountIsNegative, fmt.Sprintf("the amount is:%s", event.OpAmount))
	}
	stakeId, assetId := getStakeIDAndAssetId(event)
	//check is asset exist
	if !k.restakingStateKeeper.IsStakingAsset(ctx, assetId) {
		return errorsmod.Wrap(rtypes.ErrRewardAssetNotExist, fmt.Sprintf("the assetId is:%s", assetId))
	}

	//TODO
	changeAmount := types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: event.OpAmount,
		CanWithdrawAmountOrWantChangeValue:  event.OpAmount,
	}
	err := k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakeId, assetId, changeAmount)
	if err != nil {
		return err
	}
	if err = k.restakingStateKeeper.UpdateStakingAssetTotalAmount(ctx, assetId, event.OpAmount); err != nil {
		return err
	}
	return nil
}
