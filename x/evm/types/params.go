package types

import (
	"github.com/ExocoreNetwork/exocore/utils"
	evmtype "github.com/evmos/evmos/v14/x/evm/types"
)

// ExocoreAvailableEVMExtensions defines the default active precompiles
var (
	// DefaultEVMDenom defines the default EVM denomination on Exocore
	DefaultEVMDenom               = utils.BaseDenom
	ExocoreAvailableEVMExtensions = []string{
		"0x0000000000000000000000000000000000000800", // Staking precompile
		"0x0000000000000000000000000000000000000801", // client chains precompile
		"0x0000000000000000000000000000000000000802", // ICS20 transfer precompile
		"0x0000000000000000000000000000000000000804", // deposit precompile
		"0x0000000000000000000000000000000000000805", // delegation precompile
		"0x0000000000000000000000000000000000000806", // reward precompile
		"0x0000000000000000000000000000000000000807", // slash precompile
		"0x0000000000000000000000000000000000000808", // withdraw precompile
		"0x0000000000000000000000000000000000000809", // bls precompile
	}
)

// ExocoreEvmDefaultParams returns default evm parameters
// ExtraEIPs is empty to prevent overriding the latest hard fork instruction set
func ExocoreEvmDefaultParams() evmtype.Params {
	return evmtype.Params{
		EvmDenom:            DefaultEVMDenom,
		EnableCreate:        evmtype.DefaultEnableCreate,
		EnableCall:          evmtype.DefaultEnableCall,
		ChainConfig:         evmtype.DefaultChainConfig(),
		ExtraEIPs:           nil,
		AllowUnprotectedTxs: evmtype.DefaultAllowUnprotectedTxs,
		ActivePrecompiles:   ExocoreAvailableEVMExtensions,
	}
}
