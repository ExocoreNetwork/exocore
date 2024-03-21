package keeper_test

import (
	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
)

func (suite *DepositTestSuite) TestParams() {
	params := &deposittype.Params{}
	err := suite.App.DepositKeeper.SetParams(suite.Ctx, params)
	suite.NoError(err)

	getParams, err := suite.App.DepositKeeper.GetParams(suite.Ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
}
