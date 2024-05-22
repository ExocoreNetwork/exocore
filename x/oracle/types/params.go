package types

import (
	"errors"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	TypeChain = iota
	TypeSource
	TypeToken
	TypeRule
	TypeTokenFeeder
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
	return Params{
		Chains:       []*Chain{{}},
		Tokens:       []*Token{{}},
		Sources:      []*Source{{}},
		Rules:        []*RuleSource{{}},
		TokenFeeders: []*TokenFeeder{{}},
	}
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
				StartBaseBlock: 30,
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
	// validate tokenFeeders format
	for fID, feeder := range p.TokenFeeders {
		if fID == 0 {
			continue
		}
		if feeder.TokenID < 1 {
			return errors.New("")
		}

		if feeder.EndBlock > 0 && (feeder.StartBaseBlock >= feeder.EndBlock || feeder.StartBaseBlock == 0) {
			return errors.New("")
		}

		if feeder.Interval > 0 && feeder.Interval < 3*2 {

			return errors.New("")

		}

		if feeder.StartBaseBlock > 0 && feeder.EndBlock > 0 {
			// create a new feeder with endblock set, endblock should not whithin the window of one price round
			if feeder.Interval > 0 && (feeder.EndBlock-feeder.StartBaseBlock)%feeder.Interval < 3 {
				return errors.New("")
			}
		}
	}
	// TODO: validate chains, tokens, rules, tokenfeeders, and cross validation among these fields

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p Params) GetSourceIDByName(n string) int {
	for i, c := range p.Sources {
		if n == c.Name {
			return i
		}
	}
	return 0
}

func (p Params) GetFeederIDsByTokenID(tID uint64) []int {
	ret := make([]int, 0)
	for fID, f := range p.TokenFeeders {
		// feeder list is ordered, so the slice returned is in the order of the feeders updated for the same token
		if f.TokenID == tID {
			ret = append(ret, fID)
		}
	}
	return ret
}
