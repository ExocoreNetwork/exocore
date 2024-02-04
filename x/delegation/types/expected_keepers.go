package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var CanUndelegationDelayHeight = uint64(10)

type ISlashKeeper interface {
	IsOperatorFrozen(ctx sdk.Context, opAddr sdk.AccAddress) bool
	OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec
}

// VirtualISlashKeeper todo: When the actual keeper functionality has not been implemented yet, temporarily use the virtual keeper.
type VirtualISlashKeeper struct{}

func (VirtualISlashKeeper) IsOperatorFrozen(ctx sdk.Context, opAddr sdk.AccAddress) bool {
	return false
}

func (VirtualISlashKeeper) OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetId string, startHeight, endHeight uint64) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(0)
}

type ExpectedOperatorInterface interface {
	IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool
	GetUnbondingExpirationBlockNumber(ctx sdk.Context, OperatorAddress sdk.AccAddress, startHeight uint64) uint64

	UpdateOptedInAssetsState(ctx sdk.Context, stakerId, assetId, operatorAddr string, opAmount sdkmath.Int) error
}
