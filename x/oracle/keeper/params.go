package keeper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	startAfterBlocks = 10
	defaultInterval  = 30
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

func (k Keeper) RegisterNewTokenAndSetTokenFeeder(ctx sdk.Context, chain, token, decimal, interval, contract, assetID string) error {
	p := k.GetParams(ctx)
	if p.GetTokenIDFromAssetID(assetID) > 0 {
		return fmt.Errorf("assetID exists:%s", assetID)
	}
	chainID := uint64(0)
	for id, c := range p.Chains {
		if c.Name == chain {
			chainID = uint64(id)
			break
		}
	}
	if chainID == 0 {
		// add new chain
		p.Chains = append(p.Chains, &types.Chain{
			Name: chain,
			Desc: "registered through assets module",
		})
		chainID = uint64(len(p.Chains) - 1)
	}
	decimalInt, err := strconv.ParseInt(decimal, 10, 32)
	if err != nil {
		return err
	}
	if decimalInt < 0 {
		return fmt.Errorf("decimal can't be negative:%d", decimalInt)
	}
	intervalInt, err := strconv.ParseUint(interval, 10, 64)
	if err != nil {
		return err
	}
	if intervalInt == 0 {
		intervalInt = defaultInterval
	}

	for _, t := range p.Tokens {
		// token exists, bind assetID for this token
		// it's possible for  one price bonded with multiple assetID, like ETHUSDT from sepolia/mainnet
		if t.Name == token && t.ChainID == chainID {
			t.AssetID = strings.Join([]string{t.AssetID, assetID}, ",")
			k.SetParams(ctx, p)
			// there should have been existing tokenFeeder running(currently we register tokens from assets-module and with infinite endBlock)
			return nil
		}
	}

	// add a new token
	p.Tokens = append(p.Tokens, &types.Token{
		Name:            token,
		ChainID:         chainID,
		ContractAddress: contract,
		Decimal:         int32(decimalInt),
		Active:          true,
		AssetID:         assetID,
	})
	// set a tokenFeeder for the new token
	p.TokenFeeders = append(p.TokenFeeders, &types.TokenFeeder{
		TokenID: uint64(len(p.Tokens) - 1),
		// we support rule_1 for v1
		RuleID:         1,
		StartRoundID:   1,
		StartBaseBlock: uint64(ctx.BlockHeight() + startAfterBlocks),
		Interval:       intervalInt,
		EndBlock:       0,
	})

	k.SetParams(ctx, p)
	return nil
}
