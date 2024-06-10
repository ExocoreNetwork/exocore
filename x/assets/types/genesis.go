package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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
	for i, info := range gs.ClientChains {
		if info.Name == "" {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil Name for chain %d",
				i,
			)
		}
		// this is our primary method of cross-chain communication.
		if info.LayerZeroChainID == 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil LayerZeroChainID for chain %s",
				info.Name,
			)
		}
		// the address length is used to convert from bytes32 to address.
		if info.AddressLength == 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil AddressLength for chain %s",
				info.Name,
			)
		}
		// check for no duplicated chain, indexed by LayerZeroChainID.
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
	tokenSupplies := make(map[string]math.Int, len(gs.Tokens))
	for _, info := range gs.Tokens {
		if info.AssetBasicInfo == nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil AssetBasicInfo for token %s",
				info.AssetBasicInfo.MetaInfo,
			)
		}
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
		// ensure there are no deposits for this asset already (since they are handled in the
		// genesis exec). while it is possible to remove this field entirely (and assume 0),
		// i did not do so in order to make the genesis state more explicit.
		if info.StakingTotalAmount.IsNil() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil staking total amount for asset %s",
				assetID,
			)
		}
		if !info.StakingTotalAmount.IsZero() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"non-zero deposit amount for asset %s",
				assetID,
			)
		}
		// check that it is not a duplicate.
		if _, ok := tokenSupplies[assetID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate assetID: %s",
				assetID,
			)
		}
		// validate the amount of supply
		if info.AssetBasicInfo.TotalSupply.IsNil() ||
			!info.AssetBasicInfo.TotalSupply.IsPositive() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil total supply for token %s",
				info.AssetBasicInfo.MetaInfo,
			)
		}
		tokenSupplies[assetID] = info.AssetBasicInfo.TotalSupply
	}
	// staker_asset.go -> check deposits and withdrawals and that there is no unbonding.
	stakers := make(map[string]struct{}, len(gs.Deposits))
	for _, depositByStaker := range gs.Deposits {
		stakerID := depositByStaker.StakerID
		// validate the stakerID
		var stakerClientChainID uint64
		var err error
		if _, stakerClientChainID, err = ValidateID(stakerID, true, true); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid stakerID: %s",
				stakerID,
			)
		}
		// check that the chain is registered
		if _, ok := lzIDs[stakerClientChainID]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown LayerZeroChainID for staker %s: %d",
				stakerID, stakerClientChainID,
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
			if _, ok := tokenSupplies[assetID]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown assetID for deposit %s: %s",
					stakerID, assetID,
				)
			}
			// #nosec G703 // if it's invalid, we will not reach here.
			_, assetClientChainID, _ := ParseID(assetID)
			if assetClientChainID != stakerClientChainID {
				// we can reach here if there are multiple chains
				// and it tries to deposit assets from one chain
				// under a staker from another chain.
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"mismatched client chain IDs for staker %s and asset %s",
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
			// check that the withdrawable amount and the deposited amount are equal.
			// this is because this module's genesis only sets up free deposits.
			// the delegation module bonds them, thereby altering the withdrawable amount.
			if !info.WithdrawableAmount.Equal(info.TotalDepositAmount) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"withdrawable amount is not equal to total deposit amount for %s: %+v",
					assetID, info,
				)
			}
			// check that deposit amount does not exceed supply.
			if info.TotalDepositAmount.GT(tokenSupplies[assetID]) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"deposit amount exceeds max supply for %s: %+v",
					assetID, info,
				)
			}
		}
	}
	return gs.Params.Validate()
}
