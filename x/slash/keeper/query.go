package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/slash/types"
)

var _ types.QueryServer = Keeper{}
