package aggregator

import (
	"math/big"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregatorContext(t *testing.T) {
	Convey("init aggregatorContext with default params", t, func() {
		agc := initAggregatorContext4Test()
		var p *monkey.PatchGuard
		var ctx sdk.Context
		Convey("prepare round to gengerate round info of feeders for next block", func() {
			Convey("pepare within the window", func() {
				p = patchBlockHeight(12)
				agc.PrepareRound(ctx, 0)
				Convey("for empty round list", func() {
					So(*agc.rounds[1], ShouldResemble, roundInfo{10, 2, 1})
				})
				Convey("update already exist round info", func() {
					p = patchBlockHeight(10 + common.MaxNonce)
					agc.PrepareRound(ctx, 0)
					So(agc.rounds[1].status, ShouldEqual, 2)
				})
			})
			Convey("pepare outside the window", func() {
				Convey("for empty round list", func() {
					p = patchBlockHeight(10 + common.MaxNonce)
					agc.PrepareRound(ctx, 0)
					So(agc.rounds[1].status, ShouldEqual, 2)
				})
			})
		})
		Convey("seal existed round without any msg recieved", func() {
			p = patchBlockHeight(11)
			agc.PrepareRound(ctx, 0)
			Convey("seal when exceed the window", func() {
				So(agc.rounds[1].status, ShouldEqual, 1)
				p = patchBlockHeight(13)
				agc.SealRound(ctx, false)
				So(agc.rounds[1].status, ShouldEqual, 2)
			})
			Convey("force seal by required", func() {
				p = patchBlockHeight(12)
				agc.SealRound(ctx, false)
				So(agc.rounds[1].status, ShouldEqual, 1)
				agc.SealRound(ctx, true)
				So(agc.rounds[1].status, ShouldEqual, 2)
			})
		})

		if p != nil {
			p.Unpatch()
		}
	})
}

func initAggregatorContext4Test() *AggregatorContext {
	agc := NewAggregatorContext()

	validatorPowers := map[string]*big.Int{
		"v1": big.NewInt(1),
		"v2": big.NewInt(1),
		"v3": big.NewInt(1),
	}

	p := defaultParams
	pWrapped := common.Params(p)

	agc.SetValidatorPowers(validatorPowers)
	agc.SetParams(&pWrapped)
	return agc
}

func patchBlockHeight(h int64) *monkey.PatchGuard {
	return monkey.PatchInstanceMethod(reflect.TypeOf(sdk.Context{}), "BlockHeight", func(sdk.Context) int64 {
		return h
	})
}
