package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v16/app"
	utiltx "github.com/evmos/evmos/v16/testutil/tx"
	evm "github.com/evmos/evmos/v16/x/evm/types"

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
	avsAddress     common.Address
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
	suite.avsAddress = utiltx.GenerateAddress()
}
