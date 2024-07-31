package types_test

import (
	"testing"

	utiltx "github.com/ExocoreNetwork/exocore/testutil/tx"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
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
	accAddress1 := sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	newGen := &types.GenesisState{}

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
			name: "invalid genesis state due to non bech32 operator address",
			genState: &types.GenesisState{
				Operators: []types.OperatorInfo{
					{
						EarningsAddr: "invalid",
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis state due to duplicate operator address",
			genState: &types.GenesisState{
				Operators: []types.OperatorInfo{
					{
						EarningsAddr: accAddress1.String(),
					},
					{
						EarningsAddr: accAddress1.String(),
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis state due to duplicate lz chain id",
			genState: &types.GenesisState{
				Operators: []types.OperatorInfo{
					{
						EarningsAddr: accAddress1.String(),
						ClientChainEarningsAddr: &types.ClientChainEarningAddrList{
							EarningInfoList: []*types.ClientChainEarningAddrInfo{
								{
									LzClientChainID:        1,
									ClientChainEarningAddr: utiltx.GenerateAddress().String(),
								},
								{
									LzClientChainID:        1,
									ClientChainEarningAddr: utiltx.GenerateAddress().String(),
								},
							},
						},
					},
				},
			},
			expPass: false,
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
