package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// operators.go
	operators := make(map[string]struct{}, len(gs.Operators))
	for _, info := range gs.Operators {
		operatorAddress := info.EarningsAddr
		_, err := sdk.AccAddressFromBech32(operatorAddress)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"operator address %s is invalid: %s",
				operatorAddress, err,
			)
		}
		if _, ok := operators[operatorAddress]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate operator address: %s",
				operatorAddress,
			)
		}
		operators[operatorAddress] = struct{}{}
	}
	// consensus_keys.go
	operatorsByKeys := make(map[string]struct{}, len(gs.OperatorRecords))
	// keysByChainID stores chain id -> cons keys list. ensure that within a chain id, cons key
	// isn't repeated.
	keysByChainID := make(map[string](map[string]struct{}))
	for _, record := range gs.OperatorRecords {
		operatorAddress := record.OperatorAddress
		_, err := sdk.AccAddressFromBech32(operatorAddress)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"operator address %s is invalid: %s",
				operatorAddress, err,
			)
		}
		if _, ok := operators[operatorAddress]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"operator not registered %s",
				operatorAddress,
			)
		}
		if _, ok := operatorsByKeys[operatorAddress]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate operator address: %s",
				operatorAddress,
			)
		}
		operatorsByKeys[operatorAddress] = struct{}{}
		for _, subRecord := range record.Chains {
			consKeyString := subRecord.ConsensusKey
			if _, err := HexStringToPubKey(consKeyString); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"consensus key %s is invalid: %s",
					consKeyString, err,
				)
			}
			// validate chain id is not done, since it is not strictly enforced within Cosmos.
			// technically, it is possible to do so via ibcclienttypes.ParseChainID != 0.
			if _, ok := keysByChainID[subRecord.ChainID]; !ok {
				keysByChainID[subRecord.ChainID] = make(map[string]struct{})
			}
			if _, ok := keysByChainID[subRecord.ChainID][consKeyString]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate consensus key %s for chain %s",
					consKeyString, subRecord.ChainID,
				)
			}
			keysByChainID[subRecord.ChainID][consKeyString] = struct{}{}
		}
	}
	// it may be possible for an operator to opt into an AVS which does not have a consensus
	// key requirement, so this check could be removed if we set up the Export case. but i
	// think it is better to keep it for now.
	if len(operators) != len(operatorsByKeys) {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"operator addresses in operators and operator records do not match",
		)
	}
	// state_update.go
	// we do not know the length of this map, so we use an approximation.
	// it will auto expand anyway.
	operatorsByStakers := make(map[string]struct{}, len(operators))
	assetsByStakers := make(map[string](map[string]struct{}), len(gs.StakerRecords))
	for _, level1 := range gs.StakerRecords {
		stakerID := level1.StakerID
		if _, _, err := assetstypes.ParseID(stakerID); err != nil {
			return errorsmod.Wrapf(ErrInvalidGenesisData, "stakerID invalid: %s", err)
		}
		if _, ok := assetsByStakers[stakerID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate stakerID: %s",
				stakerID,
			)
		}
		assetsByStakers[stakerID] = make(map[string]struct{}, len(level1.StakerDetails))
		for _, level2 := range level1.StakerDetails {
			assetID := level2.AssetID
			if _, _, err := assetstypes.ParseID(assetID); err != nil {
				return errorsmod.Wrapf(ErrInvalidGenesisData, "assetID invalid: %s", err)
			}
			if _, ok := assetsByStakers[stakerID][assetID]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate assetID: %s",
					assetID,
				)
			}
			assetsByStakers[stakerID][assetID] = struct{}{}
			for _, level3 := range level2.Details {
				operatorAddress := level3.OperatorAddress
				_, err := sdk.AccAddressFromBech32(operatorAddress)
				if err != nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"operator address %s is invalid: %s",
						operatorAddress, err,
					)
				}
				if _, ok := operators[operatorAddress]; !ok {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"operator not registered %s",
						operatorAddress,
					)
				}
				// a staker may delegate different assets to multiple
				// operators, so we do not check for duplicates here.
				operatorsByStakers[operatorAddress] = struct{}{}
				amount := level3.Amount
				if amount.IsNil() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"amount is nil for %s: %s: %s",
						stakerID, assetID, operatorAddress,
					)
				}
				if amount.IsNegative() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"amount is negative for %s: %s: %s",
						stakerID, assetID, operatorAddress,
					)
				}
				// we allow 0 amounts for completeness.
			}
		}
	}
	// it is possible that a few operators do not get delegations from stakers. that means
	// operatorsByStakers may be smaller than operators.
	// operatorsByStakers can never be larger than operators anyway, since we have checked
	// that operators are already registered.
	// it may also be prudent to validate the sorted (or not) nature of these items
	// but it is not critical for the functioning. it is only used for comparison
	// of the genesis state stored across all of the validators.
	return nil
}
