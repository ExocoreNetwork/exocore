package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
	"gopkg.in/yaml.v2"

	"github.com/ExocoreNetwork/exocore/utils"
)

const (
	// DefaultMintDenom is the denomination in which the inflation is minted.
	DefaultMintDenom = utils.BaseDenom
	// DefaultEpochIdentifier is the epoch identifier which is used, by default, to identify the
	// epoch.
	DefaultEpochIdentifier = epochstypes.DayEpochID
	// DefaultEpochRewardStr is the amount of MintDenom minted at each epoch end, as a string.
	DefaultEpochRewardStr = "20000000000000000000"
)

func init() {
	// validate during init
	_, ok := sdk.NewIntFromString(DefaultEpochRewardStr)
	if !ok {
		panic(fmt.Sprintf("invalid default epoch reward: %s", DefaultEpochRewardStr))
	}
}

// NewParams creates a new Params instance
func NewParams(mintDenom string, epochReward math.Int, epochIdentifier string) Params {
	return Params{
		MintDenom:       mintDenom,
		EpochReward:     epochReward,
		EpochIdentifier: epochIdentifier,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	res, _ := sdk.NewIntFromString(DefaultEpochRewardStr)
	return NewParams(
		DefaultMintDenom, res, DefaultEpochIdentifier,
	)
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := ValidateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := ValidateEpochReward(p.EpochReward); err != nil {
		return err
	}
	return epochstypes.ValidateEpochIdentifierString(p.EpochIdentifier)
}

func ValidateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return sdk.ValidateDenom(v)
}

func ValidateEpochReward(i interface{}) error {
	v, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("mint reward must be positive: %s", v)
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
