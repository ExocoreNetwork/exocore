package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

var _ types.QueryServer = Keeper{}
