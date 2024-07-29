package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/operator/types"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
)

// This file indexes by chainID and not the avs address.
// The caller must ensure that the chainID is without the revision number.

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetOperatorConsKeyForChainID sets the (consensus) public key for the given operator address
// and chain id. If a key already exists, it will be overwritten and the edit will flow to the
// validator set at the next epoch.
// The caller must ensure that
// 1. Operator is opted in to the chain
// 2. The chain is registered with the AVS module
// 3. The wrappedKey is not nil
func (k *Keeper) SetOperatorConsKeyForChainID(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainID string,
	wrappedKey types.WrappedConsKey,
) error {
	return k.setOperatorConsKeyForChainID(ctx, opAccAddr, chainID, wrappedKey, false /* genesis */)
}

// setOperatorConsKeyForChainID is the private version of SetOperatorConsKeyForChainID.
// it is used with a boolean flag to indicate that the call is from genesis.
// if so, operator freeze status is not checked and hooks are not called.
func (k *Keeper) setOperatorConsKeyForChainID(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainID string,
	wrappedKey types.WrappedConsKey,
	genesis bool,
) error {
	// check for slashing
	if !genesis && k.slashKeeper.IsOperatorFrozen(ctx, opAccAddr) {
		return delegationtypes.ErrOperatorIsFrozen
	}
	/// in the process of opting out, do not allow key replacement
	if k.IsOperatorRemovingKeyFromChainID(ctx, opAccAddr, chainID) {
		return types.ErrAlreadyRemovingKey
	}
	// convert to bytes
	bz := k.cdc.MustMarshal(wrappedKey.ToTmKey())
	consAddr := wrappedKey.ToConsAddr()
	// check if the provided key is already in use by another operator. such use
	// also includes whether it was replaced by the same operator. this check ensures
	// that a key that has been replaced cannot be used again until it matures.
	// even if it is the same operator, do not allow calling this function twice.
	keyInUse, _ := k.GetOperatorAddressForChainIDAndConsAddr(ctx, chainID, consAddr)
	if keyInUse {
		return types.ErrConsKeyAlreadyInUse
	}
	// check that such a key is already set. if yes, we will consider it as key replacement.
	found, prevKey := k.getOperatorConsKeyForChainID(ctx, opAccAddr, chainID)
	var alreadyRecorded bool
	if found {
		// ultimately performs bytes.Equal
		if prevKey.Equals(wrappedKey) {
			// no-op
			return nil
		}
		// if this key is different, we will set the vote power of the old key to 0
		// in the validator update. but, we must only do so once in a block, since the
		// first existing key is the one to replace with 0 vote power and not any others.
		alreadyRecorded, _ = k.getOperatorPrevConsKeyForChainID(ctx, opAccAddr, chainID)
		if !alreadyRecorded {
			k.setOperatorPrevConsKeyForChainID(
				ctx, opAccAddr, chainID, prevKey,
			)
		}
	}
	k.setOperatorConsKeyForChainIDUnchecked(ctx, opAccAddr, consAddr, chainID, bz)
	// only call the hooks if this is not genesis
	if !genesis {
		if found {
			if !alreadyRecorded {
				k.Hooks().AfterOperatorKeyReplaced(ctx, opAccAddr, prevKey, wrappedKey, chainID)
			}
		} else {
			k.Hooks().AfterOperatorKeySet(ctx, opAccAddr, chainID, wrappedKey)
		}
	}
	return nil
}

