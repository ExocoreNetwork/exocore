package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	withdrawtype "github.com/ExocoreNetwork/exocore/x/withdraw/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type WithdrawParams struct {
	ClientChainLzId uint64
	Action          types.CrossChainOpType
	AssetsAddress   []byte
	WithdrawAddress []byte
	OpAmount        sdkmath.Int
}

// func getWithdrawParamsFromEventLog(log *ethtypes.Log) (*WithdrawParams, error) {
// 	// check if action is withdraw
// 	var action types.CrossChainOpType
// 	var err error
// 	readStart := 0
// 	readEnd := types.CrossChainActionLength
// 	r := bytes.NewReader(log.Data[readStart:readEnd])
// 	err = binary.Read(r, binary.BigEndian, &action)
// 	if err != nil {
// 		return nil, errorsmod.Wrap(err, "error occurred when binary read action")
// 	}
// 	if action != types.WithdrawPrinciple {
// 		// not handle the actions that isn't withdraw
// 		return nil, nil
// 	}

// 	//decode the action parameters
// 	readStart = readEnd
// 	readEnd += types.GeneralAssetsAddrLength
// 	r = bytes.NewReader(log.Data[readStart:readEnd])
// 	var assetsAddress []byte
// 	err = binary.Read(r, binary.BigEndian, assetsAddress)
// 	if err != nil {
// 		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
// 	}

// 	readStart = readEnd
// 	readEnd += types.GeneralClientChainAddrLength
// 	r = bytes.NewReader(log.Data[readStart:readEnd])
// 	var withdrawAddress []byte
// 	err = binary.Read(r, binary.BigEndian, withdrawAddress)
// 	if err != nil {
// 		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
// 	}

// 	readStart = readEnd
// 	readEnd += types.CrossChainOpAmountLength
// 	amount := sdkmath.NewIntFromBigInt(big.NewInt(0).SetBytes(log.Data[readStart:readEnd]))

// 	var clientChainLzId uint64
// 	r = bytes.NewReader(log.Topics[types.ClientChainLzIdIndexInTopics][:])
// 	err = binary.Read(r, binary.BigEndian, &clientChainLzId)
// 	if err != nil {
// 		return nil, errorsmod.Wrap(err, "error occurred when binary read clientChainLzId from topic")
// 	}

// 	return &WithdrawParams{
// 		ClientChainLzId: clientChainLzId,
// 		Action:          action,
// 		AssetsAddress:   assetsAddress,
// 		WithdrawAddress: withdrawAddress,
// 		OpAmount:        amount,
// 	}, nil
// }

func getStakeIDAndAssetId(params *WithdrawParams) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(params.ClientChainLzId)
	stakeId = strings.Join([]string{hexutil.Encode(params.WithdrawAddress), clientChainLzIdStr}, "_")
	assetId = strings.Join([]string{hexutil.Encode(params.AssetsAddress), clientChainLzIdStr}, "_")
	return
}

// Be supported by precompiled way
// func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
// 	//TODO check if contract address is valid layerZero relayer address
// 	//check if log address and topicId is valid
// 	params, err := k.GetParams(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	//filter needed logs
// 	addresses := []common.Address{common.HexToAddress(params.ExoCoreLzAppAddress)}
// 	topics := [][]common.Hash{
// 		{common.HexToHash(params.ExoCoreLzAppEventTopic)},
// 	}
// 	needLogs := filters.FilterLogs(receipt.Logs, nil, nil, addresses, topics)
// 	if len(needLogs) == 0 {
// 		log.Println("the hook message doesn't have any event needed to handle")
// 		return nil
// 	}

// 	for _, log := range needLogs {
// 		withdrawParams, err := getWithdrawParamsFromEventLog(log)
// 		if err != nil {
// 			return err
// 		}
// 		if withdrawParams != nil {
// 			err = k.Withdraw(ctx, withdrawParams)
// 			if err != nil {
// 				// todo: need to test if the changed storage state will be reverted if there is an error occurred
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

func (k Keeper) Withdraw(ctx sdk.Context, params *WithdrawParams) error {
	// check event parameter then execute deposit operation
	if params.OpAmount.IsNegative() {
		return errorsmod.Wrap(withdrawtype.ErrWithdrawAmountIsNegative, fmt.Sprintf("the amount is:%s", params.OpAmount))
	}
	stakeId, assetId := getStakeIDAndAssetId(params)

	// check if asset exist
	if !k.restakingStateKeeper.IsStakingAsset(ctx, assetId) {
		return errorsmod.Wrap(withdrawtype.ErrWithdrawAssetNotExist, fmt.Sprintf("the assetId is:%s", assetId))
	}
	changeAmount := types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: params.OpAmount.Neg(),
		CanWithdrawAmountOrWantChangeValue:  params.OpAmount.Neg(),
	}
	err := k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakeId, assetId, changeAmount)
	if err != nil {
		return err
	}
	if err = k.restakingStateKeeper.UpdateStakingAssetTotalAmount(ctx, assetId, params.OpAmount.Neg()); err != nil {
		return err
	}
	return nil
}
