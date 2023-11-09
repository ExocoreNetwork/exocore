package keeper

import (
<<<<<<< HEAD
<<<<<<< HEAD
	sdkmath "cosmossdk.io/math"
	"fmt"
=======
=======
	sdkmath "cosmossdk.io/math"
>>>>>>> 104cf78 (add some test and fix bugs)
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	delegationKeeper "github.com/exocore/x/delegation/keeper"
	depositKeeper "github.com/exocore/x/deposit/keeper"
	retakingStateKeeper "github.com/exocore/x/restaking_assets_manage/keeper"

>>>>>>> eebca7f (implement slash interface)
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"

	"github.com/exocore/x/exoslash/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

<<<<<<< HEAD
	//other keepers
	retakingStateKeeper keeper.Keeper
}

=======
>>>>>>> eebca7f (implement slash interface)
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
	OptIntoSlashing(ctx sdk.Context, event *SlashParams) error
	Slash(ctx sdk.Context, event *SlashParams) error
	FreezeOperator(ctx sdk.Context, event *SlashParams) error
	ResetFrozenStatus(ctx sdk.Context, event *SlashParams) error
<<<<<<< HEAD
<<<<<<< HEAD
	IsOperatorFrozen(ctx sdk.Context, event *SlashParams) (bool, error)
	SetParams(ctx sdk.Context, params *types.Params) error
	GetParams(ctx sdk.Context) (*types.Params, error)
	OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec
=======
	IsOperatorFrozen(ctx sdk.Context, event *SlashParams) error
=======
	IsOperatorFrozen(ctx sdk.Context, event *SlashParams) (bool, error)
>>>>>>> 5429dca (add unti test for slash and fix some  bugs)
	SetParams(ctx sdk.Context, params *types.Params) error
	GetParams(ctx sdk.Context) (*types.Params, error)
<<<<<<< HEAD
>>>>>>> eebca7f (implement slash interface)
=======
	OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec
>>>>>>> 104cf78 (add some test and fix bugs)
}
