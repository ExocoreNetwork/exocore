package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	// DefaultMintDenom = sdk.DefaultBondDenom // not a constant
	// DefaultEpochIdentifier is the epoch identifier which is used, by default, to identify the
	// epoch. Note that the options include week, day or hour.
	DefaultEpochIdentifier = epochstypes.DayEpochID
	// DefaultEpochRewardStr is the amount of MintDenom minted at each epoch end, as a string.
	DefaultEpochRewardStr = "20000000000000000000"
)

// Reflection based keys for params subspace
var (
	KeyMintDenom       = []byte("MintDenom")
	KeyEpochReward     = []byte("EpochReward")
	KeyEpochIdentifier = []byte("EpochIdentifier")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
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
	res, ok := sdk.NewIntFromString(DefaultEpochRewardStr)
	if !ok {
		panic("invalid default mint reward")
	}
	return NewParams(
		sdk.DefaultBondDenom, res, DefaultEpochIdentifier,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, p.MintDenom, ValidateMintDenom),
		paramtypes.NewParamSetPair(KeyEpochReward, p.EpochReward, ValidateEpochReward),
		paramtypes.NewParamSetPair(
			KeyEpochIdentifier, p.EpochIdentifier, epochstypes.ValidateEpochIdentifierInterface,
		),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := ValidateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := ValidateEpochReward(p.EpochReward); err != nil {
		return err
	}
	return nil
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
