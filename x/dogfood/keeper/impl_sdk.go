package keeper

import (
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

// interface guards
var (
	_ slashingtypes.StakingKeeper = Keeper{}
	_ evidencetypes.StakingKeeper = Keeper{}
	_ genutiltypes.StakingKeeper  = Keeper{}
	_ clienttypes.StakingKeeper   = Keeper{} // implemented in `validators.go`
)

// GetParams is an implementation of the staking interface expected by the SDK's evidence
// module. The module does not use it, but it is part of the interface.
func (k Keeper) GetParams(sdk.Context) stakingtypes.Params {
	return stakingtypes.Params{}
}

// IterateValidators is an implementation of the staking interface expected by the SDK's
// slashing module. The slashing module uses it for two purposes: once at genesis to
// store a mapping of pub key to cons address (which is done by our operator module),
// and then during the invariants check to ensure that the total delegated amount
// matches that of each validator. Ideally, this invariant should be implemented
// by the delegation and/or deposit module(s) instead.
func (k Keeper) IterateValidators(sdk.Context,
	func(index int64, validator stakingtypes.ValidatorI) (stop bool),
) {
	// no op
}

// Validator is an implementation of the staking interface expected by the SDK's
// slashing module. The slashing module uses it to obtain a validator's information at
// its addition to the list of validators, and then to unjail a validator. The former
// is used to create the pub key to cons address mapping, which we do in the operator module.
// The latter should also be implemented in the operator module, or maybe the slashing module
// depending upon the finalized design. We don't need to implement this function here because
// we are not calling the AfterValidatorCreated hook in our module, so this will never be
// reached.
func (k Keeper) Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI {
	panic("unimplemented on this keeper")
}

// ValidatorByConsAddr is an implementation of the staking interface expected by the SDK's
// slashing and evidence modules.
// The slashing module calls this function when it observes downtime. The only requirement on
// the returned value is that it isn't nil, and the jailed status is accurately set (to prevent
// re-jailing of the same operator).
// The evidence module calls this function when it handles equivocation evidence. The returned
// value must not be nil and must not have an UNBONDED validator status (the default is
// unspecified), or evidence will reject it.
func (k Keeper) ValidatorByConsAddr(
	ctx sdk.Context,
	addr sdk.ConsAddress,
) stakingtypes.ValidatorI {
	found, accAddr := k.operatorKeeper.GetOperatorAddressForChainIDAndConsAddr(
		ctx, ctx.ChainID(), addr,
	)
	if !found {
		// replicate the behavior of the SDK's staking module; do not panic.
		return nil
	}
	return stakingtypes.Validator{
		Jailed: k.operatorKeeper.IsOperatorJailedForChainID(ctx, accAddr, ctx.ChainID()),
	}
}

// Slash is an implementation of the staking interface expected by the SDK's slashing module.
// It forwards the call to SlashWithInfractionReason with Infraction_INFRACTION_UNSPECIFIED.
// It is not called within the slashing module, but is part of the interface.
func (k Keeper) Slash(
	ctx sdk.Context, addr sdk.ConsAddress,
	infractionHeight, power int64,
	slashFactor sdk.Dec,
) math.Int {
	return k.SlashWithInfractionReason(
		ctx, addr, infractionHeight, power,
		slashFactor, stakingtypes.Infraction_INFRACTION_UNSPECIFIED,
	)
}

// SlashWithInfractionReason is an implementation of the staking interface expected by the
// SDK's slashing module. It is called when the slashing module observes an infraction
// of either downtime or equivocation (which is via the evidence module).
func (k Keeper) SlashWithInfractionReason(
	ctx sdk.Context, addr sdk.ConsAddress, infractionHeight, power int64,
	slashFactor sdk.Dec, infraction stakingtypes.Infraction,
) math.Int {
	found, accAddress := k.operatorKeeper.GetOperatorAddressForChainIDAndConsAddr(
		ctx, ctx.ChainID(), addr,
	)
	if !found {
		// TODO(mm): already slashed and removed from the set?
		return math.NewInt(0)
	}
	// TODO(mm): add list of assets to be slashed (and not just all of them).
	// based on yet to be finalized slashing design.
	return k.slashingKeeper.SlashWithInfractionReason(
		ctx, accAddress, infractionHeight,
		power, slashFactor, infraction,
	)
}

// Jail is an implementation of the staking interface expected by the SDK's slashing module.
// It delegates the call to the operator module. Alternatively, this may be handled
// by the slashing module depending upon the design decisions.
func (k Keeper) Jail(ctx sdk.Context, addr sdk.ConsAddress) {
	k.operatorKeeper.Jail(ctx, addr, ctx.ChainID())
	// TODO(mm)
	// once the operator module jails someone, a hook should be triggered
	// and the validator removed from the set. same for unjailing.
}

// Unjail is an implementation of the staking interface expected by the SDK's slashing module.
// The function is called by the slashing module only when it receives a request from the
// operator to do so. TODO(mm): We need to use the SDK's slashing module to allow for downtime
// slashing but somehow we need to prevent its Unjail function from being called by anyone.
func (k Keeper) Unjail(sdk.Context, sdk.ConsAddress) {
	panic("unimplemented on this keeper")
}

// Delegation is an implementation of the staking interface expected by the SDK's slashing
// module. The slashing module uses it to obtain the delegation information of a validator
// before unjailing it. If the slashing module's unjail function is never called, this
// function will never be called either.
func (k Keeper) Delegation(
	sdk.Context, sdk.AccAddress, sdk.ValAddress,
) stakingtypes.DelegationI {
	panic("unimplemented on this keeper")
}

// MaxValidators is an implementation of the staking interface expected by the SDK's slashing
// module. It is not called within the slashing module, but is part of the interface.
// It returns the maximum number of validators allowed in the network.
func (k Keeper) MaxValidators(ctx sdk.Context) uint32 {
	return k.GetMaxValidators(ctx)
}

// GetAllValidators is an implementation of the staking interface expected by the SDK's
// slashing module. It is not called within the slashing module, but is part of the interface.
func (k Keeper) GetAllValidators(sdk.Context) (validators []stakingtypes.Validator) {
	return []stakingtypes.Validator{}
}

// IsValidatorJailed is an implementation of the staking interface expected by the SDK's
// slashing module. It is called by the slashing module to record validator signatures
// for downtime tracking. We delegate the call to the operator keeper.
func (k Keeper) IsValidatorJailed(ctx sdk.Context, addr sdk.ConsAddress) bool {
	found, accAddr := k.operatorKeeper.GetOperatorAddressForChainIDAndConsAddr(
		ctx, ctx.ChainID(), addr,
	)
	if !found {
		// replicate the behavior of the SDK's staking module
		return false
	}
	return k.operatorKeeper.IsOperatorJailedForChainID(ctx, accAddr, ctx.ChainID())
}

// ApplyAndReturnValidatorSetUpdates is an implementation of the staking interface expected
// by the SDK's genutil module. It is used in the gentx command, which we do not need to
// support. So this function does nothing.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(
	sdk.Context,
) (updates []abci.ValidatorUpdate, err error) {
	return
}
