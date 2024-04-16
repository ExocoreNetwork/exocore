package testdata

import "github.com/ExocoreNetwork/exocore/x/oracle/types"

func newPTD(detId, price string) *types.PriceWithTimeAndDetId {
	return &types.PriceWithTimeAndDetId{
		Price:     price,
		Decimal:   18,
		Timestamp: "-",
		DetId:     detId,
	}
}

func newPS(sourceId uint64, prices ...*types.PriceWithTimeAndDetId) *types.PriceWithSource {
	return &types.PriceWithSource{
		SourceId: sourceId,
		Prices:   prices,
	}
}

func newPTR(price string, roundId uint64) *types.PriceWithTimeAndRound {
	return &types.PriceWithTimeAndRound{
		Price:     price,
		Decimal:   18,
		Timestamp: "",
		RoundId:   roundId,
	}
}

func newPrices(tokenId uint64, nextRoundId uint64, pList ...*types.PriceWithTimeAndRound) types.Prices {
	return types.Prices{
		TokenId:     tokenId,
		NextRoundId: nextRoundId,
		PriceList:   pList,
	}
}
