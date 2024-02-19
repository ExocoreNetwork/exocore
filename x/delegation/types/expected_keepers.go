package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var CanUndelegationDelayHeight = uint64(10)

type ISlashKeeper interface {
	IsOperatorFrozen(ctx sdk.Context, opAddr sdk.AccAddress) bool
	OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetID string, startHeight, endHeight uint64) sdkmath.LegacyDec
}

// VirtualISlashKeeper todo: When the actual keeper functionality has not been implemented yet, temporarily use the virtual keeper.
type VirtualISlashKeeper struct{}

func (VirtualISlashKeeper) IsOperatorFrozen(sdk.Context, sdk.AccAddress) bool {
	return false
}

func (VirtualISlashKeeper) OperatorAssetSlashedProportion(sdk.Context, sdk.AccAddress, string, uint64, uint64) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(0)
}

type OperatorOptedInMiddlewareKeeper interface {
	GetOperatorCanUndelegateHeight(ctx sdk.Context, assetID string, opAddr sdk.AccAddress, startHeight uint64) uint64
}

type VirtualOperatorOptedInKeeper struct{}

func (VirtualOperatorOptedInKeeper) GetOperatorCanUndelegateHeight(_ sdk.Context, _ string, _ sdk.AccAddress, startHeight uint64) uint64 {
	return startHeight + CanUndelegationDelayHeight
}
