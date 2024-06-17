package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/constraints"
)

type TotalSupplyAndStaking struct {
	TotalStaking math.Int
	TotalSupply  math.Int
}

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

// CommonValidation is used to check for duplicates in the input list
// and validate the input information simultaneously.
// It might be used for validating most genesis states.
func CommonValidation[T any, V constraints.Ordered, D any](
	slice []T,
	seenFieldValue func(T) (V, D),
	validation func(int, T) error,
) (map[V]D, error) {
	seen := make(map[V]D)
	for i, v := range slice {
		field, value := seenFieldValue(v)
		// check for no duplicated element
		if _, ok := seen[field]; ok {
			return nil, errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate element: %v",
				field,
			)
		}
		// perform the validation
		if err := validation(i, v); err != nil {
			return nil, err
		}
		seen[field] = value
	}
	return seen, nil
}

// ValidateClientChains performs basic client chains validation
func (gs GenesisState) ValidateClientChains() (map[uint64]struct{}, error) {
	validationFunc := func(i int, info ClientChainInfo) error {
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
		return nil
	}
	seenFieldValueFunc := func(info ClientChainInfo) (uint64, struct{}) {
		return info.LayerZeroChainID, struct{}{}
	}
	lzIDs, err := CommonValidation(gs.ClientChains, seenFieldValueFunc, validationFunc)
	if err != nil {
		return nil, err
	}
	return lzIDs, nil
}

// ValidateTokens performs basic client chain assets validation
func (gs GenesisState) ValidateTokens(lzIDs map[uint64]struct{}) (map[string]TotalSupplyAndStaking, error) {
	validationFunc := func(i int, info StakingAssetInfo) error {
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

		// ensure there are no deposits for this asset already (since they are handled in the
		// genesis exec). while it is possible to remove this field entirely (and assume 0),
		// i did not do so in order to make the genesis state more explicit.
		if info.StakingTotalAmount.IsNil() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil staking total amount for asset %s",
				info.AssetBasicInfo.Address,
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

		// the StakingTotalAmount shouldn't be greater than the total supply when init from the general
		// exported genesis file
		if info.StakingTotalAmount.GT(info.AssetBasicInfo.TotalSupply) {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"total staking amount is greater than the total supply for token %s",
				info.AssetBasicInfo.MetaInfo,
			)
		}
		return nil
	}
	seenFieldValueFunc := func(info StakingAssetInfo) (string, TotalSupplyAndStaking) {
		// calculate the asset id.
		_, assetID := GetStakeIDAndAssetIDFromStr(
			info.AssetBasicInfo.LayerZeroChainID,
			"", info.AssetBasicInfo.Address,
		)
		return assetID, TotalSupplyAndStaking{
			TotalSupply:  info.AssetBasicInfo.TotalSupply,
			TotalStaking: info.StakingTotalAmount,
		}
	}
	totalSupplyAndStaking, err := CommonValidation(gs.Tokens, seenFieldValueFunc, validationFunc)
	if err != nil {
		return nil, err
	}
	return totalSupplyAndStaking, nil
}

