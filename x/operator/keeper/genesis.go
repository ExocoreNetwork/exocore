package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	for i := range state.Operators {
		op := state.Operators[i] // avoid implicit memory aliasing
		if op.OperatorInfo.EarningsAddr == "" {
			op.OperatorInfo.EarningsAddr = op.OperatorAddress
		}
		if err := k.SetOperatorInfo(ctx, op.OperatorAddress, &op.OperatorInfo); err != nil {
			panic(errorsmod.Wrap(err, "failed to set operator info"))
		}
	}
	for _, record := range state.OperatorRecords {
		addr := record.OperatorAddress
		// #nosec G703 // already validated
		operatorAddr, _ := sdk.AccAddressFromBech32(addr)
		for _, detail := range record.Chains {
			wrappedKey := types.NewWrappedConsKeyFromHex(detail.ConsensusKey)
			bz := k.cdc.MustMarshal(wrappedKey.ToTmProtoKey())
			k.setOperatorConsKeyForChainIDUnchecked(ctx, operatorAddr, wrappedKey.ToConsAddr(), detail.ChainID, bz)
		}
	}
	// init the state from the general exporting genesis file
	err := k.SetAllOptedInfo(ctx, state.OptStates)
	if err != nil {
		panic(errorsmod.Wrap(err, "failed to set all opted info"))
	}
	err = k.SetAllOperatorUSDValues(ctx, state.OperatorUSDValues)
	if err != nil {
		panic(errorsmod.Wrap(err, "failed to set all operator USD values"))
	}
	err = k.SetAllAVSUSDValues(ctx, state.AVSUSDValues)
	if err != nil {
		panic(errorsmod.Wrap(err, "failed to set all AVS USD values"))
	}
	err = k.SetAllSlashStates(ctx, state.SlashStates)
	if err != nil {
		panic(errorsmod.Wrap(err, "failed to set all slash info"))
	}
	err = k.SetAllPrevConsKeys(ctx, state.PreConsKeys)
	if err != nil {
		panic(errorsmod.Wrap(err, "failed to set all previous consensus keys"))
	}
	err = k.SetAllOperatorKeyRemovals(ctx, state.OperatorKeyRemovals)
	if err != nil {
		panic(errorsmod.Wrap(err, "failed to set all key removals for operators"))
	}
	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	res := types.GenesisState{}
	res.Operators = k.AllOperators(ctx)

	operatorRecords, err := k.GetAllOperatorConsKeyRecords(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to get all consensus keys for operators").Error())
	}
	res.OperatorRecords = operatorRecords

	optedInfos, err := k.GetAllOptedInfo(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to get all opted info").Error())
	}
	res.OptStates = optedInfos

	allAVSUSDValues, err := k.GetAllAVSUSDValues(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to get all AVS USD values").Error())
	}
	res.AVSUSDValues = allAVSUSDValues

	allOperatorUSDValues, err := k.GetAllOperatorUSDValues(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to get all operator USD values").Error())
	}
	res.OperatorUSDValues = allOperatorUSDValues

	slashingInfos, err := k.GetAllSlashStates(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to set all slashing info").Error())
	}
	res.SlashStates = slashingInfos

	prevConsKeys, err := k.GetAllPrevConsKeys(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to get all previous consensus keys").Error())
	}
	res.PreConsKeys = prevConsKeys

	operatorKeyRemovals, err := k.GetAllOperatorKeyRemovals(ctx)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "failed to get all key removals for operators").Error())
	}
	res.OperatorKeyRemovals = operatorKeyRemovals

	return &res
}
