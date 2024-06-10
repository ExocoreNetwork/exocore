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
func (k *Keeper) OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) error {
	if !k.IsOperator(ctx, operatorAddress) {
		return delegationtypes.ErrOperatorNotExist
	}
	// check optedIn info
	if !k.IsOptedIn(ctx, operatorAddress.String(), avsAddr) {
		return types.ErrNotOptedIn
	}
	if !common.IsHexAddress(avsAddr) {
		if avsAddr == ctx.ChainID() {
			found, _ := k.getOperatorConsKeyForChainID(ctx, operatorAddress, avsAddr)
			if found {
				// if the key exists, it should be in the process of being removed.
				// TODO: if slashing is moved to a snapshot approach, opt out should only be
				// performed if the key doesn't exist.
				if !k.IsOperatorRemovingKeyFromChainID(ctx, operatorAddress, avsAddr) {
					return types.ErrOperatorNotRemovingKey
				}
			}
		} else {
			return types.ErrInvalidAvsAddr
		}
	}

	// DeleteOperatorUSDValue, delete the operator voting power, it can facilitate to
	// update the voting powers of all opted-in operators at the end of epoch.
	// there isn't going to be any reward for the operator in this opted-out epoch.
	err := k.DeleteOperatorUSDValue(ctx, avsAddr, operatorAddress.String())
	if err != nil {
		return err
	}

	// set opted-out height
	handleFunc := func(info *types.OptedInfo) {
		// #nosec G701
		info.OptedOutHeight = uint64(ctx.BlockHeight())
	}
	err = k.HandleOptedInfo(ctx, operatorAddress.String(), avsAddr, handleFunc)
	if err != nil {
		return err
	}
	return nil
}
