package types

const (
	EventTypeCreatePrice = "create_price"

	AttributeKeyFeederID      = "feeder_id"
	AttributeKeyTokenID       = "token_id"
	AttributeKeyBasedBlock    = "based_block"
	AttributeKeyRoundID       = "round_id"
	AttributeKeyProposer      = "proposer"
	AttributeKeyFinalPrice    = "final_price"
	AttributeKeyPriceUpdated  = "price_update"
	AttributeKeyParamsUpdated = "params_update"
	AttributeKeyFeederIDs     = "feeder_ids"

	AttributeValuePriceUpdatedSuccess  = "success"
	AttributeValueParamsUpdatedSuccess = "success"
)
