package types

import (
	"encoding/hex"

	errorsmod "cosmossdk.io/errors"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
)

// NewGenesis returns a new genesis state with the given inputs.
func NewGenesis(
	delegations []DelegationsByStaker,
	associations []StakerToOperator,
) *GenesisState {
	return &GenesisState{
		Delegations:  delegations,
		Associations: associations,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis(nil, nil)
}

func ValidateIDAndOperator(stakerID, assetID, operator string) error {
	// validate the operator address
	if _, err := sdk.AccAddressFromBech32(operator); err != nil {
		return xerrors.Errorf(
			"invalid operator address for operator %s", operator,
		)
	}
	_, stakerClientChainID, err := assetstypes.ValidateID(stakerID, true, false)
	if err != nil {
		return xerrors.Errorf(
			"invalid stakerID: %s",
			stakerID,
		)
	}
	_, assetClientChainID, err := assetstypes.ValidateID(assetID, true, false)
	if err != nil {
		return xerrors.Errorf(
			"invalid assetID: %s",
			assetID,
		)
	}
	if stakerClientChainID != assetClientChainID {
		return xerrors.Errorf(
			"the client chain layerZero IDs of the staker and asset are different, stakerID:%s, assetID:%s",
			stakerID, assetID)
	}
	return nil
}

func (gs GenesisState) ValidateDelegations() error {
	// TODO(mm): this can be a very big hash table and impact system performance.
	// This is likely to be the biggest one amongst the three, and the others
	// are garbage collected within the loop anyway. Maybe reordering the genesis
	// structure could potentially help with this.
	stakers := make(map[string]struct{}, len(gs.Delegations))
	for _, level1 := range gs.Delegations {
		stakerID := level1.StakerID
		// validate staker ID
		var stakerClientChainID uint64
		var err error
		if _, stakerClientChainID, err = assetstypes.ValidateID(
			stakerID, true, true,
		); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "invalid staker ID %s: %s", stakerID, err,
			)
		}
		// check for duplicate stakers
		if _, ok := stakers[stakerID]; ok {
			return errorsmod.Wrapf(ErrInvalidGenesisData, "duplicate staker ID %s", stakerID)
		}
		stakers[stakerID] = struct{}{}
		assets := make(map[string]struct{}, len(level1.Delegations))
		for _, level2 := range level1.Delegations {
			assetID := level2.AssetID
			// check for duplicate assets
			if _, ok := assets[assetID]; ok {
				return errorsmod.Wrapf(ErrInvalidGenesisData, "duplicate asset ID %s", assetID)
			}
			assets[assetID] = struct{}{}
			// validate asset ID
			var assetClientChainID uint64
			if _, assetClientChainID, err = assetstypes.ValidateID(
				assetID, true, true,
			); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData, "invalid asset ID %s: %s", assetID, err,
				)
			}
			if assetClientChainID != stakerClientChainID {
				// a staker from chain A is delegating an asset on chain B, which is not
				// something we support right now.
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"asset %s client chain ID %d does not match staker %s client chain ID %d",
					assetID, assetClientChainID, stakerID, stakerClientChainID,
				)
			}
			operators := make(map[string]struct{}, len(level2.PerOperatorAmounts))
			for _, level3 := range level2.PerOperatorAmounts {
				operator := level3.Key
				wrappedAmount := level3.Value
				// check supplied amount
				if wrappedAmount == nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData, "nil operator amount for operator %s", operator,
					)
				}
				amount := wrappedAmount.Amount
				if amount.IsNil() || amount.IsNegative() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid operator amount %s for operator %s", amount, operator,
					)
				}
				// check operator address
				if _, err := sdk.AccAddressFromBech32(operator); err != nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid operator address for operator %s", operator,
					)
				}
				// check for duplicate operators
				if _, ok := operators[operator]; ok {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"duplicate operator %s for asset %s", operator, assetID,
					)
				}
				operators[operator] = struct{}{}
			}
		}
	}
	// for associations, one stakerID can be associated only with one operator.
	// but one operator may have multiple stakerIDs associated with it.
	associatedStakerIDs := make(map[string]struct{}, len(gs.Associations))
	for _, association := range gs.Associations {
		// check operator address
		if _, err := sdk.AccAddressFromBech32(association.Operator); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid operator address for operator %s", association.Operator,
			)
		}
		// check staker address
		if _, _, err := assetstypes.ValidateID(
			association.StakerID, true, true,
		); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "invalid staker ID %s: %s", association.StakerID, err,
			)
		}
		// check for duplicate stakerIDs
		if _, ok := associatedStakerIDs[association.StakerID]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "duplicate staker ID %s", association.StakerID,
			)
		}
		associatedStakerIDs[association.StakerID] = struct{}{}
		// we don't check that this `association.stakerID` features in `gs.Delegations`,
		// because we allow the possibility of a staker without any delegations to be associated
		// with an operator.
	}
	return nil
}

