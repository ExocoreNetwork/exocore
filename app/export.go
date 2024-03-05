package app

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/simapp"

	"github.com/ExocoreNetwork/exocore/utils"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	claimstypes "github.com/evmos/evmos/v14/x/claims/types"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	inflationtypes "github.com/evmos/evmos/v14/x/inflation/types"

	"github.com/evmos/evmos/v14/encoding"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.Codec) simapp.GenesisState {
	encCfg := encoding.MakeConfig(ModuleBasics)
	defaultGenesis := ModuleBasics.DefaultGenesis(encCfg.Codec)

	// staking module
	stakingGenesis := stakingtypes.GenesisState{}
	rawGenesis := defaultGenesis[stakingtypes.ModuleName]
	cdc.MustUnmarshalJSON(rawGenesis, &stakingGenesis)
	stakingGenesis.Params.BondDenom = utils.BaseDenom
	defaultGenesis[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenesis)

	// crisis module
	crisisGenesis := crisistypes.GenesisState{}
	rawGenesis = defaultGenesis[crisistypes.ModuleName]
	cdc.MustUnmarshalJSON(rawGenesis, &crisisGenesis)
	crisisGenesis.ConstantFee.Denom = utils.BaseDenom
	defaultGenesis[crisistypes.ModuleName] = cdc.MustMarshalJSON(&crisisGenesis)

	// gov module
	govGenesis := govtypesv1.GenesisState{}
	rawGenesis = defaultGenesis[govtypes.ModuleName]
	cdc.MustUnmarshalJSON(rawGenesis, &govGenesis)
	govGenesis.Params.MinDeposit[0].Denom = utils.BaseDenom
	defaultGenesis[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenesis)

	// evm module
	evmGenesis := evmtypes.GenesisState{}
	rawGenesis = defaultGenesis[evmtypes.ModuleName]
	cdc.MustUnmarshalJSON(rawGenesis, &evmGenesis)
	evmGenesis.Params.EvmDenom = utils.BaseDenom
	defaultGenesis[evmtypes.ModuleName] = cdc.MustMarshalJSON(&evmGenesis)

	// inflation module
	inflationGenesis := inflationtypes.GenesisState{}
	rawGenesis = defaultGenesis[inflationtypes.ModuleName]
	cdc.MustUnmarshalJSON(rawGenesis, &inflationGenesis)
	inflationGenesis.Params.MintDenom = utils.BaseDenom
	defaultGenesis[inflationtypes.ModuleName] = cdc.MustMarshalJSON(&inflationGenesis)

	// claims module
	claimsGenesis := claimstypes.GenesisState{}
	rawGenesis = defaultGenesis[claimstypes.ModuleName]
	cdc.MustUnmarshalJSON(rawGenesis, &claimsGenesis)
	claimsGenesis.Params.ClaimsDenom = utils.BaseDenom
	defaultGenesis[claimstypes.ModuleName] = cdc.MustMarshalJSON(&claimsGenesis)

	return defaultGenesis
}

// ExportAppStateAndValidators exports the state of the application for a genesis
// file.
func (app *ExocoreApp) ExportAppStateAndValidators(
	forZeroHeight bool, jailAllowedAddrs []string, modulesToExport []string,
) (servertypes.ExportedApp, error) {
	// Creates context with current height and checks txs for ctx to be usable by start of next block
	ctx := app.NewContext(true, tmproto.Header{Height: app.LastBlockHeight()})

	// We export at last height + 1, because that's the height at which
	// Tendermint will start InitChain.
	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0

		if err := app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	genState := app.mm.ExportGenesisForModules(ctx, app.appCodec, modulesToExport)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, &app.StakingKeeper)
	if err != nil {
		return servertypes.ExportedApp{}, err
	}

	return servertypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, nil
}

// prepare for fresh start at zero height
// NOTE zero height genesis is a temporary feature which will be deprecated
//
//	in favor of export at a block height
func (app *ExocoreApp) prepForZeroHeightGenesis(ctx sdk.Context, jailAllowedAddrs []string) error {
	applyAllowedAddrs := false

	// check if there is a allowed address list
	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			return err
		}
		allowedAddrsMap[addr] = true
	}

	/* Just to be safe, assert the invariants on current state. */
	app.CrisisKeeper.AssertInvariants(ctx)

	/* Handle fee distribution state. */

	// withdraw all validator commission
	app.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		_, _ = app.DistrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	// withdraw all delegator rewards
	dels := app.StakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			return err
		}

		delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			return err
		}
		_, _ = app.DistrKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	}

	// clear validator slash events
	app.DistrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.DistrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	app.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
		feePool := app.DistrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.DistrKeeper.SetFeePool(ctx, feePool)

		err := app.DistrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		// this lets us stop in case there's an error
		return err != nil
	})

	// reinitialize all delegations
	for _, del := range dels {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			return err
		}
		delAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
		if err != nil {
			return err
		}
		err = app.DistrKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr)
		if err != nil {
			return err
		}
		err = app.DistrKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr)
		if err != nil {
			return err
		}
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	app.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingtypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.StakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	app.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingtypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[stakingtypes.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := app.StakingKeeper.GetValidator(ctx, addr)
		if !found {
			return fmt.Errorf("expected validator %s not found", addr)
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		app.StakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	if err := iter.Close(); err != nil {
		return err
	}

	if _, err := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx); err != nil {
		return err
	}

	/* Handle slashing state. */

	// reset start height on signing infos
	app.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashingtypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
	return nil
}
