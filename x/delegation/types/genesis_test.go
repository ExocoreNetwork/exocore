package types_test

import (
	"testing"

	"cosmossdk.io/math"
	utiltx "github.com/ExocoreNetwork/exocore/testutil/tx"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
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
	assetAddress := utiltx.GenerateAddress()
	stakerAddress := utiltx.GenerateAddress()
	lzID := uint64(101)
	stakerID, assetID := assetstypes.GetStakeIDAndAssetID(
		lzID, stakerAddress[:], assetAddress[:],
	)
	operatorAddress := sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	delegations := []types.DelegationsByStaker{
		{
			StakerID: stakerID,
			Delegations: []types.DelegatedSingleAssetInfo{
				{
					AssetID:              assetID,
					TotalDelegatedAmount: math.NewInt(1000),
					PerOperatorAmounts: map[string]*types.ValueField{
						operatorAddress.String(): {Amount: math.NewInt(1000)},
					},
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
			name:     "base, should pass",
			genState: types.NewGenesis(delegations),
			expPass:  true,
		},
		{
			name:     "invalid staker id",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].StakerID = "invalid"
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].StakerID = stakerID
			},
		},
		{
			name:     "duplicate staker id",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations = append(gs.Delegations, gs.Delegations[0])
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations = gs.Delegations[:1]
			},
		},
		{
			name:     "duplicate asset id",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations = append(
					gs.Delegations[0].Delegations,
					gs.Delegations[0].Delegations[0],
				)
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations = gs.Delegations[0].Delegations[:1]
			},
		},
		{
			name:     "invalid asset id",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].AssetID = "invalid"
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].AssetID = assetID
			},
		},
		{
			name:     "asset id mismatch",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				stakerID, _ := assetstypes.GetStakeIDAndAssetID(
					lzID+1, stakerAddress[:], assetAddress[:],
				)
				gs.Delegations[0].StakerID = stakerID
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].StakerID = stakerID
			},
		},
		{
			name:     "invalid total amount",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].TotalDelegatedAmount = math.NewInt(-1)
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].TotalDelegatedAmount = math.NewInt(1000)
			},
		},
		{
			name:     "nil total amount",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].TotalDelegatedAmount = math.Int{}
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].TotalDelegatedAmount = math.NewInt(1000)
			},
		},
		{
			name:     "nil wrapped amount",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					nil
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{Amount: math.NewInt(1000)}
			},
		},
		{
			name:     "nil unwrapped amount",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{}
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{Amount: math.NewInt(1000)}
			},
		},
		{
			name:     "negative unwrapped amount",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{Amount: math.NewInt(-1)}
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{Amount: math.NewInt(1000)}
			},
		},
		{
			name:     "invalid operator address",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts["invalid"] =
					&types.ValueField{Amount: math.NewInt(1000)}
			},
			unmalleate: func(gs *types.GenesisState) {
				delete(gs.Delegations[0].Delegations[0].PerOperatorAmounts, "invalid")
			},
		},
		{
			name:     "total amount mismatch",
			genState: types.NewGenesis(delegations),
			expPass:  false,
			malleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{Amount: math.NewInt(2000)}
			},
			unmalleate: func(gs *types.GenesisState) {
				gs.Delegations[0].Delegations[0].PerOperatorAmounts[operatorAddress.String()] =
					&types.ValueField{Amount: math.NewInt(1000)}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		if tc.malleate != nil {
			tc.malleate(tc.genState)
			// require that unmalleate is defined
			suite.Require().NotNil(tc.unmalleate, tc.name)
		}
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
		if tc.unmalleate != nil {
			tc.unmalleate(tc.genState)
		}
		// fmt.Println(tc.name, err)
	}
}
