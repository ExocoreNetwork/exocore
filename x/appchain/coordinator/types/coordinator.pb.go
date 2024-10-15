// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: exocore/appchain/coordinator/v1/coordinator.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// PendingSubscriberChainRequests is a helper structure to store a list of
// subscriber chain requests that are pending activation.
type PendingSubscriberChainRequests struct {
	// list is the list of subscriber chain requests that are pending activation.
	List []RegisterSubscriberChainRequest `protobuf:"bytes,1,rep,name=list,proto3" json:"list"`
}

func (m *PendingSubscriberChainRequests) Reset()         { *m = PendingSubscriberChainRequests{} }
func (m *PendingSubscriberChainRequests) String() string { return proto.CompactTextString(m) }
func (*PendingSubscriberChainRequests) ProtoMessage()    {}
func (*PendingSubscriberChainRequests) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb7bb04617dc0e61, []int{0}
}
func (m *PendingSubscriberChainRequests) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PendingSubscriberChainRequests) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PendingSubscriberChainRequests.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PendingSubscriberChainRequests) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PendingSubscriberChainRequests.Merge(m, src)
}
func (m *PendingSubscriberChainRequests) XXX_Size() int {
	return m.Size()
}
func (m *PendingSubscriberChainRequests) XXX_DiscardUnknown() {
	xxx_messageInfo_PendingSubscriberChainRequests.DiscardUnknown(m)
}

var xxx_messageInfo_PendingSubscriberChainRequests proto.InternalMessageInfo

func (m *PendingSubscriberChainRequests) GetList() []RegisterSubscriberChainRequest {
	if m != nil {
		return m.List
	}
	return nil
}

// ChainIDs is a helper structure to store a list of chain IDs.
type ChainIDs struct {
	// list is the list of chain IDs.
	List []string `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (m *ChainIDs) Reset()         { *m = ChainIDs{} }
func (m *ChainIDs) String() string { return proto.CompactTextString(m) }
func (*ChainIDs) ProtoMessage()    {}
func (*ChainIDs) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb7bb04617dc0e61, []int{1}
}
func (m *ChainIDs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainIDs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainIDs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainIDs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainIDs.Merge(m, src)
}
func (m *ChainIDs) XXX_Size() int {
	return m.Size()
}
func (m *ChainIDs) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainIDs.DiscardUnknown(m)
}

var xxx_messageInfo_ChainIDs proto.InternalMessageInfo

func (m *ChainIDs) GetList() []string {
	if m != nil {
		return m.List
	}
	return nil
}

// ConsensusAddresses is a list of consensus addresses.
type ConsensusAddresses struct {
	// list is the list of consensus addresses.
	List [][]byte `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (m *ConsensusAddresses) Reset()         { *m = ConsensusAddresses{} }
func (m *ConsensusAddresses) String() string { return proto.CompactTextString(m) }
func (*ConsensusAddresses) ProtoMessage()    {}
func (*ConsensusAddresses) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb7bb04617dc0e61, []int{2}
}
func (m *ConsensusAddresses) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConsensusAddresses) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ConsensusAddresses.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ConsensusAddresses) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsensusAddresses.Merge(m, src)
}
func (m *ConsensusAddresses) XXX_Size() int {
	return m.Size()
}
func (m *ConsensusAddresses) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsensusAddresses.DiscardUnknown(m)
}

var xxx_messageInfo_ConsensusAddresses proto.InternalMessageInfo

func (m *ConsensusAddresses) GetList() [][]byte {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*PendingSubscriberChainRequests)(nil), "exocore.appchain.coordinator.v1.PendingSubscriberChainRequests")
	proto.RegisterType((*ChainIDs)(nil), "exocore.appchain.coordinator.v1.ChainIDs")
	proto.RegisterType((*ConsensusAddresses)(nil), "exocore.appchain.coordinator.v1.ConsensusAddresses")
}

func init() {
	proto.RegisterFile("exocore/appchain/coordinator/v1/coordinator.proto", fileDescriptor_fb7bb04617dc0e61)
}

