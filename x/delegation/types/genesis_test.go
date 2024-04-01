package types_test

import (
	"testing"

	"cosmossdk.io/math"
	utiltx "github.com/ExocoreNetwork/exocore/testutil/tx"
	"github.com/ExocoreNetwork/exocore/x/delegation/types"
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
	accAddress := sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
		malleate func(*types.GenesisState)
	}{
		{
			name:     "valid empty genesis",
			genState: &types.GenesisState{},
			expPass:  true,
		},
		{
			name:     "default",
			genState: types.DefaultGenesis(),
			expPass:  true,
		},
		{
			name: "invalid staker id",
			genState: &types.GenesisState{
				DelegationsByStakerAssetOperator: []types.DelegationByStakerAssetOperator{
					{
						StakerID: "asd_asd_0x64",
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid asset id",
			genState: &types.GenesisState{
				DelegationsByStakerAssetOperator: []types.DelegationByStakerAssetOperator{
					{
						StakerID: "asd_0x64",
						DelegationsByAssetOperator: []types.DelegationByAssetOperator{
							{
								AssetID: "",
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid operator address",
			genState: &types.GenesisState{
				DelegationsByStakerAssetOperator: []types.DelegationByStakerAssetOperator{
					{
						StakerID: "asd_0x64",
						DelegationsByAssetOperator: []types.DelegationByAssetOperator{
							{
								AssetID: "abcd_0x64",
								DelegationsByOperator: []types.DelegationByOperator{
									{
										OperatorAddress: "fake",
									},
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "nil amount",
			genState: &types.GenesisState{
				DelegationsByStakerAssetOperator: []types.DelegationByStakerAssetOperator{
					{
						StakerID: "asd_0x64",
						DelegationsByAssetOperator: []types.DelegationByAssetOperator{
							{
								AssetID: "abcd_0x64",
								DelegationsByOperator: []types.DelegationByOperator{
									{
										OperatorAddress: accAddress.String(),
									},
								},
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "negative amount",
			genState: &types.GenesisState{
				DelegationsByStakerAssetOperator: []types.DelegationByStakerAssetOperator{
					{
						StakerID: "asd_0x64",
						DelegationsByAssetOperator: []types.DelegationByAssetOperator{
							{
								AssetID: "abcd_0x64",
								DelegationsByOperator: []types.DelegationByOperator{
									{
										OperatorAddress: accAddress.String(),
										Amount:          math.NewInt(-1),
									},
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
				DelegationsByStakerAssetOperator: []types.DelegationByStakerAssetOperator{
					{
						StakerID: "asd_0x64",
						DelegationsByAssetOperator: []types.DelegationByAssetOperator{
							{
								AssetID: "abcd_0x64",
								DelegationsByOperator: []types.DelegationByOperator{
									{
										OperatorAddress: accAddress.String(),
										Amount:          math.NewInt(1),
									},
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
