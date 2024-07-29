package keeper

import (
	"github.com/ethereum/go-ethereum/common"

	sdkmath "cosmossdk.io/math"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AssetPriceAndDecimal struct {
	Price        sdkmath.Int
	PriceDecimal uint8
	Decimal      uint32
}

// OptIn call this function to opt in AVS
func (k *Keeper) OptIn(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) error {
	// avsAddr should be an evm contract address or a chain id.
	if !common.IsHexAddress(avsAddr) {
		if avsAddr != ctx.ChainID() { // TODO: other chain ids besides this chain's.
			return types.ErrInvalidAvsAddr
		}
	}
	// check optedIn info
	if k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrAlreadyOptedIn
	}

	// call InitOperatorUSDValue to mark the operator has been opted into the AVS
	// but the actual voting power calculation and update will be performed at the
	// end of epoch of the AVS. So there isn't any reward in the opted-in epoch for the
	// operator
	err := k.InitOperatorUSDValue(ctx, avsAddr, operatorAddress.String())
	if err != nil {
		return err
	}

	// update opted-in info
	slashContract, err := k.avsKeeper.GetAVSSlashContract(ctx, avsAddr)
	if err != nil {
		return err
	}
	optedInfo := &types.OptedInfo{
		SlashContract: slashContract,
		// #nosec G701
		OptedInHeight:  uint64(ctx.BlockHeight()),
		OptedOutHeight: types.DefaultOptedOutHeight,
	}
	err = k.SetOptedInfo(ctx, operatorAddress.String(), avsAddr, optedInfo)
	if err != nil {
		return err
	}
	return nil
}

// OptOut call this function to opt out of AVS
func (k *Keeper) OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) (err error) {
	// avsAddr should be an evm contract address or a chain id.
	if !common.IsHexAddress(avsAddr) {
		if avsAddr != ctx.ChainID() { // TODO: other chain ids besides this chain's.
			return types.ErrInvalidAvsAddr
		}
		defer func() {
			if err == nil {
				// does not fail, because it performs no validations, just stores a key-value pair.
				// all validations are instead performed by this function.
				k.InitiateOperatorKeyRemovalForChainID(ctx, operatorAddress, avsAddr)
			}
		}()
	}
	if !k.IsOperator(ctx, operatorAddress) {
		return delegationtypes.ErrOperatorNotExist
	}
	// check if the operator is active. It's not allowed to opt-out if the operator
	// isn't opted-in or is jailed.
	if !k.IsActive(ctx, operatorAddress, avsAddr) {
		return types.ErrNotOptedIn
	}
	// do not allow frozen operators to do anything meaningful
	if k.slashKeeper.IsOperatorFrozen(ctx, operatorAddress) {
		return delegationtypes.ErrOperatorIsFrozen
	}

	// DeleteOperatorUSDValue, delete the operator voting power, it can facilitate to
	// update the voting powers of all opted-in operators at the end of epoch.
	// There might still be a reward for the operator in this opted-out epoch,
	// which is determined by the reward logic.
	err = k.DeleteOperatorUSDValue(ctx, avsAddr, operatorAddress.String())
	if err != nil {
		return err
	}

	// set opted-out height
	handleFunc := func(info *types.OptedInfo) {
		// #nosec G701
		info.OptedOutHeight = uint64(ctx.BlockHeight())
		// the opt out, although is requested now, is made effective at the end of the current epoch.
		// so this is not necessarily the OptedOutHeight, rather, it is the OptOutRequestHeight.
		// the height is not directly used, beyond ascertaining whether the operator is currently opted in/out.
		// so the difference due to the epoch scheduling is not too big a concern.
	}
	err = k.HandleOptedInfo(ctx, operatorAddress.String(), avsAddr, handleFunc)
	if err != nil {
		return err
	}
	return nil
}
