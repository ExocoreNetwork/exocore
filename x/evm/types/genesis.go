// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	types2 "github.com/evmos/evmos/v14/x/evm/types"
)

// DefaultGenesisState sets default evm genesis state with empty accounts and default params and
// chain config values.
func DefaultGenesisState() *types2.GenesisState {
	return &types2.GenesisState{
		Accounts: []types2.GenesisAccount{},
		Params:   ExocoreEvmDefaultParams(),
	}
}
