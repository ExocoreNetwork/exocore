package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/exomint/types"
)

var _ types.QueryServer = Keeper{}
