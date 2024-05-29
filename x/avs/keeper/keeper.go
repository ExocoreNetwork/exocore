package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	assettypes "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		operatorKeeper delegationtypes.OperatorKeeper
		// other keepers
		assetsKeeper assettypes.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	operatorKeeper delegationtypes.OperatorKeeper,
	assetKeeper assettypes.Keeper,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		operatorKeeper: operatorKeeper,
		assetsKeeper:   assetKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AVSInfoUpdate(ctx sdk.Context, params *AVSRegisterOrDeregisterParams) error {

	avsInfo, _ := k.GetAVSInfo(ctx, params.AvsAddress)

	action := params.Action

	if action == RegisterAction && avsInfo != nil {
		return errorsmod.Wrap(types.ErrAlreadyRegistered, fmt.Sprintf("the error input arg is:%s", params.AvsAddress))
	}

	avs := &types.AVSInfo{
		Name:               params.AvsName,
		AvsAddress:         params.AvsAddress,
		SlashAddr:          params.SlashContractAddr,
		AvsOwnerAddress:    params.AvsOwnerAddress,
		AssetId:            params.AssetID,
		AvsUnbondingEpochs: uint32(params.MinimumDelegation),
		MinimumDelegation:  sdk.NewIntFromUint64(params.UnbondingEpochs),
		AvsEpoch:           nil,
		OperatorAddress:    nil,
	}
	if action == RegisterAction && avsInfo == nil {
		return k.SetAVSInfo(ctx, avs)
	}

	if avsInfo == nil {
		return errorsmod.Wrap(types.ErrUnregisterNonExistent, fmt.Sprintf("the error input arg is:%s", avsInfo))

	}

	//TODO:if avs DeRegisterAction check UnbondingEpochs
	//if avsInfo.Info.AvsUnbondingEpochs < currenUnbondingEpoch - regUnbondingEpoch {
	//	return errorsmod.Wrap(err, fmt.Sprintf("not qualified to deregister %s", avsInfo))
	//}
	return k.DeleteAVSInfo(ctx, params.AvsAddress)
}

func (k Keeper) AVSInfoUpdateWithOperator(ctx sdk.Context, params *OperatorOptParams) error {
	operatorAddress := params.OperatorAddress
	opAccAddr, err := sdk.AccAddressFromBech32(operatorAddress)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse acc address from Bech32,the addr is:%s", operatorAddress))
	}
	if !k.operatorKeeper.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(delegationtypes.ErrOperatorNotExist, fmt.Sprintf("AVSInfoUpdate: invalid operator address:%s", operatorAddress))
	}
	avsInfo, err := k.GetAVSInfo(ctx, params.AvsAddress)
	if err != nil || avsInfo.GetInfo() == nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when get avs info %s", avsInfo))
	}
	avs := avsInfo.GetInfo()
	addresses := avs.OperatorAddress

	if params.Action == RegisterAction && types.ContainsString(avs.OperatorAddress, operatorAddress) {
		return errorsmod.Wrap(types.ErrAlreadyRegistered, fmt.Sprintf("Error: Already registered，operatorAddress %s", operatorAddress))
	}

	if params.Action == RegisterAction && !types.ContainsString(avs.OperatorAddress, operatorAddress) {
		addresses = append(addresses, operatorAddress)
		avs.OperatorAddress = addresses
		return k.SetAVSInfo(ctx, avs)
	}

	if params.Action == DeRegisterAction && types.ContainsString(avs.OperatorAddress, operatorAddress) {
		avs.OperatorAddress = types.RemoveOperatorAddress(addresses, operatorAddress)
		return k.SetAVSInfo(ctx, avs)
	}
	return errorsmod.Wrap(types.ErrUnregisterNonExistent, fmt.Sprintf("No available operatorAddress to DeRegisterAction ,operatorAddress: %s", operatorAddress))

}
func (k Keeper) SetAVSInfo(ctx sdk.Context, avs *types.AVSInfo) (err error) {
	avsAddr, err := sdk.AccAddressFromBech32(avs.AvsAddress)
	if err != nil {
		return errorsmod.Wrap(err, "SetAVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)

	AVSInfo := avs
	bz := k.cdc.MustMarshal(AVSInfo)
	store.Set(avsAddr, bz)
	return nil

}

func (k Keeper) GetAVSInfo(ctx sdk.Context, addr string) (*types.QueryAVSInfoResponse, error) {
	avsAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetAVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(avsAddr) {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSInfo: key is %s", avsAddr))
	}

	value := store.Get(avsAddr)
	ret := types.AVSInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	res := &types.QueryAVSInfoResponse{
		Info: &ret,
	}
	return res, nil
}

func (k Keeper) DeleteAVSInfo(ctx sdk.Context, addr string) error {
	avsAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return errorsmod.Wrap(err, "AVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(avsAddr) {
		return errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("AVSInfo didn't exist: key is %s", avsAddr))
	}
	store.Delete(avsAddr)
	return nil
}
