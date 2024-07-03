package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/evmos/v14/app"
	evm "github.com/evmos/evmos/v14/x/evm/types"
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type AVSTestSuite struct {
	testutil.BaseTestSuite

	ctx            sdk.Context
	app            *app.Evmos
	queryClientEvm evm.QueryClient
	consAddress    sdk.ConsAddress
}

var s *AVSTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(AVSTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *AVSTestSuite) SetupTest() {
	suite.DoSetupTest()
}
