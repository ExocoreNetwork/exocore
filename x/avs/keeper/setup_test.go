package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	blscommon "github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v14/app"
	utiltx "github.com/evmos/evmos/v14/testutil/tx"
	evm "github.com/evmos/evmos/v14/x/evm/types"

	"github.com/ExocoreNetwork/exocore/testutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type AVSTestSuite struct {
	testutil.BaseTestSuite

	// needed by test
	operatorAddr          sdk.AccAddress
	avsAddr               string
	assetID               string
	stakerID              string
	assetAddr             common.Address
	assetDecimal          uint32
	clientChainLzID       uint64
	depositAmount         sdkmath.Int
	delegationAmount      sdkmath.Int
	updatedAmountForOptIn sdkmath.Int

	ctx            sdk.Context
	app            *app.Evmos
	queryClientEvm evm.QueryClient
	consAddress    sdk.ConsAddress
	avsAddress     common.Address
	taskAddress    common.Address
	taskId         uint64
	blsKey         blscommon.SecretKey
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
	suite.taskAddress = utiltx.GenerateAddress()

}
