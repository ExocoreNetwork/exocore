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
	testCases := []struct {
		name     string
		genState *types.GenesisState
		expPass  bool
		malleate func(*types.GenesisState)
	}{
		{
			name: "constructor",
			genState: &types.GenesisState{
				Params: params,
			},
			expPass: true,
		},
		{
			name:     "default",
			genState: types.DefaultGenesis(),
			expPass:  true,
		},
		{
			name: "NewGenesis call",
			genState: types.NewGenesis(
				params, []types.GenesisValidator{},
			),
			expPass: true,
		},
		{
			name: "too many validators",
			genState: &types.GenesisState{
				Params: params,
			},
			malleate: func(gs *types.GenesisState) {
				// note the plus 1
				for i := 0; i < int(gs.Params.MaxValidators)+1; i++ {
					// generate a  new key each time
					key := hexutil.Encode(
						ed25519.GenPrivKey().PubKey().Bytes(),
					)
					gs.InitialValSet = append(gs.InitialValSet,
						types.GenesisValidator{
							PublicKey: key,
							Power:     5,
						},
					)
				}
			},
			expPass: false,
		},
		{
			name: "duplicate keys",
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
			name: "key with missing 0x prefix",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						// remove 2 chars and add 2 chars
						PublicKey: sharedKey[2:] + "ab",
						Power:     5,
					},
				},
			},
			expPass: false,
		},
		{
			// also covers empty key
			name: "key with the wrong length",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						PublicKey: sharedKey + "ab",
						Power:     5,
					},
				},
			},
			expPass: false,
		},
		{
			name: "non hex key",
			genState: &types.GenesisState{
				Params: params,
				InitialValSet: []types.GenesisValidator{
					{
						// replace last 2 chars with non-hex
						PublicKey: sharedKey[:64] + "ss",
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
		// fmt.Println(tc.name, err)
	}
}
