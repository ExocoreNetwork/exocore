package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func NewGenesisState(
	operators []OperatorDetail,
	records []OperatorConsKeyRecord,
) *GenesisState {
	return &GenesisState{
		Operators:       operators,
		OperatorRecords: records,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState([]OperatorDetail{}, []OperatorConsKeyRecord{})
}

// ValidateOperators rationale for the validation:
//  1. since this function should support chain restarts and upgrades, we cannot require
//     the format of the earnings address be EVM only.
func (gs GenesisState) ValidateOperators() (map[string]struct{}, error) {
	// checks list:
	// - no duplicate addresses in `gs.Operators`.
	// - correct bech32 format for each address in `gs.Operators`
	// - no `chainID` duplicates for earnings addresses list in `gs.Operators`.
	operators := make(map[string]struct{}, len(gs.Operators))
	for _, op := range gs.Operators {
		address := op.OperatorAddress
		if _, found := operators[address]; found {
			return nil, errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate operator address %s", address,
			)
		}
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return nil, errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid operator address %s: %s", address, err,
			)
		}
		if op.OperatorInfo.EarningsAddr != "" {
			_, err := sdk.AccAddressFromBech32(op.OperatorInfo.EarningsAddr)
			if err != nil {
				return nil, errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid operator earning address %s: %s", op.OperatorInfo.EarningsAddr, err,
				)
			}
		}
		operators[address] = struct{}{}
		if op.OperatorInfo.ClientChainEarningsAddr != nil {
			lzIDs := make(map[uint64]struct{}, len(op.OperatorInfo.ClientChainEarningsAddr.EarningInfoList))
			for _, info := range op.OperatorInfo.ClientChainEarningsAddr.EarningInfoList {
				lzID := info.LzClientChainID
				if _, found := lzIDs[lzID]; found {
					return nil, errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"duplicate lz client chain id %d", lzID,
					)
				}
				lzIDs[lzID] = struct{}{}
				// TODO: when moving to support non-EVM chains, this check should be modified
				// to work based on the `lzID` or possibly removed.
				if !common.IsHexAddress(info.ClientChainEarningAddr) {
					return nil, errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid client chain earning address %s", info.ClientChainEarningAddr,
					)
				}
			}
		}
		if op.OperatorInfo.Commission.CommissionRates.Rate.IsNil() ||
			op.OperatorInfo.Commission.CommissionRates.MaxRate.IsNil() ||
			op.OperatorInfo.Commission.CommissionRates.MaxChangeRate.IsNil() {
			return nil, errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"missing commission for operator %s", address,
			)
		}
		if err := op.OperatorInfo.Commission.Validate(); err != nil {
			return nil, errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid commission for operator %s: %s", address, err,
			)
		}
	}
	return operators, nil
}

// ValidateOperatorConsKeyRecords rationale for the validation:
//  2. since the operator module is not meant to handle dogfooding, we should not check
//     whether an operator has keys defined for our chainID. this is left for the dogfood
//     module.
func (gs GenesisState) ValidateOperatorConsKeyRecords(operators map[string]struct{}) error {
	// - correct bech32 format for each address in `gs.OperatorRecords`.
	// - no duplicate addresses in `gs.OperatorRecords`.
	// - no operator that is in `gs.OperatorRecords` but not in `gs.Operators`.
	// - validity of consensus key format for each entry in `gs.OperatorRecords`.
	// - within each chainID, no duplicate consensus keys.
	operatorRecords := make(map[string]struct{}, len(gs.OperatorRecords))
	keysByChainID := make(map[string]map[string]struct{})
	for _, record := range gs.OperatorRecords {
		addr := record.OperatorAddress
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid operator address %s: %s", record.OperatorAddress, err,
			)
		}
		if _, found := operatorRecords[addr]; found {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate operator record for operator %s", addr,
			)
		}
		operatorRecords[addr] = struct{}{}
		if _, opFound := operators[addr]; !opFound {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"operator record for un-registered operator %s", addr,
			)
		}
		for _, chain := range record.Chains {
			consensusKeyString := chain.ConsensusKey
			chainID := chain.ChainID
			// Cosmos does not describe a specific `chainID` format, so can't validate it.
			if _, found := keysByChainID[chainID]; !found {
				keysByChainID[chainID] = make(map[string]struct{})
			}
			if _, err := HexStringToPubKey(consensusKeyString); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid consensus key for operator %s: %s", addr, err,
				)
			}
			// within a chain id, there should not be duplicate consensus keys
			if _, found := keysByChainID[chainID][consensusKeyString]; found {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate consensus key for operator %s on chain %s", addr, chainID,
				)
			}
			keysByChainID[chainID][consensusKeyString] = struct{}{}
		}
	}
	return nil
}

