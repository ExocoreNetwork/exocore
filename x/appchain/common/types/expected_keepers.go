package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

type ClientKeeper interface {
	CreateClient(
		sdk.Context, ibcexported.ClientState, ibcexported.ConsensusState,
	) (string, error)
	GetClientState(sdk.Context, string) (ibcexported.ClientState, bool)
	GetLatestClientConsensusState(
		sdk.Context, string,
	) (ibcexported.ConsensusState, bool)
	GetSelfConsensusState(
		sdk.Context, ibcexported.Height,
	) (ibcexported.ConsensusState, error)
}
