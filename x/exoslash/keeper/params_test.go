package keeper_test

import (
<<<<<<< HEAD
<<<<<<< HEAD
	types2 "github.com/exocore/x/exoslash/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := &types2.Params{
		ExoCoreLzAppAddress:    "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		ExoCoreLzAppEventTopic: "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec",
	}
	err := suite.app.ExoslashKeeper.SetParams(suite.ctx, params)
	suite.NoError(err)

	getParams, err := suite.app.ExoslashKeeper.GetParams(suite.ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
=======
	"testing"
)

func TestGetParams(t *testing.T) {
>>>>>>> eebca7f (implement slash interface)
=======
	types2 "github.com/exocore/x/exoslash/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := &types2.Params{
		ExoCoreLzAppAddress:    "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
		ExoCoreLzAppEventTopic: "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec",
	}
	err := suite.app.ExoslashKeeper.SetParams(suite.ctx, params)
	suite.NoError(err)

	getParams, err := suite.app.ExoslashKeeper.GetParams(suite.ctx)
	suite.NoError(err)
	suite.Equal(*params, *getParams)
>>>>>>> 593c3a5 (add param unit test)
}
