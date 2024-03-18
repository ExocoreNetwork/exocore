package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

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
type AvshKeeper struct{}

func (k AvshKeeper) IsAVS(sdk.Context, sdk.AccAddress) bool {
	return true
}

func (k AvshKeeper) getAvsAddress(sdk.Context, sdk.AccAddress) string {
	return "0x00000000000000000000000000000000000009999"
}
