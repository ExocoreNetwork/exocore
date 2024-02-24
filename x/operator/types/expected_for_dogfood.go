package types

import (
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SlashKeeper interface {
	IsOperatorFrozen(ctx sdk.Context, addr sdk.AccAddress) bool
}

type RedelegationKeeper interface {
	AppChainInfoIsExist(ctx sdk.Context, chainId string) bool
}

type OperatorConsentHooks interface {
	// This hook is called when an operator opts in to a chain.
	AfterOperatorOptIn(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainId string,
		pubKey tmprotocrypto.PublicKey,
	)
	// This hook is called when an operator's consensus key is replaced for
	// a chain.
	AfterOperatorKeyReplacement(
		ctx sdk.Context,
		addr sdk.AccAddress,
		oldKey tmprotocrypto.PublicKey,
		newKey tmprotocrypto.PublicKey,
		chainId string,
	)
	// This hook is called when an operator opts out of a chain.
	AfterOperatorOptOutInitiated(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainId string,
		key tmprotocrypto.PublicKey,
	)
}
