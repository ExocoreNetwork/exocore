package keeper

import (
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	delegationKeeper "github.com/exocore/x/delegation/keeper"
	depositKeeper "github.com/exocore/x/deposit/keeper"
	depositTypes "github.com/exocore/x/deposit/types"
	retakingStateKeeper "github.com/exocore/x/restaking_assets_manage/keeper"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/core"

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

type IEXOSlash interface {
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error
	Slash(ctx sdk.Context, event *SlashParams) error
	FreezeOperator(ctx sdk.Context, event *SlashParams) error
	ResetFrozenStatus(ctx sdk.Context, event *SlashParams) error
	IsOperatorFrozen(ctx sdk.Context, event *SlashParams) error
	SetParams(ctx sdk.Context, params *depositTypes.Params) error
	GetParams(ctx sdk.Context) (*depositTypes.Params, error)
}
