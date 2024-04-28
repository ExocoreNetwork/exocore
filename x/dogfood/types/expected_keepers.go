package types

import (
	"cosmossdk.io/math"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// OperatorKeeper represents the expected keeper interface for the operator module.
type OperatorKeeper interface {
	GetOperatorConsKeyForChainID(
		sdk.Context, sdk.AccAddress, string,
	) (bool, *tmprotocrypto.PublicKey, error)
	IsOperatorRemovingKeyFromChainID(
		sdk.Context, sdk.AccAddress, string,
	) bool
	CompleteOperatorKeyRemovalForChainID(sdk.Context, sdk.AccAddress, string) error
	GetOperatorAddressForChainIDAndConsAddr(
		sdk.Context, string, sdk.ConsAddress,
	) (bool, sdk.AccAddress)
	IsOperatorJailedForChainID(sdk.Context, sdk.ConsAddress, string) bool
	Jail(sdk.Context, sdk.ConsAddress, string)
	Unjail(sdk.Context, sdk.ConsAddress, string)
	// GetActiveOperatorsForChainID should return a list of operators and their public keys.
	// These operators should not be in the process of opting out, and should not be jailed
	// whether permanently or temporarily.
	GetActiveOperatorsForChainID(
		sdk.Context, string,
	) ([]sdk.AccAddress, []*tmprotocrypto.PublicKey)
	GetAvgDelegatedValue(
		sdk.Context, []sdk.AccAddress, string, string,
	) ([]int64, error)
	SlashWithInfractionReason(
		sdk.Context, sdk.AccAddress, int64,
		int64, sdk.Dec, stakingtypes.Infraction,
	) math.Int
	ValidatorByConsAddrForChainID(
		ctx sdk.Context, consAddr sdk.ConsAddress, chainID string,
	) stakingtypes.ValidatorI
}

// DelegationKeeper represents the expected keeper interface for the delegation module.
type DelegationKeeper interface {
	IncrementUndelegationHoldCount(sdk.Context, []byte) error
	DecrementUndelegationHoldCount(sdk.Context, []byte) error
}

// AssetsKeeper represents the expected keeper interface for the assets module.
type AssetsKeeper interface {
	IsStakingAsset(sdk.Context, string) bool
}

// SlashingKeeper represents the expected keeper interface for the (exo-)slashing module.
type SlashingKeeper interface{}
