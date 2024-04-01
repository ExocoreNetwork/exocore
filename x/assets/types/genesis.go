package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesis returns a new genesis state with the given inputs.
func NewGenesis(
	params Params, chains []ClientChainInfo, tokens []StakingAssetInfo,
	opAssets []OpAssetIDAndInfos, stAssets []StAssetIDAndInfos,
) *GenesisState {
	return &GenesisState{
		Params:             params,
		ClientChains:       chains,
		Tokens:             tokens,
		OperatorAssetInfos: opAssets,
		StakerAssetInfos:   stAssets,
	}
}

// DefaultGenesis returns the default genesis state. It intentionally
// does not have any supported assets or deposits, since these must
// be supplied manually before bootstrapping the chain.
func DefaultGenesis() *GenesisState {
	return NewGenesis(
		DefaultParams(), []ClientChainInfo{}, []StakingAssetInfo{},
		[]OpAssetIDAndInfos{}, []StAssetIDAndInfos{},
	)
}

// Validate performs basic genesis state validation returning an error
// upon any failure.
func (gs GenesisState) Validate() error {
	// TODO: validate the other items as well.
	// such validations could potentially include that the sum
	// of all balances for an asset equals the deposits, etc.
	// operator_asset.go

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
	deposits := make(map[string]math.Int, len(gs.Tokens))
	for _, info := range gs.Tokens {
		id := info.AssetBasicInfo.LayerZeroChainID
		if _, ok := lzIDs[id]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown LayerZeroChainID for token %s: %d",
				info.AssetBasicInfo.MetaInfo, id,
			)
		}
		_, assetID := GetStakeIDAndAssetIDFromStr(
			info.AssetBasicInfo.LayerZeroChainID,
			"", info.AssetBasicInfo.Address,
		)
		if _, ok := deposits[assetID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate assetID: %s",
				assetID,
			)
		}
		if info.StakingTotalAmount.IsNegative() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"negative staking total amount for asset %s",
				assetID,
			)
		}
		deposits[assetID] = info.StakingTotalAmount
	}
	// operator_asset.go => check for duplicates and that operator address
	// is bech32 encoded sdk.AccAddress and note the delegated amount.
	delegationsToOperator := make(map[string](map[string]math.Int), len(gs.OperatorAssetInfos))
	delegationsByAssetToOperator := make(map[string]math.Int, len(deposits))
	for _, level1 := range gs.OperatorAssetInfos {
		addr := level1.OperatorAddress
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return err
		}
		if _, ok := delegationsToOperator[addr]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate operator address: %s",
				addr,
			)
		}
		delegationsToOperator[addr] = make(map[string]math.Int, len(level1.AssetIdAndInfos))
		for _, level2 := range level1.AssetIdAndInfos {
			// this is compared against deposits assetID, so we don't need to validate it
			assetID := level2.AssetID
			info := level2.Info
			// check that all amounts are supplied; do not accept nil values
			if info.TotalAmount.IsNil() || info.OperatorAmount.IsNil() ||
				info.WaitUnbondingAmount.IsNil() || info.OperatorUnbondingAmount.IsNil() ||
				info.OperatorUnbondableAmountAfterSlash.IsNil() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"missing amounts for operator %s: %s",
					addr, assetID,
				)
			}
			// for a genesis bootstrap (not a chain restart), this condition should hold.
			if !info.OperatorAmount.IsZero() ||
				!info.WaitUnbondingAmount.IsZero() ||
				!info.OperatorUnbondingAmount.IsZero() ||
				!info.OperatorUnbondableAmountAfterSlash.IsZero() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unexpected non-zero amounts for operator %s: %s",
					addr, assetID,
				)
			}
			amount := info.TotalAmount
			// check that deposits against this asset actually exist
			if _, ok := deposits[assetID]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown assetID for operator %s: %s",
					addr, assetID,
				)
			}
			// check for duplicates
			if _, ok := delegationsToOperator[addr][assetID]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate assetID entry for operator %s: %s",
					addr, assetID,
				)
			}
			if amount.IsNegative() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"negative assetID entry for operator %s: %s",
					addr, assetID,
				)
			}
			delegationsToOperator[addr][assetID] = amount
			if _, ok := delegationsByAssetToOperator[assetID]; ok {
				delegationsByAssetToOperator[assetID] = delegationsByAssetToOperator[assetID].Add(
					amount,
				)
			} else {
				delegationsByAssetToOperator[assetID] = amount
			}
		}
	}
	// staker_asset.go => check for duplicates and that staker ID is not empty
	// and note the delegated amount.
	delegationsByStaker := make(map[string](map[string]math.Int), len(gs.StakerAssetInfos))
	delegationsByAssetByStaker := make(map[string]math.Int, len(deposits))
	depositsByAssetByStaker := make(map[string]math.Int, len(deposits))
	for _, level1 := range gs.StakerAssetInfos {
		staker := level1.StakerID
		if staker == "" {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"empty staker ID",
			)
		}
		if _, id, err := ParseID(staker); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid staker ID %s: %s", staker, err,
			)
		} else {
			if _, ok := lzIDs[id]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown LayerZeroChainID for staker %s: %d",
					staker, id,
				)
			}
		}
		if _, ok := delegationsByStaker[staker]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate staker ID: %s",
				staker,
			)
		}
		delegationsByStaker[staker] = make(map[string]math.Int, len(level1.AssetIdAndInfos))
		for _, level2 := range level1.AssetIdAndInfos {
			// this is compared against deposits assetID, so we don't need to validate it
			assetID := level2.AssetID
			info := level2.Info
			// check that all amounts are supplied; do not accept nil values
			if info.TotalDepositAmount.IsNil() || info.WithdrawableAmount.IsNil() ||
				info.WaitUnbondingAmount.IsNil() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"missing amounts for staker %s: %s",
					staker, assetID,
				)
			}
			if !info.WaitUnbondingAmount.IsZero() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unexpected non-zero amounts for staker %s: %s",
					staker, assetID,
				)
			}
			// delegated amount is the difference between total deposit and withdrawable amount
			amount := info.TotalDepositAmount.Sub(
				info.WithdrawableAmount,
			)
			// check that deposits against this asset actually exist
			if _, ok := deposits[assetID]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown assetID for staker %s: %s",
					staker, assetID,
				)
			}
			// check for duplicates
			if _, ok := delegationsByStaker[staker][assetID]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate assetID entry for staker %s: %s",
					staker, assetID,
				)
			}
			if amount.IsNegative() || info.TotalDepositAmount.IsNegative() ||
				info.WithdrawableAmount.IsNegative() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"negative assetID entry for staker %s: %s",
					staker, assetID,
				)
			}
			delegationsByStaker[staker][assetID] = amount
			if _, ok := delegationsByAssetByStaker[assetID]; ok {
				delegationsByAssetByStaker[assetID] = delegationsByAssetByStaker[assetID].Add(
					amount,
				)
			} else {
				delegationsByAssetByStaker[assetID] = amount
			}
			if _, ok := depositsByAssetByStaker[assetID]; ok {
				depositsByAssetByStaker[assetID] = depositsByAssetByStaker[assetID].Add(
					level2.Info.TotalDepositAmount,
				)
			} else {
				depositsByAssetByStaker[assetID] = level2.Info.TotalDepositAmount
			}
		}
	}
	// check that the total amount delegated to operators equals the amount delegated by stakers
	if !areMapsIdentical(delegationsByAssetToOperator, delegationsByAssetByStaker) {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"delegations to operators and by stakers don't match",
		)
	}
	// check that the total deposits by stakers equals the deposits by assets
	if !areMapsIdentical(deposits, depositsByAssetByStaker) {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"deposits and deposits by stakers don't match",
		)
	}
	return gs.Params.Validate()
}

func areMapsIdentical(m1, m2 map[string]math.Int) bool {
	if len(m1) != len(m2) {
		return false
	}

	// this code is not consensus critical, so we can loop.
	// the codeQL warning can be ignored.
	
	for k1, v1 := range m1 {
		v2, ok := m2[k1]
		if !ok || !v1.Equal(v2) {
			return false
		}
	}

	return true
}
