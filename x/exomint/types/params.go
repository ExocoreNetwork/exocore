package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Reflection based keys for params subspace
var (
	KeyMintDenom   = []byte("MintDenom")
	KeyBlockReward = []byte("BlockReward")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(mintDenom string, blockReward math.Int) Params {
	return Params{
		MintDenom:   mintDenom,
		BlockReward: blockReward,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		sdk.DefaultBondDenom,
		sdk.NewInt(1),
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, p.MintDenom, ValidateMintDenom),
		paramtypes.NewParamSetPair(KeyBlockReward, p.BlockReward, ValidateBlockReward),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := ValidateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := ValidateBlockReward(p.BlockReward); err != nil {
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

func ValidateBlockReward(i interface{}) error {
	v, ok := i.(sdk.Int)
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
