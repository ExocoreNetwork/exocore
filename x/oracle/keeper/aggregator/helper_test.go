package aggregator

import "github.com/ExocoreNetwork/exocore/x/oracle/types"

func newPTD(detId, price string) *types.PriceWithTimeAndDetId {
	return &types.PriceWithTimeAndDetId{
		Price:     price,
		Decimal:   1,
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
