package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
)

// NewGenesis returns a new genesis state with the given inputs.
func NewGenesis(
	params Params, chains []ClientChainInfo,
	tokens []StakingAssetInfo, deposits []DepositsByStaker,
) *GenesisState {
	return &GenesisState{
		Params:       params,
		ClientChains: chains,
		Tokens:       tokens,
		Deposits:     deposits,
	}
}

// DefaultGenesis returns the default genesis state. It intentionally
// does not have any supported assets or deposits, since these must
// be supplied manually before bootstrapping the chain. The same is true
// for any unit / integration tests.
func DefaultGenesis() *GenesisState {
	return NewGenesis(
		DefaultParams(), []ClientChainInfo{}, []StakingAssetInfo{}, []DepositsByStaker{},
	)
}

// Validate performs basic genesis state validation returning an error
// upon any failure.
func (gs GenesisState) Validate() error {
	// client_chain.go -> check no repeat chain
	lzIDs := make(map[uint64]struct{}, len(gs.ClientChains))
	for _, info := range gs.ClientChains {
		if _, ok := lzIDs[info.LayerZeroChainID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate LayerZeroChainID: %d",
				info.LayerZeroChainID,
			)
		}
		lzIDs[info.LayerZeroChainID] = struct{}{}
	}
	// client_chain_asset.go -> check presence of client chain
	// for all assets and no duplicates
	tokens := make(map[string]struct{}, len(gs.Tokens))
	for _, info := range gs.Tokens {
		id := info.AssetBasicInfo.LayerZeroChainID
		// check that the chain is registered
		if _, ok := lzIDs[id]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown LayerZeroChainID for token %s: %d",
				info.AssetBasicInfo.MetaInfo, id,
			)
		}
		address := info.AssetBasicInfo.Address
		// build for 0x addresses only.
		// TODO: consider removing this check for non-EVM client chains.
		if !common.IsHexAddress(address) {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"not hex address for token %s: %s",
				info.AssetBasicInfo.MetaInfo, address,
			)
		}
		// calculate the asset id.
		_, assetID := GetStakeIDAndAssetIDFromStr(
			info.AssetBasicInfo.LayerZeroChainID,
			"", address,
		)
		// check that it is not a duplicate.
		if _, ok := tokens[assetID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate assetID: %s",
				assetID,
			)
		}
		// ensure there are no deposits for this asset already (since they are handled in the
		// genesis exec). while it is possible to remove this field entirely (and assume 0),
		// i did not do so in order to make the genesis state more explicit.
		if !info.StakingTotalAmount.IsZero() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"non-zero deposit amount for asset %s",
				assetID,
			)
		}
		tokens[assetID] = struct{}{}
	}
	// staker_asset.go -> check deposits and withdrawals and that there is no unbonding.
	stakers := make(map[string]struct{}, len(gs.Deposits))
	for _, depositByStaker := range gs.Deposits {
		stakerID := depositByStaker.StakerID
		// validate the stakerID
		if stakerID != strings.ToLower(stakerID) {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"stakerID not lowercase: %s",
				stakerID,
			)
		}
		var stakerClientAddress string
		var lzID uint64
		var err error
		if stakerClientAddress, lzID, err = ParseID(stakerID); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid stakerID: %s",
				stakerID,
			)
		}
		// check that the chain is registered
		if _, ok := lzIDs[lzID]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown LayerZeroChainID for staker %s: %d",
				stakerID, lzID,
			)
		}
		// build for 0x addresses only.
		// TODO: consider removing this check for non-EVM client chains.
		if !common.IsHexAddress(stakerClientAddress) {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"not hex staker address for staker %s: %s",
				stakerID, stakerClientAddress,
			)
		}
		// check that it is not a duplicate
		if _, ok := stakers[stakerID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate stakerID: %s",
				stakerID,
			)
		}
		stakers[stakerID] = struct{}{}
		// map to check for duplicate tokens for the staker.
		tokensForStaker := make(map[string]struct{}, len(depositByStaker.Deposits))
		for _, deposit := range depositByStaker.Deposits {
			assetID := deposit.AssetID
			// check that the asset is registered
			// no need to check for the validity of the assetID, since
			// an invalid assetID cannot be in the tokens map.
			if _, ok := tokens[assetID]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown assetID for deposit %s: %s",
					stakerID, assetID,
				)
			}
			// check that it is not a duplicate
			if _, ok := tokensForStaker[assetID]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate assetID for staker %s: %s",
					stakerID, assetID,
				)
			}
			tokensForStaker[assetID] = struct{}{}
			info := deposit.Info
			// check that there is no nil value provided.
			if info.TotalDepositAmount.IsNil() || info.WithdrawableAmount.IsNil() ||
				info.WaitUnbondingAmount.IsNil() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"nil deposit info for %s: %+v",
					assetID, info,
				)
			}
			// at genesis (not chain restart), there is no unbonding amount.
			if !info.WaitUnbondingAmount.IsZero() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"non-zero unbonding amount for %s: %s",
					assetID, info.WaitUnbondingAmount,
				)
			}
			// check for negative values.
			if info.TotalDepositAmount.IsNegative() || info.WithdrawableAmount.IsNegative() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"negative deposit amount for %s: %+v",
					assetID, info,
				)
			}
			// check that the withdrawable amount is not greater than the total deposit amount.
			// since withdrawable amount should be less than or equal to the amount deposited.
			if info.WithdrawableAmount.GT(info.TotalDepositAmount) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"withdrawable amount exceeds total deposit amount for %s: %+v",
					assetID, info,
				)
			}
		}
	}
	return gs.Params.Validate()
}
