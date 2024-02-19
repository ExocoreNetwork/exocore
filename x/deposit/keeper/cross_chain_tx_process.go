package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	despoittypes "github.com/ExocoreNetwork/exocore/x/deposit/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DepositParams struct {
	ClientChainLzID uint64
	// The action field might need to be removed,it will be used when called from event hook.
	Action        types.CrossChainOpType
	AssetsAddress []byte
	StakerAddress []byte
	OpAmount      sdkmath.Int
}

// The event hook process has been deprecated, now we use precompile contract to trigger the calls.
/*func (k Keeper) getDepositParamsFromEventLog(ctx sdk.Context, log *ethtypes.Log) (*DepositParams, error) {
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
	if action != types.Deposit {
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
	readEnd += clientChainInfo.AddressLength
	r = bytes.NewReader(log.Data[readStart:readEnd])
	depositAddress := make([]byte, clientChainInfo.AddressLength)
	err = binary.Read(r, binary.BigEndian, depositAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when binary read assets address")
	}

	readStart = readEnd
	readEnd += types.CrossChainOpAmountLength
	amount := sdkmath.NewIntFromBigInt(big.NewInt(0).SetBytes(log.Data[readStart:readEnd]))

	return &DepositParams{
		ClientChainLzID: clientChainLzID,
		Action:          action,
		AssetsAddress:   assetsAddress,
		StakerAddress:   depositAddress,
		OpAmount:        amount,
	}, nil
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
	//TODO check if contract address is valid layerZero relayer address

	needLogs, err := k.FilterCrossChainEventLogs(ctx, msg, receipt)
	if err != nil {
		return err
	}
	if len(needLogs) == 0 {
		log.Println("the hook message doesn't have any event needed to handle")
		return nil
	}

	for _, log := range needLogs {
		depositParams, err := k.getDepositParamsFromEventLog(ctx, log)
		if err != nil {
			return err
		}
		if depositParams != nil {
			err = k.Deposit(ctx, depositParams)
			if err != nil {
				// todo: need to test if the changed storage state will be reverted when there is an error occurred
				return err
			}
		}
	}
	return nil
}*/

// Deposit the deposit precompile contract will call this function to update asset state when there is a deposit.
func (k Keeper) Deposit(ctx sdk.Context, params *DepositParams) error {
	// check params parameter before executing deposit operation
	if params.OpAmount.IsNegative() {
		return errorsmod.Wrap(despoittypes.ErrDepositAmountIsNegative, fmt.Sprintf("the amount is:%s", params.OpAmount))
	}
	stakeID, assetID := types.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)
	// check if asset exist
	if !k.restakingStateKeeper.IsStakingAsset(ctx, assetID) {
		return errorsmod.Wrap(despoittypes.ErrDepositAssetNotExist, fmt.Sprintf("the assetID is:%s", assetID))
	}
	changeAmount := types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: params.OpAmount,
		CanWithdrawAmountOrWantChangeValue:  params.OpAmount,
	}
	// update asset state of the specified staker
	err := k.restakingStateKeeper.UpdateStakerAssetState(ctx, stakeID, assetID, changeAmount)
	if err != nil {
		return err
	}

	// update total amount of the deposited asset
	err = k.restakingStateKeeper.UpdateStakingAssetTotalAmount(ctx, assetID, params.OpAmount)
	if err != nil {
		return err
	}
	return nil
}