var fileDescriptor_fb7bb04617dc0e61 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x13, 0x51, 0x21, 0x08, 0x4c, 0x11, 0x03, 0xea, 0xe0, 0xa2, 0x4e, 0x9d, 0x6c, 0x05,
	0x76, 0x10, 0x2d, 0x0c, 0x2c, 0x08, 0x85, 0x05, 0xd8, 0x12, 0xe7, 0xe4, 0x5a, 0x80, 0x2f, 0xf8,
	0x9c, 0x12, 0xc4, 0x9f, 0xe0, 0x67, 0x75, 0xec, 0xc8, 0x84, 0x50, 0xf2, 0x47, 0x50, 0x93, 0x56,
	0x0a, 0x52, 0xa5, 0x6e, 0x67, 0xdf, 0x7b, 0xdf, 0x3d, 0xbd, 0x20, 0x82, 0x12, 0x25, 0x5a, 0x10,
	0x49, 0x9e, 0xcb, 0x69, 0xa2, 0x8d, 0x90, 0x88, 0x36, 0xd3, 0x26, 0x71, 0x68, 0xc5, 0x2c, 0xea,
	0x3e, 0x79, 0x6e, 0xd1, 0x61, 0x38, 0x58, 0x59, 0xf8, 0xda, 0xc2, 0xbb, 0x9a, 0x59, 0xd4, 0x1f,
	0x6d, 0x63, 0xba, 0xb2, 0x45, 0xf5, 0x8f, 0x14, 0x2a, 0x6c, 0x46, 0xb1, 0x9c, 0xda, 0xdf, 0xe1,
	0x67, 0xc0, 0xee, 0xc0, 0x64, 0xda, 0xa8, 0xfb, 0x22, 0x25, 0x69, 0x75, 0x0a, 0x76, 0xb2, 0xe4,
	0xc4, 0xf0, 0x56, 0x00, 0x39, 0x0a, 0x1f, 0x83, 0xde, 0x8b, 0x26, 0x77, 0xec, 0x9f, 0xec, 0x8c,
	0x0e, 0x4e, 0x2f, 0xf8, 0x96, 0x44, 0x3c, 0x06, 0xa5, 0xc9, 0x81, 0xdd, 0xcc, 0x1b, 0xf7, 0xe6,
	0x3f, 0x03, 0x2f, 0x6e, 0x90, 0x43, 0x16, 0xec, 0x35, 0xbb, 0x9b, 0x2b, 0x0a, 0xc3, 0xce, 0x99,
	0xfd, 0xd5, 0x7e, 0x14, 0x84, 0x13, 0x34, 0x04, 0x86, 0x0a, 0xba, 0xcc, 0x32, 0x0b, 0x44, 0xf0,
	0x5f, 0x79, 0xd8, 0x2a, 0xc7, 0x0f, 0xf3, 0x8a, 0xf9, 0x8b, 0x8a, 0xf9, 0xbf, 0x15, 0xf3, 0xbf,
	0x6a, 0xe6, 0x2d, 0x6a, 0xe6, 0x7d, 0xd7, 0xcc, 0x7b, 0x3a, 0x57, 0xda, 0x4d, 0x8b, 0x94, 0x4b,
	0x7c, 0x15, 0xd7, 0x6d, 0xf4, 0x5b, 0x70, 0xef, 0x68, 0x9f, 0xc5, 0xba, 0xba, 0x72, 0x73, 0x79,
	0xee, 0x23, 0x07, 0x4a, 0x77, 0x9b, 0x9e, 0xce, 0xfe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x12, 0xc7,
	0x0f, 0x6b, 0xbd, 0x01, 0x00, 0x00,
}

func (m *PendingSubscriberChainRequests) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PendingSubscriberChainRequests) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PendingSubscriberChainRequests) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.List) > 0 {
		for iNdEx := len(m.List) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.List[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCoordinator(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ChainIDs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainIDs) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainIDs) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.List) > 0 {
		for iNdEx := len(m.List) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.List[iNdEx])
			copy(dAtA[i:], m.List[iNdEx])
			i = encodeVarintCoordinator(dAtA, i, uint64(len(m.List[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ConsensusAddresses) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConsensusAddresses) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConsensusAddresses) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.List) > 0 {
		for iNdEx := len(m.List) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.List[iNdEx])
			copy(dAtA[i:], m.List[iNdEx])
			i = encodeVarintCoordinator(dAtA, i, uint64(len(m.List[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintCoordinator(dAtA []byte, offset int, v uint64) int {
	offset -= sovCoordinator(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PendingSubscriberChainRequests) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.List) > 0 {
		for _, e := range m.List {
			l = e.Size()
			n += 1 + l + sovCoordinator(uint64(l))
		}
	}
	return n
}

func (m *ChainIDs) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.List) > 0 {
		for _, s := range m.List {
			l = len(s)
			n += 1 + l + sovCoordinator(uint64(l))
		}
	}
	return n
}

func (m *ConsensusAddresses) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.List) > 0 {
		for _, b := range m.List {
			l = len(b)
			n += 1 + l + sovCoordinator(uint64(l))
		}
	}
	return n
}

func sovCoordinator(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCoordinator(x uint64) (n int) {
	return sovCoordinator(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PendingSubscriberChainRequests) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCoordinator
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
			return fmt.Errorf("proto: PendingSubscriberChainRequests: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PendingSubscriberChainRequests: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field List", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCoordinator
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
				return ErrInvalidLengthCoordinator
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCoordinator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.List = append(m.List, RegisterSubscriberChainRequest{})
			if err := m.List[len(m.List)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCoordinator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCoordinator
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
func (m *ChainIDs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCoordinator
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
			return fmt.Errorf("proto: ChainIDs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainIDs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field List", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCoordinator
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
				return ErrInvalidLengthCoordinator
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCoordinator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.List = append(m.List, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCoordinator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCoordinator
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
func (m *ConsensusAddresses) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCoordinator
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
			return fmt.Errorf("proto: ConsensusAddresses: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConsensusAddresses: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field List", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCoordinator
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCoordinator
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCoordinator
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.List = append(m.List, make([]byte, postIndex-iNdEx))
			copy(m.List[len(m.List)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCoordinator(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCoordinator
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
func skipCoordinator(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCoordinator
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
					return 0, ErrIntOverflowCoordinator
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
					return 0, ErrIntOverflowCoordinator
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
				return 0, ErrInvalidLengthCoordinator
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCoordinator
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCoordinator
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCoordinator        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCoordinator          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCoordinator = fmt.Errorf("proto: unexpected end of group")
)
