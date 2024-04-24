package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	"github.com/stretchr/testify/suite"
)

type DelegationTestSuite struct {
	testutil.BaseTestSuite
	assetAddr        common.Address
	opAccAddr        types.AccAddress
	clientChainLzID  uint64
	depositAmount    sdkmath.Int
	delegationAmount sdkmath.Int
}

var s *DelegationTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(DelegationTestSuite)
	suite.Run(t, s)

}

func (suite *DelegationTestSuite) SetupTest() {
	suite.DoSetupTest()
}
