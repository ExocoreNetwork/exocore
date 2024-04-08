package keeper_test

import (
	"context"
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type KeeperSuite struct {
	testutil.BaseTestSuite

	t        *testing.T
	k        keeper.Keeper
	ctx      sdk.Context
	ms       types.MsgServer
	ctrl     *gomock.Controller
	valAddr1 sdk.ValAddress
	valAddr2 sdk.ValAddress
}

var ks *KeeperSuite

func TestKeeper(t *testing.T) {
	var ctxW context.Context
	ks = &KeeperSuite{}
	ks.ms, ctxW, ks.k = setupMsgServer(t)
	ks.ctx = sdk.UnwrapSDKContext(ctxW)
	ks.t = t

	suite.Run(t, ks)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *KeeperSuite) Reset() {
	var ctxW context.Context
	suite.ms, ctxW, suite.k = setupMsgServer(suite.t)
	suite.ctx = sdk.UnwrapSDKContext(ctxW)
	suite.ctrl = gomock.NewController(suite.t)
}

func (suite *KeeperSuite) SetupTest() {
	suite.DoSetupTest()
	suite.valAddr1, _ = sdk.ValAddressFromBech32(suite.Validators[0].OperatorAddress)
	suite.valAddr2, _ = sdk.ValAddressFromBech32(suite.Validators[1].OperatorAddress)
	keeper.ResetAggregatorContext()
	keeper.ResetCache()
}
