package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

var _ evmtypes.StakingKeeper = Keeper{}

// GetValidatorByConsAddr is an implementation of the StakingKeeper interface
// expected by the EVM module. It returns a validator given a consensus address.
// The EVM module uses it to determine the proposer's AccAddress for the block.
// ConsAddress -> lookup AccAddress -> convert to ValAddress -> bech32ify.
// ValAddress string is then decoded back by the EVM module to get bytes, which
// make the 0x address of the coinbase.
func (k Keeper) GetValidatorByConsAddr(
	ctx sdk.Context, consAddr sdk.ConsAddress,
) (validator stakingtypes.Validator, found bool) {
	val := k.ValidatorByConsAddr(ctx, consAddr)
	if val == nil {
		return stakingtypes.Validator{}, false
	}
	return val.(stakingtypes.Validator), true
}
