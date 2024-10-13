package types

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

const maxLogSize = 1024

// NewResultAcknowledgementWithLog creates a result acknowledgement with a log message.
func NewResultAcknowledgementWithLog(ctx sdk.Context, res []byte) channeltypes.Acknowledgement {
	if len(res) > maxLogSize {
		ctx.Logger().Info(
			"IBC ResultAcknowledgement constructed",
			"res_size", len(res),
			"res_preview", string(res[:maxLogSize]),
		)
	} else {
		ctx.Logger().Info(
			"IBC ResultAcknowledgement constructed",
			"res_size", len(res),
			"res", string(res),
		)
	}
	return channeltypes.NewResultAcknowledgement(res)
}

// NewErrorAcknowledgementWithLog creates an error acknowledgement with a log message.
func NewErrorAcknowledgementWithLog(ctx sdk.Context, err error) channeltypes.Acknowledgement {
	ctx.Logger().Error("IBC ErrorAcknowledgement constructed", "error", err)
	return channeltypes.NewErrorAcknowledgement(err)
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

// NewVscPacketData creates a new VscMaturedPacketData instance.
func NewVscMaturedPacketData(
	valsetUpdateID uint64,
) *VscMaturedPacketData {
	return &VscMaturedPacketData{
		ValsetUpdateID: valsetUpdateID,
	}
}
