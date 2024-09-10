package types

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

// NewErrorAcknowledgementWithLog creates an error acknowledgement with a log message.
func NewErrorAcknowledgementWithLog(ctx sdk.Context, err error) channeltypes.Acknowledgement {
	ctx.Logger().Error("IBC ErrorAcknowledgement constructed", "error", err)
	return channeltypes.NewErrorAcknowledgement(err)
}

// NewResultAcknowledgementWithLog creates a result acknowledgement with a log message.
func NewResultAcknowledgementWithLog(ctx sdk.Context, res []byte) channeltypes.Acknowledgement {
	ctx.Logger().Info("IBC ResultAcknowledgement constructed", "res", res)
	return channeltypes.NewResultAcknowledgement(res)
}

// NewVscPacketData creates a new ValidatorSetChangePacketData instance.
func NewVscPacketData(
	updates []abci.ValidatorUpdate,
	valsetUpdateID uint64,
	slashAcks [][]byte,
) ValidatorSetChangePacketData {
	return ValidatorSetChangePacketData{
		ValidatorUpdates: updates,
		ValsetUpdateID:   valsetUpdateID,
		SlashAcks:        slashAcks,
	}
}

func NewVscMaturedPacketData(
	valsetUpdateID uint64,
) *VscMaturedPacketData {
	return &VscMaturedPacketData{
		ValsetUpdateID: valsetUpdateID,
	}
}
