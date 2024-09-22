package types

import (
	fmt "fmt"
	"time"

	ibchost "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// CoordinatorFeePoolAddrStr is from the coordinator params.
	// DistributionTransmissionChannel is currently blank.
	// RewardDenom should be fetched from the add-avs tx.
	// IBCTimeoutPeriod is from the coordinator params.

	// DefaultSubscriberUnbondingPeriod is the default unbonding period for a subscriber chain.
	// There is no particular rationale to choose this number, since an
	// operator can opt-in to or opt-out from all AVSs independently.
	DefaultSubscriberUnbondingPeriod = 24 * 7 * time.Hour
	// DefaultBlocksPerDistributionTransmission is the default number of blocks after which a reward distribution is
	// transmitted to the coordinator. TODO: This can be replaced with epochs.
	DefaultBlocksPerDistributionTransmission = 1000
	// DefaultSubscriberRedistributionFraction is the default percentage of rewards that are distributed
	// to the coordinator.
	DefaultSubscriberRedistributionFraction = "0.1"
	// DefaultTransferTimeoutPeriod is the default timeout period for IBC transfers.
	DefaultTransferTimeoutPeriod = time.Hour
)

// NewSubscriberParams returns a new subscriber params with the given values.
func NewSubscriberParams(
	coordinatorFeePoolAddrStr string,
	distributionTransmissionChannel string,
	blocksPerDistributionTransmission int64,
	subscriberRedistributionFraction string,
	rewardDenom string,
	ibcTimeoutPeriod time.Duration,
	transferTimeoutPeriod time.Duration,
	unbondingPeriod time.Duration,
	historicalEntries int64,
	slashFractionDowntime string,
	downtimeJailDuration time.Duration,
	slashFractionDoubleSign string,
) *SubscriberParams {
	return (&SubscriberParams{}).
		WithCoordinatorFeePoolAddrStr(coordinatorFeePoolAddrStr).
		WithDistributionTransmissionChannel(distributionTransmissionChannel).
		WithBlocksPerDistributionTransmission(blocksPerDistributionTransmission).
		WithSubscriberRedistributionFraction(subscriberRedistributionFraction).
		WithRewardDenom(rewardDenom).
		WithIBCTimeoutPeriod(ibcTimeoutPeriod).
		WithTransferTimeoutPeriod(transferTimeoutPeriod).
		WithUnbondingPeriod(unbondingPeriod).
		WithHistoricalEntries(historicalEntries).
		WithSlashFractionDowntime(slashFractionDowntime).
		WithDowntimeJailDuration(downtimeJailDuration).
		WithSlashFractionDoubleSign(slashFractionDoubleSign)
}

// DefaultSubscriberParams returns the default subscriber params.
func DefaultSubscriberParams() *SubscriberParams {
	return NewSubscriberParams(
		"", // this is set by the coordinator, not by the subscriber
		"", // this is set by the coordinator, not by the subscriber
		DefaultBlocksPerDistributionTransmission,
		DefaultSubscriberRedistributionFraction,
		"", // this should be set in the transaction to register the subscriber
		DefaultIBCTimeoutPeriod,
		DefaultTransferTimeoutPeriod,
		DefaultSubscriberUnbondingPeriod,
		int64(stakingtypes.DefaultHistoricalEntries),
		slashingtypes.DefaultSlashFractionDowntime.String(),
		DefaultSubscriberUnbondingPeriod,
		slashingtypes.DefaultSlashFractionDoubleSign.String(),
	)
}

// Validate the subscriber params.
func (p SubscriberParams) Validate() error {
	// CoordinatorFeePoolAddrStr needs no validations.
	if err := ValidateDistributionTransmissionChannel(p.DistributionTransmissionChannel); err != nil {
		return fmt.Errorf("distribution transmission channel: %w", err)
	}
	if err := ValidatePositiveInt64(p.BlocksPerDistributionTransmission); err != nil {
		return fmt.Errorf("blocks per distribution transmission: %w", err)
	}
	if err := ValidateStringFraction(p.SubscriberRedistributionFraction); err != nil {
		return fmt.Errorf("subscriber redistribution fraction: %w", err)
	}
	if err := ValidateDuration(p.TransferTimeoutPeriod); err != nil {
		return fmt.Errorf("transfer timeout period: %w", err)
	}
	if err := ValidatePositiveInt64(p.HistoricalEntries); err != nil {
		return fmt.Errorf("historical entries: %w", err)
	}
	if err := ValidateDenomination(p.RewardDenom); err != nil {
		return fmt.Errorf("reward denom: %w", err)
	}
	if err := ValidateDuration(p.UnbondingPeriod); err != nil {
		return fmt.Errorf("unbonding period: %w", err)
	}
	if err := ValidateDuration(p.IBCTimeoutPeriod); err != nil {
		return fmt.Errorf("IBC timeout period: %w", err)
	}
	if err := ValidateStringFraction(p.SlashFractionDowntime); err != nil {
		return fmt.Errorf("slash fraction downtime: %w", err)
	}
	if err := ValidateStringFraction(p.SlashFractionDoubleSign); err != nil {
		return fmt.Errorf("slash fraction double sign: %w", err)
	}
	// technically duration could still be 0 for downtime, but
	// it's not a good idea.
	if err := ValidateDuration(p.DowntimeJailDuration); err != nil {
		return fmt.Errorf("downtime jail duration: %w", err)
	}
	return nil
}

