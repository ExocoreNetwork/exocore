package testdata

import "github.com/ExocoreNetwork/exocore/x/oracle/types"

func newPTD(detID, price string) *types.PriceWithTimeAndDetId {
	return &types.PriceWithTimeAndDetId{
		Price:     price,
		Decimal:   18,
		Timestamp: "-",
		DetID:     detID,
	}
}

func newPS(sourceID uint64, prices ...*types.PriceWithTimeAndDetId) *types.PriceWithSource {
	return &types.PriceWithSource{
		SourceID: sourceID,
		Prices:   prices,
	}
}

func newPTR(price string, roundID uint64) *types.PriceWithTimeAndRound {
	return &types.PriceWithTimeAndRound{
		Price:     price,
		Decimal:   18,
		Timestamp: "",
		RoundID:   roundID,
	}
}

func newPrices(tokenID uint64, nextRoundID uint64, pList ...*types.PriceWithTimeAndRound) types.Prices {
	return types.Prices{
		TokenID:     tokenID,
		NextRoundID: nextRoundID,
		PriceList:   pList,
	}
}
