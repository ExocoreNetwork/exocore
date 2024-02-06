package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochsTypes "github.com/evmos/evmos/v14/x/epochs/types"
)

// EpochsKeeper represents the expected keeper interface for the epochs module.
type EpochsKeeper interface {
	GetEpochInfo(sdk.Context, string) (epochsTypes.EpochInfo, bool)
}

// DogfoodHooks represent the event hooks for dogfood module. Ideally, these should
// match those of the staking module but for now it is only a subset of them. The side effects
// of calling the other hooks are not relevant to running the chain, so they can be skipped.
type DogfoodHooks interface {
	AfterValidatorBonded(
		sdk.Context, sdk.ConsAddress, sdk.ValAddress,
	) error
}
