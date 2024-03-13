package common

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperOracle interface {
	GetParams(sdk.Context) types.Params

	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingTypes.ValidatorI) bool)
	GetLastTotalPower(sdk.Context) *big.Int

	GetIndexRecentMsg(sdk.Context) (types.IndexRecentMsg, bool)
	GetAllRecentMsgAsMap(sdk.Context) map[uint64][]*types.MsgItem

	GetIndexRecentParams(sdk.Context) (types.IndexRecentParams, bool)
	GetAllRecentParamsAsMap(sdk.Context) map[uint64]*types.Params

	GetValidatorUpdateBlock(sdk.Context) (types.ValidatorUpdateBlock, bool)

	SetIndexRecentMsg(sdk.Context, types.IndexRecentMsg)
	SetRecentMsg(sdk.Context, types.RecentMsg)

	SetIndexRecentParams(sdk.Context, types.IndexRecentParams)
	SetRecentParams(sdk.Context, types.RecentParams)

	SetValidatorUpdateBlock(sdk.Context, types.ValidatorUpdateBlock)

	RemoveRecentParams(sdk.Context, uint64)
	RemoveRecentMsg(sdk.Context, uint64)
}
