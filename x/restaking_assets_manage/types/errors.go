// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:LGPL-3.0-only

package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrNoClientChainKey      = errorsmod.Register(ModuleName, 0, "there is no stored key for the input chain index")
	ErrNoClientChainAssetKey = errorsmod.Register(ModuleName, 1, "there is no stored key for the input assetId")

	ErrInputUpdateStateIsZero = errorsmod.Register(ModuleName, 2, "all of the input parameter value are zero")

	ErrSubDepositAmountIsMoreThanOrigin = errorsmod.Register(ModuleName, 3, "the staker deposit amount that want to decrease is more than the original state")

	ErrSubCanWithdrawAmountIsMoreThanOrigin = errorsmod.Register(ModuleName, 4, "the staker canWithdraw amount that want to decrease is more than the original state")

	ErrNoStakerAssetKey = errorsmod.Register(ModuleName, 5, "there is no stored key for the input staker and assetId")

	ErrSubOperatorTotalAmountIsMoreThanOrigin = errorsmod.Register(ModuleName, 6, "the operator total amount that want to decrease is more than the original state")

	ErrSubOperatorOwnAmountIsMoreThanOrigin = errorsmod.Register(ModuleName, 7, "the operator ownAmount that want to decrease is more than the original state")

	ErrNoOperatorAssetKey = errorsmod.Register(ModuleName, 8, "there is no stored key for the input operator address and assetId")
)
