// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package restaking_assets_manage

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/exocore/x/deposit/keeper"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
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

type AppModuleBasic struct {
}

func (b AppModuleBasic) Name() string {
	return types2.ModuleName
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	//TODO implement me
	panic("implement me")
}

// DefaultGenesis returns default genesis state as raw bytes for the auth
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the auth module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(data)
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	//TODO implement me
	panic("implement me")
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	//TODO implement me
	panic("implement me")
}

type AppModule struct {
	AppModuleBasic
	keeper *keeper.Keeper
}

func (am AppModule) GenerateGenesisState(input *module.SimulationState) {
	//TODO implement me
	panic("implement me")
}

func (am AppModule) RegisterStoreDecoder(registry sdk.StoreDecoderRegistry) {
	//TODO implement me
	panic("implement me")
}

func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	//TODO implement me
	panic("implement me")
}

type ReStakingChainInfo struct {
	ChainName string
	ChainId   uint64
}

type ReStakingTokenInfo struct {
	TokenAddress string
	TokenName    string
}

// IReStakingAssetsManage interface will be implemented by restaking_assets_manage keeper
type IReStakingAssetsManage interface {
	SetClientChainInfo(info *types2.ClientChainInfo) (exoCoreChainIndex uint64, err error)
	GetClientChainInfoByIndex(exoCoreChainIndex uint64) (info types2.ClientChainInfo, err error)
	GetAllClientChainInfo() (infos map[uint64]types2.ClientChainInfo, err error)

	SetReStakingAssetInfo(info *types2.ReStakingAssetInfo) (exoCoreAssetIndex uint64, err error)
	GetReStakingAssetInfo(assetId string) (info types2.ReStakingAssetInfo, err error)
	GetAllReStakingAssetsInfo() (allAssets map[string]types2.ReStakingAssetInfo, err error)

	GetReStakerAssetInfos(reStakerId string) (assetsInfo map[string]math.Uint, err error)
	GetReStakerSpecifiedAssetAmount(reStakerId string, assetId string) (amount math.Uint, err error)
	IncreaseReStakerAssetsAmount(reStakerId string, assetsAddAmount map[string]math.Uint) (err error)
	DecreaseReStakerAssetsAmount(reStakerId string, assetsSubAmount map[string]math.Uint) (err error)

	GetOperatorAssetInfos(operatorAddr sdk.Address) (assetsInfo map[string]math.Uint, err error)
	GetOperatorSpecifiedAssetAmount(operatorAddr sdk.Address, assetId string) (amount math.Uint, err error)
	IncreaseOperatorAssetsAmount(operatorAddr sdk.Address, assetsAddAmount map[string]math.Uint) (err error)
	DecreaseOperatorAssetsAmount(operatorAddr sdk.Address, assetsSubAmount map[string]math.Uint) (err error)
	GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error)

	// GetAllOperatorAssetOptedInMiddleWare can also be implemented in operator optedIn module
	GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error)
	SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error)
}
