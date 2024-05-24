package oracle

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	// this line is used by starport scaffolding # 1

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/ExocoreNetwork/exocore/x/oracle/client/cli"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/cache"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface that defines the independent methods a Cosmos SDK module needs to implement.
type AppModuleBasic struct {
	cdc codec.BinaryCodec
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the name of the module as a string
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns a default GenesisState for the module, marshaled to json.RawMessage. The default GenesisState need to be defined by the module developer and is primarily used for testing
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root Tx command for the module. The subcommands of this root command are used by end-users to generate new transactions containing messages defined in the module
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for the module. The subcommands of this root command are used by end-users to generate new queries to the subset of the state defined by the module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper

	// used for simulation
	accountKeeper types.AccountKeeper

	// used for simulation
	bankKeeper types.BankKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the module. It should be incremented on each consensus-breaking change introduced by the module. To avoid wrong/empty versions, the initial version should be set to 1
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	cs := keeper.GetCaches()
	validatorUpdates := am.keeper.GetValidatorUpdates(ctx)
	forceSeal := false
	agc := keeper.GetAggregatorContext(ctx, am.keeper)

	logger := am.keeper.Logger(ctx)
	if len(validatorUpdates) > 0 {
		validatorList := make(map[string]*big.Int)
		for _, vu := range validatorUpdates {
			pubKey, _ := cryptocodec.FromTmProtoPublicKey(vu.PubKey)
			validator, _ := am.keeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pubKey))
			validatorList[validator.OperatorAddress] = big.NewInt(vu.Power)
		}
		// update validator set information in cache
		cs.AddCache(cache.ItemV(validatorList))
		validatorPowers := make(map[string]*big.Int)
		cs.GetCache(cache.ItemV(validatorPowers))
		// update validatorPowerList in aggregatorContext
		agc.SetValidatorPowers(validatorPowers)
		// TODO: seal all alive round since validatorSet changed here
		forceSeal = true
		logger.Info("validator set changed, force seal all active rounds")
	}

	// TODO: for v1 use mode==1, just check the failed feeders
	_, failed := agc.SealRound(ctx, forceSeal)
	// append new round with previous price for fail-seal token
	for _, tokenID := range failed {
		event := sdk.NewEvent(
			types.EventTypeCreatePrice,
			sdk.NewAttribute(types.AttributeKeyTokenID, strconv.FormatUint(tokenID, 10)),
			sdk.NewAttribute(types.AttributeKeyPriceUpdated, types.AttributeValuePriceUpdatedFail),
		)
		logInfo := fmt.Sprintf("add new round with previous price under fail aggregation, tokenID:%d", tokenID)
		if pTR, ok := am.keeper.GetPriceTRLatest(ctx, tokenID); ok {
			pTR.RoundID++
			am.keeper.AppendPriceTR(ctx, tokenID, pTR)
			logger.Info("add new round with previous price under fail aggregation", "tokenID", tokenID, "roundID", pTR.RoundID)
			logInfo += fmt.Sprintf(", roundID:%d, price:%s", pTR.RoundID, pTR.Price)
			event.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyRoundID, strconv.FormatUint(pTR.RoundID, 10)),
				sdk.NewAttribute(types.AttributeKeyFinalPrice, pTR.Price),
			)
		} else {
			nextRoundID := am.keeper.GetNextRoundID(ctx, tokenID)
			am.keeper.AppendPriceTR(ctx, tokenID, types.PriceTimeRound{
				RoundID: nextRoundID,
			})
			logInfo += fmt.Sprintf(", roundID:%d, price:-", nextRoundID)
			event.AppendAttributes(
				sdk.NewAttribute(types.AttributeKeyRoundID, strconv.FormatUint(nextRoundID, 10)),
				sdk.NewAttribute(types.AttributeKeyFinalPrice, "-"),
			)
		}
		logger.Info(logInfo)
		ctx.EventManager().EmitEvent(event)
	}
	// TODO: emit events for success sealed rounds(could ignore for v1)

	logger.Info("prepare for next oracle round of each tokenFeeder")
	agc.PrepareRound(ctx, 0)

	// TODO: update params happened during this block for cache, for the case: agc is recached from history and cache'params is set with the latest params by recache, but parmas changed during this block as well. or force agc to be GET first before cache be GET(later approch is better)
	cs.CommitCache(ctx, false, am.keeper)
	return []abci.ValidatorUpdate{}
}
