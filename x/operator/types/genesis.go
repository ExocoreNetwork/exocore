package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func NewGenesisState(
	operators []OperatorInfo,
	records []OperatorConsKeyRecord,
) *GenesisState {
	return &GenesisState{
		Operators:       operators,
		OperatorRecords: records,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState([]OperatorInfo{}, []OperatorConsKeyRecord{})
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// checks list:
	// - no duplicate addresses in `gs.Operators`.
	// - correct bech32 format for each address in `gs.Operators`
	// - no `chainID` duplicates for earnings addresses list in `gs.Operators`.
	operators := make(map[string]struct{}, len(gs.Operators))
	for _, op := range gs.Operators {
		address := op.EarningsAddr
		if _, found := operators[address]; found {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate operator address %s", address,
			)
		}
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid operator address %s: %s", address, err,
			)
		}
		operators[address] = struct{}{}
		if op.ClientChainEarningsAddr != nil {
			lzIDs := make(map[uint64]struct{}, len(op.ClientChainEarningsAddr.EarningInfoList))
			for _, info := range op.ClientChainEarningsAddr.EarningInfoList {
				lzID := info.LzClientChainID
				if _, found := lzIDs[lzID]; found {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"duplicate lz client chain id %d", lzID,
					)
				}
				lzIDs[lzID] = struct{}{}
				// TODO: when moving to support non-EVM chains, this check should be modified
				// to work based on the `lzID` or possibly removed.
				if !common.IsHexAddress(info.ClientChainEarningAddr) {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid client chain earning address %s", info.ClientChainEarningAddr,
					)
				}
			}
		}
		if op.Commission.CommissionRates.Rate.IsNil() ||
			op.Commission.CommissionRates.MaxRate.IsNil() ||
			op.Commission.CommissionRates.MaxChangeRate.IsNil() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"missing commission for operator %s", address,
			)
		}
		if err := op.Commission.Validate(); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid commission for operator %s: %s", address, err,
			)
		}
	}
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
	// rationale for the validations above:
	// 1. since this function should support chain restarts and upgrades, we cannot require
	//    the format of the earnings address be EVM only.
	// 2. since the operator module is not meant to handle dogfooding, we should not check
	//    whether an operator has keys defined for our chainID. this is left for the dogfood
	//    module.
	return nil
}
