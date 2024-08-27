package types

const (
	// EventTypeTimeout       = "timeout"
	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "error"
	AttributeChainID          = "chain_id"
	AttributeValidatorAddress = "validator_address"
	AttributeValSetUpdateId   = "valset_update_id"
	AttributeInfractionType   = "infraction_type"

	EventTypeChannelEstablished = "channel_established"
	EventTypePacket             = "common_packet"
	EventTypeTimeout            = "common_timeout"
)
