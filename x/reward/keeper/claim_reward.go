package keeper

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/evmos/evmos/v14/rpc/namespaces/ethereum/eth/filters"
	"github.com/exocore/x/restaking_assets_manage/types"
)

type RewardParams struct {
	clientChainLzId uint64
	action          types.CrossChainOpType
	assetsAddress   types.GeneralAssetsAddr
	depositAddress  types.GeneralClientChainAddr
	opAmount        sdkmath.Int
}

func getRewardParamsFromEventLog(log *ethtypes.Log) (*RewardParams, error) {
	// check if action is deposit
	var action types.CrossChainOpType
	var err error
	readStart := 0
	readEnd := types.CrossChainActionLength
	r := bytes.NewReader(log.Data[readStart:readEnd])
	err = binary.Read(r, binary.BigEndian, &action)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read action")
	}
	if action != types.DepositAction {
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
	var depositAddress types.GeneralClientChainAddr
	err = binary.Read(r, binary.BigEndian, depositAddress)
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

	return &DepositParams{
		clientChainLzId: clientChainLzId,
		action:          action,
		assetsAddress:   assetsAddress,
		depositAddress:  depositAddress,
		opAmount:        amount,
	}, nil
}

// To be decided here!
// func (k Keeper) sendRewards(ctx sdk.Context, rewards []*types.MsgRewardDetail, addr string, rewardProgramDenom string) (sdk.Coin, error) {
// 	amount := sdk.NewInt64Coin(rewardProgramDenom, 0)

// 	if len(rewards) == 0 {
// 		return amount, nil
// 	}

// 	for _, reward := range rewards {
// 		amount.Denom = reward.GetClaimPeriodReward().Denom
// 		amount = amount.Add(reward.GetClaimPeriodReward())
// 	}

// 	return k.sendCoinsToAccount(ctx, amount, addr)
// }

// // sendCoinsToAccount is mainly for `SendCoinsFromModuleToAccount`
// func (k Keeper) sendCoinsToAccount(ctx sdk.Context, amount sdk.Coin, addr string) (sdk.Coin, error) {
// 	if amount.IsZero() {
// 		return sdk.NewInt64Coin(amount.GetDenom(), 0), nil
// 	}

// 	acc, err := sdk.AccAddressFromBech32(addr)
// 	if err != nil {
// 		return sdk.NewInt64Coin(amount.Denom, 0), err
// 	}

// 	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc, sdk.NewCoins(amount))
// 	if err != nil {
// 		return sdk.NewInt64Coin(amount.Denom, 0), err
// 	}

// 	return amount, nil
// }

// // Get reward value of the corresponding address in the rewards pool
// func (p rewardPool) getRewards(address sdk.ValAddress) (sdk.Coins, bool) {
// 	for _, reward := range p.Rewards {
// 		if reward.Validator.Equals(address) {
// 			return reward.Coins, true
// 		}
// 	}

// 	return sdk.Coins{}, false
// }

// // Add and record for the corresponding reward in the rewards pool
// func (p *rewardPool) AddReward(address sdk.ValAddress, coin sdk.Coin) {
// 	defer func() {
// 		p.k.Logger(p.ctx).Debug("adding rewards in pool", "pool", p.Name, "validator", address.String(), "coin", coin.String())

// 		p.k.setPool(p.ctx, p.Pool)
// 	}()

// 	if coin.Amount.IsZero() {
// 		return
// 	}

// 	for i, reward := range p.Rewards {
// 		if reward.Validator.Equals(address) {
// 			p.Rewards[i].Coins = reward.Coins.Add(coin)
// 			return
// 		}
// 	}

// 	p.Rewards = append(p.Rewards, types.Pool_Reward{
// 		Validator: address,
// 		Coins:     sdk.NewCoins(coin),
// 	})
// }

// // Clear rewards of the specific address
// func (p *rewardPool) ClearRewards(address sdk.ValAddress) {
// 	for i, reward := range p.Rewards {
// 		if reward.Validator.Equals(address) {
// 			p.k.Logger(p.ctx).Info("clearing rewards in pool", "pool", p.Name, "validator", address.String())

// 			p.Rewards = append(p.Rewards[:i], p.Rewards[i+1:]...)
// 			p.k.setPool(p.ctx, p.Pool)
// 			return
// 		}
// 	}
// }

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
		depositParams, err := getDepositParamsFromEventLog(log)
		if err != nil {
			return err
		}
		if depositParams != nil {
			err = k.Deposit(ctx, depositParams)
			if err != nil {
				// todo: need to test if the changed storage state will be reverted if there is an error occurred
				return err
			}
		}
	}
	return nil
}
