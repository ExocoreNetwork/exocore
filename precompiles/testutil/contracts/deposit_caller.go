package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"

	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

var (
	//go:embed DepositCaller.json
	DepositCallerJSON []byte

	// DepositCallerContract is the compiled contract calling the deposit precompile
	DepositCallerContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(DepositCallerJSON, &DepositCallerContract)
	if err != nil {
		panic(err)
	}

	if len(DepositCallerContract.Bin) == 0 {
		panic("failed to load smart contract that calls deposit precompile")
	}
}
