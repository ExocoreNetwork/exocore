package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoKeyInTheStore = errorsmod.Register(
		ModuleName, 2,
		"there is not the key for in the store",
	)
	ErrOperatorIsFrozen = errorsmod.Register(
		ModuleName, 3,
		"the operator has been frozen",
	)

	ErrOperatorNotExist = errorsmod.Register(
		ModuleName, 4,
		"the operator has not been registered",
	)

	ErrAmountIsNotPositive = errorsmod.Register(
		ModuleName, 5,
		"the amount isn't positive",
	)

	OperatorAddrIsNotAccAddr = errorsmod.Register(
		ModuleName, 6,
		"the operator address isn't a valid acc addr",
	)

	ErrSubAmountIsGreaterThanOriginal = errorsmod.Register(
		ModuleName, 7,
		"the sub amount is greater than the original amount",
	)

	ErrParseDelegationKey = errorsmod.Register(
		ModuleName, 8,
		"delegation state key can't be parsed",
	)

	ErrStakerGetRecordType = errorsmod.Register(
		ModuleName, 9,
		"the input getType is error when get staker Undelegation records",
	)

	ErrUndelegationAmountTooBig = errorsmod.Register(
		ModuleName, 10,
		"the Undelegation amount is bigger than the delegated amount",
	)

	ErrNotSupportYet = errorsmod.Register(
		ModuleName, 11,
		"don't have supported it yet",
	)

	ErrDelegationAmountTooBig = errorsmod.Register(
		ModuleName, 12,
		"the delegation amount is bigger than the canWithdraw amount",
	)

	ErrCannotIncHoldCount = errorsmod.Register(
		ModuleName, 13,
		"cannot increment undelegation hold count above max uint64",
	)

	ErrCannotDecHoldCount = errorsmod.Register(
		ModuleName, 14,
		"cannot decrement undelegation hold count below zero",
	)

	ErrInvalidGenesisData = errorsmod.Register(
		ModuleName, 15,
		"the genesis data supplied is invalid",
	)

	ErrDivisorIsZero = errorsmod.Register(
		ModuleName, 16,
		"the divisor is zero")

	ErrInsufficientShares = errorsmod.Register(
		ModuleName, 17,
		"insufficient delegation shares")

	ErrInsufficientAssetAmount = errorsmod.Register(
		ModuleName, 18,
		"insufficient asset amount")

	ErrInvalidHash = errorsmod.Register(
		ModuleName, 19,
		"invalid hash",
	)

	ErrOperatorAlreadyAssociated = errorsmod.Register(
		ModuleName, 20,
		"the operator is already associated by this staker",
	)

	ErrNoAssociatedOperatorByStaker = errorsmod.Register(
		ModuleName, 21,
		"there isn't any operator marked by the staker",
	)
)
