package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/ethereum/go-ethereum/common"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	// DefaultExocoreLzAppAddress is the default address of ExocoreGateway.sol.
	// When it is automatically deployed within the genesis block (along with
	// supporting contracts) by the EVM module, this address should be changed
	// to point to that contract.
	DefaultExocoreLzAppAddress = "0x0000000000000000000000000000000000000000"
	// DefaultExocoreLzAppEventTopic is the default topic of the exocore lz app
	// event. TODO: Set this to a sane default?
	DefaultExocoreLzAppEventTopic = "0x000000000000000000000000000000000000000000000000000000000000000"
)

// Reflection based keys for params subspace.
var (
	KeyExocoreLzAppAddress    = []byte("ExocoreLzAppAddress")
	KeyExocoreLzAppEventTopic = []byte("ExocoreLzAppEventTopic")
)

// ParamKeyTable returns a key table with the necessary registered params.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance.
func NewParams(
	exocoreLzAppAddress string,
	exocoreLzAppEventTopic string,
) Params {
	return Params{
		ExocoreLzAppAddress:    exocoreLzAppAddress,
		ExocoreLzAppEventTopic: exocoreLzAppEventTopic,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultExocoreLzAppAddress,
		DefaultExocoreLzAppEventTopic,
	)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			KeyExocoreLzAppAddress,
			&p.ExocoreLzAppAddress,
			ValidateHexAddress,
		),
		paramtypes.NewParamSetPair(
			KeyExocoreLzAppEventTopic,
			&p.ExocoreLzAppEventTopic,
			ValidateHexHash,
		),
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := ValidateHexAddress(p.ExocoreLzAppAddress); err != nil {
		return fmt.Errorf("exocore lz app address: %w", err)
	}
	if err := ValidateHexHash(p.ExocoreLzAppEventTopic); err != nil {
		return fmt.Errorf("exocore lz app event topic: %w", err)
	}
	return nil
}

// ValidateHexAddress validates a hex address.
func ValidateHexAddress(i interface{}) error {
	addr, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("invalid hex address: %s", addr)
	}
	return nil
}

// ValidateHexHash validates a hex hash.
func ValidateHexHash(i interface{}) error {
	hash, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(common.FromHex(hash)) != common.HashLength {
		return fmt.Errorf("invalid hex hash: %s", hash)
	}
	return nil
}
