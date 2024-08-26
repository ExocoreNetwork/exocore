package types

import (
	"math/big"

	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OracleKeeper interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetID string) (oracletypes.Price, error)
	UpdateNativeTokenByDepositOrWithdraw(ctx sdk.Context, assetID, stakerAddr, amount string) *big.Int
}
