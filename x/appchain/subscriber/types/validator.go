package types

import (
	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewSubscriberChainValidator creates a new SubscriberChainValidator instance.
func NewSubscriberChainValidator(
	address []byte, power int64, pubKey cryptotypes.PubKey,
) (SubscriberChainValidator, error) {
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return SubscriberChainValidator{}, err
	}

	return SubscriberChainValidator{
		ConsAddress: address,
		Power:       power,
		Pubkey:      pkAny,
	}, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces.
// It is required to ensure that ConsPubKey below works.
func (ocv SubscriberChainValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(ocv.Pubkey, &pk)
}

// ConsPubKey returns the validator PubKey as a cryptotypes.PubKey.
func (ocv SubscriberChainValidator) ConsPubKey() (cryptotypes.PubKey, error) {
	pk, ok := ocv.Pubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidType,
			"expecting cryptotypes.PubKey, got %T",
			pk,
		)
	}

	return pk, nil
}
