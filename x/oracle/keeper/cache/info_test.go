package cache

import "github.com/ExocoreNetwork/exocore/x/oracle/types"

var defaultParams = types.Params{
	Chains:       []*types.Chain{{Name: "-", Desc: "-"}, {Name: "Ethereum", Desc: "-"}},
	Tokens:       []*types.Token{{}, {Name: "eth", ChainId: 1, ContractAddress: "0xabc", Decimal: 18, Active: true}},
	Sources:      []*types.Source{{}, {Name: "chainLink", Entry: &types.Endpoint{}, Valid: true, Deterministic: true}},
	Rules:        []*types.RuleWithSource{{}, {SourceIds: []int32{1}}},
	TokenFeeders: []*types.TokenFeeder{{}, {TokenId: 1, RuleId: 1, StartRoundId: 1, StartBaseBlock: 0, Interval: 10, EndBlock: 0}},
}
