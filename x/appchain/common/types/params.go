package types

import "time"

const (
	// Within the dogfood module, the default unbonding duration is 7 epochs, where 1 epoch = 1 day. This means a maximum of 8 days to unbond. We go lower than that here.
	DefaultSubscriberUnbondingPeriod = 24 * 7 * time.Hour
)
