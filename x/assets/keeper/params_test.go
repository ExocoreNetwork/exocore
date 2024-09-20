package keeper_test

import (
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
)

func (suite *StakingAssetsTestSuite) TestParams() {
	params := &assetstype.Params{
		ExocoreLzAppAddress:    "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		ExocoreLzAppEventTopic: "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec",
	}
	err := suite.App.AssetsKeeper.SetParams(suite.Ctx, params)
	suite.NoError(err)

	getParams, err := suite.App.AssetsKeeper.GetParams(suite.Ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
}

func (suite *StakingAssetsTestSuite) TestNullFlag() {
	genesisState := assetstype.GenesisState{
		Params: assetstype.Params{
			ExocoreLzAppAddress:    "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
			ExocoreLzAppEventTopic: "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec",
		},
	}
	bytes := suite.App.AppCodec().MustMarshalJSON(&genesisState)
	var unmarshalResult assetstype.GenesisState
	suite.App.AppCodec().MustUnmarshalJSON(bytes, &unmarshalResult)
	suite.False(unmarshalResult.NotInitFromBootStrap)
}