// ValidateDeposits performs basic deposits validation
func (gs GenesisState) ValidateDeposits(lzIDs map[uint64]struct{}, tokenState map[string]TotalSupplyAndStaking) error {
	validationFunc := func(i int, depositByStaker DepositsByStaker) error {
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

		// validate the deposits
		depositValidationF := func(i int, deposit DepositByAsset) error {
			assetID := deposit.AssetID
			// check that the asset is registered
			// no need to check for the validity of the assetID, since
			// an invalid assetID cannot be in the tokens map.
			if _, ok := tokenState[assetID]; !ok {
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

			// check for negative values.
			if info.TotalDepositAmount.IsNegative() || info.WithdrawableAmount.IsNegative() ||
				info.WaitUnbondingAmount.IsNegative() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"negative deposit amount for %s: %+v",
					assetID, info,
				)
			}
			return nil
		}
		depositFieldValueF := func(deposit DepositByAsset) (string, struct{}) {
			return deposit.AssetID, struct{}{}
		}
		_, err = CommonValidation(depositByStaker.Deposits, depositFieldValueF, depositValidationF)
		if err != nil {
			return err
		}
		return nil
	}
	seenFieldValueFunc := func(deposits DepositsByStaker) (string, struct{}) {
		return deposits.StakerID, struct{}{}
	}
	_, err := CommonValidation(gs.Deposits, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

// ValidateOperatorAssets performs basic operator assets validation
func (gs GenesisState) ValidateOperatorAssets(tokenState map[string]TotalSupplyAndStaking) error {
	validationFunc := func(i int, assets AssetsByOperator) error {
		_, err := sdk.AccAddressFromBech32(assets.Operator)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid operator address %s: %s", assets.Operator, err,
			)
		}

		// validate the assets list for the specified operator
		assetValidationFunc := func(i int, asset AssetByID) error {
			// check that the asset is registered
			// no need to check for the validity of the assetID, since
			// an invalid assetID cannot be in the tokens map.
			if _, ok := tokenState[asset.AssetID]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown assetID for operator assets %s: %s",
					assets.Operator, asset.AssetID,
				)
			}
			// the sum amount of operators shouldn't be greater than the total staking amount of this asset
			if asset.Info.TotalAmount.Add(asset.Info.WaitUnbondingAmount).GT(tokenState[asset.AssetID].TotalStaking) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"operator's sum amount exceeds the total staking amount for %s: %+v",
					assets.Operator, asset,
				)
			}

			if asset.Info.OperatorAmount.GT(asset.Info.TotalAmount) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"operator's amount exceeds the total amount for %s: %+v",
					assets.Operator, asset,
				)
			}

			if asset.Info.OperatorUnbondingAmount.GT(asset.Info.WaitUnbondingAmount) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"operator's unbonding amount exceeds the total unbonding amount for %s: %+v",
					assets.Operator, asset,
				)
			}

			if asset.Info.OperatorShare.GT(asset.Info.TotalShare) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"operator's share exceeds the total share for %s: %+v",
					assets.Operator, asset,
				)
			}
			return nil
		}
		assetFieldValueFunc := func(asset AssetByID) (string, struct{}) {
			return asset.AssetID, struct{}{}
		}
		_, err = CommonValidation(assets.AssetsState, assetFieldValueFunc, assetValidationFunc)
		if err != nil {
			return err
		}
		return nil
	}
	seenFieldValueFunc := func(assets AssetsByOperator) (string, struct{}) {
		return assets.Operator, struct{}{}
	}
	_, err := CommonValidation(gs.OperatorAssets, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

// Validate performs basic genesis state validation returning an error
// upon any failure.
func (gs GenesisState) Validate() error {
	lzIDs, err := gs.ValidateClientChains()
	if err != nil {
		return err
	}
	totalSupplyAndStaking, err := gs.ValidateTokens(lzIDs)
	if err != nil {
		return err
	}
	err = gs.ValidateDeposits(lzIDs, totalSupplyAndStaking)
	if err != nil {
		return err
	}
	err = gs.ValidateOperatorAssets(totalSupplyAndStaking)
	if err != nil {
		return err
	}
	return gs.Params.Validate()
}

// todo: This should be removed if the refactored validate is fine.
// Validate performs basic genesis state validation returning an error
// upon any failure.
/*func (gs GenesisState) Validate() error {
	// client_chain.go -> check no repeat chain
	lzIDs := make(map[uint64]struct{}, len(gs.ClientChains))
	for i, info := range gs.ClientChains {
		if info.Name == "" || len(info.Name) > MaxChainTokenNameLength {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil Name or too long for chain %d, name:%s, maxLength:%d",
				i, info.Name, MaxChainTokenNameLength,
			)
		}
		if info.MetaInfo == "" || len(info.MetaInfo) > MaxChainTokenMetaInfoLength {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil meta info or too long for chain %d, metaInfo:%s, maxLength:%d",
				i, info.MetaInfo, MaxChainTokenMetaInfoLength,
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
	tokens := make(map[string]struct{}, len(gs.Tokens))
	for _, info := range gs.Tokens {
		if info.AssetBasicInfo == nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil AssetBasicInfo for token %s",
				info.AssetBasicInfo.MetaInfo,
			)
		}
		if info.AssetBasicInfo.Name == "" || len(info.AssetBasicInfo.Name) > MaxChainTokenNameLength {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil Name or too long for token, name:%s, maxLength:%d",
				info.AssetBasicInfo.Name, MaxChainTokenNameLength,
			)
		}
		if info.AssetBasicInfo.MetaInfo == "" || len(info.AssetBasicInfo.MetaInfo) > MaxChainTokenMetaInfoLength {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil meta info or too long for token, metaInfo:%s, maxLength:%d",
				info.AssetBasicInfo.MetaInfo, MaxChainTokenMetaInfoLength,
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
		if _, ok := tokens[assetID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate assetID: %s",
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
			if _, ok := tokens[assetID]; !ok {
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
				info.PendingUndelegationAmount.IsNil() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"nil deposit info for %s: %+v",
					assetID, info,
				)
			}
			// at genesis (not chain restart), there is no unbonding amount.
			if !info.PendingUndelegationAmount.IsZero() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"non-zero unbonding amount for %s: %s",
					assetID, info.PendingUndelegationAmount,
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
		}
	}
	return gs.Params.Validate()
}*/
