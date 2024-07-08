package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

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
