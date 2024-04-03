package types_test

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	utiltx "github.com/ExocoreNetwork/exocore/testutil/tx"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite
}

func (suite *GenesisTestSuite) SetupTest() {
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestValidateGenesis() {
	params := types.DefaultParams()
	params.ExocoreLzAppAddress = "0x0000000000000000000000000000000000000001"
	newGen := types.NewGenesis(
		params, []types.ClientChainInfo{},
		[]types.StakingAssetInfo{}, []types.DepositsByStaker{},
	)
	// genesis data that is hardcoded for use in the tests
	ethClientChain := types.ClientChainInfo{
		Name:               "ethereum",
		MetaInfo:           "ethereum blockchain",
		ChainId:            1,
		FinalizationBlocks: 10,
		LayerZeroChainID:   101,
		AddressLength:      20,
	}
	tokenAddress := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	usdtClientChainAsset := types.AssetInfo{
		Name:             "Tether USD",
		Symbol:           "USDT",
		Address:          tokenAddress,
		Decimals:         6,
		LayerZeroChainID: ethClientChain.LayerZeroChainID,
		MetaInfo:         "Tether USD token",
	}
	totalSupply, _ := sdk.NewIntFromString("40022689732746729")
	usdtClientChainAsset.TotalSupply = totalSupply
	stakingInfo := types.StakingAssetInfo{
		AssetBasicInfo:     &usdtClientChainAsset,
		StakingTotalAmount: math.NewInt(0),
	}
	// generated information
	ethAddress := utiltx.GenerateAddress()
	// csmAddress := sdk.AccAddress(ethAddress.Bytes())
	stakerID, assetID := types.GetStakeIDAndAssetIDFromStr(
		usdtClientChainAsset.LayerZeroChainID, ethAddress.String(), usdtClientChainAsset.Address,
	)
	genesisDeposit := types.DepositsByStaker{
		StakerID: stakerID,
		Deposits: []types.DepositByAsset{
			{
				AssetID: assetID,
				Info: types.StakerAssetInfo{
					TotalDepositAmount:  math.NewInt(100),
					WithdrawableAmount:  math.NewInt(0),
					WaitUnbondingAmount: math.NewInt(0),
				},
			},
		},
	}

	testCases := []struct {
		name       string
		genState   *types.GenesisState
		expPass    bool
		malleate   func(*types.GenesisState)
		unmalleate func(*types.GenesisState)
	}{
		{
			name:     "valid genesis constructor",
			genState: newGen,
			expPass:  true,
		},
		{
			name:     "default",
			genState: types.DefaultGenesis(),
			expPass:  true,
		},
		{
			name: "valid genesis created here",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
			},
			expPass: true,
		},
		{
			name: "invalid genesis due to duplicate client chain",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain, ethClientChain,
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to missing client chain",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to wrong token address format",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
			},
			malleate: func(gs *types.GenesisState) {
				gs.Tokens[0].AssetBasicInfo.Address = "fakeTokenAddress"
			},
			unmalleate: func(gs *types.GenesisState) {
				// gs.Tokens[0].AssetBasicInfo is a pointer, so we undo the change manually
				gs.Tokens[0].AssetBasicInfo.Address = tokenAddress
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to duplicate token",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo, stakingInfo,
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to non zero deposit",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Tokens[0].StakingTotalAmount = math.NewInt(1)
			},
		},
		{
			name: "invalid genesis due to upper case staker id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].StakerID = strings.ToUpper(gs.Deposits[0].StakerID)
			},
		},
		{
			name: "invalid genesis due to invalid staker id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].StakerID = "fakeStaker"
			},
		},
		{
			name: "invalid genesis due to staker id from unknown chain",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].StakerID = "fakeStaker_0x63"
			},
		},
		{
			name: "invalid genesis due to non hex staker id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].StakerID = "fakeNonHexStaker_0x65"
			},
		},
		{
			name: "invalid genesis due to duplicate staker id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit, genesisDeposit},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to unknown asset id for staker",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].AssetID = "fakeAssetID"
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].AssetID = assetID
			},
		},
		{
			name: "invalid genesis due to duplicate asset id for staker",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits = append(
					gs.Deposits[0].Deposits,
					genesisDeposit.Deposits[0],
				)
			},
		},
		{
			name: "invalid genesis due to nil values for staker",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].Info = types.StakerAssetInfo{}
			},
			unmalleate: func(gs *types.GenesisState) {
				genesisDeposit.Deposits[0].Info = types.StakerAssetInfo{
					TotalDepositAmount:  math.NewInt(100),
					WithdrawableAmount:  math.NewInt(0),
					WaitUnbondingAmount: math.NewInt(0),
				}
				gs.Deposits[0].Deposits[0].Info = genesisDeposit.Deposits[0].Info
			},
		},
		{
			name: "invalid genesis due to non zero unbonding amount for staker",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].Info.WaitUnbondingAmount = math.NewInt(1)
			},
			unmalleate: func(gs *types.GenesisState) {
				genesisDeposit.Deposits[0].Info = types.StakerAssetInfo{
					TotalDepositAmount:  math.NewInt(100),
					WithdrawableAmount:  math.NewInt(0),
					WaitUnbondingAmount: math.NewInt(0),
				}
				gs.Deposits[0].Deposits[0].Info = genesisDeposit.Deposits[0].Info
			},
		},
		{
			name: "invalid genesis due to negative amount for staker",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].Info.TotalDepositAmount = math.NewInt(-1)
			},
			unmalleate: func(gs *types.GenesisState) {
				genesisDeposit.Deposits[0].Info = types.StakerAssetInfo{
					TotalDepositAmount:  math.NewInt(100),
					WithdrawableAmount:  math.NewInt(0),
					WaitUnbondingAmount: math.NewInt(0),
				}
				gs.Deposits[0].Deposits[0].Info = genesisDeposit.Deposits[0].Info
			},
		},
		{
			name: "invalid genesis due to excess withdrawable amount for staker",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				Deposits: []types.DepositsByStaker{genesisDeposit},
			},
			expPass: false,
			malleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].Info.TotalDepositAmount = math.NewInt(1)
				gs.Deposits[0].Deposits[0].Info.WithdrawableAmount = math.NewInt(2)
			},
			unmalleate: func(gs *types.GenesisState) {
				genesisDeposit.Deposits[0].Info = types.StakerAssetInfo{
					TotalDepositAmount:  math.NewInt(100),
					WithdrawableAmount:  math.NewInt(0),
					WaitUnbondingAmount: math.NewInt(0),
				}
				gs.Deposits[0].Deposits[0].Info = genesisDeposit.Deposits[0].Info
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		if tc.malleate != nil {
			tc.malleate(tc.genState)
		}
		err := tc.genState.Validate()
		fmt.Println(tc.name, ":", err)
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
		if tc.unmalleate != nil {
			tc.unmalleate(tc.genState)
		}
	}
}
