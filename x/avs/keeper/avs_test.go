package keeper_test

func (suite *AVSTestSuite) TestAVS() {
	avsName, avsAddres, operatorAddress, assetID := "avsTest", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo18h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", ""
	err := suite.App.AVSManagerKeeper.SetAVSInfo(suite.Ctx, avsName, avsAddres, operatorAddress, assetID)
	suite.NoError(err)

	info, err := suite.App.AVSManagerKeeper.GetAVSInfo(suite.Ctx, avsAddres)

	suite.NoError(err)
	suite.Equal(avsAddres, info.GetInfo().AvsAddress)

}
