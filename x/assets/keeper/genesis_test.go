package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	utiltx "github.com/evmos/evmos/v14/testutil/tx"
)

func (suite *StakingAssetsTestSuite) TestInitGenesisFromBootStrap() {
	params := types.DefaultParams()
	params.ExocoreLzAppAddress = "0x0000000000000000000000000000000000000001"
	// genesis data that is hardcoded for use in the tests
	ethClientChain := types.ClientChainInfo{
		Name:               "ethereum",
		MetaInfo:           "ethereum blockchain",
		ChainId:            1,
		FinalizationBlocks: 10,
		LayerZeroChainID:   101,
		AddressLength:      20,
	}
	// do not hardcode the address to avoid gitleaks complaining.
	tokenAddress := utiltx.GenerateAddress().String()
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
		AssetBasicInfo:     usdtClientChainAsset,
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
					WithdrawableAmount:  math.NewInt(100),
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
			unmalleate: func(gs *types.GenesisState) {
				gs.Tokens[0].StakingTotalAmount = math.NewInt(0)
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
		{
			name: "invalid genesis due to excess deposited amount for staker",
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
				gs.Deposits[0].Deposits[0].Info.TotalDepositAmount = stakingInfo.AssetBasicInfo.TotalSupply.Add(math.NewInt(1))
				gs.Deposits[0].Deposits[0].Info.WithdrawableAmount = stakingInfo.AssetBasicInfo.TotalSupply.Add(math.NewInt(1))
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Deposits[0].Deposits[0].Info.TotalDepositAmount = math.NewInt(100)
				gs.Deposits[0].Deposits[0].Info.WithdrawableAmount = math.NewInt(100)
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		if tc.malleate != nil {
			tc.malleate(tc.genState)
			// check that unmalleate is not nil
			suite.Require().NotNil(tc.unmalleate, tc.name)
		}
		// Add defer and recover to handle panic
		func() {
			defer func() {
				r := recover()
				if tc.expPass {
					suite.Require().Nil(r, tc.name)
				} else {
					suite.Require().NotNil(r, tc.name)
				}
			}()
			suite.App.AssetsKeeper.InitGenesis(suite.Ctx, tc.genState)
		}()
		if tc.unmalleate != nil {
			tc.unmalleate(tc.genState)
		}
	}
}
