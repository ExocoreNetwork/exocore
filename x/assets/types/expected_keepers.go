package types

import (
	sdkmath "cosmossdk.io/math"
	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OracleKeeper interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetID string) (oracletypes.Price, error)
	RegisterNewTokenAndSetTokenFeeder(ctx sdk.Context, oInfo *oracletypes.OracleInfo) error
	UpdateNativeTokenByDepositOrWithdraw(ctx sdk.Context, assetID, stakerAddr string, amount sdkmath.Int) sdkmath.Int
}

type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
