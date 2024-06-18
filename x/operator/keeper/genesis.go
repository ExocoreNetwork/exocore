package keeper

import (
	errorsmod "cosmossdk.io/errors"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidationBeforeInit
// some validations can't be performed in the `GenesisState.Validate` because the `isGeneralInit`
// flag is set in the keeper, so perform the validation here.
func (k Keeper) ValidationBeforeInit(_ sdk.Context, state types.GenesisState) error {
	if !k.isGeneralInit {
		if len(state.OptStates) != 0 {
			return errorsmod.Wrap(
				delegationtype.ErrInvalidGenesisData,
				"the opted states should be null when initializing from the bootStrap contract",
			)
		}
		if len(state.VotingPowers) != 0 {
			return errorsmod.Wrap(
				delegationtype.ErrInvalidGenesisData,
				"the voting powers should be null when initializing from the bootStrap contract",
			)
		}
		if len(state.SlashStates) != 0 {
			return errorsmod.Wrap(
				delegationtype.ErrInvalidGenesisData,
				"the slashing states should be null when initializing from the bootStrap contract",
			)
		}
		if len(state.PreConsKeys) != 0 {
			return errorsmod.Wrap(
				delegationtype.ErrInvalidGenesisData,
				"the previous consensus key should be null when initializing from the bootStrap contract",
			)
		}
		if len(state.OperatorKeyRemovals) != 0 {
			return errorsmod.Wrap(
				delegationtype.ErrInvalidGenesisData,
				"the operator key removals should be null when initializing from the bootStrap contract",
			)
		}
	}
	return nil
}

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	err := k.ValidationBeforeInit(ctx, state)
	if err != nil {
		panic(err)
	}
	for i := range state.Operators {
		op := state.Operators[i] // avoid implicit memory aliasing
		if !k.isGeneralInit {
			if op.OperatorInfo.EarningsAddr == "" {
				op.OperatorInfo.EarningsAddr = op.OperatorAddress
			}
		} else {
			if op.OperatorInfo.EarningsAddr == "" {
				panic(errorsmod.Wrap(delegationtype.ErrInvalidGenesisData, "earning addr is empty when init genesis from the general exporting genesis file"))
			}
		}
		if err := k.SetOperatorInfo(ctx, op.OperatorAddress, &op.OperatorInfo); err != nil {
			panic(err)
		}
	}
	for _, record := range state.OperatorRecords {
		addr := record.OperatorAddress
		// #nosec G703 // already validated
		operatorAddr, _ := sdk.AccAddressFromBech32(addr)

		for _, detail := range record.Chains {
			chainID := detail.ChainID
			// validate that the chain exists
			// TODO: move this check to the avs keeper when implemented.
			if chainID != ctx.ChainID() {
				panic("unknown chain id")
			}
			// #nosec G703 // already validated
			key, _ := types.HexStringToPubKey(detail.ConsensusKey)
			if k.isGeneralInit {
				// convert to bytes
				bz := k.cdc.MustMarshal(key)
				// convert to address for reverse lookup
				consAddr, err := types.TMCryptoPublicKeyToConsAddr(key)
				if err != nil {
					panic(errorsmod.Wrap(
						err, "SetOperatorConsKeyForChainID: cannot convert pub key to consensus address",
					))
				}
				k.setOperatorConsKeyForChainIDUnchecked(ctx, operatorAddr, consAddr, chainID, bz)
			} else {
				// opt into the specified chain (TODO: avs address format)
				if err := k.OptIn(ctx, operatorAddr, chainID); err != nil {
					panic(err)
				}
				// then set pub key
				if err := k.setOperatorConsKeyForChainID(
					ctx, operatorAddr, chainID, key, true,
				); err != nil {
					panic(err)
				}
			}
		}
	}
	// init the state from the general exporting genesis file
	err = k.SetAllOptedInfo(ctx, state.OptStates)
	if err != nil {
		panic(err)
	}
	err = k.SetAllUSDValues(ctx, state.VotingPowers)
	if err != nil {
		panic(err)
	}
	err = k.SetAllSlashInfo(ctx, state.SlashStates)
	if err != nil {
		panic(err)
	}
	err = k.SetAllPrevConsKeys(ctx, state.PreConsKeys)
	if err != nil {
		panic(err)
	}
	err = k.SetAllOperatorKeyRemovals(ctx, state.OperatorKeyRemovals)
	if err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	res := types.GenesisState{}
	res.Operators = k.AllOperators(ctx)

	operatorRecords, err := k.GetAllOperatorConsKeyRecords(ctx)
	if err != nil {
		panic(err)
	}
	res.OperatorRecords = operatorRecords

	optedInfos, err := k.GetAllOptedInfo(ctx)
	if err != nil {
		panic(err)
	}
	res.OptStates = optedInfos

	votingPowers, err := k.GetAllUSDValues(ctx)
	if err != nil {
		panic(err)
	}
	res.VotingPowers = votingPowers

	slashingInfos, err := k.GetAllSlashInfo(ctx)
	if err != nil {
		panic(err)
	}
	res.SlashStates = slashingInfos

	prevConsKeys, err := k.GetAllPrevConsKeys(ctx)
	if err != nil {
		panic(err)
	}
	res.PreConsKeys = prevConsKeys

	operatorKeyRemovals, err := k.GetAllOperatorKeyRemovals(ctx)
	if err != nil {
		panic(err)
	}
	res.OperatorKeyRemovals = operatorKeyRemovals

	return &res
}
