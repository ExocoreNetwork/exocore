package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	withdrawtype "github.com/ExocoreNetwork/exocore/x/withdraw/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type WithdrawParams struct {
	ClientChainLzID uint64
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

// 	var clientChainLzID uint64
// 	r = bytes.NewReader(log.Topics[types.ClientChainLzIDIndexInTopics][:])
// 	err = binary.Read(r, binary.BigEndian, &clientChainLzID)
// 	if err != nil {
// 		return nil, errorsmod.Wrap(err, "error occurred when binary read clientChainLzID from topic")
// 	}

// 	return &WithdrawParams{
// 		ClientChainLzID: clientChainLzID,
// 		Action:          action,
// 		AssetsAddress:   assetsAddress,
// 		WithdrawAddress: withdrawAddress,
// 		OpAmount:        amount,
// 	}, nil
// }

func getStakeIDAndAssetID(params *WithdrawParams) (stakeID string, assetID string) {
	clientChainLzIDStr := hexutil.EncodeUint64(params.ClientChainLzID)
	stakeID = strings.Join([]string{hexutil.Encode(params.WithdrawAddress), clientChainLzIDStr}, "_")
	assetID = strings.Join([]string{hexutil.Encode(params.AssetsAddress), clientChainLzIDStr}, "_")
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
	stakeID, assetID := getStakeIDAndAssetID(params)

	// check if asset exist
	if !k.assetsKeeper.IsStakingAsset(ctx, assetID) {
		return errorsmod.Wrap(withdrawtype.ErrWithdrawAssetNotExist, fmt.Sprintf("the assetID is:%s", assetID))
	}
	changeAmount := types.DeltaStakerSingleAsset{
		TotalDepositAmount: params.OpAmount.Neg(),
		WithdrawableAmount: params.OpAmount.Neg(),
	}
	err := k.assetsKeeper.UpdateStakerAssetState(ctx, stakeID, assetID, changeAmount)
	if err != nil {
		return err
	}
	if err = k.assetsKeeper.UpdateStakingAssetTotalAmount(ctx, assetID, params.OpAmount.Neg()); err != nil {
		return err
	}
	return nil
}
