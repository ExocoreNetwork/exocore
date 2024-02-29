package keeper_test

import (
	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
)

func (suite *DepositTestSuite) TestParams() {
	params := &deposittype.Params{
		ExoCoreLzAppAddress:    "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		ExoCoreLzAppEventTopic: "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec",
	}
	err := suite.App.DepositKeeper.SetParams(suite.Ctx, params)
	suite.NoError(err)

	getParams, err := suite.App.DepositKeeper.GetParams(suite.Ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
}
