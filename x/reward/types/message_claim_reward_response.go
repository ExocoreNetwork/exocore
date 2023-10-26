package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimRewardResponse = "claim_reward_response"

var _ sdk.Msg = &MsgClaimRewardResponse{}

func NewMsgClaimRewardResponse(creator string) *MsgClaimRewardResponse {
  return &MsgClaimRewardResponse{
		Creator: creator,
	}
}

func (msg *MsgClaimRewardResponse) Route() string {
  return RouterKey
}

func (msg *MsgClaimRewardResponse) Type() string {
  return TypeMsgClaimRewardResponse
}

func (msg *MsgClaimRewardResponse) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgClaimRewardResponse) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimRewardResponse) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

