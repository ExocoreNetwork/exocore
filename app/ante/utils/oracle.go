package utils

import (
	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const CreatePriceGas = 20000

func IsOracleCreatePriceTx(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return false
	}
	for _, msg := range msgs {
		if _, ok := msg.(*oracletypes.MsgCreatePrice); ok {
			continue
		}
		return false
	}
	return true
}
