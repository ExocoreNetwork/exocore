package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgClaimRewardRequest{}, "reward/ClaimRewardRequest", nil)
cdc.RegisterConcrete(&MsgRewardDetail{}, "reward/RewardDetail", nil)
cdc.RegisterConcrete(&MsgClaimRewardResponse{}, "reward/ClaimRewardResponse", nil)
// this line is used by starport scaffolding # 2
} 

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgClaimRewardRequest{},
)
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgRewardDetail{},
)
registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgClaimRewardResponse{},
)
// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