// setOperatorConsKeyForChainIDUnchecked is the internal private version. It performs
// no error checking of the input. The caller must do the error checking
// and then call this function.
func (k Keeper) setOperatorConsKeyForChainIDUnchecked(
	ctx sdk.Context, opAccAddr sdk.AccAddress, consAddr sdk.ConsAddress,
	chainID string, bz []byte,
) {
	store := ctx.KVStore(k.storeKey)
	// forward lookup
	// given operator address and chain id, find the consensus key,
	// since it is sorted by operator address, it helps for faster indexing by operator
	// for example, when an operator is delegated to, we can find all impacted
	// chain ids and their respective consensus keys
	store.Set(types.KeyForOperatorAndChainIDToConsKey(opAccAddr, chainID), bz)
	// reverse lookups
	// 1. given chain id and operator address, find the consensus key,
	// at initial onboarding of an app chain, it will allow us to find all
	// operators that have opted in and their consensus keys
	store.Set(types.KeyForChainIDAndOperatorToConsKey(chainID, opAccAddr), bz)
	// 2. given a chain id and a consensus addr, find the operator address,
	// the slashing module asks for an operator to be slashed by their consensus
	// address, so this will allow us to find the operator address to slash.
	// however, we do not want to retain this information forever, so we will
	// prune it once the validator set update id matures (if key replacement).
	// this pruning will be triggered by the app chain module and will not be
	// recorded here.
	store.Set(types.KeyForChainIDAndConsKeyToOperator(chainID, consAddr), opAccAddr.Bytes())
}

// setOperatorPrevConsKeyForChainID sets the previous (consensus) public key for the given
// operator address and chain id. This is used to track the previous key when a key is replaced.
// It is internal-only because such a key must only be set upon key replacement. So it does
// not perform any meaningful error checking of the input beyond its ability to be marshaled.
func (k *Keeper) setOperatorPrevConsKeyForChainID(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainID string,
	prevKey types.WrappedConsKey,
) {
	bz := k.cdc.MustMarshal(prevKey.ToTmKey())
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForChainIDAndOperatorToPrevConsKey(chainID, opAccAddr), bz)
}

// GetOperatorPrevConsKeyForChainID gets the previous (consensus) public key for the given
// operator address and chain id. When such a key is returned, callers should set its vote power
// to 0 in the validator update. It checks whether the chainID is registered in the AVS module
// and whether the operator is registered in this module.
func (k *Keeper) GetOperatorPrevConsKeyForChainID(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainID string,
) (bool, types.WrappedConsKey, error) {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		return false, nil, delegationtypes.ErrOperatorNotExist
	}
	// check if the chain exists as an AVS
	if isAvs, _ := k.avsKeeper.IsAVSByChainID(ctx, chainID); !isAvs {
		return false, nil, types.ErrUnknownChainID
	}
	found, key := k.getOperatorPrevConsKeyForChainID(ctx, opAccAddr, chainID)
	return found, key, nil
}

// getOperatorPrevConsKeyForChainID is the internal version of GetOperatorPrevConsKeyForChainID.
// It performs no error checking of the input.
func (k *Keeper) getOperatorPrevConsKeyForChainID(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainID string,
) (bool, types.WrappedConsKey) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.KeyForChainIDAndOperatorToPrevConsKey(chainID, opAccAddr))
	if res == nil {
		return false, nil
	}
	key := &tmprotocrypto.PublicKey{}
	k.cdc.MustUnmarshal(res, key)
	return true, types.NewWrappedConsKeyFromTmKey(key)
}

// GetOperatorConsKeyForChainID gets the (consensus) public key for the given operator address
// and chain id. This should be exposed via the query surface. If there is no such key,
// false and a nil key are returned.
func (k Keeper) GetOperatorConsKeyForChainID(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainID string,
) (bool, types.WrappedConsKey, error) {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		return false, nil, delegationtypes.ErrOperatorNotExist
	}
	// check if the chain exists as an AVS
	if isAvs, _ := k.avsKeeper.IsAVSByChainID(ctx, chainID); !isAvs {
		return false, nil, types.ErrUnknownChainID
	}
	found, key := k.getOperatorConsKeyForChainID(ctx, opAccAddr, chainID)
	return found, key, nil
}

// getOperatorConsKeyForChainID is the internal version of GetOperatorConsKeyForChainID. It
// performs no error checking of the input.
func (k *Keeper) getOperatorConsKeyForChainID(
	ctx sdk.Context,
	opAccAddr sdk.AccAddress,
	chainID string,
) (bool, types.WrappedConsKey) {
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.KeyForOperatorAndChainIDToConsKey(opAccAddr, chainID))
	if res == nil {
		return false, nil
	}
	key := &tmprotocrypto.PublicKey{}
	k.cdc.MustUnmarshal(res, key)
	return true, types.NewWrappedConsKeyFromTmKey(key)
}

