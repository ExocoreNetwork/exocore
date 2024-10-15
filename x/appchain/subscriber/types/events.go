package types

const (
	EventTypeFeeDistribution          = "fee_distribution"
	EventTypeFeeTransferChannelOpened = "fee_transfer_channel_opened"
	EventTypeSubscriberSlashRequest   = "subscriber_slash_request"
	EventTypeVSCMatured               = "vsc_matured"

	AttributeDistributionCurrentHeight = "distribution_current_height"
	AttributeDistributionDenom         = "distribution_denom"
	AttributeDistributionFraction      = "distribution_fraction"
	AttributeDistributionNextHeight    = "distribution_next_height"
	AttributeDistributionValue         = "distribution_value"
	AttributeSubscriberHeight          = "subscriber_height"
	AttributeTimestamp                 = "timestamp"
)
