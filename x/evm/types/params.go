// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	types2 "github.com/evmos/evmos/v14/x/evm/types"
)

var (
	// ExocoreAvailableEVMExtensions defines the default active precompiles
	ExocoreAvailableEVMExtensions = []string{
		"0x0000000000000000000000000000000000000800", // Staking precompile
		"0x0000000000000000000000000000000000000801", // Distribution precompile
		"0x0000000000000000000000000000000000000802", // ICS20 transfer precompile
		"0x0000000000000000000000000000000000000803", // Vesting precompile
		"0x0000000000000000000000000000000000000804", // deposit precompile
		"0x0000000000000000000000000000000000000805", // delegation precompile
		"0x0000000000000000000000000000000000000806", // reward precompile
		"0x0000000000000000000000000000000000000807", // slash precompile
		"0x0000000000000000000000000000000000000808", // withdraw precompile
	}
)

// ExocoreEvmDefaultParams returns default evm parameters
// ExtraEIPs is empty to prevent overriding the latest hard fork instruction set
// ActivePrecompiles is empty to prevent overriding the default precompiles
// from the EVM configuration.
func ExocoreEvmDefaultParams() types2.Params {
	return types2.Params{
		EvmDenom:            types2.DefaultEVMDenom,
		EnableCreate:        types2.DefaultEnableCreate,
		EnableCall:          types2.DefaultEnableCall,
		ChainConfig:         types2.DefaultChainConfig(),
		ExtraEIPs:           nil,
		AllowUnprotectedTxs: types2.DefaultAllowUnprotectedTxs,
		ActivePrecompiles:   ExocoreAvailableEVMExtensions,
	}
}
