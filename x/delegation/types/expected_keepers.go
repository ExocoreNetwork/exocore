// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	CanUnDelegationDelayHeight = uint64(10)
)

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

type OperatorOptedInMiddlewareKeeper interface {
	GetOperatorCanUnDelegateHeight(ctx sdk.Context, assetId string, opAddr sdk.AccAddress, startHeight uint64) uint64
}

type VirtualOperatorOptedInKeeper struct{}

func (VirtualOperatorOptedInKeeper) GetOperatorCanUnDelegateHeight(ctx sdk.Context, assetId string, opAddr sdk.AccAddress, startHeight uint64) uint64 {
	return startHeight + CanUnDelegationDelayHeight
}
