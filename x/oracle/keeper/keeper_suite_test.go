package keeper_test

import (
	"context"
	"testing"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

type KeeperSuite struct {
	t    *testing.T
	k    keeper.Keeper
	ctx  sdk.Context
	ms   types.MsgServer
	ctrl *gomock.Controller
}

var ks *KeeperSuite

func TestKeeper(t *testing.T) {
	var ctxW context.Context
	ks = &KeeperSuite{}
	ks.ms, ctxW, ks.k = setupMsgServer(t)
	ks.ctx = sdk.UnwrapSDKContext(ctxW)
	ks.t = t

	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (k *KeeperSuite) Reset() {
	var ctxW context.Context
	k.ms, ctxW, k.k = setupMsgServer(k.t)
	k.ctx = sdk.UnwrapSDKContext(ctxW)
	k.ctrl = gomock.NewController(k.t)
}
