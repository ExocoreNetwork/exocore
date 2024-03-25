package cache

import (
	"math/big"
	"testing"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	. "github.com/smartystreets/goconvey/convey"
	// "go.uber.org/mock/gomock"
)

func TestCache(t *testing.T) {
	c := NewCache()
	p := defaultParams
	pWrapped := common.Params(p)

	//	ctrl := gomock.NewController(t)
	//	defer ctrl.Finish()
	//ko := common.NewMockKeeperOracle(ctrl)
	//c.AddCache(CacheItemP(&pWrapped), ko)

	Convey("test cache", t, func() {
		Convey("add pramams item", func() {
			c.AddCache(CacheItemP(&pWrapped))
			pReturn := &common.Params{}
			c.GetCache(CacheItemP(pReturn))
			So(*pReturn, ShouldResemble, pWrapped)
		})

		Convey("add validatorPower item", func() {
			validatorPowers := map[string]*big.Int{
				"v1": big.NewInt(100),
				"v2": big.NewInt(109),
				"v3": big.NewInt(119),
			}
			c.AddCache(CacheItemV(validatorPowers))
			vpReturn := make(map[string]*big.Int)
			Convey("for empty cache", func() {
				c.GetCache(CacheItemV(vpReturn))
				So(vpReturn, ShouldResemble, validatorPowers)
			})
			Convey("then update validatorPower item for this cache", func() {
				validaotrPowers := map[string]*big.Int{
					//add v5
					"v5": big.NewInt(123),
					//remove v1
					"v1": big.NewInt(0),
					//update v2
					"v2": big.NewInt(199),
				}
				c.AddCache(CacheItemV(validaotrPowers))
				c.GetCache(CacheItemV(vpReturn))
				So(vpReturn, ShouldNotContainKey, "v1")
				So(vpReturn, ShouldContainKey, "v5")
				So(vpReturn["v2"], ShouldResemble, big.NewInt(199))
			})
		})

		Convey("add msg item", func() {
			msgItems := []*CacheItemM{
				{
					FeederId: 1,
					PSources: []*types.PriceWithSource{
						{
							SourceId: 1,
							Prices: []*types.PriceWithTimeAndDetId{
								{Price: "600000", Decimal: 1, Timestamp: "-", DetId: "1"}, {Price: "620000", Decimal: 1, Timestamp: "-", DetId: "2"},
							},
						},
					},
					Validator: "v1",
				},
				{
					FeederId: 1,
					PSources: []*types.PriceWithSource{
						{SourceId: 1, Prices: []*types.PriceWithTimeAndDetId{{Price: "600000", Decimal: 1, Timestamp: "-", DetId: "4"}, {Price: "620000", Decimal: 1, Timestamp: "-", DetId: "3"}}}},
					Validator: "v1",
				},
				{
					FeederId:  2,
					PSources:  []*types.PriceWithSource{{SourceId: 1, Prices: []*types.PriceWithTimeAndDetId{{Price: "30000", Decimal: 1, Timestamp: "-", DetId: "4"}, {Price: "32000", Decimal: 1, Timestamp: "-", DetId: "3"}}}},
					Validator: "v2",
				},
			}
			c.AddCache(msgItems[0])
			msgItemsReturn := make([]*CacheItemM, 0, 3)
			Convey("add single item", func() {
				c.GetCache(&msgItemsReturn)
				So(msgItemsReturn, ShouldContain, msgItems[0])
			})
			Convey("add more items", func() {
				c.AddCache(msgItems[1])
				c.AddCache(msgItems[2])
				c.GetCache(&msgItemsReturn)
				So(msgItemsReturn, ShouldContain, msgItems[2])
				So(msgItemsReturn, ShouldNotContain, msgItems[0])
			})
		})
	})
}