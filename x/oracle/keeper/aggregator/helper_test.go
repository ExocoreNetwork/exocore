package aggregator

import "github.com/ExocoreNetwork/exocore/x/oracle/types"

func newPTD(detID, price string) *types.PriceTimeDetID {
	return &types.PriceTimeDetID{
		Price:     price,
		Decimal:   1,
		Timestamp: "-",
		DetID:     detID,
	}
}

func newPS(sourceID uint64, prices ...*types.PriceTimeDetID) *types.PriceWithSource {
	return &types.PriceWithSource{
		SourceID: sourceID,
		Prices:   prices,
	}
}
