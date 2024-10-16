// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: exocore/oracle/v1/recent_msg.proto

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

// RecentMsg represent the messages to be cached for recent blocks
type RecentMsg struct {
	// block height these messages from
	Block uint64 `protobuf:"varint,1,opt,name=block,proto3" json:"block,omitempty"`
	// cached messages
	Msgs []*MsgItem `protobuf:"bytes,2,rep,name=msgs,proto3" json:"msgs,omitempty"`
}

func (m *RecentMsg) Reset()         { *m = RecentMsg{} }
func (m *RecentMsg) String() string { return proto.CompactTextString(m) }
func (*RecentMsg) ProtoMessage()    {}
func (*RecentMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe72ca6e7b5df271, []int{0}
}
func (m *RecentMsg) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RecentMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RecentMsg.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RecentMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecentMsg.Merge(m, src)
}
func (m *RecentMsg) XXX_Size() int {
	return m.Size()
}
func (m *RecentMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_RecentMsg.DiscardUnknown(m)
}

var xxx_messageInfo_RecentMsg proto.InternalMessageInfo

func (m *RecentMsg) GetBlock() uint64 {
	if m != nil {
		return m.Block
	}
	return 0
}

func (m *RecentMsg) GetMsgs() []*MsgItem {
	if m != nil {
		return m.Msgs
	}
	return nil
}

// MsgItem represents the message info of createPrice
type MsgItem struct {
	// feeder_id tells of wich feeder this price if corresponding to
	FeederID uint64 `protobuf:"varint,2,opt,name=feeder_id,json=feederId,proto3" json:"feeder_id,omitempty"`
	// p_source price with its source info
	PSources []*PriceSource `protobuf:"bytes,3,rep,name=p_sources,json=pSources,proto3" json:"p_sources,omitempty"`
	// validator tells which validator create this price
	Validator string `protobuf:"bytes,4,opt,name=validator,proto3" json:"validator,omitempty"`
}

func (m *MsgItem) Reset()         { *m = MsgItem{} }
func (m *MsgItem) String() string { return proto.CompactTextString(m) }
func (*MsgItem) ProtoMessage()    {}
func (*MsgItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_fe72ca6e7b5df271, []int{1}
}
func (m *MsgItem) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgItem.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgItem.Merge(m, src)
}
func (m *MsgItem) XXX_Size() int {
	return m.Size()
}
func (m *MsgItem) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgItem.DiscardUnknown(m)
}

var xxx_messageInfo_MsgItem proto.InternalMessageInfo

func (m *MsgItem) GetFeederID() uint64 {
	if m != nil {
		return m.FeederID
	}
	return 0
}

func (m *MsgItem) GetPSources() []*PriceSource {
	if m != nil {
		return m.PSources
	}
	return nil
}

func (m *MsgItem) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func init() {
	proto.RegisterType((*RecentMsg)(nil), "exocore.oracle.v1.RecentMsg")
	proto.RegisterType((*MsgItem)(nil), "exocore.oracle.v1.MsgItem")
}

func init() {
	proto.RegisterFile("exocore/oracle/v1/recent_msg.proto", fileDescriptor_fe72ca6e7b5df271)
}

