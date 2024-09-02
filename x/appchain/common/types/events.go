package types

const (
	AttributeChainID       = "chain_id"
	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "ack_error"

	EventTypeChannelEstablished = "channel_established"
	EventTypePacket             = "common_packet"
	EventTypeTimeout            = "common_timeout"
)
