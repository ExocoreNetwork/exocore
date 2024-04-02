package types_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	sharedKey := hexutil.Encode(ed25519.GenPrivKey().PubKey().Bytes())
	params := types.DefaultParams()
	newGen := &types.GenesisState{
		Params: params,
	}

	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
		malleate func(*types.GenesisState)
	}{
		{
			name:     "valid genesis constructor",
			genState: newGen,
			expPass:  false, // no validators
		},
		{
			name:     "default",
			genState: types.DefaultGenesis(),
			expPass:  false,
		},
		{
			name: "invalid genesis since it is missing 0x prefix",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						PublicKey: sharedKey,
						Power:     5,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis since it has the wrong length",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						PublicKey: sharedKey,
						Power:     5,
					},
				},
			},
			expPass: false,
		},
		{
			name: "valid genesis with one validator",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						PublicKey: sharedKey,
						Power:     5,
					},
				},
			},
			expPass: true,
		},
		{
			name: "invalid genesis with duplicate validators",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						PublicKey: sharedKey,
						Power:     5,
					},
					{
						PublicKey: sharedKey,
						Power:     10,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis with too many validators",
			genState: &types.GenesisState{
				Params: params,
			},
			malleate: func(gs *types.GenesisState) {
				for i := 0; i < int(gs.Params.MaxValidators)+1; i++ {
					key := hexutil.Encode(ed25519.GenPrivKey().PubKey().Bytes())
					gs.InitialValSet = append(gs.InitialValSet, types.GenesisValidator{
						PublicKey: key,
						Power:     5,
					})
				}
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
