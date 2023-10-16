package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewPool is the constructor of Pool
func NewPool(name string) Pool {
	return Pool{
		Name:    name,
		Rewards: []Pool_Reward{},
	}
}

// ValidateBasic returns an error if the Pool is not valid; nil otherwise
func (m Pool) ValidateBasic() error {

	validatorSeen := make(map[string]bool)
	for _, reward := range m.Rewards {
		validatorAddr := reward.Validator.String()
		if validatorSeen[validatorAddr] {
			return fmt.Errorf("duplicate validator %s found in pool %s", validatorAddr, m.Name)
		}

		if err := sdk.VerifyAddressFormat(reward.Validator); err != nil {
			return fmt.Errorf("invalid validator %s found in pool %s", validatorAddr, m.Name)
		}

		validatorSeen[validatorAddr] = true
	}

	return nil
}

// ValidateBasic returns an error if the Refund is not valid; nil otherwise
func (m Refund) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(m.Payer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, sdkerrors.Wrap(err, "payer").Error())
	}

	return nil
}