func (gs GenesisState) ValidateOptedStates(operators map[string]struct{}) (map[string]struct{}, error) {
	avs := make(map[string]struct{})
	validationFunc := func(i int, state OptedState) error {
		stringList, err := assetstypes.ParseJoinedStoreKey([]byte(state.Key), 3)
		if err != nil {
			return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
		}
		operator, avsAddr := stringList[0], stringList[1]
		// check that the operator is registered
		if _, ok := operators[operator]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown operator address for the opted state, %+v",
				state,
			)
		}
		if state.OptInfo.OptedOutHeight < state.OptInfo.OptedInHeight {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"the opted-out height should be greater than the opted-in height, %+v",
				state,
			)
		}
		// todo: check the AVS address when the format is finalized.
		avs[avsAddr] = struct{}{}
		return nil
	}
	seenFieldValueFunc := func(state OptedState) (string, struct{}) {
		return state.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.OptStates, seenFieldValueFunc, validationFunc)
	if err != nil {
		return nil, err
	}
	return avs, nil
}

func (gs GenesisState) ValidateVotingPowers(operators, avs map[string]struct{}) error {
	var avsAddr string
	var avsVP DecValueField
	validationFunc := func(i int, vp VotingPower) error {
		if vp.Value.Amount.IsNil() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil voting power for :%+v",
				vp,
			)
		}
		if vp.Value.Amount.IsNegative() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"negative voting power for :%+v",
				vp,
			)
		}
		// the firs field should be the voting power of the AVS, and the operators that
		// opted into this AVS should follow it. because the iterator
		// iterates all the keys in ascending order when exporting the genesis file
		if !assetstypes.IsJoinedStoreKey(vp.Key) {
			avsAddr = vp.Key
			avsVP = vp.Value
			// check whether the AVS is in the opted states.
			// This check might be removed if the opted-in states are deleted when
			// the operator opts out of the AVS.
			if _, ok := avs[avsAddr]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown AVS address for the voting power, %+v",
					vp,
				)
			}
		} else {
			stringList, err := assetstypes.ParseJoinedStoreKey([]byte(vp.Key), 2)
			if err != nil {
				return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
			}
			operator, avsAddress := stringList[0], stringList[1]
			// check that the operator is registered
			if _, ok := operators[operator]; !ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"unknown operator address for the voting power, %+v",
					vp,
				)
			}
			if avsAddress != avsAddr {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"the operator should follows the opted-in AVS, AVS: %s, vp: %+v",
					avsAddr, vp,
				)
			}
			if vp.Value.Amount.GT(avsVP.Amount) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"the operator's voting power shouldn't be greater than the voting power of the AVS, avsVP: %s, vp: %+v",
					avsVP.Amount.String(), vp,
				)
			}
		}
		return nil
	}
	seenFieldValueFunc := func(vp VotingPower) (string, struct{}) {
		return vp.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.VotingPowers, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

func (gs GenesisState) ValidateSlashStates(operators, avs map[string]struct{}) error {
	validationFunc := func(i int, slash OperatorSlashState) error {
		stringList, err := assetstypes.ParseJoinedStoreKey([]byte(slash.Key), 3)
		if err != nil {
			return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
		}
		operator, avsAddr := stringList[0], stringList[1]
		// check that the operator is registered
		if _, ok := operators[operator]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown operator address for the slashing state, %+v",
				slash,
			)
		}
		// check whether the AVS is in the opted states.
		// This check might be removed if the opted-in states are deleted when
		// the operator opts out of the AVS.
		if _, ok := avs[avsAddr]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown AVS address for the slashing state, %+v",
				slash,
			)
		}
		if slash.Info.EventHeight > slash.Info.SubmittedHeight {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"the submitted height shouldn't be greater than the event height for a slashing record, %+v",
				slash,
			)
		}
		if slash.Info.SlashProportion.IsNil() || slash.Info.SlashProportion.LTE(sdkmath.LegacyNewDec(0)) {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid slash proportion, it's nil, zero, or negative: %+v",
				slash,
			)
		}

		// validate the slashing execution information
		// the actual executed proportion and value might be zero because of the rounding in an extreme case
		if slash.Info.ExecutionInfo.SlashProportion.IsNil() || slash.Info.ExecutionInfo.SlashProportion.IsNegative() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid slashing execution proportion, it's nil, or negative: %+v",
				slash,
			)
		}
		if slash.Info.ExecutionInfo.SlashValue.IsNil() || slash.Info.ExecutionInfo.SlashValue.IsNegative() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid slashing execution value, it's nil, or negative: %+v",
				slash,
			)
		}
		// validate the slashing record regarding undelegation
		SlashFromUndelegationVal := func(i int, slashFromUndelegation SlashFromUndelegation) error {
			if slashFromUndelegation.Amount.IsNil() || slashFromUndelegation.Amount.LTE(sdkmath.NewInt(0)) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid slashing amount from the undelegation, it's nil, zero, or negative: %+v",
					slash,
				)
			}
			return nil
		}
		seenFieldValueFunc := func(slashFromUndelegation SlashFromUndelegation) (string, struct{}) {
			key := assetstypes.GetJoinedStoreKey(slashFromUndelegation.StakerID, slashFromUndelegation.AssetID)
			return string(key), struct{}{}
		}
		_, err = assetstypes.CommonValidation(slash.Info.ExecutionInfo.SlashUndelegations, seenFieldValueFunc, SlashFromUndelegationVal)
		if err != nil {
			return err
		}
		// validate the slashing record regarding assets pool
		SlashFromAssetsPoolVal := func(i int, slashFromAssetsPool SlashFromAssetsPool) error {
			if slashFromAssetsPool.Amount.IsNil() || slashFromAssetsPool.Amount.LTE(sdkmath.NewInt(0)) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid slashing amount from the assets pool, it's nil, zero, or negative: %+v",
					slash,
				)
			}
			return nil
		}
		SlashFromAssetsPooLSeenFunc := func(slashFromAssetsPool SlashFromAssetsPool) (string, struct{}) {
			return slashFromAssetsPool.AssetID, struct{}{}
		}
		_, err = assetstypes.CommonValidation(slash.Info.ExecutionInfo.SlashAssetsPool, SlashFromAssetsPooLSeenFunc, SlashFromAssetsPoolVal)
		if err != nil {
			return err
		}
		return nil
	}
	seenFieldValueFunc := func(slash OperatorSlashState) (string, struct{}) {
		return slash.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.SlashStates, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

func (gs GenesisState) ValidatePrevConsKeys(operators map[string]struct{}) error {
	validationFunc := func(i int, prevConsKey PrevConsKey) error {
		keyBytes, err := hexutil.Decode(prevConsKey.Key)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"ValidatePrevConsKeys can't decode the key with hexutil.Decode, %+v",
				prevConsKey,
			)
		}
		_, operatorAddr, err := ParsePrevConsKey(keyBytes)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"ValidatePrevConsKeys can't parse the key, %+v",
				prevConsKey,
			)
		}
		// check that the operator is registered
		if _, ok := operators[operatorAddr.String()]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown operator address for the previous consensus key, %+v",
				prevConsKey,
			)
		}
		if _, err := HexStringToPubKey(prevConsKey.ConsensusKey); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid previous consensus key for operator %v: %s", prevConsKey, err,
			)
		}
		// todo: not sure if the duplication of previous consensus keys needs to be checked
		return nil
	}
	seenFieldValueFunc := func(prevConsKey PrevConsKey) (string, struct{}) {
		return prevConsKey.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.PreConsKeys, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

func (gs GenesisState) ValidateOperatorKeyRemovals(operators map[string]struct{}) error {
	validationFunc := func(i int, operatorKeyRemoval OperatorKeyRemoval) error {
		keyBytes, err := hexutil.Decode(operatorKeyRemoval.Key)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"ValidateOperatorKeyRemovals can't decode the key with hexutil.Decode, %+v",
				operatorKeyRemoval,
			)
		}
		operatorAddr, _, err := ParseKeyForOperatorKeyRemoval(keyBytes)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"ValidateOperatorKeyRemovals can't parse the key, %+v",
				operatorKeyRemoval,
			)
		}
		// check that the operator is registered
		if _, ok := operators[operatorAddr.String()]; !ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"unknown operator address for the operator key removal, %+v",
				operatorKeyRemoval,
			)
		}
		return nil
	}
	seenFieldValueFunc := func(operatorKeyRemoval OperatorKeyRemoval) (string, struct{}) {
		return operatorKeyRemoval.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.OperatorKeyRemovals, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	operators, err := gs.ValidateOperators()
	if err != nil {
		return err
	}
	err = gs.ValidateOperatorConsKeyRecords(operators)
	if err != nil {
		return err
	}
	avsMap, err := gs.ValidateOptedStates(operators)
	if err != nil {
		return err
	}
	err = gs.ValidateVotingPowers(operators, avsMap)
	if err != nil {
		return err
	}
	err = gs.ValidateSlashStates(operators, avsMap)
	if err != nil {
		return err
	}
	err = gs.ValidatePrevConsKeys(operators)
	if err != nil {
		return err
	}
	err = gs.ValidateOperatorKeyRemovals(operators)
	if err != nil {
		return err
	}
	return nil
}
