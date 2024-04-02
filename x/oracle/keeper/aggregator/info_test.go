package aggregator

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

var (
	one     = big.NewInt(1)
	zero    = big.NewInt(0)
	ten     = big.NewInt(10)
	eleven  = big.NewInt(11)
	fifteen = big.NewInt(15)
	twenty  = big.NewInt(20)
)

var (
	pTD1  = newPTD("1", "10")
	pTD2  = newPTD("2", "12")
	pTD3  = newPTD("3", "15")
	pTD2M = newPTD("2", "11")
	pTD3M = newPTD("3", "19")
	//1-10, 2-12
	pS1 = []*types.PriceWithSource{newPS(1, pTD1, pTD2)}
	//2-12, 3-1
	pS2 = []*types.PriceWithSource{newPS(1, pTD3, pTD2)}
	//1-10, 2-11(m)
	pS3 = []*types.PriceWithSource{newPS(1, pTD1, pTD2M)}
	//2-12, 3-19(m)
	pS4 = []*types.PriceWithSource{newPS(1, pTD2, pTD3M)}
	//1-10, 3-19(m)
	pS5 = []*types.PriceWithSource{newPS(1, pTD1, pTD3M)}

	pS6 = []*types.PriceWithSource{newPS(2, pTD1)}

	//1-10, 2-12
	pS21 = []*types.PriceWithSource{newPS(1, pTD1, pTD2), newPS(2, pTD1, pTD3)}
	//2-12, 3-15
	pS22 = []*types.PriceWithSource{newPS(1, pTD3, pTD2), newPS(2, pTD2, pTD3)}
	//1-10, 2-11(m)
	pS23 = []*types.PriceWithSource{newPS(1, pTD1, pTD2M), newPS(2, pTD2M, pTD1)}
	//2-12, 3-19(m)
	pS24 = []*types.PriceWithSource{newPS(1, pTD2, pTD3M), newPS(2, pTD3, pTD2M)}
	//1-10, 3-19(m)
	pS25 = []*types.PriceWithSource{newPS(1, pTD1, pTD3M), newPS(2, pTD2M, pTD3M)}
)

var defaultParams = types.Params{
	Chains:       []*types.Chain{{Name: "-", Desc: "-"}, {Name: "Ethereum", Desc: "-"}},
	Tokens:       []*types.Token{{}, {Name: "eth", ChainId: 1, ContractAddress: "0xabc", Decimal: 18, Active: true}},
	Sources:      []*types.Source{{}, {Name: "chainLink", Entry: &types.Endpoint{}, Valid: true, Deterministic: true}},
	Rules:        []*types.RuleWithSource{{}, {SourceIds: []int32{1}}},
	TokenFeeders: []*types.TokenFeeder{{}, {TokenId: 1, RuleId: 1, StartRoundId: 1, StartBaseBlock: 0, Interval: 10, EndBlock: 0}},
}
