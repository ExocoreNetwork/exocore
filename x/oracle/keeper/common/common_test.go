package common

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination mock_test.go -package common github.com/ExocoreNetwork/exocore/x/oracle/keeper/common KeeperOracle

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ko := NewMockKeeperOracle(ctrl)

	ko.EXPECT().GetLastTotalPower(gomock.Any()).Return(big.NewInt(99))

	x := ko.GetLastTotalPower(sdk.Context{})
	_ = x

	Convey("mock oracle keeper", t, func() {
		Convey("GetLastTotalPower", func() { So(x, ShouldResemble, big.NewInt(99)) })
	})
}
