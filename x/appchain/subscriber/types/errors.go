package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/subscriber module sentinel errors
var (
	ErrNoProposerChannelID = errorsmod.Register(ModuleName, 2, "no established channel")
)
