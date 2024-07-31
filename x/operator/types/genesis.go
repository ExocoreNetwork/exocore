package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func NewGenesisState(
	operators []OperatorInfo,
) *GenesisState {
	return &GenesisState{
		Operators: operators,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(nil)
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
	return nil
}
