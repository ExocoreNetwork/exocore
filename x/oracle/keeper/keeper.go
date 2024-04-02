package keeper

import (
	"fmt"
	"math/big"

	//	"cosmossdk.io/api/tendermint/abci"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		paramstore    paramtypes.Subspace
		stakingKeeper stakingkeeper.Keeper
	}
)

var _ common.KeeperOracle = Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	sKeeper stakingkeeper.Keeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		stakingKeeper: sKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetLastTotalPower(ctx sdk.Context) *big.Int {
	return k.stakingKeeper.GetLastTotalPower(ctx).BigInt()
}

func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context, f func(index int64, validator stakingtypes.ValidatorI) bool) {
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, f)
}

func (k Keeper) GetValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	return k.stakingKeeper.GetValidatorUpdates(ctx)
}

func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) (stakingtypes.Validator, bool) {
	return k.stakingKeeper.GetValidatorByConsAddr(ctx, addr)
}
