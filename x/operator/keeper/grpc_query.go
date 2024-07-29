package keeper

import (
	"context"
	"errors"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

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
		OptingOut: k.IsOperatorRemovingKeyFromChainID(ctx, addr, req.Chain),
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
		ConsAddr:  consAddr.String(),
		OptingOut: k.IsOperatorRemovingKeyFromChainID(ctx, addr, req.Chain),
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
			OptingOut:       k.IsOperatorRemovingKeyFromChainID(ctx, addr, req.Chain),
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
			OptingOut:       k.IsOperatorRemovingKeyFromChainID(ctx, addr, req.Chain),
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

func (k *Keeper) QueryOperatorUSDValue(ctx context.Context, req *operatortypes.QueryOperatorUSDValueRequest) (*operatortypes.QueryOperatorUSDValueResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	optedUSDValues, err := k.GetOperatorOptedUSDValue(c, req.Details.AVSAddress, req.Details.OperatorAddr)
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryOperatorUSDValueResponse{
		USDValues: &optedUSDValues,
	}, nil
}

func (k *Keeper) QueryAVSUSDValue(ctx context.Context, req *operatortypes.QueryAVSUSDValueRequest) (*operatortypes.DecValueField, error) {
	c := sdk.UnwrapSDKContext(ctx)
	usdValue, err := k.GetAVSUSDValue(c, req.AVSAddress)
	if err != nil {
		return nil, err
	}
	return &operatortypes.DecValueField{
		Amount: usdValue,
	}, nil
}

func (k *Keeper) QueryOperatorSlashInfo(goCtx context.Context, req *operatortypes.QueryOperatorSlashInfoRequest) (*operatortypes.QueryOperatorSlashInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	res := make([]*operatortypes.OperatorSlashInfoByID, 0)

	slashPrefix := operatortypes.AppendMany(operatortypes.KeyPrefixOperatorSlashInfo, assetstype.GetJoinedStoreKeyForPrefix(req.Details.OperatorAddr, req.Details.AVSAddress))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), slashPrefix)
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		ret := &operatortypes.OperatorSlashInfo{}
		// don't use MustUnmarshal to not panic for queries
		if err := ret.Unmarshal(value); err != nil {
			return err
		}

		res = append(res, &operatortypes.OperatorSlashInfoByID{
			SlashID: string(key),
			Info:    ret,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operatortypes.QueryOperatorSlashInfoResponse{
		AllSlashInfo: res,
		Pagination:   pageRes,
	}, nil
}