// GetOperatorAddressForChainIDAndConsAddr returns the operator address for the given chain id
// and consensus address. This is used during slashing to find the operator address to slash.
func (k Keeper) GetOperatorAddressForChainIDAndConsAddr(
	ctx sdk.Context, chainID string, consAddr sdk.ConsAddress,
) (bool, sdk.AccAddress) {
	// check if the chain exists as an AVS
	if isAvs, _ := k.avsKeeper.IsAVSByChainID(ctx, chainID); !isAvs {
		return false, nil
	}
	store := ctx.KVStore(k.storeKey)
	res := store.Get(types.KeyForChainIDAndConsKeyToOperator(chainID, consAddr))
	if res == nil {
		return false, sdk.AccAddress{}
	}
	return true, sdk.AccAddress(res)
}

// InitiateOperatorKeyRemovalForChainID initiates an operator removing their key from the
// chain id. The caller must validate that the chainID is registered and that the address
// is an operator, that is not frozen, and that the operator is currently opted in.
func (k *Keeper) InitiateOperatorKeyRemovalForChainID(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainID string,
) {
	// found will always be true, since the operator has registered into the chain
	// and during registration a key must be set.
	_, key := k.getOperatorConsKeyForChainID(ctx, opAccAddr, chainID)
	// we don't check if the operator is already opted out, because this function
	// can only be called if the operator is currently opted in.
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForOperatorKeyRemovalForChainID(opAccAddr, chainID), []byte{})
	k.Hooks().AfterOperatorKeyRemovalInitiated(ctx, opAccAddr, chainID, key)
}

// IsOperatorRemovingKeyFromChainID returns true if the operator is removing the consensus
// key from the given chain id.
func (k Keeper) IsOperatorRemovingKeyFromChainID(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainID string,
) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForOperatorKeyRemovalForChainID(opAccAddr, chainID))
	return bz != nil
}

// CompleteOperatorKeyRemovalForChainID completes the operator key removal from the given
// chain id.
func (k Keeper) CompleteOperatorKeyRemovalForChainID(
	ctx sdk.Context, opAccAddr sdk.AccAddress, chainID string,
) error {
	// check if we are an operator
	if !k.IsOperator(ctx, opAccAddr) {
		return delegationtypes.ErrOperatorNotExist
	}
	// validate chain id
	avsAddr := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
	if isAvs, _ := k.avsKeeper.IsAVS(ctx, avsAddr); !isAvs {
		return types.ErrUnknownChainID
	}
	// check if the operator is opting out as we speak
	if !k.IsOperatorRemovingKeyFromChainID(ctx, opAccAddr, chainID) {
		return types.ErrOperatorNotRemovingKey
	}
	store := ctx.KVStore(k.storeKey)
	// get previous key to calculate consensus address
	_, prevKey := k.getOperatorConsKeyForChainID(ctx, opAccAddr, chainID)
	consAddr := prevKey.ToConsAddr()
	store.Delete(types.KeyForOperatorAndChainIDToConsKey(opAccAddr, chainID))
	store.Delete(types.KeyForChainIDAndOperatorToConsKey(chainID, opAccAddr))
	store.Delete(types.KeyForChainIDAndConsKeyToOperator(chainID, consAddr))
	store.Delete(types.KeyForOperatorKeyRemovalForChainID(opAccAddr, chainID))
	return nil
}

