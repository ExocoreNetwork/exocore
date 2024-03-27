package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		memKey         storetypes.StoreKey
		paramstore     paramtypes.Subspace
		operatorKeeper delegationtypes.OperatorKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AVSInfoUpdate(ctx sdk.Context, params *AVSParams) error {
	operatorAddress := params.OperatorAddress
	opAccAddr, err := sdk.AccAddressFromBech32(string(operatorAddress))
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse acc address from Bech32,the addr is:%s", operatorAddress))
	}
	if !k.operatorKeeper.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(delegationtypes.ErrOperatorNotExist, "AVSInfoUpdate: invalid operator address")
	}
	action := params.Action

	if action == RegisterAction {
		return k.SetAVSInfo(ctx, params.AVSName, string(params.AVSAddress), string(operatorAddress))
	} else {
		return k.DeleteAVSInfo(ctx, string(params.AVSAddress))
	}
}

func (k Keeper) SetAVSInfo(ctx sdk.Context, AVSName string, AVSAddress string, OperatorAddress string) (err error) {
	avsAddr, err := sdk.AccAddressFromBech32(AVSAddress)
	if err != nil {
		return errorsmod.Wrap(err, "GetAVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(avsAddr) {
		AVSInfo := &types.AVSInfo{Name: AVSName, AVSAddress: AVSAddress, OperatorAddress: []string{OperatorAddress}}
		bz := k.cdc.MustMarshal(AVSInfo)
		store.Set(avsAddr, bz)
		return nil
	} else {
		value := store.Get(avsAddr)
		ret := &types.AVSInfo{}
		k.cdc.MustUnmarshal(value, ret)
		ret.OperatorAddress = append(ret.OperatorAddress, OperatorAddress)
		return nil
	}
}

func (k Keeper) GetAVSInfo(ctx sdk.Context, addr string) (*types.AVSInfo, error) {
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
	return &ret, nil
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
