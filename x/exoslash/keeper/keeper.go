package keeper

import (
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	delegationKeeper "github.com/exocore/x/delegation/keeper"
	depositKeeper "github.com/exocore/x/deposit/keeper"
	retakingStateKeeper "github.com/exocore/x/restaking_assets_manage/keeper"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"

	"github.com/exocore/x/exoslash/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		memKey              storetypes.StoreKey
		paramstore          paramtypes.Subspace
		retakingStateKeeper retakingStateKeeper.Keeper
		depositKeeper       depositKeeper.Keeper
		delegationKeeper    delegationKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	retakingStateKeeper keeper.Keeper,
) Keeper {

	return Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		retakingStateKeeper: retakingStateKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

type IEXOSlash interface {
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error
<<<<<<< HEAD
	OptIntoSlashing(ctx sdk.Context, event *SlashParams) error
	Slash(ctx sdk.Context, event *SlashParams) error
	FreezeOperator(ctx sdk.Context, event *SlashParams) error
	ResetFrozenStatus(ctx sdk.Context, event *SlashParams) error
	IsOperatorFrozen(ctx sdk.Context, event *SlashParams) (bool, error)
	SetParams(ctx sdk.Context, params *types.Params) error
	GetParams(ctx sdk.Context) (*types.Params, error)
	OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec
=======
	Slash(ctx sdk.Context, event *SlashParams) error
	FreezeOperator(ctx sdk.Context, event *SlashParams) error
	ResetFrozenStatus(ctx sdk.Context, event *SlashParams) error
	IsOperatorFrozen(ctx sdk.Context, event *SlashParams) error
<<<<<<< HEAD
	SetParams(ctx sdk.Context, params *depositTypes.Params) error
	GetParams(ctx sdk.Context) (*depositTypes.Params, error)
>>>>>>> 74ed4f2 (update exoslash module)
=======
	SetParams(ctx sdk.Context, params *types.Params) error
	GetParams(ctx sdk.Context) (*types.Params, error)
>>>>>>> f3ac4aa (implement slash interface)
}
