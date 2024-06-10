package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NewGenesis creates a new genesis state with the provided parameters and
// data.
func NewGenesis(
	params Params,
	vals []GenesisValidator,
	expiries []EpochToOperatorAddrs,
	consAddrs []EpochToConsensusAddrs,
	recordKeys []EpochToUndelegationRecordKeys,
	power math.Int,
) *GenesisState {
	return &GenesisState{
		Params:                 params,
		ValSet:                 vals,
		EpochsOptOutExpiries:   expiries,
		EpochsConsensusAddrs:   consAddrs,
		UndelegationMaturities: recordKeys,
		LastTotalPower:         power,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesis(
		DefaultParams(),
		[]GenesisValidator{},
		[]EpochToOperatorAddrs{},
		[]EpochToConsensusAddrs{},
		[]EpochToUndelegationRecordKeys{},
		math.ZeroInt(),
	)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// #nosec G701 // ok on 64-bit systems.
	maxValidators := int(gs.Params.MaxValidators)
	if len(gs.ValSet) > maxValidators {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"too many validators %d",
			len(gs.ValSet),
		)
	}
	// do not complain about 0 validators, let Tendermint do that.
	vals := make(map[string]struct{}, len(gs.ValSet))
	totalPower := int64(0)
	for _, val := range gs.ValSet {
		// check for duplicates
		if _, ok := vals[val.PublicKey]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate public key %s", val.PublicKey,
			)
		}
		vals[val.PublicKey] = struct{}{}
		// HexStringToPubKey checks the size and returns a tmprotocrypto type.
		// and since its specific type (ed25519) is already set, it converts
		// easily to the sdk Key format as well.
		if _, err := operatortypes.HexStringToPubKey(
			val.PublicKey,
		); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid public key %s: %s",
				val.PublicKey, err,
			)
		}
		power := val.Power
		if power <= 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid power %d",
				power,
			)
		}
		totalPower += power
	}

	// we don't know the current epoch, since this is stateless validation.
	// to check epoochs aren't duplicated.
	epochs := make(map[int64]struct{}, len(gs.EpochsOptOutExpiries))
	// to check that there is no duplicate address - not by per epoch but overall.
	addrsMap := make(map[string]struct{})
	for _, obj := range gs.EpochsOptOutExpiries {
		epoch := obj.Epoch
		if _, ok := epochs[epoch]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate epoch %d", epoch,
			)
		}
		// the first epoch in the epochs module is 1. this epoch is first
		// incremented, and then AfterEpochEnd is called with a value of 2.
		// therefore, the first epoch in the dogfood module is 2. all expiries
		// must happen at the end of this epoch or any epoch thereafter.
		if epoch <= 1 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"epoch %d should be > 1", epoch,
			)
		}
		epochs[epoch] = struct{}{}
		addrs := obj.OperatorAccAddrs
		if len(addrs) == 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"empty operator addresses for epoch %d", epoch,
			)
		}
		for _, addr := range addrs {
			if _, err := sdk.AccAddressFromBech32(addr); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid operator address %s: %s",
					addr, err,
				)
			}
			if _, ok := addrsMap[addr]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate operator address %s", addr,
				)
			}
			addrsMap[addr] = struct{}{}
		}
	}

	epochs = make(map[int64]struct{}, len(gs.EpochsConsensusAddrs))
	addrsMap = make(map[string]struct{})
	for _, obj := range gs.EpochsConsensusAddrs {
		epoch := obj.Epoch
		if _, ok := epochs[epoch]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate epoch %d", epoch,
			)
		}
		epochs[epoch] = struct{}{}
		if epoch <= 1 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"epoch %d should be > 1", epoch,
			)
		}
		addrs := obj.ConsAddrs
		if len(addrs) == 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"empty consensus addresses for epoch %d", epoch,
			)
		}
		for _, addr := range addrs {
			if _, err := sdk.ConsAddressFromBech32(addr); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid consensus address %s: %s",
					addr, err,
				)
			}
			if _, ok := addrsMap[addr]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate consensus address %s", addr,
				)
			}
			addrsMap[addr] = struct{}{}
		}
	}

	epochs = make(map[int64]struct{}, len(gs.UndelegationMaturities))
	recordKeysMap := make(map[string]struct{})
	for _, obj := range gs.UndelegationMaturities {
		epoch := obj.Epoch
		if _, ok := epochs[epoch]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate epoch %d", epoch,
			)
		}
		if epoch <= 1 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"epoch %d should be > 1", epoch,
			)
		}
		epochs[epoch] = struct{}{}
		recordKeys := obj.UndelegationRecordKeys
		if len(recordKeys) == 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"empty record keys for epoch %d", epoch,
			)
		}
		for _, recordKey := range recordKeys {
			if _, ok := recordKeysMap[recordKey]; ok {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"duplicate record key %s", recordKey,
				)
			}
			if recordBytes, err := hexutil.Decode(recordKey); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid record key (non hex) %s: %s",
					recordKey, err,
				)
			} else if _, err := delegationtypes.ParseUndelegationRecordKey(recordBytes); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid record key (parse) %s: %s",
					recordKey, err,
				)
			}
			recordKeysMap[recordKey] = struct{}{}
		}
	}

	if gs.LastTotalPower.IsNil() {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"nil last total power",
		)
	}

	if !gs.LastTotalPower.IsPositive() {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"non-positive last total power %s",
			gs.LastTotalPower,
		)
	}

	if !gs.LastTotalPower.Equal(math.NewInt(totalPower)) {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"last total power mismatch %s, expected %d",
			gs.LastTotalPower, totalPower,
		)
	}

	return gs.Params.Validate()
}
