package native_token

import (
	"context"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/exocore/x/deposit/keeper"
	"github.com/exocore/x/native_token/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
)

const consensusVersion = 0

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModuleBasic struct{}

func (b AppModuleBasic) Name() string {
	return types.ModuleName
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	// TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {
	// TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	// TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	// TODO implement me
	panic("implement me")
}

type AppModule struct {
	AppModuleBasic
	keeper *keeper.Keeper
}

func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
	// TODO implement me
	panic("implement me")
}

func (am AppModule) RegisterStoreDecoder(registry sdk.StoreDecoderRegistry) {
	// TODO implement me
	panic("implement me")
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	// TODO implement me
	panic("implement me")
}

// IDeposit interface will be implemented by deposit keeper
type IDeposit interface {
	// PostTxProcessing automatically call PostTxProcessing to update deposit state after receiving deposit event tx from layerZero protocol
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// SetReStakerExoCoreAddr handle the SetReStakerExoCoreAddr txs from msg service
	SetReStakerExoCoreAddr(ctx context.Context, reStakerId string) (err error)
	GetReStakerExoCoreAddr(reStakerId string) (addr sdk.Address, err error)

	// Deposit internal func for PostTxProcessing
	Deposit(reStakerId string, assetsInfo map[string]math.Uint) error
}
