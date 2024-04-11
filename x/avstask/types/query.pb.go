// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: exocore/avstask/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// QueryAVSTaskInfoReq is the request to obtain the task information.
type GetAVSTaskInfoReq struct {
	// task_addr is the task contract address,its type should be a sdk.AccAddress
	TaskAddr string `protobuf:"bytes,1,opt,name=task_addr,json=taskAddr,proto3" json:"task_addr,omitempty"`
}

func (m *GetAVSTaskInfoReq) Reset()         { *m = GetAVSTaskInfoReq{} }
func (m *GetAVSTaskInfoReq) String() string { return proto.CompactTextString(m) }
func (*GetAVSTaskInfoReq) ProtoMessage()    {}
func (*GetAVSTaskInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb88546bd6602038, []int{0}
}
func (m *GetAVSTaskInfoReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GetAVSTaskInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GetAVSTaskInfoReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GetAVSTaskInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAVSTaskInfoReq.Merge(m, src)
}
func (m *GetAVSTaskInfoReq) XXX_Size() int {
	return m.Size()
}
func (m *GetAVSTaskInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAVSTaskInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetAVSTaskInfoReq proto.InternalMessageInfo

func (m *GetAVSTaskInfoReq) GetTaskAddr() string {
	if m != nil {
		return m.TaskAddr
	}
	return ""
}

func init() {
	proto.RegisterType((*GetAVSTaskInfoReq)(nil), "exocore.avstask.v1.GetAVSTaskInfoReq")
}

func init() { proto.RegisterFile("exocore/avstask/v1/query.proto", fileDescriptor_fb88546bd6602038) }

var fileDescriptor_fb88546bd6602038 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xad, 0xc8, 0x4f,
	0xce, 0x2f, 0x4a, 0xd5, 0x4f, 0x2c, 0x2b, 0x2e, 0x49, 0x2c, 0xce, 0xd6, 0x2f, 0x33, 0xd4, 0x2f,
	0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x82, 0xca, 0xeb, 0x41,
	0xe5, 0xf5, 0xca, 0x0c, 0xa5, 0x24, 0x93, 0xf3, 0x8b, 0x73, 0xf3, 0x8b, 0xe3, 0xc1, 0x2a, 0xf4,
	0x21, 0x1c, 0x88, 0x72, 0x29, 0x69, 0x2c, 0xc6, 0x95, 0x54, 0x40, 0x25, 0x65, 0xd2, 0xf3, 0xf3,
	0xd3, 0x73, 0x52, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12, 0x4b, 0x32,
	0xf3, 0xf3, 0xa0, 0x5a, 0x95, 0xbc, 0xb8, 0x04, 0xdd, 0x53, 0x4b, 0x1c, 0xc3, 0x82, 0x43, 0x12,
	0x8b, 0xb3, 0x3d, 0xf3, 0xd2, 0xf2, 0x83, 0x52, 0x0b, 0x85, 0x4c, 0xb9, 0x38, 0x41, 0xc6, 0xc4,
	0x27, 0xa6, 0xa4, 0x14, 0x49, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x3a, 0x49, 0x5c, 0xda, 0xa2, 0x2b,
	0x02, 0xb5, 0xd4, 0x31, 0x25, 0xa5, 0x28, 0xb5, 0xb8, 0x38, 0xb8, 0xa4, 0x28, 0x33, 0x2f, 0x3d,
	0x88, 0x03, 0xa4, 0x14, 0x24, 0x64, 0x34, 0x8d, 0x91, 0x8b, 0x35, 0x10, 0xe4, 0x0b, 0xa1, 0x1e,
	0x46, 0x2e, 0x3e, 0x54, 0x63, 0x85, 0x54, 0xf5, 0x30, 0xfd, 0xa4, 0x87, 0x61, 0xb5, 0x94, 0x0a,
	0x36, 0x65, 0x20, 0x05, 0xce, 0xf9, 0x79, 0x25, 0x45, 0x89, 0xc9, 0x25, 0x20, 0x85, 0x4a, 0xba,
	0x4d, 0x97, 0x9f, 0x4c, 0x66, 0x52, 0x17, 0x52, 0xd5, 0xc7, 0xe2, 0x73, 0x0c, 0x43, 0x9d, 0xbc,
	0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5,
	0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x30, 0x3d, 0xb3, 0x24, 0xa3,
	0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0xdf, 0x15, 0x62, 0x94, 0x5f, 0x6a, 0x49, 0x79, 0x7e, 0x51,
	0x36, 0xdc, 0xe4, 0x0a, 0xb8, 0xd9, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x49, 0x6c, 0xe0, 0x80, 0x33,
	0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x2a, 0x78, 0xf2, 0xc4, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// TaskInfo queries the task information.
	GetAVSTaskInfo(ctx context.Context, in *GetAVSTaskInfoReq, opts ...grpc.CallOption) (*TaskContractInfo, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetAVSTaskInfo(ctx context.Context, in *GetAVSTaskInfoReq, opts ...grpc.CallOption) (*TaskContractInfo, error) {
	out := new(TaskContractInfo)
	err := c.cc.Invoke(ctx, "/exocore.avstask.v1.Query/GetAVSTaskInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// TaskInfo queries the task information.
	GetAVSTaskInfo(context.Context, *GetAVSTaskInfoReq) (*TaskContractInfo, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) GetAVSTaskInfo(ctx context.Context, req *GetAVSTaskInfoReq) (*TaskContractInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAVSTaskInfo not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetAVSTaskInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAVSTaskInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetAVSTaskInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exocore.avstask.v1.Query/GetAVSTaskInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetAVSTaskInfo(ctx, req.(*GetAVSTaskInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "exocore.avstask.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAVSTaskInfo",
			Handler:    _Query_GetAVSTaskInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "exocore/avstask/v1/query.proto",
}

func (m *GetAVSTaskInfoReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GetAVSTaskInfoReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GetAVSTaskInfoReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TaskAddr) > 0 {
		i -= len(m.TaskAddr)
		copy(dAtA[i:], m.TaskAddr)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.TaskAddr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GetAVSTaskInfoReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TaskAddr)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GetAVSTaskInfoReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: GetAVSTaskInfoReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetAVSTaskInfoReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TaskAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TaskAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