// GetOperatorsForChainID returns a list of {operatorAddr, pubKey} for the given
// chainID. This is used to create or update the validator set. It includes
// jailed operators, frozen operators and those in the process of opting out.
func (k *Keeper) GetOperatorsForChainID(
	ctx sdk.Context, chainID string,
) ([]sdk.AccAddress, []types.WrappedConsKey) {
	avsAddr := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
	if isAvs, _ := k.avsKeeper.IsAVS(ctx, avsAddr); !isAvs {
		return nil, nil
	}
	// prefix is the byte prefix and then chainID with length
	prefix := types.ChainIDAndAddrKey(
		types.BytePrefixForChainIDAndOperatorToConsKey,
		chainID, nil,
	)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(
		store, prefix,
	)
	defer iterator.Close()
	var addrs []sdk.AccAddress
	var pubKeys []types.WrappedConsKey
	for ; iterator.Valid(); iterator.Next() {
		// this key is of the format prefix | len | chainID | addr
		// and our prefix is of the format prefix | len | chainID
		// so just drop it and convert to sdk.AccAddress
		addr := iterator.Key()[len(prefix):]
		res := iterator.Value()
		ret := &tmprotocrypto.PublicKey{}
		k.cdc.MustUnmarshal(res, ret)
		addrs = append(addrs, addr)
		pubKeys = append(pubKeys, types.NewWrappedConsKeyFromTmKey(ret))
	}
	return addrs, pubKeys
}

// GetActiveOperatorsForChainID should return a list of operators and their public keys.
// These operators are neither jailed, nor frozen, nor opted out, and nor in the process
// of doing so.
func (k Keeper) GetActiveOperatorsForChainID(
	ctx sdk.Context, chainID string,
) ([]sdk.AccAddress, []types.WrappedConsKey) {
	operatorsAddr, pks := k.GetOperatorsForChainID(ctx, chainID)
	avsAddr := k.avsKeeper.GetAVSAddrByChainID(ctx, chainID)
	activeOperator := make([]sdk.AccAddress, 0)
	activePks := make([]types.WrappedConsKey, 0)
	// check if the operator is active
	for i, operator := range operatorsAddr {
		if k.IsActive(ctx, operator, avsAddr) {
			activeOperator = append(activeOperator, operator)
			activePks = append(activePks, pks[i])
		}
	}
	return activeOperator, activePks
}

// ValidatorByConsAddrForChainID returns a stakingtypes.ValidatorI for the given consensus
// address and chain id.
func (k Keeper) ValidatorByConsAddrForChainID(
	ctx sdk.Context, consAddr sdk.ConsAddress, chainID string,
) (stakingtypes.Validator, bool) {
	found, operatorAddr := k.GetOperatorAddressForChainIDAndConsAddr(
		ctx, chainID, consAddr,
	)
	if !found {
		// TODO(mm): create tests for the case where a validator
		// changes their pub key in the middle of an epoch, resulting
		// in a different key. work around that issue here.
		return stakingtypes.Validator{}, false
	}
	found, wrappedKey, err := k.GetOperatorConsKeyForChainID(ctx, operatorAddr, chainID)
	if !found || err != nil {
		return stakingtypes.Validator{}, false
	}
	// since we are sending the address, we have to send the consensus key as well.
	// this is because the presence of a non-empty address triggers a call to Validator
	// which triggers a call to fetch the consensus key, in the slashing module.
	val, err := stakingtypes.NewValidator(
		sdk.ValAddress(operatorAddr), wrappedKey.ToSdkKey(), stakingtypes.Description{},
	)
	if err != nil {
		return stakingtypes.Validator{}, false
	}
	val.Jailed = k.IsOperatorJailedForChainID(ctx, consAddr, chainID)
	return val, true
}

// DeleteOperatorAddressForChainIDAndConsAddr is a pruning method used to delete the
// mapping from chain id and consensus address to operator address. This mapping is used
// to obtain the operator address from its consensus public key when slashing.
func (k Keeper) DeleteOperatorAddressForChainIDAndConsAddr(
	ctx sdk.Context, chainID string, consAddr sdk.ConsAddress,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForChainIDAndConsKeyToOperator(chainID, consAddr))
}

// ClearPreviousConsensusKeys clears the previous consensus public key for all operators
// of the specified chain.
func (k Keeper) ClearPreviousConsensusKeys(ctx sdk.Context, chainID string) {
	partialKey := types.ChainIDWithLenKey(chainID)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(
		store,
		types.AppendMany(
			[]byte{types.BytePrefixForOperatorAndChainIDToPrevConsKey},
			partialKey,
		),
	)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
