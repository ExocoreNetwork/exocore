package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey) // return types.NewParams()
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &params)
	}
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	// TODO: validation check
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}

//func (k Keeper) SetUpdateParams(ctx sdk.Context, params types.Params) {
//	// validation has been done, keeper just update params
//	p := k.GetParams(ctx)
//
//	// update chains
//	for _, c := range params.Chains {
//		p.Chains = append(p.Chains, c)
//	}
//
//	// update tokens
//	for _, t := range params.Tokens {
//		p.Tokens = append(p.Tokens, t)
//	}
//
//	// update  sources
//	for _, s := range params.Sources {
//		if !s.Valid {
//			sID := p.GetSourceIDByName(s.Name)
//			p.Sources[sID].Valid = s.Valid
//		} else {
//			p.Sources = append(p.Sources, s)
//		}
//	}
//
//	// update rules
//	for _, r := range params.Rules {
//		p.Rules = append(p.Rules, r)
//	}
//
//	//update tokenFeeder
//	for _, f := range params.TokenFeeders {
//		if fID := p.GetFeederIDsByTokenID(f.TokenID); fID > 0 {
//			feeder := p.TokenFeeders[fID]
//			// all fiels has been verified before update params, just set all values
//			if f.StartBaseBlock > 0 {
//				feeder.StartBaseBlock = f.StartBaseBlock
//			}
//			if f.EndBlock > 0 {
//				feeder.EndBlock = f.EndBlock
//			}
//		} else {
//			p.TokenFeeders = append(p.TokenFeeders, f)
//		}
//	}
//	k.SetParams(ctx, p)
//}
