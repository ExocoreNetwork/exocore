package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
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
) *Keeper {
	return &Keeper{
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
	operatorAddress := params.OperatorAddress
	opAccAddr, err := sdk.AccAddressFromBech32(operatorAddress)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse acc address from Bech32,the addr is:%s", operatorAddress))
	}
	if !k.operatorKeeper.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(delegationtypes.ErrOperatorNotExist, "AVSInfoUpdate: invalid operator address")
	}

	action := params.Action

	if action == RegisterAction {
		return k.SetAVSInfo(ctx, params.AvsName, params.AvsAddress, operatorAddress, params.AssetID)
	}
	avsInfo, err := k.GetAVSInfo(ctx, params.AvsAddress)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when get avs info %s", avsInfo))
	}
	if avsInfo.Info.AvsOwnerAddress != params.AvsOwnerAddress {
		return errorsmod.Wrap(err, fmt.Sprintf("not qualified to deregister %s", avsInfo))
	}
	return k.DeleteAVSInfo(ctx, params.AvsAddress)
}

func (k Keeper) SetAVSInfo(ctx sdk.Context, avsName, avsAddress, operatorAddress, assetID string) (err error) {
	avsAddr, err := sdk.AccAddressFromBech32(avsAddress)
	if err != nil {
		return errorsmod.Wrap(err, "SetAVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(avsAddr) {
		AVSInfo := &types.AVSInfo{Name: avsName, AvsAddress: avsAddress, OperatorAddress: []string{operatorAddress}, AssetId: []string{assetID}}
		bz := k.cdc.MustMarshal(AVSInfo)
		store.Set(avsAddr, bz)
		return nil
	}
	value := store.Get(avsAddr)
	ret := &types.AVSInfo{}
	k.cdc.MustUnmarshal(value, ret)
	ret.OperatorAddress = append(ret.OperatorAddress, operatorAddress)
	ret.AssetId = append(ret.AssetId, assetID)
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
