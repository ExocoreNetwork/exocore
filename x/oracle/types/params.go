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
	return Params{
		MaxNonce:     1,
		MaxDetId:     1,
		Mode:         1,
		ThresholdA:   2,
		ThresholdB:   3,
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
				StartBaseBlock: 1000000,
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
	// Some basic configure params validation:
	// Maxnonce: This tells how many transactions for one round can a validator send, This also restrict how many blocks a window lasts for one round to collect transactions
	// MaxDetID: This only works for DS, to tell how many continuous roundID_from_DS could be accept at most for one round of exorcore_oracle
	// ThresholdA/ThresholdB: represents the threshold of voting power to confirm a price as final price
	// Mode: tells how and when to confirm a final price, expect for voting power threshold, v1 set this value to 1 means final price will be confirmed as soon as it has reached the threshold of total voting power, and just ignore any remaining transactions followed for current round.
	if p.MaxNonce < 1 || p.MaxDetId < 1 || p.ThresholdA < 1 || p.ThresholdB < p.ThresholdA || p.Mode != 1 {
		return ErrInvalidParams.Wrapf("invalid maxNonce/maxDetID/Threshold/Mode: %d, %d, %d, %d, %d", p.MaxNonce, p.MaxDetId, p.ThresholdA, p.ThresholdB, p.Mode)
	}

	// validate tokenFeeders
	for fID, feeder := range p.TokenFeeders {
		// id==0 is reserved
		if fID == 0 {
			continue
		}
		if err := feeder.validate(); err != nil {
			return err
		}
		// If Endblock is set, it should not be in the window of one round
		if feeder.EndBlock > 0 && (feeder.EndBlock-feeder.StartBaseBlock)%feeder.Interval < uint64(p.MaxNonce) {
			return ErrInvalidParams.Wrap("invalid tokenFeeder, invalid EndBlock")
		}
		// Interval should be long enough, make it more than twice pricing window of one round
		if feeder.Interval < 2*uint64(p.MaxNonce) {
			return ErrInvalidParams.Wrap("invalid tokenFeeder, invalid interval")
		}
		// cross validation with tokens
		if feeder.TokenID >= uint64(len(p.Tokens)) {
			return ErrInvalidParams.Wrap("invalid tokenFeeder, non-exist tokenID referred")
		}
		// cross validation with rules
		if feeder.RuleID >= uint64(len(p.Rules)) {
			return ErrInvalidParams.Wrap("invalid tokenFeeder, non-exist ruleID referred")
		}
	}

	// validate chain
	for cID, chain := range p.Chains {
		// id==0 is reserved
		if cID == 0 {
			continue
		}
		if err := chain.validate(); err != nil {
			return err
		}
	}

	// validate token
	for tID, token := range p.Tokens {
		// id==0 is reserved
		if tID == 0 {
			continue
		}
		if err := token.validate(); err != nil {
			return err
		}
		// cross validation with chain
		if token.ChainID >= uint64(len(p.Chains)) {
			return ErrInvalidParams.Wrap("invalid token, non-exist chainID referred")
		}
	}

	// validate rules
	for rID, rule := range p.Rules {
		if rID == 0 {
			continue
		}
		if err := rule.validate(); err != nil {
			return err
		}
		// cross validation with sources
		for _, id := range rule.SourceIDs {
			if id >= uint64(len(p.Rules)) {
				return ErrInvalidParams.Wrap("invalid rule")
			}
		}
		if rule.Nom != nil {
			for _, id := range rule.Nom.SourceIDs {
				if id < 1 || id >= uint64(len(p.Rules)) {
					return ErrInvalidParams.Wrap("invalid rule")
				}
			}
		}
	}
	// validete sources
	for sID, source := range p.Sources {
		if sID == 0 {
			continue
		}
		if err := source.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p Params) AddSources(sources ...*Source) (Params, error) {
	sNames := make(map[string]struct{})
	for _, source := range p.Sources {
		sNames[source.Name] = struct{}{}
	}
	for _, s := range sources {
		if !s.Valid {
			return p, ErrInvalidParams.Wrap("invalid source to add, new source should be valid")
		}
		if _, exists := sNames[s.Name]; exists {
			return p, ErrInvalidParams.Wrap("invalid source to add, duplicated")
		}
		sNames[s.Name] = struct{}{}
		p.Sources = append(p.Sources, s)
	}
	return p, nil
}

func (p Params) AddChains(chains ...*Chain) (Params, error) {
	cNames := make(map[string]struct{})
	for _, chain := range p.Chains {
		cNames[chain.Name] = struct{}{}
	}
	for _, c := range chains {
		if _, exists := cNames[c.Name]; exists {
			return p, ErrInvalidParams.Wrap("invalid chain to add, duplicated")
		}
		p.Chains = append(p.Chains, c)
	}
	return p, nil
}

func (p Params) UpdateTokens(tokens ...*Token) (Params, error) {
	for _, t := range tokens {
		update := false
		for _, token := range p.Tokens {
			if token.ChainID == t.ChainID && token.Name == t.Name {
				// AssetID is only field can be updated/modified
				if len(t.AssetID) > 0 {
					token.AssetID = t.AssetID
					update = true
					break
				}
				return p, ErrInvalidParams.Wrap("invalid token to add, duplicated")
			}
		}
		if !update {
			p.Tokens = append(p.Tokens, t)
		}
	}
	return p, nil
}

func (p Params) AddRules(rules ...*RuleSource) (Params, error) {
	p.Rules = append(p.Rules, rules...)
	return p, nil
}

func (p Params) UpdateTokenFeeder(tf *TokenFeeder, currentHeight uint64) (Params, error) {
	tfIDs := p.GetFeederIDsByTokenID(tf.TokenID)
	if len(tfIDs) == 0 {
		// first tokenfeeder for this token
		p.TokenFeeders = append(p.TokenFeeders, tf)
		return p, nil
	}
	tfIdx := tfIDs[len(tfIDs)-1]
	tokenFeeder := p.TokenFeeders[tfIdx]

	// latest feeder is not started yet
	if tokenFeeder.StartBaseBlock > currentHeight {
		// fields can be modified: startBaseBlock, interval, endBlock
		update := false
		if tf.StartBaseBlock > 0 {
			if tf.StartBaseBlock <= currentHeight {
				return p, ErrInvalidParams.Wrap("invalid tokenFeeder to update, invalid StartBaseBlock")
			}
			update = true
			tokenFeeder.StartBaseBlock = tf.StartBaseBlock
		}
		if tf.Interval > 0 {
			tokenFeeder.Interval = tf.Interval
			update = true
		}
		if tf.EndBlock > 0 {
			if tf.EndBlock <= currentHeight {
				return p, ErrInvalidParams.Wrap("invalid tokenFeeder to update, invalid EndBlock")
			}
			tokenFeeder.EndBlock = tf.EndBlock
			update = true
		}
		if !update {
			return p, ErrInvalidParams.Wrap("invalid tokenFeeder to update, no valid field set")
		}
		p.TokenFeeders[tfIdx] = tokenFeeder
		return p, nil
	}

	// latest feeder is running
	if tokenFeeder.EndBlock == 0 || tokenFeeder.EndBlock > currentHeight {
		// fields can be modified: endBlock
		if tf.EndBlock == 0 || tf.EndBlock <= currentHeight {
			return p, ErrInvalidParams.Wrap("invalid tokenFeeder to update, invalid EndBlock")
		}
		tokenFeeder.EndBlock = tf.EndBlock
		p.TokenFeeders[tfIdx] = tokenFeeder
		return p, nil
	}

	// latest feeder is stopped, this is actually a new feeder to resume the latest one for the same token
	latestRoundID := tokenFeeder.StartRoundID + (tokenFeeder.EndBlock-tokenFeeder.StartBaseBlock)/tokenFeeder.Interval
	if tf.StartBaseBlock <= currentHeight || tf.StartRoundID != latestRoundID+1 {
		return p, ErrInvalidParams.Wrap("invalid tokenFeeder to update")
	}
	p.TokenFeeders = append(p.TokenFeeders, tf)

	return p, nil
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

func (c Chain) validate() error {
	if len(c.Name) == 0 {
		return ErrInvalidParams.Wrap("invalid Chain")
	}
	return nil
}

func (t Token) validate() error {
	if len(t.Name) == 0 || t.ChainID < 1 {
		return ErrInvalidParams.Wrap("invalid Token")
	}
	return nil
}

func (r RuleSource) validate() error {
	if len(r.SourceIDs) == 0 && (r.Nom == nil || len(r.Nom.SourceIDs) == 0) {
		return ErrInvalidParams.Wrap("invalid RuleSource")
	}
	if r.Nom != nil &&
		r.Nom.Minimum > uint64(len(r.Nom.SourceIDs)) {
		return ErrInvalidParams.Wrap("invalid RuleSource")
	}
	return nil
}

func (s Source) validate() error {
	if len(s.Name) == 0 {
		return ErrInvalidParams.Wrap("invalid Source, duplicated name")
	}
	return nil
}

func (f TokenFeeder) validate() error {
	if f.TokenID < 1 ||
		f.StartRoundID < 1 ||
		f.Interval < 1 ||
		f.StartBaseBlock < 1 {
		return ErrInvalidParams.Wrap("invalid TokenFeeder")
	}

	if f.EndBlock > 0 && f.StartBaseBlock >= f.EndBlock {
		return ErrInvalidParams.Wrap("invalid TokenFeeder")
	}

	return nil
}
