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

func (VirtualISlashKeeper) IsOperatorFrozen(ctx sdk.Context, opAddr sdk.AccAddress) bool {
	return false
}

func (VirtualISlashKeeper) OperatorAssetSlashedProportion(ctx sdk.Context, opAddr sdk.AccAddress, assetID string, startHeight, endHeight uint64) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(0)
}

// DelegationHooks add for dogfood
type DelegationHooks interface {
	//AfterDelegation we don't want the ability to cancel delegation or undelegation so no return type for
	// either
	// for delegation, we only care about the address of the operator to cache the event
	AfterDelegation(ctx sdk.Context, operator sdk.AccAddress)
	//AfterUndelegationStarted for undelegation, we use the address of the operator to figure out the list of impacted
	// chains for that operator. and we need the identifier to hold it until confirmed by subscriber
	AfterUndelegationStarted(ctx sdk.Context, addr sdk.AccAddress, recordKey []byte)
	//AfterUndelegationCompleted whenever an undelegation completes, we should update the vote power of the associated operator
	// on all of the chain ids that are impacted
	AfterUndelegationCompleted(ctx sdk.Context, addr sdk.AccAddress)
}

type ExpectedOperatorInterface interface {
	IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool
	GetUnbondingExpirationBlockNumber(ctx sdk.Context, OperatorAddress sdk.AccAddress, startHeight uint64) uint64

	UpdateOptedInAssetsState(ctx sdk.Context, stakerID, assetID, operatorAddr string, opAmount sdkmath.Int) error
}
