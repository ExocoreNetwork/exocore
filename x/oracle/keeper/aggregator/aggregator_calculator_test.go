package aggregator

import (
	"math/big"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

/*
	1-10, 	2-12, 	3-15

ps1: 1-10,	2-12
ps2: 		2-12,	3-15
ps3: 1-10,	2-11(m)
---
ps4: 		2-12,	3-19(m)
ps5: 1-10, 			3-19(m)
----
ps1, ps2, ps3, ps4 ---> 2-12
ps2, ps2, ps3, ps5 ---> 1-10
*/
func TestCalculator(t *testing.T) {
	//pTD1 := newPTD("1", "10")
	//pTD2 := newPTD("2", "12")
	//pTD3 := newPTD("3", "15")
	//pTD2M := newPTD("2", "11")
	//pTD3M := newPTD("3", "19")
	////1-10, 2-12
	//ps1 := []*types.PriceWithSource{newPS(1, pTD1, pTD2)}
	////2-12, 3-15
	//ps2 := []*types.PriceWithSource{newPS(1, pTD3, pTD2)}
	////1-10, 2-11(m)
	//ps3 := []*types.PriceWithSource{newPS(1, pTD1, pTD2M)}
	////2-12, 3-19(m)
	//ps4 := []*types.PriceWithSource{newPS(1, pTD2, pTD3M)}
	////1-10, 3-19(m)
	//ps5 := []*types.PriceWithSource{newPS(1, pTD1, pTD3M)}

	////1-10, 2-12
	//ps21 := []*types.PriceWithSource{newPS(1, pTD1, pTD2), newPS(2, pTD1, pTD3)}
	////2-12, 3-15
	//ps22 := []*types.PriceWithSource{newPS(1, pTD3, pTD2), newPS(2, pTD2, pTD3)}
	////1-10, 2-11(m)
	//ps23 := []*types.PriceWithSource{newPS(1, pTD1, pTD2M), newPS(2, pTD2M, pTD1)}
	////2-12, 3-19(m)
	//ps24 := []*types.PriceWithSource{newPS(1, pTD2, pTD3M), newPS(2, pTD3, pTD2M)}
	////1-10, 3-19(m)
	//ps25 := []*types.PriceWithSource{newPS(1, pTD1, pTD3M), newPS(2, pTD2M, pTD3M)}

	one := big.NewInt(1)
	Convey("fill prices into calculator", t, func() {
		c := newCalculator(5, big.NewInt(4))
		Convey("fill prices from single deterministic source", func() {
			c.fillPrice(pS1, "v1", one) //1-10, 2-12
			c.fillPrice(pS2, "v2", one) //2-12, 3-15
			c.fillPrice(pS3, "v3", one) //1-10, 2-11
			Convey("consensus on detid=2 and price=12", func() {
				confirmed := c.fillPrice(pS4, "v4", one) //2-12, 3-19
				So(confirmed[0].detId, ShouldEqual, "2")
				So(confirmed[0].price, ShouldResemble, big.NewInt(12))
			})
			Convey("consensus on detid=1 and price=10", func() {
				confirmed := c.fillPrice(pS5, "v5", one) //1-10, 3-19
				So(confirmed[0].detId, ShouldEqual, "1")
				So(confirmed[0].price, ShouldResemble, big.NewInt(10))

				confirmed = c.fillPrice(pS4, "v4", one)
				So(confirmed, ShouldBeNil)
			})
		})
		Convey("fill prices from multiple deterministic sources", func() {
			c.fillPrice(pS21, "v1", one)
			c.fillPrice(pS22, "v2", one)
			c.fillPrice(pS23, "v3", one)
			Convey("consensus on both source 1 and source 2", func() {
				confirmed := c.fillPrice(pS24, "v4", one)
				So(len(confirmed), ShouldEqual, 2)
				i := 0
				if confirmed[0].sourceId == 2 {
					i = 1
				}
				So(confirmed[i].detId, ShouldEqual, "2")
				So(confirmed[i].price, ShouldResemble, big.NewInt(12))

				So(confirmed[1-i].detId, ShouldEqual, "3")
				So(confirmed[1-i].price, ShouldResemble, big.NewInt(15))

			})
			Convey("consenus on source 1 only", func() {
				confirmed := c.fillPrice(pS25, "v5", one)
				So(len(confirmed), ShouldEqual, 1)
				So(confirmed[0].detId, ShouldEqual, "1")
				So(confirmed[0].price, ShouldResemble, big.NewInt(10))

			})
		})
	})
}
