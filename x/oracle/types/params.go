package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var (
	KeyChains       = []byte("Chains")
	KeyTokens       = []byte("Tokens")
	KeySources      = []byte("Sources")
	KeyRules        = []byte("Rules")
	KeyTokenFeeders = []byte("TokenFeeders")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		Chains: []*Chain{
			{Name: "-", Desc: "-"},
			{Name: "Ethereum", Desc: "-"},
		},
		Tokens: []*Token{
			{},
			{
				Name:            "ETH",
				ChainID:         1,
				ContractAddress: "0x",
				Decimal:         18,
				Active:          true,
				AssetID:         "",
			},
		},
		Sources: []*Source{
			{
				Name: "0 position is reserved",
			},
			{
				Name: "Chainlink",
				Entry: &Endpoint{
					Offchain: map[uint64]string{0: ""},
				},
				Valid:         true,
				Deterministic: true,
			},
		},
		Rules: []*RuleSource{
			// 0 is reserved
			{},
			{
				// all sources math
				SourceIDs: []uint64{0},
			},
		},
		TokenFeeders: []*TokenFeeder{
			{},
			{
				TokenID:        1,
				RuleID:         1,
				StartRoundID:   1,
				StartBaseBlock: 20,
				Interval:       10,
			},
		},
		MaxNonce:   3,
		ThresholdA: 2,
		ThresholdB: 3,
		Mode:       1,
		MaxDetId:   5,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyChains, &p.Chains, func(_ interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeyTokens, &p.Tokens, func(_ interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeySources, &p.Sources, func(_ interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeyRules, &p.Rules, func(_ interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeyTokenFeeders, &p.TokenFeeders, func(_ interface{}) error { return nil }),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
