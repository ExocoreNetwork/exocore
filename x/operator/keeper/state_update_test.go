package keeper_test

func (s *KeeperTestSuite) prepare() {
	s.avsAddr = "avsTestAddr"
	//staking assets
	//register operator
	//delegate to operator
}
func (s *KeeperTestSuite) OptIn() {
	s.prepare()
	s.app.OperatorKeeper.OptIn(s.ctx)
}
func (s *KeeperTestSuite) UpdateOptedInAssetsState() {

}
