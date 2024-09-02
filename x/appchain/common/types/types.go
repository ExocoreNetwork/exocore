package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

// NewErrorAcknowledgementWithLog creates an error acknowledgement with a log message.
func NewErrorAcknowledgementWithLog(ctx sdk.Context, err error) channeltypes.Acknowledgement {
	ctx.Logger().Error("IBC ErrorAcknowledgement constructed", "error", err)
	return channeltypes.NewErrorAcknowledgement(err)
}
