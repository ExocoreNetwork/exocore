package types

import (
	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OracleKeeper interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetID string) (oracletypes.Price, error)
	RegisterNewTokenAndSetTokenFeeder(ctx sdk.Context, oInfo *oracletypes.OracleInfo) error
}
