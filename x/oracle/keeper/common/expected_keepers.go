package common

import (
	"math/big"

	//	"cosmossdk.io/api/tendermint/abci"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type KeeperOracle interface {
	GetParams(sdk.Context) types.Params

	IterateBondedValidatorsByPower(sdk.Context, func(index int64, validator stakingTypes.ValidatorI) bool)
	GetLastTotalPower(sdk.Context) *big.Int
	GetValidatorUpdates(sdk.Context) []abci.ValidatorUpdate
	GetValidatorByConsAddr(sdk.Context, sdk.ConsAddress) (stakingTypes.Validator, bool)

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
