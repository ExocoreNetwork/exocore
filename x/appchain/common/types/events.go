package types

const (
	AttributeChainID          = "chain_id"
	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "ack_error"
	AttributeInfractionType   = "infraction_type"
	AttributeValidatorAddress = "validator_address"
	AttributeValSetUpdateID   = "valset_update_id"

	EventTypeChannelEstablished = "channel_established"
	EventTypePacket             = "common_packet"
	EventTypeTimeout            = "common_timeout"
)
