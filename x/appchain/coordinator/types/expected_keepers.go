package types

import (
	time "time"

	exocoretypes "github.com/ExocoreNetwork/exocore/types"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// AVSKeeper represents the expected keeper interface for the AVS module.
type AVSKeeper interface {
	RegisterAVSWithChainID(sdk.Context, *avstypes.AVSRegisterOrDeregisterParams) (common.Address, error)
	IsAVSByChainID(sdk.Context, string) (bool, common.Address)
	DeleteAVSInfo(sdk.Context, common.Address) error
}

// EpochsKeeper represents the expected keeper interface for the epochs module.
type EpochsKeeper interface {
	GetEpochInfo(sdk.Context, string) (epochstypes.EpochInfo, bool)
}

// StakingKeeper represents the expected keeper interface for the staking module.
type StakingKeeper interface {
	UnbondingTime(sdk.Context) time.Duration
}

// OperatorKeeper represents the expected keeper interface for the operator module.
type OperatorKeeper interface {
	GetActiveOperatorsForChainID(sdk.Context, string) ([]sdk.AccAddress, []exocoretypes.WrappedConsKey)
	GetVotePowerForChainID(sdk.Context, []sdk.AccAddress, string) ([]int64, error)
}
