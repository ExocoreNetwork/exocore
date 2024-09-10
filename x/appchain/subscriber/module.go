package subscriber

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/client/cli"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/keeper"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
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

// AppModuleBasic implements the AppModuleBasic interface that defines the independent methods a
// Cosmos SDK module needs to implement.
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

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used to marshal
// and unmarshal structs to/from []byte in order to persist them in the module's KVStore
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers a module's interface types and their concrete implementations as
// proto.Message
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns a default GenesisState for the module, marshaled to json.RawMessage.
// The default GenesisState need to be defined by the module developer and is primarily used for
// testing
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form
func (AppModuleBasic) ValidateGenesis(
	cdc codec.JSONCodec,
	_ client.TxEncodingConfig,
	bz json.RawMessage,
) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(
	clientCtx client.Context,
	mux *runtime.ServeMux,
) {
	if err := types.RegisterQueryHandlerClient(
		context.Background(), mux, types.NewQueryClient(clientCtx),
	); err != nil {
		// this panic is safe to do because it means an error in setting up the module.
		panic(err)
	}
}

// GetTxCmd returns the root Tx command for the module. The subcommands of this root command are
// used by end-users to generate new transactions containing messages defined in the module
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for the module. The subcommands of this root
// command are used by end-users to generate new queries to the subset of the state defined by
// the module
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that
// modules need to implement
type AppModule struct {
	AppModuleBasic
	// keeper of the module receives the cdc codec separately.
	keeper keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC
// queries
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its
// predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will
// be halted)
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	gs json.RawMessage,
) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	return am.keeper.InitGenesis(ctx, genState)
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the module. It should be
// incremented on each consensus-breaking change introduced by the module. To avoid wrong/empty
// versions, the initial version should be set to 1
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	channelId, found := am.keeper.GetCoordinatorChannel(ctx)
	if found && am.keeper.IsChannelClosed(ctx, channelId) {
		// we are now PoA
		am.keeper.Logger(ctx).
			Error("coordinator channel is closed, we are now PoA", "channelId", channelId)
	}

	// get height of the yet-to-be-made block
	height := ctx.BlockHeight()
	// this should either be the last known vscId
	// or the one set by the last processed vsc packet
	// since that processing applies to the next block
	vscID := am.keeper.GetValsetUpdateIDForHeight(ctx, height)
	am.keeper.SetValsetUpdateIDForHeight(ctx, height+1, vscID)
	am.keeper.Logger(ctx).Debug(
		"block height was mapped to vscID",
		"height", height, "vscID", vscID,
	)

	am.keeper.TrackHistoricalInfo(ctx)
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// send rewards to coordinator
	am.keeper.EndBlockSendRewards(ctx)

	// queue maturity packets to coordinator
	am.keeper.QueueVscMaturedPackets(ctx)
	// remember that slash packets are queued in the subscriber module
	// by the slashing and evidence modules when a slashing event is observed by them

	// broadcast queued packets to coordinator
	am.keeper.SendPackets(ctx)

	// apply validator changes and then delete them
	data := am.keeper.GetPendingChanges(ctx)
	if len(data.ValidatorUpdates) == 0 {
		return []abci.ValidatorUpdate{}
	}
	updates := am.keeper.ApplyValidatorChanges(ctx, data.ValidatorUpdates)
	am.keeper.DeletePendingChanges(ctx)
	if len(updates) > 0 {
		am.keeper.Logger(ctx).Info("applying validator updates", "updates", updates)
	}
	return updates
}
