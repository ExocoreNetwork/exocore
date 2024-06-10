package bls

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls/blst"
	"github.com/prysmaticlabs/prysm/v4/crypto/bls/common"
)

const (
	// MethodFastAggregateVerify defines the ABI method name to fast verify aggregated signature and its
	// corresponding public keys
	MethodFastAggregateVerify = "fastAggregateVerify"
	// MethodVerify defines the ABI method name to verify aggregated signature and aggregated public key.
	MethodVerify              = "verify"
	MethodGeneratePrivateKey  = "generatePrivateKey"
	MethodPublicKey           = "publicKey"
	MethodSign                = "sign"
	MethodAggregatePubkeys    = "aggregatePubkeys"
	MethodAggregateSignatures = "aggregateSignatures"
	MethodAddTwoPubkeys       = "addTwoPubkeys"
)

// Verify checks the validity of an aggregated signature against msg and aggregated public keys.
func (p Precompile) Verify(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	sigBz, ok := args[1].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	sig, err := bls.SignatureFromBytes(sigBz)
	if err != nil {
		return nil, ErrInvalidArg
	}

	pubkeyBz, ok := args[2].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	pubkey, err := bls.PublicKeyFromBytes(pubkeyBz)
	if err != nil {
		return nil, ErrInvalidArg
	}

	msg, ok := args[0].([32]byte)
	if !ok {
		return nil, ErrInvalidArg
	}

	return method.Outputs.Pack(sig.Verify(pubkey, msg[:]))
}

// Verify checks the validity of an aggregated signature against msg and aggregated public keys.
func (p Precompile) FastAggregateVerify(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	sigBz, ok := args[1].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	sig, err := bls.SignatureFromBytes(sigBz)
	if err != nil {
		return nil, ErrInvalidArg
	}

	pubkeysBz, ok := args[2].([][]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	pubkeys := make([]common.PublicKey, len(pubkeysBz))
	for i, pubkeyBz := range pubkeysBz {
		pubkey, err := bls.PublicKeyFromBytes(pubkeyBz)
		if err != nil {
			return nil, ErrInvalidArg
		}
		pubkeys[i] = pubkey
	}

	msg, ok := args[0].([32]byte)
	if !ok {
		return nil, ErrInvalidArg
	}

	return method.Outputs.Pack(sig.FastAggregateVerify(pubkeys, msg))
}

func (p Precompile) GeneratePrivateKey(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	privkey, err := blst.RandKey()
	if err != nil {
		return nil, err
	}
	pri := privkey.Marshal()
	return method.Outputs.Pack(pri)
}

func (p Precompile) PublicKey(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	privkeyBz, ok := args[0].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	privkey, err := blst.SecretKeyFromBytes(privkeyBz)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(privkey.PublicKey().Marshal())
}

func (p Precompile) Sign(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	privkeyBz, ok := args[0].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	privkey, err := blst.SecretKeyFromBytes(privkeyBz)
	if err != nil {
		return nil, err
	}

	msg, ok := args[1].([32]byte)
	if !ok {
		return nil, ErrInvalidArg
	}

	return method.Outputs.Pack(privkey.Sign(msg[:]).Marshal())
}

func (p Precompile) AggregatePubkeys(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	pubkeysBz, ok := args[0].([][]byte)
	if !ok {
		return nil, ErrInvalidArg
	}

	aggregatedPubkey, err := blst.AggregatePublicKeys(pubkeysBz)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate public keys")
	}

	return method.Outputs.Pack(aggregatedPubkey.Marshal())
}

func (p Precompile) AggregateSignatures(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	sigsBz, ok := args[0].([][]byte)
	if !ok {
		return nil, ErrInvalidArg
	}

	aggregatedSig, err := blst.AggregateCompressedSignatures(sigsBz)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate public keys")
	}

	return method.Outputs.Pack(aggregatedSig.Marshal())
}

func (p Precompile) AddTwoPubkeys(
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	pubkeyOneBz, ok := args[0].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}
	pubkeyTwoBz, ok := args[1].([]byte)
	if !ok {
		return nil, ErrInvalidArg
	}

	pubkeyOne, err := blst.PublicKeyFromBytes(pubkeyOneBz)
	if err != nil {
		return nil, ErrInvalidArg
	}
	pubkeyTwo, err := blst.PublicKeyFromBytes(pubkeyTwoBz)
	if err != nil {
		return nil, ErrInvalidArg
	}
	newPubkey := pubkeyOne.Aggregate(pubkeyTwo)

	return method.Outputs.Pack(newPubkey.Marshal())
}
