package types

import (
	"errors"

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
				AssetID:         "0x0b34c4d876cd569129cf56bafabb3f9e97a4ff42_0x9ce1",
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
				StartBaseBlock: 100000000,
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

func (p Params) GetTokenIDFromAssetID(assetID string) int {
	for id, token := range p.Tokens {
		if token.AssetID == assetID {
			return id
		}
	}
	return 0
}

func (p Params) IsDeterministicSource(sourceID uint64) bool {
	return p.Sources[sourceID].Deterministic
}

func (p Params) IsValidSource(sourceID uint64) bool {
	if sourceID == 0 {
		// custom defined source
		return true
	}
	return p.Sources[sourceID].Valid
}

func (p Params) GetTokenFeeder(feederID uint64) *TokenFeeder {
	for k, v := range p.TokenFeeders {
		if uint64(k) == feederID {
			return v
		}
	}
	return nil
}

func (p Params) GetTokenInfo(feederID uint64) *Token {
	for k, v := range p.TokenFeeders {
		if uint64(k) == feederID {
			return p.Tokens[v.TokenID]
		}
	}
	return nil
}

func (p Params) CheckRules(feederID uint64, prices []*PriceSource) (bool, error) {
	feeder := p.TokenFeeders[feederID]
	rule := p.Rules[feeder.RuleID]
	// specified sources set, v1 use this rule to set `chainlink` as official source
	if rule.SourceIDs != nil && len(rule.SourceIDs) > 0 {
		if len(rule.SourceIDs) != len(prices) {
			return false, errors.New("count prices should match rule")
		}
		notFound := false
		if rule.SourceIDs[0] == 0 {
			// match all sources listed
			for sID, source := range p.Sources {
				if sID == 0 {
					continue
				}
				if source.Valid {
					notFound = true
					for _, p := range prices {
						if p.SourceID == uint64(sID) {
							notFound = false
							break
						}
					}

				}
			}
		} else {
			for _, source := range rule.SourceIDs {
				notFound = true
				for _, p := range prices {
					if p.SourceID == source {
						notFound = false
						break
					}
				}
			}
		}
		if notFound {
			return false, errors.New("price source not match with rule")
		}
	}

	// TODO: check NOM
	// return true if no rule set, we will accept any source
	return true, nil
}
