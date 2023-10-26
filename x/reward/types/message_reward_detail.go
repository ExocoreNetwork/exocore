package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRewardDetail = "reward_detail"

var _ sdk.Msg = &MsgRewardDetail{}

func NewMsgRewardDetail(creator string, id uint64) *MsgRewardDetail {
  return &MsgRewardDetail{
		Creator: creator,
    Id: id,
	}
}

func (msg *MsgRewardDetail) Route() string {
  return RouterKey
}

func (msg *MsgRewardDetail) Type() string {
  return TypeMsgRewardDetail
}

func (msg *MsgRewardDetail) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgRewardDetail) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgRewardDetail) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

