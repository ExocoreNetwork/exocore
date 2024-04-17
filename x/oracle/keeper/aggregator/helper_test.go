package aggregator

import "github.com/ExocoreNetwork/exocore/x/oracle/types"

func newPTD(detID, price string) *types.PriceWithTimeAndDetId {
	return &types.PriceWithTimeAndDetId{
		Price:     price,
		Decimal:   1,
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
