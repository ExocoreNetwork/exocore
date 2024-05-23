package common

import (
	"cosmossdk.io/math"
	stakingkeeper "github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperOracle interface {
	KeeperStaking

	GetParams(sdk.Context) types.Params

	GetIndexRecentMsg(sdk.Context) (types.IndexRecentMsg, bool)
	GetAllRecentMsgAsMap(sdk.Context) map[int64][]*types.MsgItem

	GetIndexRecentParams(sdk.Context) (types.IndexRecentParams, bool)
	GetAllRecentParamsAsMap(sdk.Context) map[int64]*types.Params

	GetValidatorUpdateBlock(sdk.Context) (types.ValidatorUpdateBlock, bool)

	SetIndexRecentMsg(sdk.Context, types.IndexRecentMsg)
	SetRecentMsg(sdk.Context, types.RecentMsg)

	SetIndexRecentParams(sdk.Context, types.IndexRecentParams)
	SetRecentParams(sdk.Context, types.RecentParams)

	SetValidatorUpdateBlock(sdk.Context, types.ValidatorUpdateBlock)

	RemoveRecentParams(sdk.Context, uint64)
	RemoveRecentMsg(sdk.Context, uint64)
}

var _ KeeperStaking = stakingkeeper.Keeper{}

type KeeperStaking interface {
	GetLastTotalPower(ctx sdk.Context) math.Int
	IterateBondedValidatorsByPower(ctx sdk.Context, fn func(index int64, validator stakingTypes.ValidatorI) (stop bool))
	GetValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingTypes.Validator, found bool)

	GetAllExocoreValidators(ctx sdk.Context) (validators []dogfoodtypes.ExocoreValidator)
}
