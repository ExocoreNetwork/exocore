package types

import (
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	"golang.org/x/net/context"
)

// EpochsKeeper represents the expected keeper interface for the epochs module.
type EpochsKeeper interface {
	GetEpochInfo(sdk.Context, string) (epochstypes.EpochInfo, bool)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// EVMKeeper defines the expected EVM keeper interface used on erc20
type EVMKeeper interface {
	SetAccount(ctx sdk.Context, addr common.Address, account statedb.Account) error
	SetCode(ctx sdk.Context, codeHash, code []byte)
	GetCode(ctx sdk.Context, codeHash common.Hash) []byte
	GetParams(ctx sdk.Context) evmtypes.Params
	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
}

// OperatorKeeper represents the expected keeper interface for the operator module.
type OperatorKeeper interface {
	IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool
	OptIn(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) error
	OptOut(ctx sdk.Context, operatorAddress sdk.AccAddress, avsAddr string) (err error)
}

// AssetsKeeper represents the expected keeper interface for the assets module.
type AssetsKeeper interface {
	GetStakingAssetInfo(
		ctx sdk.Context, assetID string,
	) (info *assetstype.StakingAssetInfo, err error)
}
