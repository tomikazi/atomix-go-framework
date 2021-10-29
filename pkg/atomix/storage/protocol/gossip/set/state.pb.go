// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: atomix/storage/protocol/gossip/set/state.proto

package set

import (
	fmt "fmt"
	meta "github.com/atomix/atomix-api/go/atomix/primitive/meta"
	_ "github.com/atomix/atomix-go-sdk/pkg/atomix/storage/protocol/gossip/primitive"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type SetElement struct {
	meta.ObjectMeta `protobuf:"bytes,1,opt,name=meta,proto3,embedded=meta" json:"meta"`
	Value           string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *SetElement) Reset()         { *m = SetElement{} }
func (m *SetElement) String() string { return proto.CompactTextString(m) }
func (*SetElement) ProtoMessage()    {}
func (*SetElement) Descriptor() ([]byte, []int) {
	return fileDescriptor_d4b8461bca3afde6, []int{0}
}
func (m *SetElement) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetElement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetElement.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetElement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetElement.Merge(m, src)
}
func (m *SetElement) XXX_Size() int {
	return m.Size()
}
func (m *SetElement) XXX_DiscardUnknown() {
	xxx_messageInfo_SetElement.DiscardUnknown(m)
}

var xxx_messageInfo_SetElement proto.InternalMessageInfo

func (m *SetElement) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*SetElement)(nil), "atomix.storage.protocol.gossip.set.SetElement")
}

func init() {
	proto.RegisterFile("atomix/storage/protocol/gossip/set/state.proto", fileDescriptor_d4b8461bca3afde6)
}

var fileDescriptor_d4b8461bca3afde6 = []byte{
	// 254 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x8f, 0xb1, 0x4a, 0xc4, 0x40,
	0x10, 0x86, 0xb3, 0x72, 0x0a, 0xae, 0x56, 0x87, 0x45, 0x48, 0xb1, 0x77, 0x5e, 0x95, 0x6a, 0x07,
	0xb4, 0xb2, 0x8d, 0x58, 0x8a, 0xa0, 0x4f, 0xb0, 0x77, 0x0c, 0x61, 0x25, 0x9b, 0x09, 0xd9, 0xf1,
	0xb8, 0xd2, 0x46, 0xb0, 0xf4, 0xb1, 0xae, 0x4c, 0x69, 0x75, 0x48, 0xf2, 0x22, 0x92, 0xdd, 0xa8,
	0xa5, 0xdd, 0x30, 0xff, 0xf7, 0xed, 0xfc, 0x2b, 0xb5, 0x61, 0x72, 0x76, 0x07, 0x9e, 0xa9, 0x35,
	0x25, 0x42, 0xd3, 0x12, 0xd3, 0x86, 0x2a, 0x28, 0xc9, 0x7b, 0xdb, 0x80, 0x47, 0x06, 0xcf, 0x86,
	0x51, 0x87, 0x64, 0xbe, 0x8a, 0xbc, 0x9e, 0x78, 0xfd, 0xc3, 0xeb, 0xc8, 0x6b, 0x8f, 0x9c, 0x4d,
	0x0c, 0x34, 0xad, 0x75, 0x96, 0xed, 0x16, 0xc1, 0x21, 0x1b, 0xa0, 0xf5, 0x33, 0x6e, 0x38, 0x1a,
	0xd9, 0xcd, 0x3f, 0x77, 0xff, 0x5c, 0xdc, 0x31, 0xd6, 0xde, 0x52, 0xed, 0x27, 0xf5, 0xa2, 0xa4,
	0x92, 0xc2, 0x08, 0xe3, 0x14, 0xb7, 0x2b, 0x27, 0xe5, 0x13, 0xf2, 0x5d, 0x85, 0x0e, 0x6b, 0x9e,
	0xdf, 0xca, 0xd9, 0x78, 0x33, 0x15, 0x4b, 0x91, 0x9f, 0x5d, 0x5d, 0x4e, 0xbf, 0xd4, 0xbf, 0xaf,
	0xea, 0x31, 0xd5, 0x0f, 0xa1, 0xd1, 0x3d, 0xb2, 0x29, 0xce, 0xf7, 0x87, 0x45, 0xd2, 0x1d, 0x16,
	0xe2, 0xfd, 0x2d, 0x17, 0x8f, 0x41, 0x9e, 0x67, 0xf2, 0x78, 0x6b, 0xaa, 0x17, 0x4c, 0x8f, 0x96,
	0x22, 0x3f, 0x2d, 0x66, 0xaf, 0x63, 0x14, 0x57, 0x45, 0xba, 0xef, 0x95, 0xe8, 0x7a, 0x25, 0xbe,
	0x7a, 0x25, 0x3e, 0x06, 0x95, 0x74, 0x83, 0x4a, 0x3e, 0x07, 0x95, 0xac, 0x4f, 0x42, 0x9f, 0xeb,
	0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc0, 0x85, 0x3c, 0xb6, 0x5a, 0x01, 0x00, 0x00,
}

func (m *SetElement) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetElement) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetElement) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintState(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.ObjectMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintState(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintState(dAtA []byte, offset int, v uint64) int {
	offset -= sovState(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SetElement) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ObjectMeta.Size()
	n += 1 + l + sovState(uint64(l))
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovState(uint64(l))
	}
	return n
}

func sovState(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozState(x uint64) (n int) {
	return sovState(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SetElement) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowState
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
			return fmt.Errorf("proto: SetElement: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetElement: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObjectMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
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
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ObjectMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowState
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
				return ErrInvalidLengthState
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthState
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipState(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthState
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthState
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
func skipState(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowState
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
					return 0, ErrIntOverflowState
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
					return 0, ErrIntOverflowState
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
				return 0, ErrInvalidLengthState
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupState
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthState
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthState        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowState          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupState = fmt.Errorf("proto: unexpected end of group")
)
