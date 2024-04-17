// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: exocore/oracle/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgCreatePrice provide the price updating message
type MsgCreatePrice struct {
	// creator tells which is the message sender and should sign this message
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	//refer to id from Params.TokenFeeders, 0 is reserved, invalid to use
	FeederID uint64 `protobuf:"varint,2,opt,name=feeder_id,json=feederId,proto3" json:"feeder_id,omitempty"`
	// prices price with its corresponding source
	Prices []*PriceWithSource `protobuf:"bytes,3,rep,name=prices,proto3" json:"prices,omitempty"`
	//on which block commit does this message be built on
	BasedBlock uint64 `protobuf:"varint,4,opt,name=based_block,json=basedBlock,proto3" json:"based_block,omitempty"`
	// nonce represents the unique number to disginguish duplicated messages
	Nonce int32 `protobuf:"varint,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *MsgCreatePrice) Reset()         { *m = MsgCreatePrice{} }
func (m *MsgCreatePrice) String() string { return proto.CompactTextString(m) }
func (*MsgCreatePrice) ProtoMessage()    {}
func (*MsgCreatePrice) Descriptor() ([]byte, []int) {
	return fileDescriptor_02cf64aff79d2288, []int{0}
}
func (m *MsgCreatePrice) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreatePrice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreatePrice.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreatePrice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreatePrice.Merge(m, src)
}
func (m *MsgCreatePrice) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreatePrice) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreatePrice.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreatePrice proto.InternalMessageInfo

func (m *MsgCreatePrice) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgCreatePrice) GetFeederID() uint64 {
	if m != nil {
		return m.FeederID
	}
	return 0
}

func (m *MsgCreatePrice) GetPrices() []*PriceWithSource {
	if m != nil {
		return m.Prices
	}
	return nil
}

func (m *MsgCreatePrice) GetBasedBlock() uint64 {
	if m != nil {
		return m.BasedBlock
	}
	return 0
}

func (m *MsgCreatePrice) GetNonce() int32 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

// MsgCreatePriceResponse
type MsgCreatePriceResponse struct {
}

func (m *MsgCreatePriceResponse) Reset()         { *m = MsgCreatePriceResponse{} }
func (m *MsgCreatePriceResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreatePriceResponse) ProtoMessage()    {}
func (*MsgCreatePriceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_02cf64aff79d2288, []int{1}
}
func (m *MsgCreatePriceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreatePriceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreatePriceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreatePriceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreatePriceResponse.Merge(m, src)
}
func (m *MsgCreatePriceResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreatePriceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreatePriceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreatePriceResponse proto.InternalMessageInfo

// MsgUpdateParms
type MsgUpdateParams struct {
	// authority is the address that controls the module (defaults to x/gov unless overwritten).
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// params defines the x/staking parameters to update.
	//
	// NOTE: All parameters must be supplied.
	Params Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdateParams) Reset()         { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParams) ProtoMessage()    {}
func (*MsgUpdateParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_02cf64aff79d2288, []int{2}
}
func (m *MsgUpdateParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParams.Merge(m, src)
}
func (m *MsgUpdateParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParams proto.InternalMessageInfo

func (m *MsgUpdateParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdateParams) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// MsgUpdateParamsResponse
type MsgUpdateParamsResponse struct {
}

func (m *MsgUpdateParamsResponse) Reset()         { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParamsResponse) ProtoMessage()    {}
func (*MsgUpdateParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_02cf64aff79d2288, []int{3}
}
func (m *MsgUpdateParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParamsResponse.Merge(m, src)
}
func (m *MsgUpdateParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgCreatePrice)(nil), "exocore.oracle.MsgCreatePrice")
	proto.RegisterType((*MsgCreatePriceResponse)(nil), "exocore.oracle.MsgCreatePriceResponse")
	proto.RegisterType((*MsgUpdateParams)(nil), "exocore.oracle.MsgUpdateParams")
	proto.RegisterType((*MsgUpdateParamsResponse)(nil), "exocore.oracle.MsgUpdateParamsResponse")
}

func init() { proto.RegisterFile("exocore/oracle/tx.proto", fileDescriptor_02cf64aff79d2288) }

var fileDescriptor_02cf64aff79d2288 = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb3, 0xa4, 0x09, 0xcd, 0xa6, 0x2a, 0x62, 0x15, 0x35, 0x6e, 0x10, 0x76, 0x08, 0x12,
	0x84, 0x48, 0xb1, 0x21, 0x48, 0x45, 0xc0, 0x09, 0xf3, 0x47, 0x2a, 0x52, 0x10, 0x72, 0x55, 0x40,
	0x5c, 0x22, 0xc7, 0x5e, 0x1c, 0x2b, 0xb1, 0xd7, 0xda, 0xdd, 0x40, 0x7a, 0xe5, 0xc8, 0x89, 0xc7,
	0xe0, 0x98, 0x43, 0xc5, 0x89, 0x07, 0xe8, 0xb1, 0xea, 0x89, 0x53, 0x85, 0x92, 0x43, 0x78, 0x0c,
	0xe4, 0x5d, 0x9b, 0xd6, 0x06, 0xc1, 0x25, 0xf1, 0xcc, 0xef, 0xdb, 0xd9, 0x99, 0x6f, 0x16, 0xd6,
	0xf1, 0x8c, 0x38, 0x84, 0x62, 0x83, 0x50, 0xdb, 0x99, 0x60, 0x83, 0xcf, 0xf4, 0x88, 0x12, 0x4e,
	0xd0, 0x66, 0x02, 0x74, 0x09, 0x1a, 0x97, 0xed, 0xc0, 0x0f, 0x89, 0x21, 0x7e, 0xa5, 0xa4, 0x51,
	0x77, 0x08, 0x0b, 0x08, 0x33, 0x02, 0xe6, 0x19, 0xef, 0xef, 0xc4, 0x7f, 0x09, 0xd8, 0x96, 0x60,
	0x20, 0x22, 0x43, 0x06, 0x09, 0xba, 0x92, 0xbb, 0x2f, 0xb2, 0xa9, 0x1d, 0xa4, 0xb0, 0x91, 0x87,
	0xd4, 0x77, 0x70, 0xc2, 0x6a, 0x1e, 0xf1, 0x88, 0x2c, 0x18, 0x7f, 0xc9, 0x6c, 0xeb, 0x27, 0x80,
	0x9b, 0x7d, 0xe6, 0x3d, 0xa6, 0xd8, 0xe6, 0xf8, 0x65, 0x2c, 0x47, 0x0f, 0xe1, 0x45, 0x27, 0x0e,
	0x09, 0x55, 0x40, 0x13, 0xb4, 0x2b, 0xe6, 0xb5, 0x93, 0xc3, 0xee, 0xd5, 0xa4, 0x89, 0x57, 0xf6,
	0xc4, 0x77, 0x63, 0xf6, 0xc8, 0x75, 0x29, 0x66, 0x6c, 0x8f, 0x53, 0x3f, 0xf4, 0xac, 0xf4, 0x04,
	0xba, 0x05, 0x2b, 0xef, 0x30, 0x76, 0x31, 0x1d, 0xf8, 0xae, 0x72, 0xa1, 0x09, 0xda, 0x6b, 0xe6,
	0xc6, 0xe2, 0x54, 0x5b, 0x7f, 0x26, 0x92, 0xbb, 0x4f, 0xac, 0x75, 0x89, 0x77, 0x5d, 0x74, 0x0f,
	0x96, 0x45, 0x7f, 0x4c, 0x29, 0x36, 0x8b, 0xed, 0x6a, 0x4f, 0xd3, 0xb3, 0x8e, 0xe9, 0xa2, 0x9d,
	0xd7, 0x3e, 0x1f, 0xed, 0x91, 0x29, 0x75, 0xb0, 0x95, 0xc8, 0x91, 0x06, 0xab, 0x43, 0x9b, 0x61,
	0x77, 0x30, 0x9c, 0x10, 0x67, 0xac, 0xac, 0xc5, 0xb7, 0x58, 0x50, 0xa4, 0xcc, 0x38, 0x83, 0x6a,
	0xb0, 0x14, 0x92, 0xd0, 0xc1, 0x4a, 0xa9, 0x09, 0xda, 0x25, 0x4b, 0x06, 0x2d, 0x05, 0x6e, 0x65,
	0x27, 0xb5, 0x30, 0x8b, 0x48, 0xc8, 0x70, 0xeb, 0x1b, 0x80, 0x97, 0xfa, 0xcc, 0xdb, 0x8f, 0xdc,
	0x18, 0x09, 0x43, 0xd1, 0x0e, 0xac, 0xd8, 0x53, 0x3e, 0x22, 0xd4, 0xe7, 0x07, 0x89, 0x0f, 0xca,
	0xc9, 0x61, 0xb7, 0x96, 0xf8, 0x90, 0x1d, 0xff, 0x4c, 0x8a, 0xee, 0xc3, 0xb2, 0x5c, 0x89, 0x98,
	0xbe, 0xda, 0xdb, 0xfa, 0x63, 0x2a, 0x41, 0xcd, 0xca, 0xd1, 0xa9, 0x56, 0xf8, 0xb2, 0x9a, 0x77,
	0x80, 0x95, 0x1c, 0x78, 0xb0, 0xf3, 0x71, 0x35, 0xef, 0x9c, 0x95, 0xfa, 0xb4, 0x9a, 0x77, 0xae,
	0xcb, 0xeb, 0xba, 0xcc, 0x1d, 0x1b, 0xb3, 0x74, 0xab, 0xb9, 0x56, 0x5b, 0xdb, 0xb0, 0x9e, 0x4b,
	0xa5, 0x93, 0xf5, 0xbe, 0x02, 0x58, 0xec, 0x33, 0x0f, 0xed, 0xc3, 0xea, 0xf9, 0x15, 0xab, 0xf9,
	0xa6, 0xb2, 0xc6, 0x34, 0x6e, 0xfc, 0x9b, 0xa7, 0xe5, 0xd1, 0x1b, 0xb8, 0x91, 0x31, 0x4d, 0xfb,
	0xcb, 0xb9, 0xf3, 0x82, 0xc6, 0xcd, 0xff, 0x08, 0xd2, 0xca, 0xe6, 0xf3, 0xa3, 0x85, 0x0a, 0x8e,
	0x17, 0x2a, 0xf8, 0xb1, 0x50, 0xc1, 0xe7, 0xa5, 0x5a, 0x38, 0x5e, 0xaa, 0x85, 0xef, 0x4b, 0xb5,
	0xf0, 0xf6, 0xb6, 0xe7, 0xf3, 0xd1, 0x74, 0xa8, 0x3b, 0x24, 0x30, 0x9e, 0xca, 0x62, 0x2f, 0x30,
	0xff, 0x40, 0xe8, 0xd8, 0x48, 0x5f, 0xff, 0x6f, 0xa7, 0xf8, 0x41, 0x84, 0xd9, 0xb0, 0x2c, 0x9e,
	0xfa, 0xdd, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x97, 0xa1, 0x04, 0x18, 0xab, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// CreatePrice creates price for a new oracle round
	CreatePrice(ctx context.Context, in *MsgCreatePrice, opts ...grpc.CallOption) (*MsgCreatePriceResponse, error)
	// UpdateParams update params value
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreatePrice(ctx context.Context, in *MsgCreatePrice, opts ...grpc.CallOption) (*MsgCreatePriceResponse, error) {
	out := new(MsgCreatePriceResponse)
	err := c.cc.Invoke(ctx, "/exocore.oracle.Msg/CreatePrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/exocore.oracle.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// CreatePrice creates price for a new oracle round
	CreatePrice(context.Context, *MsgCreatePrice) (*MsgCreatePriceResponse, error)
	// UpdateParams update params value
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) CreatePrice(ctx context.Context, req *MsgCreatePrice) (*MsgCreatePriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePrice not implemented")
}
func (*UnimplementedMsgServer) UpdateParams(ctx context.Context, req *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreatePrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreatePrice)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreatePrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exocore.oracle.Msg/CreatePrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreatePrice(ctx, req.(*MsgCreatePrice))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exocore.oracle.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "exocore.oracle.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePrice",
			Handler:    _Msg_CreatePrice_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "exocore/oracle/tx.proto",
}

func (m *MsgCreatePrice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreatePrice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreatePrice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Nonce != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x28
	}
	if m.BasedBlock != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.BasedBlock))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Prices) > 0 {
		for iNdEx := len(m.Prices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.FeederID != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.FeederID))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreatePriceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreatePriceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreatePriceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgCreatePrice) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.FeederID != 0 {
		n += 1 + sovTx(uint64(m.FeederID))
	}
	if len(m.Prices) > 0 {
		for _, e := range m.Prices {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	if m.BasedBlock != 0 {
		n += 1 + sovTx(uint64(m.BasedBlock))
	}
	if m.Nonce != 0 {
		n += 1 + sovTx(uint64(m.Nonce))
	}
	return n
}

func (m *MsgCreatePriceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgUpdateParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgUpdateParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgCreatePrice) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreatePrice: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreatePrice: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeederID", wireType)
			}
			m.FeederID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeederID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prices = append(m.Prices, &PriceWithSource{})
			if err := m.Prices[len(m.Prices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BasedBlock", wireType)
			}
			m.BasedBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BasedBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgCreatePriceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreatePriceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreatePriceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpdateParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpdateParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
