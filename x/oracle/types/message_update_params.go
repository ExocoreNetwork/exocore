package types

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/libs/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateParams = "update_params"

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic executes sanity validation on the provided data
// MsgUpdateParams is used to update params, the validation will mostly be stateful which is done by service
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}
	return nil
}

// GetSigners returns the expected signers for a MsgUpdateParams message
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func NewMsgUpdateParams(creator, paramsJSON string) *MsgUpdateParams {
	var p Params
	if err := json.Unmarshal([]byte(paramsJSON), &p); err != nil {
		panic("invalid json for params")
	}
	return &MsgUpdateParams{
		Authority: creator,
		Params:    p,
	}
}
