package types

import (
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochsTypes "github.com/evmos/evmos/v14/x/epochs/types"
)

// EpochsKeeper represents the expected keeper interface for the epochs module.
type EpochsKeeper interface {
	GetEpochInfo(sdk.Context, string) (epochsTypes.EpochInfo, bool)
}

// DogfoodHooks represents the event hooks for dogfood module. Ideally, these should
// match those of the staking module but for now it is only a subset of them. The side effects
// of calling the other hooks are not relevant to running the chain, so they can be skipped.
type DogfoodHooks interface {
	AfterValidatorBonded(
		sdk.Context, sdk.ConsAddress, sdk.ValAddress,
	) error
}

// OperatorHooks is the interface for the operator module's hooks. The functions are called
// whenever an operator opts in to a Cosmos chain, opts out of a Cosmos chain, or replaces their
// public key with another one.
type OperatorHooks interface {
	AfterOperatorOptIn(sdk.Context, sdk.AccAddress, string, tmprotocrypto.PublicKey)
	AfterOperatorKeyReplacement(
		sdk.Context, sdk.AccAddress, tmprotocrypto.PublicKey, tmprotocrypto.PublicKey, string,
	)
	AfterOperatorOptOutInitiated(sdk.Context, sdk.AccAddress, string, tmprotocrypto.PublicKey)
}

// DelegationHooks represent the event hooks for delegation module.
type DelegationHooks interface {
	AfterDelegation(sdk.Context, sdk.AccAddress)
	AfterUndelegationStarted(sdk.Context, sdk.AccAddress, []byte)
	AfterUndelegationCompleted(sdk.Context, sdk.AccAddress)
}

// OperatorKeeper represents the expected keeper interface for the operator module.
type OperatorKeeper interface {
	GetOperatorConsKeyForChainId(
		sdk.Context, sdk.AccAddress, string,
	) (bool, tmprotocrypto.PublicKey, error)
	IsOperatorOptingOutFromChainId(
		sdk.Context, sdk.AccAddress, string,
	) bool
}

// DelegationKeeper represents the expected keeper interface for the delegation module.
type DelegationKeeper interface {
	IncrementUndelegationHoldCount(sdk.Context, []byte)
	DecrementUndelegationHoldCount(sdk.Context, []byte)
}