// ValidateDistributionTransmissionChannel validates the distribution transmission channel.
// The channel must be a valid IBC channel identifier, or unset.
func ValidateDistributionTransmissionChannel(i string) error {
	if i == "" {
		// if it is unset, it means that the channel is yet to be
		// created, so we don't need to validate it.
		// this is why the function is called DistributionChannel
		// and not just Channel
		return nil
	}
	return ibchost.ChannelIdentifierValidator(i)
}

// ValidatePositiveInt64 validates that the given int64 is positive.
func ValidatePositiveInt64(i int64) error {
	if i <= 0 {
		return fmt.Errorf("int must be positive")
	}
	return nil
}

// ValidateDenomination validates that the given denomination is valid.
func ValidateDenomination(denom string) error {
	return sdk.ValidateDenom(denom)
}

// option pattern implementation, for clarity and ease of instantiation

// WithCoordinatorFeePoolAddrStr sets the CoordinatorFeePoolAddrStr field
func (p *SubscriberParams) WithCoordinatorFeePoolAddrStr(
	coordinatorFeePoolAddrStr string,
) *SubscriberParams {
	p.CoordinatorFeePoolAddrStr = coordinatorFeePoolAddrStr
	return p
}

// WithDistributionTransmissionChannel sets the DistributionTransmissionChannel field
func (p *SubscriberParams) WithDistributionTransmissionChannel(
	distributionTransmissionChannel string,
) *SubscriberParams {
	p.DistributionTransmissionChannel = distributionTransmissionChannel
	return p
}

// WithBlocksPerDistributionTransmission sets the BlocksPerDistributionTransmission field
func (p *SubscriberParams) WithBlocksPerDistributionTransmission(
	blocksPerDistributionTransmission int64,
) *SubscriberParams {
	p.BlocksPerDistributionTransmission = blocksPerDistributionTransmission
	return p
}

// WithSubscriberRedistributionFraction sets the SubscriberRedistributionFraction field
func (p *SubscriberParams) WithSubscriberRedistributionFraction(
	subscriberRedistributionFraction string,
) *SubscriberParams {
	p.SubscriberRedistributionFraction = subscriberRedistributionFraction
	return p
}

// WithRewardDenom sets the RewardDenom field
func (p *SubscriberParams) WithRewardDenom(
	rewardDenom string,
) *SubscriberParams {
	p.RewardDenom = rewardDenom
	return p
}

// WithIBCTimeoutPeriod sets the IBCTimeoutPeriod field
func (p *SubscriberParams) WithIBCTimeoutPeriod(
	ibcTimeoutPeriod time.Duration,
) *SubscriberParams {
	p.IBCTimeoutPeriod = ibcTimeoutPeriod
	return p
}

// WithTransferTimeoutPeriod sets the TransferTimeoutPeriod field
func (p *SubscriberParams) WithTransferTimeoutPeriod(
	transferTimeoutPeriod time.Duration,
) *SubscriberParams {
	p.TransferTimeoutPeriod = transferTimeoutPeriod
	return p
}

// WithUnbondingPeriod sets the UnbondingPeriod field
func (p *SubscriberParams) WithUnbondingPeriod(
	unbondingPeriod time.Duration,
) *SubscriberParams {
	p.UnbondingPeriod = unbondingPeriod
	return p
}

// WithHistoricalEntries sets the HistoricalEntries field
func (p *SubscriberParams) WithHistoricalEntries(
	historicalEntries int64,
) *SubscriberParams {
	p.HistoricalEntries = historicalEntries
	return p
}

// WithSlashFractionDowntime sets the SlashFractionDowntime field
func (p *SubscriberParams) WithSlashFractionDowntime(
	slashFractionDowntime string,
) *SubscriberParams {
	p.SlashFractionDowntime = slashFractionDowntime
	return p
}

// WithDowntimeJailDuration sets the DowntimeJailDuration field
func (p *SubscriberParams) WithDowntimeJailDuration(
	downtimeJailDuration time.Duration,
) *SubscriberParams {
	p.DowntimeJailDuration = downtimeJailDuration
	return p
}

// WithSlashFractionDoubleSign sets the SlashFractionDoubleSign field
func (p *SubscriberParams) WithSlashFractionDoubleSign(
	slashFractionDoubleSign string,
) *SubscriberParams {
	p.SlashFractionDoubleSign = slashFractionDoubleSign
	return p
}
