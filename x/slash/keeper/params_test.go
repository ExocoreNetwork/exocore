package keeper_test

import (
	slashtype "github.com/exocore/x/slash/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := &slashtype.Params{
		ExoCoreLzAppAddress:    "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		ExoCoreLzAppEventTopic: "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec",
	}
	err := suite.app.ExoSlashKeeper.SetParams(suite.ctx, params)
	suite.NoError(err)

	getParams, err := suite.app.ExoSlashKeeper.GetParams(suite.ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
}
