package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/keeper"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	// other keepers
	restakingStateKeeper keeper.Keeper
	delegationKeeper     operatortypes.ExpectDelegationInterface
	oracleKeeper         operatortypes.ExpectOracleInterface
	avsKeeper            operatortypes.ExpectAvsInterface

	// add for dogfood
	hooks       operatortypes.OperatorConsentHooks // set separately via call to SetHooks
	slashKeeper operatortypes.SlashKeeper          // for jailing and unjailing check TODO(mm)
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	restakingStateKeeper keeper.Keeper,
	oracleKeeper operatortypes.ExpectOracleInterface,
	avsKeeper operatortypes.ExpectAvsInterface,
	slashKeeper operatortypes.SlashKeeper,
) Keeper {
	return Keeper{
		storeKey:             storeKey,
		cdc:                  cdc,
		restakingStateKeeper: restakingStateKeeper,
		oracleKeeper:         oracleKeeper,
		avsKeeper:            avsKeeper,
		slashKeeper:          slashKeeper,
	}
}

func (k *Keeper) RegisterExpectDelegationInterface(delegationKeeper operatortypes.ExpectDelegationInterface) {
	k.delegationKeeper = delegationKeeper
}

func (k *Keeper) OracleInterface() operatortypes.ExpectOracleInterface {
	return k.oracleKeeper
}

func (k *Keeper) GetUnbondingExpirationBlockNumber(_ sdk.Context, _ sdk.AccAddress, startHeight uint64) uint64 {
	return startHeight + operatortypes.UnbondingExpiration
}

// OperatorKeeper interface will be implemented by deposit keeper
type OperatorKeeper interface {
	// RegisterOperator handle the registerOperator txs from msg service
	RegisterOperator(ctx context.Context, req *operatortypes.RegisterOperatorReq) (*operatortypes.RegisterOperatorResponse, error)

	IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool

	GetUnbondingExpirationBlockNumber(ctx sdk.Context, OperatorAddress sdk.AccAddress, startHeight uint64) uint64

	UpdateOptedInAssetsState(ctx sdk.Context, stakerID, assetID, operatorAddr string, opAmount sdkmath.Int) error

	OptIn(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr string) error

	OptOut(ctx sdk.Context, OperatorAddress sdk.AccAddress, AVSAddr string) error

	Slash(ctx sdk.Context, operatorAddress sdk.AccAddress, AVSAddr, slashContract, slashID string, occurredSateHeight int64, slashProportion sdkmath.LegacyDec) error
}
