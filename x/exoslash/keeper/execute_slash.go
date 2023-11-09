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
	ClientChainLzId           uint64
	Action                    types.CrossChainOpType
	AssetsAddress             []byte
	OperatorAddress           sdk.AccAddress
	StakerAddress             []byte
<<<<<<< HEAD
<<<<<<< HEAD
	MiddlewareContractAddress []byte
	Proportion                sdkmath.LegacyDec
	OpAmount                  sdkmath.Int
	Proof                     []byte
=======
	OpAmount                  sdkmath.Int
	MiddlewareContractAddress []byte
	Proportion                sdkmath.LegacyDec
	Evidence                  string
>>>>>>> 104cf78 (add some test and fix bugs)
=======
	MiddlewareContractAddress []byte
	Proportion                sdkmath.LegacyDec
	OpAmount                  sdkmath.Int
	Proof                     []byte
>>>>>>> 5429dca (add unti test for slash and fix some  bugs)
}
type OperatorFrozenStatus struct {
	operatorAddress sdk.AccAddress
	status          bool
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
		ClientChainLzId: clientChainLzId,
		Action:          action,
		AssetsAddress:   assetsAddress,
		StakerAddress:   stakerAddress,
		OperatorAddress: opAccAddr,
		OpAmount:        amount,
	}, nil
}

func getStakeIDAndAssetId(params *SlashParams) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(params.ClientChainLzId)
	stakeId = strings.Join([]string{hexutil.Encode(params.StakerAddress[:]), clientChainLzIdStr}, "_")
	assetId = strings.Join([]string{hexutil.Encode(params.AssetsAddress[:]), clientChainLzIdStr}, "_")
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

func (k Keeper) OptIntoSlashing(ctx sdk.Context, event *SlashParams) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Slash(ctx sdk.Context, event *SlashParams) error {
	//the stakes are frozen for the impacted middleware, and deposits and withdrawals are disabled as well.
	//All pending deposits and withdrawals for the current epoch will be invalidated.
	//check event parameter then execute slash operation
<<<<<<< HEAD
<<<<<<< HEAD
	_ = k.SetFrozenStatus(ctx, string(event.OperatorAddress), true)

=======
>>>>>>> 104cf78 (add some test and fix bugs)
=======
	_ = k.SetFrozenStatus(ctx, string(event.OperatorAddress), true)

>>>>>>> 5429dca (add unti test for slash and fix some  bugs)
	if event.OpAmount.IsNegative() {
		return errorsmod.Wrap(rtypes.ErrSlashAmountIsNegative, fmt.Sprintf("the amount is:%s", event.OpAmount))
	}
	stakeId, assetId := getStakeIDAndAssetId(event)
	//check is asset exist
	if !k.retakingStateKeeper.StakingAssetIsExist(ctx, assetId) {
		return errorsmod.Wrap(rtypes.ErrSlashAssetNotExist, fmt.Sprintf("the assetId is:%s", assetId))
	}

	//TODO Processing Slash Core Logic
	changeAmount := types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: event.OpAmount,
		CanWithdrawAmountOrWantChangeValue:  event.OpAmount,
	}
	err := k.retakingStateKeeper.UpdateStakerAssetState(ctx, stakeId, assetId, changeAmount)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) FreezeOperator(ctx sdk.Context, event *SlashParams) error {
	k.SetFrozenStatus(ctx, string(event.OperatorAddress), true)
	return nil
}

func (k Keeper) ResetFrozenStatus(ctx sdk.Context, event *SlashParams) error {
	k.SetFrozenStatus(ctx, string(event.OperatorAddress), true)
	return nil
}
func (k Keeper) IsOperatorFrozen(ctx sdk.Context, event *SlashParams) (bool, error) {
	return k.GetFrozenStatus(ctx, string(event.OperatorAddress))

<<<<<<< HEAD
}
func (k Keeper) OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec {
	//TODO
	return sdkmath.LegacyNewDec(3)
=======
>>>>>>> 5429dca (add unti test for slash and fix some  bugs)
}
func (k Keeper) OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec {
	//TODO
	return sdkmath.LegacyNewDec(3)
}
