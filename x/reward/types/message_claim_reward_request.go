package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimRewardRequest = "claim_reward_request"

var _ sdk.Msg = &MsgClaimRewardRequest{}

func NewMsgClaimRewardRequest(creator string, id uint64, rewardaddress string) *MsgClaimRewardRequest {
	return &MsgClaimRewardRequest{
		Creator:       creator,
		Id:            id,
		Rewardaddress: rewardaddress,
	}
}

func (msg *MsgClaimRewardRequest) Route() string {
	return RouterKey
}

func (msg *MsgClaimRewardRequest) Type() string {
	return TypeMsgClaimRewardRequest
}

func (msg *MsgClaimRewardRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimRewardRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimRewardRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
