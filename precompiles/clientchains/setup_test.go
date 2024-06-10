package clientchains_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/clientchains"
	"github.com/ExocoreNetwork/exocore/testutil"

	"github.com/stretchr/testify/suite"
)

var s *ClientChainsPrecompileSuite

type ClientChainsPrecompileSuite struct {
	testutil.BaseTestSuite

	precompile *clientchains.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(ClientChainsPrecompileSuite)
	suite.Run(t, s)
}

func (s *ClientChainsPrecompileSuite) SetupTest() {
	s.DoSetupTest()
	precompile, err := clientchains.NewPrecompile(
		s.App.AssetsKeeper, s.App.AuthzKeeper,
	)
	s.Require().NoError(err)
	s.precompile = precompile
}
