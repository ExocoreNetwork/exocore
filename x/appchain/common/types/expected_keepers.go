package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
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

// ScopedKeeper defines the expected IBC capability keeper
type ScopedKeeper interface {
	GetCapability(sdk.Context, string) (*capabilitytypes.Capability, bool)
	AuthenticateCapability(sdk.Context, *capabilitytypes.Capability, string) bool
	ClaimCapability(sdk.Context, *capabilitytypes.Capability, string) error
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
}
