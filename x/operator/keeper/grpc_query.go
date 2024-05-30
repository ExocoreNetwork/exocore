package keeper

import (
	"context"
	"errors"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ operatortypes.QueryServer = &Keeper{}

// QueryOperatorInfo queries the operator information for the given address.
func (k *Keeper) QueryOperatorInfo(
	ctx context.Context, req *operatortypes.GetOperatorInfoReq,
) (*operatortypes.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.OperatorInfo(c, req.OperatorAddr)
}

// QueryAllOperators queries all operators on the chain.
func (k *Keeper) QueryAllOperators(
	goCtx context.Context, req *operatortypes.QueryAllOperatorsRequest,
) (*operatortypes.QueryAllOperatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	res := make([]string, 0)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorInfo)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, _ []byte) error {
		addr := sdk.AccAddress(key)
		res = append(res, addr.String())
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryAllOperatorsResponse{
		OperatorAccAddrs: res,
		Pagination:       pageRes,
	}, nil
}

// QueryOperatorConsKeyForChainID queries the consensus key for the operator on the given chain.
func (k *Keeper) QueryOperatorConsKeyForChainID(
	goCtx context.Context,
	req *operatortypes.QueryOperatorConsKeyRequest,
) (*operatortypes.QueryOperatorConsKeyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAccAddr)
	if err != nil {
		return nil, err
	}
	found, key, err := k.GetOperatorConsKeyForChainID(
		ctx, addr, req.Chain,
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
	addr, err := sdk.AccAddressFromBech32(req.OperatorAccAddr)
	if err != nil {
		return nil, err
	}
	found, key, err := k.GetOperatorConsKeyForChainID(
		ctx, addr, req.Chain,
	)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("no key assigned")
	}
	consAddr, err := operatortypes.TMCryptoPublicKeyToConsAddr(key)
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryOperatorConsAddressResponse{
		ConsAddr: consAddr.String(),
	}, nil
}

// QueryAllOperatorConsKeysByChainID queries all operators for the given chain and returns
// their consensus keys.
func (k Keeper) QueryAllOperatorConsKeysByChainID(
	goCtx context.Context,
	req *operatortypes.QueryAllOperatorConsKeysByChainIDRequest,
) (*operatortypes.QueryAllOperatorConsKeysByChainIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	res := make([]*operatortypes.OperatorConsKeyPair, 0)
	chainPrefix := operatortypes.ChainIDAndAddrKey(
		operatortypes.BytePrefixForChainIDAndOperatorToConsKey,
		req.Chain, nil,
	)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), chainPrefix)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		addr := sdk.AccAddress(key)
		ret := &tmprotocrypto.PublicKey{}
		// don't use MustUnmarshal to not panic for queries
		if err := ret.Unmarshal(value); err != nil {
			return err
		}
		res = append(res, &operatortypes.OperatorConsKeyPair{
			OperatorAccAddr: addr.String(),
			PublicKey:       ret,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryAllOperatorConsKeysByChainIDResponse{
		OperatorConsKeys: res,
		Pagination:       pageRes,
	}, nil
}

// QueryAllOperatorConsAddrsByChainID queries all operators for the given chain and returns
// their consensus addresses.
func (k Keeper) QueryAllOperatorConsAddrsByChainID(
	goCtx context.Context,
	req *operatortypes.QueryAllOperatorConsAddrsByChainIDRequest,
) (*operatortypes.QueryAllOperatorConsAddrsByChainIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	res := make([]*operatortypes.OperatorConsAddrPair, 0)
	chainPrefix := operatortypes.ChainIDAndAddrKey(
		operatortypes.BytePrefixForChainIDAndOperatorToConsKey,
		req.Chain, nil,
	)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), chainPrefix)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		addr := sdk.AccAddress(key)
		ret := &tmprotocrypto.PublicKey{}
		// don't use MustUnmarshal to not panic for queries
		if err := ret.Unmarshal(value); err != nil {
			return err
		}
		consAddr, err := operatortypes.TMCryptoPublicKeyToConsAddr(ret)
		if err != nil {
			return err
		}
		res = append(res, &operatortypes.OperatorConsAddrPair{
			OperatorAccAddr: addr.String(),
			ConsAddr:        consAddr.String(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryAllOperatorConsAddrsByChainIDResponse{
		OperatorConsAddrs: res,
		Pagination:        pageRes,
	}, nil
}
