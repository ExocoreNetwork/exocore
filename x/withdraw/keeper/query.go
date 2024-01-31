package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/withdraw/types"
)

var _ types.QueryServer = Keeper{}
