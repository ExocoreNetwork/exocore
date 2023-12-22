package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	delegationtype "github.com/exocore/x/delegation/types"
	depositkeeper "github.com/exocore/x/deposit/keeper"
	"github.com/exocore/x/restaking_assets_manage/keeper"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	//other keepers
	restakingStateKeeper  keeper.Keeper
	depositKeeper         depositkeeper.Keeper
	slashKeeper           delegationtype.ISlashKeeper
	operatorOptedInKeeper delegationtype.OperatorOptedInMiddlewareKeeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	restakingStateKeeper keeper.Keeper,
	depositKeeper depositkeeper.Keeper,
	slashKeeper delegationtype.ISlashKeeper,
	operatorOptedInKeeper delegationtype.OperatorOptedInMiddlewareKeeper,
) Keeper {
	return Keeper{
		storeKey:              storeKey,
		cdc:                   cdc,
		restakingStateKeeper:  restakingStateKeeper,
		depositKeeper:         depositKeeper,
		slashKeeper:           slashKeeper,
		operatorOptedInKeeper: operatorOptedInKeeper,
	}
}

// SetOperatorInfo This function is used to register to be an operator in exoCore, the provided info will be stored on the chain.
// Once an address has become an operator,the operator can't return to a normal address.But the operator can update the info through this function
// As for the operator opt-in function,it needs to be implemented in operator opt-in or AVS module
func (k Keeper) SetOperatorInfo(ctx sdk.Context, addr string, info *delegationtype.OperatorInfo) (err error) {
	opAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return errorsmod.Wrap(err, "SetOperatorInfo: error occurred when parse acc address from Bech32")
	}
	// todo: to check the validation of input info
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixOperatorInfo)
	// todo: think about the difference between init and update in future

	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(info)

	store.Set(opAccAddr, bz)
	return nil
}

func (k Keeper) GetOperatorInfo(ctx sdk.Context, addr string) (info *delegationtype.OperatorInfo, err error) {
	opAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetOperatorInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixOperatorInfo)
	//key := common.HexToAddress(incentive.Contract)
	ifExist := store.Has(opAccAddr)
	if !ifExist {
		return nil, errorsmod.Wrap(delegationtype.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorInfo: key is %s", opAccAddr))
	}

	value := store.Get(opAccAddr)

	ret := delegationtype.OperatorInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k Keeper) IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), delegationtype.KeyPrefixOperatorInfo)
	return store.Has(addr)
}

// GetExoCoreLzAppAddress Get exoCoreLzAppAddr from deposit keeper,it will be used when check the caller of precompile contract.
// This function needs to be moved to `restaking_assets_manage` module,which will facilitate its use for the other modules
func (k Keeper) GetExoCoreLzAppAddress(ctx sdk.Context) (common.Address, error) {
	return k.depositKeeper.GetExoCoreLzAppAddress(ctx)
}

// IDelegation interface will be implemented by deposit keeper
type IDelegation interface {
	// PostTxProcessing automatically call PostTxProcessing to update delegation state after receiving delegation event tx from layerZero protocol
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// RegisterOperator handle the registerOperator txs from msg service
	RegisterOperator(ctx context.Context, req *delegationtype.RegisterOperatorReq) (*delegationtype.RegisterOperatorResponse, error)
	// DelegateAssetToOperator handle the DelegateAssetToOperator txs from msg service
	DelegateAssetToOperator(ctx context.Context, delegation *delegationtype.MsgDelegation) (*delegationtype.DelegationResponse, error)
	// UndelegateAssetFromOperator handle the UndelegateAssetFromOperator txs from msg service
	UndelegateAssetFromOperator(ctx context.Context, delegation *delegationtype.MsgUndelegation) (*delegationtype.UndelegationResponse, error)

	GetSingleDelegationInfo(ctx sdk.Context, stakerId, assetId, operatorAddr string) (*delegationtype.DelegationAmounts, error)

	GetDelegationInfo(ctx sdk.Context, stakerId, assetId string) (*delegationtype.QueryDelegationInfoResponse, error)
}
