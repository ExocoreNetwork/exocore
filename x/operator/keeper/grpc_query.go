package keeper

import (
	"context"
	"errors"

	"github.com/ExocoreNetwork/exocore/x/operator/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ operatortypes.QueryServer = &Keeper{}

func (k *Keeper) QueryOperatorInfo(
	ctx context.Context, req *operatortypes.GetOperatorInfoReq,
) (*operatortypes.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.OperatorInfo(c, req.OperatorAddr)
}

// QueryOperatorConsKeyForChainID queries the consensus key for the operator on the given chain.
func (k *Keeper) QueryOperatorConsKeyForChainID(
	goCtx context.Context,
	req *operatortypes.QueryOperatorConsKeyRequest,
) (*operatortypes.QueryOperatorConsKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Addr)
	if err != nil {
		return nil, err
	}
	found, key, err := k.GetOperatorConsKeyForChainID(
		ctx, addr, req.ChainID,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("no key assigned")
	}
	return &operatortypes.QueryOperatorConsKeyResponse{
		PublicKey: *key,
	}, nil
}

// QueryOperatorConsAddressForChainID queries the consensus address for the operator on
// the given chain.
func (k Keeper) QueryOperatorConsAddressForChainID(
	goCtx context.Context,
	req *operatortypes.QueryOperatorConsAddressRequest,
) (*operatortypes.QueryOperatorConsAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Addr)
	if err != nil {
		return nil, err
	}
	found, key, err := k.GetOperatorConsKeyForChainID(
		ctx, addr, req.ChainID,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("no key assigned")
	}
	consAddr, err := types.TMCryptoPublicKeyToConsAddr(key)
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryOperatorConsAddressResponse{
		Address: consAddr.String(),
	}, nil
}

// QueryAllOperatorKeysByChainID queries all operators for the given chain.
func (k Keeper) QueryAllOperatorKeysByChainID(
	goCtx context.Context,
	req *operatortypes.QueryAllOperatorsByChainIDRequest,
) (*operatortypes.QueryAllOperatorsByChainIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	res := make([]*operatortypes.OperatorConsKeyPair, 0)
	chainPrefix := operatortypes.ChainIDAndAddrKey(
		operatortypes.BytePrefixForChainIDAndOperatorToConsKey,
		req.ChainID, nil,
	)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), chainPrefix)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		addr := sdk.AccAddress(key[:])
		ret := &tmprotocrypto.PublicKey{}
		// don't use MustUnmarshal to not panic for queries
		if err := ret.Unmarshal(value); err != nil {
			return err
		}
		res = append(res, &operatortypes.OperatorConsKeyPair{
			OperatorAddr: addr.String(),
			PublicKey:    ret,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryAllOperatorsByChainIDResponse{
		OperatorConsKeys: res,
		Pagination:       pageRes,
	}, nil
}
