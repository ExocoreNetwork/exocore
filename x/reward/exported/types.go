package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RewardPool represents a pool of rewards
type RewardPool interface {
	AddReward(sdk.ValAddress, sdk.Coin)
	ClearRewards(sdk.ValAddress)
}

// Refundable message
type Refundable interface {
	sdk.Msg
}
