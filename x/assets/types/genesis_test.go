package types_test

import (
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
		params, []types.ClientChainInfo{}, []types.StakingAssetInfo{},
		[]types.OpAssetIDAndInfos{}, []types.StAssetIDAndInfos{},
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
	usdtClientChainAsset := types.AssetInfo{
		Name:             "Tether USD",
		Symbol:           "USDT",
		Address:          "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Decimals:         6,
		LayerZeroChainID: ethClientChain.LayerZeroChainID,
		MetaInfo:         "Tether USD token",
	}
	totalSupply, _ := sdk.NewIntFromString("40022689732746729")
	usdtClientChainAsset.TotalSupply = totalSupply
	stakingInfo := types.StakingAssetInfo{
		AssetBasicInfo:     &usdtClientChainAsset,
		StakingTotalAmount: math.NewInt(1),
	}
	// generated information
	accAddress := sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	stakerID, assetID := types.GetStakeIDAndAssetIDFromStr(
		usdtClientChainAsset.LayerZeroChainID, "0x123456789", usdtClientChainAsset.Address,
	)

	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
		malleate func(*types.GenesisState)
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
			name: "invalid genesis due to negative deposit",
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
				gs.Tokens[0].StakingTotalAmount = math.NewInt(-1)
			},
		},
		{
			name: "invalid genesis due to invalid bech32 operator address",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: "0x",
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to duplicate operator address",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
					},
					{
						OperatorAddress: accAddress.String(),
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to nil operator values",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: "1",
								Info:    &types.OperatorAssetInfo{},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to non-zero operator values",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: "1",
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(1),
									WaitUnbondingAmount:                math.NewInt(1),
									OperatorUnbondingAmount:            math.NewInt(1),
									OperatorUnbondableAmountAfterSlash: math.NewInt(1),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to duplicate asset id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(2),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to negative amount",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(-1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
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
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: "0x",
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to invalid lzID within staker id",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: "0x123_0x666",
					},
				},
			},
			expPass: false,
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
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: "0x123_0x65",
					},
					{
						StakerID: "0x123_0x65",
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to nil staker amounts",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info:    &types.StakerAssetInfo{},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to unexpected staker amounts",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(0),
									WaitUnbondingAmount: math.NewInt(1),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to unexpected staker assetID",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID + "fake",
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(0),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to duplicate staker assetID",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(0),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(0),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to negative staker amount",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(-1),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to amount mismatch between operator and stakers",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(2),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(1),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis due to amount mismatch on deposits",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(2),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(2),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "valid genesis",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				ClientChains: []types.ClientChainInfo{
					ethClientChain,
				},
				Tokens: []types.StakingAssetInfo{
					stakingInfo,
				},
				OperatorAssetInfos: []types.OpAssetIDAndInfos{
					{
						OperatorAddress: accAddress.String(),
						AssetIdAndInfos: []*types.OpAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.OperatorAssetInfo{
									TotalAmount:                        math.NewInt(1),
									OperatorAmount:                     math.NewInt(0),
									WaitUnbondingAmount:                math.NewInt(0),
									OperatorUnbondingAmount:            math.NewInt(0),
									OperatorUnbondableAmountAfterSlash: math.NewInt(0),
								},
							},
						},
					},
				},
				StakerAssetInfos: []types.StAssetIDAndInfos{
					{
						StakerID: stakerID,
						AssetIdAndInfos: []*types.StAssetIDAndInfo{
							{
								AssetID: assetID,
								Info: &types.StakerAssetInfo{
									TotalDepositAmount:  math.NewInt(1),
									WaitUnbondingAmount: math.NewInt(0),
									WithdrawableAmount:  math.NewInt(0),
								},
							},
						},
					},
				},
			},
			expPass: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		if tc.malleate != nil {
			tc.malleate(tc.genState)
		}
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
