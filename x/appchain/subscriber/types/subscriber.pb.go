// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: exocore/appchain/subscriber/v1/subscriber.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
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

type OmniChainValidator struct {
	// The address is derived from the consenus key. It has no relation with the operator
	// address on Exocore.
	Address []byte `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Last known
	Power int64 `protobuf:"varint,2,opt,name=power,proto3" json:"power,omitempty"`
	// pubkey is the consensus public key of the validator, as a Protobuf Any.
	Pubkey *types.Any `protobuf:"bytes,3,opt,name=pubkey,proto3" json:"pubkey,omitempty" yaml:"consensus_pubkey"`
}

func (m *OmniChainValidator) Reset()         { *m = OmniChainValidator{} }
func (m *OmniChainValidator) String() string { return proto.CompactTextString(m) }
func (*OmniChainValidator) ProtoMessage()    {}
func (*OmniChainValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c15efc73dd829e9, []int{0}
}
func (m *OmniChainValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OmniChainValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OmniChainValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OmniChainValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OmniChainValidator.Merge(m, src)
}
func (m *OmniChainValidator) XXX_Size() int {
	return m.Size()
}
func (m *OmniChainValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_OmniChainValidator.DiscardUnknown(m)
}

var xxx_messageInfo_OmniChainValidator proto.InternalMessageInfo

func (m *OmniChainValidator) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *OmniChainValidator) GetPower() int64 {
	if m != nil {
		return m.Power
	}
	return 0
}

func (m *OmniChainValidator) GetPubkey() *types.Any {
	if m != nil {
		return m.Pubkey
	}
	return nil
}

func init() {
	proto.RegisterType((*OmniChainValidator)(nil), "exocore.appchain.subscriber.v1.OmniChainValidator")
}

func init() {
	proto.RegisterFile("exocore/appchain/subscriber/v1/subscriber.proto", fileDescriptor_1c15efc73dd829e9)
}

var fileDescriptor_1c15efc73dd829e9 = []byte{
	// 320 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4e, 0x02, 0x31,
	0x14, 0x86, 0xa9, 0x44, 0x4c, 0x46, 0x57, 0x93, 0x49, 0x1c, 0x58, 0x54, 0xc2, 0x8a, 0x8d, 0x6d,
	0x90, 0x9d, 0x89, 0x0b, 0x31, 0xae, 0x4c, 0xd4, 0xb0, 0xd0, 0xc4, 0x0d, 0x69, 0x4b, 0x1d, 0x26,
	0x30, 0xf3, 0x9a, 0x76, 0x0a, 0xf4, 0x16, 0xde, 0xc2, 0x0b, 0x78, 0x08, 0xe3, 0x8a, 0xa5, 0x2b,
	0x63, 0xe0, 0x06, 0x9e, 0xc0, 0x40, 0x07, 0xc3, 0xc2, 0xdd, 0xfb, 0xfb, 0xfe, 0xbf, 0xef, 0xcb,
	0x1f, 0x50, 0x39, 0x07, 0x01, 0x5a, 0x52, 0xa6, 0x94, 0x18, 0xb1, 0x34, 0xa7, 0xc6, 0x72, 0x23,
	0x74, 0xca, 0xa5, 0xa6, 0xd3, 0xce, 0x8e, 0x22, 0x4a, 0x43, 0x01, 0x21, 0x2e, 0x03, 0x64, 0x1b,
	0x20, 0x3b, 0x96, 0x69, 0xa7, 0x51, 0x4f, 0x00, 0x92, 0x89, 0xa4, 0x1b, 0x37, 0xb7, 0xcf, 0x94,
	0xe5, 0xce, 0x47, 0x1b, 0x51, 0x02, 0x09, 0x6c, 0x46, 0xba, 0x9e, 0xca, 0xd7, 0xba, 0x00, 0x93,
	0x81, 0x19, 0xf8, 0x85, 0x17, 0x7e, 0xd5, 0x7a, 0x45, 0x41, 0x78, 0x97, 0xe5, 0xe9, 0xd5, 0xfa,
	0xce, 0x03, 0x9b, 0xa4, 0x43, 0x56, 0x80, 0x0e, 0xe3, 0xe0, 0x80, 0x0d, 0x87, 0x5a, 0x1a, 0x13,
	0xa3, 0x26, 0x6a, 0x1f, 0xf5, 0xb7, 0x32, 0x8c, 0x82, 0x7d, 0x05, 0x33, 0xa9, 0xe3, 0xbd, 0x26,
	0x6a, 0x57, 0xfb, 0x5e, 0x84, 0x2c, 0xa8, 0x29, 0xcb, 0xc7, 0xd2, 0xc5, 0xd5, 0x26, 0x6a, 0x1f,
	0x9e, 0x45, 0xc4, 0x33, 0x92, 0x2d, 0x23, 0xb9, 0xcc, 0x5d, 0xaf, 0xfb, 0xf3, 0x75, 0x72, 0xec,
	0x58, 0x36, 0x39, 0x6f, 0x09, 0xc8, 0x8d, 0xcc, 0x8d, 0x35, 0x03, 0x9f, 0x6b, 0x7d, 0xbc, 0x9d,
	0x46, 0x25, 0x99, 0xd0, 0x4e, 0x15, 0x40, 0xee, 0x2d, 0xbf, 0x91, 0xae, 0x5f, 0x7e, 0xdc, 0x7b,
	0x7c, 0x5f, 0x62, 0xb4, 0x58, 0x62, 0xf4, 0xbd, 0xc4, 0xe8, 0x65, 0x85, 0x2b, 0x8b, 0x15, 0xae,
	0x7c, 0xae, 0x70, 0xe5, 0xe9, 0x22, 0x49, 0x8b, 0x91, 0xe5, 0x44, 0x40, 0x46, 0xaf, 0x7d, 0x75,
	0xb7, 0xb2, 0x98, 0x81, 0x1e, 0xff, 0x55, 0x3f, 0xff, 0xb7, 0xfc, 0xc2, 0x29, 0x69, 0x78, 0x6d,
	0xc3, 0xd8, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xee, 0x66, 0xe1, 0x6d, 0xa8, 0x01, 0x00, 0x00,
}

func (m *OmniChainValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OmniChainValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OmniChainValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pubkey != nil {
		{
			size, err := m.Pubkey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSubscriber(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.Power != 0 {
		i = encodeVarintSubscriber(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintSubscriber(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSubscriber(dAtA []byte, offset int, v uint64) int {
	offset -= sovSubscriber(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *OmniChainValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovSubscriber(uint64(l))
	}
	if m.Power != 0 {
		n += 1 + sovSubscriber(uint64(m.Power))
	}
	if m.Pubkey != nil {
		l = m.Pubkey.Size()
		n += 1 + l + sovSubscriber(uint64(l))
	}
	return n
}

func sovSubscriber(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSubscriber(x uint64) (n int) {
	return sovSubscriber(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *OmniChainValidator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSubscriber
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
			return fmt.Errorf("proto: OmniChainValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OmniChainValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubscriber
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
				return ErrInvalidLengthSubscriber
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSubscriber
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubscriber
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pubkey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubscriber
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
				return ErrInvalidLengthSubscriber
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubscriber
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pubkey == nil {
				m.Pubkey = &types.Any{}
			}
			if err := m.Pubkey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSubscriber(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSubscriber
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
func skipSubscriber(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSubscriber
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
					return 0, ErrIntOverflowSubscriber
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
					return 0, ErrIntOverflowSubscriber
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
				return 0, ErrInvalidLengthSubscriber
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSubscriber
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSubscriber
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSubscriber        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSubscriber          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSubscriber = fmt.Errorf("proto: unexpected end of group")
)
