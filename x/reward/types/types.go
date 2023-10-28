package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Constructor of the Pool
func NewPool(name string) Pool {
	return Pool{
		Name:    name,
		Rewards: []Pool_Reward{},
	}
}

// ValidateBasic for the pool
func (m Pool) ValidateBasic() error {
	validatorDic := make(map[string]bool)
	for _, reward := range m.Rewards {
		validatorAddr := reward.Validator.String()
		if validatorDic[validatorAddr] {
			return fmt.Errorf("duplicate validator %s in pool %s", validatorAddr, m.Name)
		}

		if err := sdk.VerifyAddressFormat(reward.Validator); err != nil {
			return fmt.Errorf("invalid validator %s in pool %s", validatorAddr, m.Name)
		}

		if reward.Coins == nil || reward.Coins.Empty() {
			return fmt.Errorf("empty rewards found for validator %s in pool %s", validatorAddr, m.Name)
		}

		if err := reward.Coins.Validate(); err != nil {
			return fmt.Errorf("invalid rewards for validator %s found in pool %s", validatorAddr, m.Name)
		}
		validatorDic[validatorAddr] = true
	}

	return nil
}
