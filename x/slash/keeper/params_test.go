package keeper_test

import (
	slashtype "github.com/ExocoreNetwork/exocore/x/slash/types"
)

func (suite *SlashTestSuite) TestParams() {
	params := &slashtype.Params{}
	err := suite.App.ExoSlashKeeper.SetParams(suite.Ctx, params)
	suite.NoError(err)

	getParams, err := suite.App.ExoSlashKeeper.GetParams(suite.Ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
}