func (gs GenesisState) ValidateDelegationStates() error {
	validationFunc := func(i int, info DelegationStates) error {
		keys, err := ParseStakerAssetIDAndOperator([]byte(info.Key))
		if err != nil {
			return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
		}

		err = ValidateIDAndOperator(keys.StakerID, keys.AssetID, keys.OperatorAddr)
		if err != nil {
			return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
		}

		// check that there is no nil value provided.
		if info.States.UndelegatableShare.IsNil() || info.States.WaitUndelegationAmount.IsNil() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"nil delegation state for %s: %+v",
				info.Key, info,
			)
		}

		// check for negative values.
		if info.States.UndelegatableShare.IsNegative() || info.States.WaitUndelegationAmount.IsNegative() {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"negative delegation state  for %s: %+v",
				info.Key, info,
			)
		}

		return nil
	}
	seenFieldValueFunc := func(info DelegationStates) (string, struct{}) {
		return info.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.DelegationStates, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

func (gs GenesisState) ValidateStakerList() error {
	validationFunc := func(i int, stakersByOperator StakersByOperator) error {
		// validate the key
		stringList, err := assetstypes.ParseJoinedStoreKey([]byte(stakersByOperator.Key), 2)
		if err != nil {
			return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
		}
		// validate the operator address
		if _, err := sdk.AccAddressFromBech32(stringList[0]); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid operator address for operator %s", stringList[0],
			)
		}
		// validate the assetID
		_, assetClientChainID, err := assetstypes.ValidateID(stringList[1], true, false)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid assetID: %s",
				stringList[1],
			)
		}
		// validate the staker list
		stakerValidationFunc := func(i int, stakerID string) error {
			_, stakerClientChainID, err := assetstypes.ValidateID(stakerID, true, false)
			if err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"invalid stakerID: %s",
					stakerID,
				)
			}
			if stakerClientChainID != assetClientChainID {
				return errorsmod.Wrapf(ErrInvalidGenesisData, "the client chain layerZero IDs of the staker and asset are different,key:%s stakerID:%s", stakersByOperator.Key, stakerID)
			}
			return nil
		}
		seenStakerFunc := func(stakerID string) (string, struct{}) {
			return stakerID, struct{}{}
		}
		_, err = assetstypes.CommonValidation(stakersByOperator.Stakers.Stakers, seenStakerFunc, stakerValidationFunc)
		if err != nil {
			return err
		}
		return nil
	}
	seenFieldValueFunc := func(info StakersByOperator) (string, struct{}) {
		return info.Key, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.StakersByOperator, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

func (gs GenesisState) ValidateUndelegations() error {
	validationFunc := func(i int, undelegaion UndelegationRecord) error {
		err := ValidateIDAndOperator(undelegaion.StakerID, undelegaion.AssetID, undelegaion.OperatorAddr)
		if err != nil {
			return errorsmod.Wrap(ErrInvalidGenesisData, err.Error())
		}

		bytes, err := hex.DecodeString(undelegaion.TxHash)
		if err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "TxHash isn't a hex string, TxHash: %s",
				undelegaion.TxHash,
			)
		}
		if len(bytes) != common.HashLength {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "invalid length of TxHash ,TxHash:%s length: %d, should:%d",
				undelegaion.TxHash, len(bytes), common.HashLength,
			)
		}
		if !undelegaion.IsPending {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "all undelegations should be pending, undelegation:%v",
				undelegaion,
			)
		}
		if undelegaion.CompleteBlockNumber < undelegaion.BlockNumber {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "the block number to complete shouldn't be less than the submitted , undelegation：%v",
				undelegaion,
			)
		}
		if undelegaion.ActualCompletedAmount.GT(undelegaion.Amount) {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "the completed amount shouldn't be greater than the submitted amount , undelegation：%v",
				undelegaion,
			)
		}
		return nil
	}
	seenFieldValueFunc := func(undelegaion UndelegationRecord) (string, struct{}) {
		return undelegaion.TxHash, struct{}{}
	}
	_, err := assetstypes.CommonValidation(gs.Undelegations, seenFieldValueFunc, validationFunc)
	if err != nil {
		return err
	}
	return nil
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateDelegations()
	if err != nil {
		return err
	}
	err = gs.ValidateDelegationStates()
	if err != nil {
		return err
	}
	err = gs.ValidateStakerList()
	if err != nil {
		return err
	}
	err = gs.ValidateUndelegations()
	if err != nil {
		return err
	}
	return nil
}
