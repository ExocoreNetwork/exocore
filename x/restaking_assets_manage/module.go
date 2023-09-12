// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package restaking_assets_manage

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/exocore/x/deposit/keeper"
	"github.com/exocore/x/deposit/types"
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
	return types.ModuleName
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	//TODO implement me
	panic("implement me")
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
	keeper keeper.Keeper
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

// IReStakingAssetsManage interface provided by restaking_assets_manage
/*
	Eigenlayer:
	@notice Mapping: staker => Strategy => number of shares which they currently hold
    mapping(address => mapping(IStrategy => uint256)) public stakerStrategyShares;
	@notice Mapping: staker => array of strategies in which they have nonzero shares
    mapping(address => IStrategy[]) public stakerStrategyList;


	exoCore stored info:

	//stored info in restaking_assets_manage module
	//used to record supported client chain and reStaking token info
	chainIndex->ChainInfo
	tokenIndex->tokenInfo
	chainList ?
	tokenList ?

	//record restaker reStaking info
	restaker->mapping(tokenIndex->amount)
	restaker->ReStakingTokenList ?
	restakerList?

	//record operator reStaking info
	operator->mapping(tokenIndex->amount)
	operator->ReStakingTokenList ?
	operator->mapping(tokenIndex->middleWareAddress) ?


	//stored info in delegation module
	//record the operator info which restaker delegate to
	restaker->mapping(operator->mapping(tokenIndex->amount))
	restaker->operatorList
	operator->operatorInfo

	//stored info in middleWare module
	middleWareAddr->middleWareInfo
	middleWareAddr->OptedInOperatorInfo
*/
type IReStakingAssetsManage interface {
	SetClientChainInfo(info *types2.ClientChainInfo) (exoCoreChainIndex uint64, err error)
	GetClientChainInfoByIndex(exoCoreChainIndex uint64) (info types2.ClientChainInfo, err error)
	GetAllClientChainInfo() (infos map[uint64]types2.ClientChainInfo, err error)

	SetReStakingAssetInfo(info *types2.ReStakingAssetInfo) (exoCoreAssetIndex uint64, err error)
	GetReStakingAssetInfo(assetId string) (info types2.ReStakingAssetInfo, err error)
	GetReStakingAssetsInfoList() (list []types2.ReStakingAssetInfo, err error)

	GetReStakerAssetInfos(reStakerAddr string) (assetsInfo map[uint64]math.Uint, err error)
	GetReStakerSpecifiedAssetAmount(reStakerAddr string, assetId string) (amount math.Uint, err error)
}