var fileDescriptor_fe72ca6e7b5df271 = []byte{
	// 310 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4e, 0x32, 0x31,
	0x14, 0x85, 0x29, 0xf0, 0xff, 0x32, 0xd5, 0x8d, 0x13, 0x16, 0x13, 0xa2, 0x95, 0xb0, 0xc2, 0x4d,
	0x47, 0x74, 0xe9, 0x8e, 0xa8, 0x09, 0x26, 0x18, 0xad, 0x3b, 0x37, 0x04, 0x3a, 0xd7, 0x3a, 0x81,
	0xf1, 0x4e, 0xda, 0x82, 0xf8, 0x14, 0xfa, 0x58, 0x2e, 0x59, 0xba, 0x32, 0x66, 0x78, 0x11, 0x43,
	0x0b, 0x71, 0x81, 0xbb, 0xd3, 0xd3, 0x73, 0xef, 0x97, 0x7b, 0x68, 0x0b, 0xe6, 0x28, 0x51, 0x43,
	0x8c, 0x7a, 0x28, 0x27, 0x10, 0xcf, 0x3a, 0xb1, 0x06, 0x09, 0xcf, 0x76, 0x90, 0x19, 0xc5, 0x73,
	0x8d, 0x16, 0xc3, 0xfd, 0x75, 0x86, 0xfb, 0x0c, 0x9f, 0x75, 0x1a, 0x87, 0xdb, 0x63, 0xb9, 0x4e,
	0x25, 0xf8, 0x89, 0x46, 0x5d, 0xa1, 0x42, 0x27, 0xe3, 0x95, 0xf2, 0x6e, 0xeb, 0x8e, 0x06, 0xc2,
	0xed, 0xee, 0x1b, 0x15, 0xd6, 0xe9, 0xbf, 0xd1, 0x04, 0xe5, 0x38, 0x22, 0x4d, 0xd2, 0xae, 0x0a,
	0xff, 0x08, 0x39, 0xad, 0x66, 0x46, 0x99, 0xa8, 0xdc, 0xac, 0xb4, 0x77, 0x4f, 0x1b, 0x7c, 0x8b,
	0xcc, 0xfb, 0x46, 0xf5, 0x2c, 0x64, 0xc2, 0xe5, 0x5a, 0x6f, 0x84, 0xee, 0xac, 0x9d, 0xf0, 0x98,
	0x06, 0x8f, 0x00, 0x09, 0xe8, 0x41, 0x9a, 0x44, 0xe5, 0xd5, 0xd6, 0xee, 0x5e, 0xf1, 0x75, 0x54,
	0xbb, 0x72, 0x66, 0xef, 0x42, 0xd4, 0xfc, 0x77, 0x2f, 0x09, 0xcf, 0x69, 0x90, 0x0f, 0x0c, 0x4e,
	0xb5, 0x04, 0x13, 0x55, 0x1c, 0x8b, 0xfd, 0xc1, 0xba, 0x5d, 0x9d, 0x74, 0xef, 0x62, 0xa2, 0x96,
	0x7b, 0x61, 0xc2, 0x03, 0x1a, 0xcc, 0x86, 0x93, 0x34, 0x19, 0x5a, 0xd4, 0x51, 0xb5, 0x49, 0xda,
	0x81, 0xf8, 0x35, 0xba, 0xd7, 0x1f, 0x05, 0x23, 0x8b, 0x82, 0x91, 0xef, 0x82, 0x91, 0xf7, 0x25,
	0x2b, 0x2d, 0x96, 0xac, 0xf4, 0xb9, 0x64, 0xa5, 0x87, 0x13, 0x95, 0xda, 0xa7, 0xe9, 0x88, 0x4b,
	0xcc, 0xe2, 0x4b, 0xcf, 0xba, 0x01, 0xfb, 0x82, 0x7a, 0x1c, 0x6f, 0xda, 0x9c, 0x6f, 0xfa, 0xb4,
	0xaf, 0x39, 0x98, 0xd1, 0x7f, 0xd7, 0xdb, 0xd9, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5e, 0x69,
	0x23, 0x94, 0xa5, 0x01, 0x00, 0x00,
}

func (m *RecentMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RecentMsg) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RecentMsg) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Msgs) > 0 {
		for iNdEx := len(m.Msgs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Msgs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRecentMsg(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Block != 0 {
		i = encodeVarintRecentMsg(dAtA, i, uint64(m.Block))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MsgItem) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgItem) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgItem) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintRecentMsg(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.PSources) > 0 {
		for iNdEx := len(m.PSources) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PSources[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRecentMsg(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.FeederID != 0 {
		i = encodeVarintRecentMsg(dAtA, i, uint64(m.FeederID))
		i--
		dAtA[i] = 0x10
	}
	return len(dAtA) - i, nil
}

func encodeVarintRecentMsg(dAtA []byte, offset int, v uint64) int {
	offset -= sovRecentMsg(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RecentMsg) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Block != 0 {
		n += 1 + sovRecentMsg(uint64(m.Block))
	}
	if len(m.Msgs) > 0 {
		for _, e := range m.Msgs {
			l = e.Size()
			n += 1 + l + sovRecentMsg(uint64(l))
		}
	}
	return n
}

func (m *MsgItem) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.FeederID != 0 {
		n += 1 + sovRecentMsg(uint64(m.FeederID))
	}
	if len(m.PSources) > 0 {
		for _, e := range m.PSources {
			l = e.Size()
			n += 1 + l + sovRecentMsg(uint64(l))
		}
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovRecentMsg(uint64(l))
	}
	return n
}

func sovRecentMsg(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRecentMsg(x uint64) (n int) {
	return sovRecentMsg(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RecentMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRecentMsg
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
			return fmt.Errorf("proto: RecentMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RecentMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Block", wireType)
			}
			m.Block = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecentMsg
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Block |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecentMsg
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
				return ErrInvalidLengthRecentMsg
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRecentMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msgs = append(m.Msgs, &MsgItem{})
			if err := m.Msgs[len(m.Msgs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRecentMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRecentMsg
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
func (m *MsgItem) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRecentMsg
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
			return fmt.Errorf("proto: MsgItem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgItem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeederID", wireType)
			}
			m.FeederID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecentMsg
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
				return fmt.Errorf("proto: wrong wireType = %d for field PSources", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecentMsg
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
				return ErrInvalidLengthRecentMsg
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRecentMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PSources = append(m.PSources, &PriceSource{})
			if err := m.PSources[len(m.PSources)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRecentMsg
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
				return ErrInvalidLengthRecentMsg
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRecentMsg
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRecentMsg(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRecentMsg
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
func skipRecentMsg(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRecentMsg
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
					return 0, ErrIntOverflowRecentMsg
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
					return 0, ErrIntOverflowRecentMsg
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
				return 0, ErrInvalidLengthRecentMsg
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRecentMsg
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRecentMsg
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRecentMsg        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRecentMsg          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRecentMsg = fmt.Errorf("proto: unexpected end of group")
)
