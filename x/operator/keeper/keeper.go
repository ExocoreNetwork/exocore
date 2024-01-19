package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	operatortypes "github.com/exocore/x/operator/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	//other keepers
	restakingStateKeeper keeper.Keeper
	delegationKeeper     operatortypes.ExpectDelegationInterface
	oracleKeeper         operatortypes.ExpectOracleInterface
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	restakingStateKeeper keeper.Keeper,
	oracleKeeper operatortypes.ExpectOracleInterface,
) Keeper {
	return Keeper{
		storeKey:             storeKey,
		cdc:                  cdc,
		restakingStateKeeper: restakingStateKeeper,
		oracleKeeper:         oracleKeeper,
	}
}

func (k Keeper) RegisterExpectDelegationInterface(delegationKeeper operatortypes.ExpectDelegationInterface) {
	k.delegationKeeper = delegationKeeper
}

func (k Keeper) GetUnBondingExpirationBlockNumber(ctx sdk.Context, OperatorAddress sdk.AccAddress, startHeight uint64) uint64 {
	return startHeight + 10
}

// IOperator interface will be implemented by deposit keeper
type IOperator interface {
	// RegisterOperator handle the registerOperator txs from msg service
	RegisterOperator(ctx context.Context, req *operatortypes.RegisterOperatorReq) (*operatortypes.RegisterOperatorResponse, error)

	IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool

	GetUnBondingExpirationBlockNumber(ctx sdk.Context, OperatorAddress sdk.AccAddress, startHeight uint64) uint64

	OptIn(ctx sdk.Context, OperatorAddress sdk.AccAddress, AVSAddr string) error

	OptOut(ctx sdk.Context, OperatorAddress sdk.AccAddress, AVSAddr string) error

	IncreasedOptedInAssets(ctx sdk.Context, stakerId, assetId, operatorAddr string) error

	DecreaseOptedInAssets(ctx sdk.Context, stakerId, assetId, operatorAddr string) error

	SlashOperator()

	GetOperatorValueShare()
}
