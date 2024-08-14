package exocore

import (
	"errors"

	"github.com/ExocoreNetwork/exocore/app/ante/utils"
	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type IncrementNonceDecorator struct {
	oracleKeeper utils.OracleKeeper
}

func NewIncrementNonceDecorator() IncrementNonceDecorator {
	return IncrementNonceDecorator{}
}

func (ind IncrementNonceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Increment the nonce of the account
	// This is done to prevent replay attacks
	// The nonce is incremented even if the transaction is invalid
	if !utils.IsOracleCreatePriceTx(tx) {
		return next(ctx, tx, simulate)
	}
	for _, msg := range tx.GetMsgs() {
		msg := msg.(*oracletypes.MsgCreatePrice)
		if accAddress, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
			return ctx, errors.New("invalid address")
		} else if _, err := ind.oracleKeeper.CheckAndIncreaseNonce(ctx, sdk.ConsAddress(accAddress).String(), msg.FeederID, uint32(msg.Nonce)); err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
