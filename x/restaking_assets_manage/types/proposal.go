package types

import (
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// constants
const (
	ProposalTypeRegisterClientChain   string = "RegisterClientChain"
	ProposalTypeDeregisterClientChain string = "DeregisterClientChain"
	ProposalTypeRegisterAsset         string = "RegisterAsset"
	ProposalTypeDeregisterAsset       string = "DeregisterAsset"
)

// Implements Proposal Interface
var (
	_ govv1beta1.Content = &RegisterClientChainProposal{}
	_ govv1beta1.Content = &DeregisterClientChainProposal{}
	_ govv1beta1.Content = &RegisterAssetProposal{}
	_ govv1beta1.Content = &DeregisterAssetProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeRegisterClientChain)
	govv1beta1.RegisterProposalType(ProposalTypeDeregisterClientChain)
	govv1beta1.RegisterProposalType(ProposalTypeRegisterAsset)
	govv1beta1.RegisterProposalType(ProposalTypeDeregisterAsset)
}

// NewRegisterClientChainProposal returns new instance of RegisterClientChainProposal
func NewRegisterClientChainProposal(
	title, description, contract string,
	clientChain *ClientChainInfo,
) govv1beta1.Content {
	return &RegisterClientChainProposal{
		Title:       title,
		Description: description,
		ClientChain: clientChain,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterClientChainProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterClientChainProposal) ProposalType() string {
	return ProposalTypeRegisterClientChain
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *RegisterClientChainProposal) ValidateBasic() error {
	// todo: simply check the client chain info

	return govv1beta1.ValidateAbstract(rip)
}

// NewDeregisterClientChainProposal returns new instance of DeregisterClientChainProposal
func NewDeregisterClientChainProposal(
	title, description, contract string,
	clientChainID string,
) govv1beta1.Content {
	return &DeregisterClientChainProposal{
		Title:         title,
		Description:   description,
		ClientChainID: clientChainID,
	}
}

// ProposalRoute returns router key for this proposal
func (*DeregisterClientChainProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*DeregisterClientChainProposal) ProposalType() string {
	return ProposalTypeRegisterClientChain
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *DeregisterClientChainProposal) ValidateBasic() error {
	// todo: check the clientChainID

	return govv1beta1.ValidateAbstract(rip)
}

// NewRegisterAssetProposal returns new instance of NewRegisterAssetProposal
func NewRegisterAssetProposal(
	title, description, contract string,
	asset *ClientChainTokenInfo,
) govv1beta1.Content {
	return &RegisterAssetProposal{
		Title:       title,
		Description: description,
		Asset:       asset,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterAssetProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterAssetProposal) ProposalType() string {
	return ProposalTypeRegisterClientChain
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *RegisterAssetProposal) ValidateBasic() error {
	// todo: simply check the asset info

	return govv1beta1.ValidateAbstract(rip)
}

// NewDeregisterAssetProposal returns new instance of NewDeregisterAssetProposal
func NewDeregisterAssetProposal(
	title, description, contract string,
	assetID string,
) govv1beta1.Content {
	return &DeregisterAssetProposal{
		Title:       title,
		Description: description,
		AssetID:     assetID,
	}
}

// ProposalRoute returns router key for this proposal
func (*DeregisterAssetProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*DeregisterAssetProposal) ProposalType() string {
	return ProposalTypeRegisterClientChain
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *DeregisterAssetProposal) ValidateBasic() error {
	// todo: check the assetID

	return govv1beta1.ValidateAbstract(rip)
}
