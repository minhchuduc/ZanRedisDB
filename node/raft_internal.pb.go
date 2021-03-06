// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: raft_internal.proto

package node

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ReqSourceType int32

const (
	FromAPI           ReqSourceType = 0
	FromClusterSyncer ReqSourceType = 1
)

var ReqSourceType_name = map[int32]string{
	0: "FromAPI",
	1: "FromClusterSyncer",
}

var ReqSourceType_value = map[string]int32{
	"FromAPI":           0,
	"FromClusterSyncer": 1,
}

func (x ReqSourceType) String() string {
	return proto.EnumName(ReqSourceType_name, int32(x))
}

func (ReqSourceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b4c9a9be0cfca103, []int{0}
}

type SchemaChangeType int32

const (
	SchemaChangeAddHsetIndex    SchemaChangeType = 0
	SchemaChangeUpdateHsetIndex SchemaChangeType = 1
	SchemaChangeDeleteHsetIndex SchemaChangeType = 2
)

var SchemaChangeType_name = map[int32]string{
	0: "SchemaChangeAddHsetIndex",
	1: "SchemaChangeUpdateHsetIndex",
	2: "SchemaChangeDeleteHsetIndex",
}

var SchemaChangeType_value = map[string]int32{
	"SchemaChangeAddHsetIndex":    0,
	"SchemaChangeUpdateHsetIndex": 1,
	"SchemaChangeDeleteHsetIndex": 2,
}

func (x SchemaChangeType) String() string {
	return proto.EnumName(SchemaChangeType_name, int32(x))
}

func (SchemaChangeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b4c9a9be0cfca103, []int{1}
}

type RequestHeader struct {
	ID        uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	DataType  int32  `protobuf:"varint,2,opt,name=data_type,json=dataType,proto3" json:"data_type,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (m *RequestHeader) Reset()         { *m = RequestHeader{} }
func (m *RequestHeader) String() string { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()    {}
func (*RequestHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c9a9be0cfca103, []int{0}
}
func (m *RequestHeader) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RequestHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RequestHeader.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RequestHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestHeader.Merge(m, src)
}
func (m *RequestHeader) XXX_Size() int {
	return m.Size()
}
func (m *RequestHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestHeader.DiscardUnknown(m)
}

var xxx_messageInfo_RequestHeader proto.InternalMessageInfo

type InternalRaftRequest struct {
	Header RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header"`
	Data   []byte        `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *InternalRaftRequest) Reset()         { *m = InternalRaftRequest{} }
func (m *InternalRaftRequest) String() string { return proto.CompactTextString(m) }
func (*InternalRaftRequest) ProtoMessage()    {}
func (*InternalRaftRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c9a9be0cfca103, []int{1}
}
func (m *InternalRaftRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InternalRaftRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InternalRaftRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InternalRaftRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InternalRaftRequest.Merge(m, src)
}
func (m *InternalRaftRequest) XXX_Size() int {
	return m.Size()
}
func (m *InternalRaftRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InternalRaftRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InternalRaftRequest proto.InternalMessageInfo

type BatchInternalRaftRequest struct {
	ReqNum    int32                 `protobuf:"varint,1,opt,name=req_num,json=reqNum,proto3" json:"req_num,omitempty"`
	Reqs      []InternalRaftRequest `protobuf:"bytes,2,rep,name=reqs,proto3" json:"reqs"`
	Timestamp int64                 `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Type      ReqSourceType         `protobuf:"varint,4,opt,name=type,proto3,enum=node.ReqSourceType" json:"type,omitempty"`
	ReqId     uint64                `protobuf:"varint,5,opt,name=req_id,json=reqId,proto3" json:"req_id,omitempty"`
	// used for cluster log syncer
	OrigTerm    uint64 `protobuf:"varint,6,opt,name=orig_term,json=origTerm,proto3" json:"orig_term,omitempty"`
	OrigIndex   uint64 `protobuf:"varint,7,opt,name=orig_index,json=origIndex,proto3" json:"orig_index,omitempty"`
	OrigCluster string `protobuf:"bytes,8,opt,name=orig_cluster,json=origCluster,proto3" json:"orig_cluster,omitempty"`
}

func (m *BatchInternalRaftRequest) Reset()         { *m = BatchInternalRaftRequest{} }
func (m *BatchInternalRaftRequest) String() string { return proto.CompactTextString(m) }
func (*BatchInternalRaftRequest) ProtoMessage()    {}
func (*BatchInternalRaftRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c9a9be0cfca103, []int{2}
}
func (m *BatchInternalRaftRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BatchInternalRaftRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BatchInternalRaftRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BatchInternalRaftRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchInternalRaftRequest.Merge(m, src)
}
func (m *BatchInternalRaftRequest) XXX_Size() int {
	return m.Size()
}
func (m *BatchInternalRaftRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchInternalRaftRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BatchInternalRaftRequest proto.InternalMessageInfo

type SchemaChange struct {
	Type       SchemaChangeType `protobuf:"varint,1,opt,name=Type,proto3,enum=node.SchemaChangeType" json:"Type,omitempty"`
	Table      string           `protobuf:"bytes,2,opt,name=Table,proto3" json:"Table,omitempty"`
	SchemaData []byte           `protobuf:"bytes,3,opt,name=SchemaData,proto3" json:"SchemaData,omitempty"`
}

func (m *SchemaChange) Reset()         { *m = SchemaChange{} }
func (m *SchemaChange) String() string { return proto.CompactTextString(m) }
func (*SchemaChange) ProtoMessage()    {}
func (*SchemaChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c9a9be0cfca103, []int{3}
}
func (m *SchemaChange) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SchemaChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SchemaChange.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SchemaChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SchemaChange.Merge(m, src)
}
func (m *SchemaChange) XXX_Size() int {
	return m.Size()
}
func (m *SchemaChange) XXX_DiscardUnknown() {
	xxx_messageInfo_SchemaChange.DiscardUnknown(m)
}

var xxx_messageInfo_SchemaChange proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("node.ReqSourceType", ReqSourceType_name, ReqSourceType_value)
	proto.RegisterEnum("node.SchemaChangeType", SchemaChangeType_name, SchemaChangeType_value)
	proto.RegisterType((*RequestHeader)(nil), "node.RequestHeader")
	proto.RegisterType((*InternalRaftRequest)(nil), "node.InternalRaftRequest")
	proto.RegisterType((*BatchInternalRaftRequest)(nil), "node.BatchInternalRaftRequest")
	proto.RegisterType((*SchemaChange)(nil), "node.SchemaChange")
}

func init() { proto.RegisterFile("raft_internal.proto", fileDescriptor_b4c9a9be0cfca103) }

var fileDescriptor_b4c9a9be0cfca103 = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xcd, 0x6e, 0xda, 0x40,
	0x18, 0xf4, 0x1a, 0xf3, 0xf7, 0x41, 0x11, 0x5d, 0x92, 0xd6, 0x6d, 0x52, 0xc7, 0xe5, 0x52, 0x8b,
	0x03, 0x55, 0xe1, 0x09, 0x42, 0x50, 0x15, 0x5f, 0xaa, 0x6a, 0xa1, 0x97, 0xaa, 0x12, 0xda, 0xe0,
	0x2f, 0x80, 0x84, 0x7f, 0x58, 0x16, 0xa9, 0xbc, 0x45, 0x5f, 0xa2, 0xef, 0xc2, 0x91, 0x63, 0x4f,
	0x55, 0x03, 0x2f, 0x52, 0xed, 0xda, 0x12, 0x6e, 0x14, 0xf5, 0xb6, 0x3b, 0x33, 0xf2, 0xcc, 0x37,
	0x9f, 0x17, 0x5a, 0x82, 0xdf, 0xcb, 0xc9, 0x22, 0x92, 0x28, 0x22, 0xbe, 0xec, 0x26, 0x22, 0x96,
	0x31, 0xb5, 0xa2, 0x38, 0xc0, 0xd7, 0x67, 0xb3, 0x78, 0x16, 0x6b, 0xe0, 0xbd, 0x3a, 0xa5, 0x5c,
	0xfb, 0x2b, 0x3c, 0x63, 0xb8, 0xda, 0xe0, 0x5a, 0xde, 0x22, 0x0f, 0x50, 0xd0, 0x06, 0x98, 0xfe,
	0xd0, 0x26, 0x2e, 0xf1, 0x2c, 0x66, 0xfa, 0x43, 0x7a, 0x01, 0xd5, 0x80, 0x4b, 0x3e, 0x91, 0xdb,
	0x04, 0x6d, 0xd3, 0x25, 0x5e, 0x91, 0x55, 0x14, 0x30, 0xde, 0x26, 0x48, 0x2f, 0xa1, 0x2a, 0x17,
	0x21, 0xae, 0x25, 0x0f, 0x13, 0xbb, 0xe0, 0x12, 0xaf, 0xc0, 0x4e, 0x40, 0xfb, 0x1b, 0xb4, 0xfc,
	0x2c, 0x09, 0xe3, 0xf7, 0x32, 0xf3, 0xa1, 0x1f, 0xa0, 0x34, 0xd7, 0x5e, 0xda, 0xa5, 0xd6, 0x6b,
	0x75, 0x55, 0xbe, 0xee, 0x3f, 0x31, 0x06, 0xd6, 0xee, 0xf7, 0x95, 0xc1, 0x32, 0x21, 0xa5, 0x60,
	0x29, 0x4f, 0xed, 0x5f, 0x67, 0xfa, 0xdc, 0xfe, 0x69, 0x82, 0x3d, 0xe0, 0x72, 0x3a, 0x7f, 0xca,
	0xe3, 0x25, 0x94, 0x05, 0xae, 0x26, 0xd1, 0x26, 0xd4, 0x26, 0x45, 0x56, 0x12, 0xb8, 0xfa, 0xb4,
	0x09, 0x69, 0x1f, 0x2c, 0x81, 0xab, 0xb5, 0x6d, 0xba, 0x05, 0xaf, 0xd6, 0x7b, 0x95, 0x5a, 0x3f,
	0xf1, 0x85, 0x2c, 0x80, 0x16, 0xff, 0x7f, 0x4c, 0xfa, 0x0e, 0x2c, 0x5d, 0x8e, 0xe5, 0x12, 0xaf,
	0x91, 0x9b, 0x66, 0x14, 0x6f, 0xc4, 0x14, 0x55, 0x4f, 0x4c, 0x0b, 0xe8, 0x39, 0xa8, 0x14, 0x93,
	0x45, 0x60, 0x17, 0x75, 0xbd, 0x45, 0x81, 0x2b, 0x3f, 0x50, 0x0d, 0xc7, 0x62, 0x31, 0x9b, 0x48,
	0x14, 0xa1, 0x5d, 0xd2, 0x4c, 0x45, 0x01, 0x63, 0x14, 0x21, 0x7d, 0x03, 0xa0, 0xc9, 0x45, 0x14,
	0xe0, 0x77, 0xbb, 0xac, 0x59, 0x2d, 0xf7, 0x15, 0x40, 0xdf, 0x42, 0x5d, 0xd3, 0xd3, 0xe5, 0x66,
	0x2d, 0x51, 0xd8, 0x15, 0x97, 0x78, 0x55, 0x56, 0x53, 0xd8, 0x4d, 0x0a, 0xb5, 0x13, 0xa8, 0x8f,
	0xa6, 0x73, 0x0c, 0xf9, 0xcd, 0x9c, 0x47, 0x33, 0xa4, 0x1d, 0xb0, 0x54, 0x26, 0xdd, 0x4b, 0xa3,
	0xf7, 0x22, 0x8d, 0x9b, 0x57, 0xa4, 0x89, 0xf5, 0x7e, 0xcf, 0xa0, 0x38, 0xe6, 0x77, 0xcb, 0x74,
	0xf1, 0x55, 0x96, 0x5e, 0xa8, 0x03, 0x90, 0xea, 0x87, 0x6a, 0x27, 0x05, 0xbd, 0x93, 0x1c, 0xd2,
	0xe9, 0xeb, 0x7f, 0xea, 0x34, 0x3e, 0xad, 0x41, 0xf9, 0xa3, 0x88, 0xc3, 0xeb, 0xcf, 0x7e, 0xd3,
	0xa0, 0xe7, 0xf0, 0x5c, 0x5d, 0xb2, 0x78, 0xa3, 0x6d, 0x34, 0x45, 0xd1, 0x24, 0x1d, 0x01, 0xcd,
	0xc7, 0x21, 0xe8, 0x25, 0xd8, 0x79, 0xec, 0x3a, 0x08, 0x6e, 0xd7, 0x28, 0xf5, 0xe4, 0x4d, 0x83,
	0x5e, 0xc1, 0x45, 0x9e, 0xfd, 0x92, 0x04, 0x5c, 0xe2, 0x49, 0x40, 0x1e, 0x0b, 0x86, 0xb8, 0xc4,
	0xbc, 0xc0, 0x1c, 0xb8, 0xbb, 0x07, 0xc7, 0xd8, 0x3f, 0x38, 0xc6, 0xee, 0xe0, 0x90, 0xfd, 0xc1,
	0x21, 0x7f, 0x0e, 0x0e, 0xf9, 0x71, 0x74, 0x8c, 0xfd, 0xd1, 0x31, 0x7e, 0x1d, 0x1d, 0xe3, 0xae,
	0xa4, 0x5f, 0x49, 0xff, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3a, 0x85, 0x6c, 0x1a, 0x58, 0x03,
	0x00, 0x00,
}

func (m *RequestHeader) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RequestHeader) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ID != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.ID))
	}
	if m.DataType != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.DataType))
	}
	if m.Timestamp != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.Timestamp))
	}
	return i, nil
}

func (m *InternalRaftRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InternalRaftRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintRaftInternal(dAtA, i, uint64(m.Header.Size()))
	n1, err := m.Header.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if len(m.Data) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	return i, nil
}

func (m *BatchInternalRaftRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BatchInternalRaftRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ReqNum != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.ReqNum))
	}
	if len(m.Reqs) > 0 {
		for _, msg := range m.Reqs {
			dAtA[i] = 0x12
			i++
			i = encodeVarintRaftInternal(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Timestamp != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.Timestamp))
	}
	if m.Type != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.Type))
	}
	if m.ReqId != 0 {
		dAtA[i] = 0x28
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.ReqId))
	}
	if m.OrigTerm != 0 {
		dAtA[i] = 0x30
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.OrigTerm))
	}
	if m.OrigIndex != 0 {
		dAtA[i] = 0x38
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.OrigIndex))
	}
	if len(m.OrigCluster) > 0 {
		dAtA[i] = 0x42
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(len(m.OrigCluster)))
		i += copy(dAtA[i:], m.OrigCluster)
	}
	return i, nil
}

func (m *SchemaChange) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SchemaChange) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Type != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(m.Type))
	}
	if len(m.Table) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(len(m.Table)))
		i += copy(dAtA[i:], m.Table)
	}
	if len(m.SchemaData) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintRaftInternal(dAtA, i, uint64(len(m.SchemaData)))
		i += copy(dAtA[i:], m.SchemaData)
	}
	return i, nil
}

func encodeVarintRaftInternal(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RequestHeader) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovRaftInternal(uint64(m.ID))
	}
	if m.DataType != 0 {
		n += 1 + sovRaftInternal(uint64(m.DataType))
	}
	if m.Timestamp != 0 {
		n += 1 + sovRaftInternal(uint64(m.Timestamp))
	}
	return n
}

func (m *InternalRaftRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Header.Size()
	n += 1 + l + sovRaftInternal(uint64(l))
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovRaftInternal(uint64(l))
	}
	return n
}

func (m *BatchInternalRaftRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ReqNum != 0 {
		n += 1 + sovRaftInternal(uint64(m.ReqNum))
	}
	if len(m.Reqs) > 0 {
		for _, e := range m.Reqs {
			l = e.Size()
			n += 1 + l + sovRaftInternal(uint64(l))
		}
	}
	if m.Timestamp != 0 {
		n += 1 + sovRaftInternal(uint64(m.Timestamp))
	}
	if m.Type != 0 {
		n += 1 + sovRaftInternal(uint64(m.Type))
	}
	if m.ReqId != 0 {
		n += 1 + sovRaftInternal(uint64(m.ReqId))
	}
	if m.OrigTerm != 0 {
		n += 1 + sovRaftInternal(uint64(m.OrigTerm))
	}
	if m.OrigIndex != 0 {
		n += 1 + sovRaftInternal(uint64(m.OrigIndex))
	}
	l = len(m.OrigCluster)
	if l > 0 {
		n += 1 + l + sovRaftInternal(uint64(l))
	}
	return n
}

func (m *SchemaChange) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovRaftInternal(uint64(m.Type))
	}
	l = len(m.Table)
	if l > 0 {
		n += 1 + l + sovRaftInternal(uint64(l))
	}
	l = len(m.SchemaData)
	if l > 0 {
		n += 1 + l + sovRaftInternal(uint64(l))
	}
	return n
}

func sovRaftInternal(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRaftInternal(x uint64) (n int) {
	return sovRaftInternal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RequestHeader) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRaftInternal
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
			return fmt.Errorf("proto: RequestHeader: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RequestHeader: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataType", wireType)
			}
			m.DataType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DataType |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRaftInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRaftInternal
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
func (m *InternalRaftRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRaftInternal
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
			return fmt.Errorf("proto: InternalRaftRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InternalRaftRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Header", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
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
				return ErrInvalidLengthRaftInternal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Header.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
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
				return ErrInvalidLengthRaftInternal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRaftInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRaftInternal
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
func (m *BatchInternalRaftRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRaftInternal
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
			return fmt.Errorf("proto: BatchInternalRaftRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BatchInternalRaftRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReqNum", wireType)
			}
			m.ReqNum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReqNum |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reqs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
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
				return ErrInvalidLengthRaftInternal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Reqs = append(m.Reqs, InternalRaftRequest{})
			if err := m.Reqs[len(m.Reqs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ReqSourceType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReqId", wireType)
			}
			m.ReqId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReqId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigTerm", wireType)
			}
			m.OrigTerm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrigTerm |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigIndex", wireType)
			}
			m.OrigIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrigIndex |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrigCluster", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
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
				return ErrInvalidLengthRaftInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OrigCluster = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRaftInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRaftInternal
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
func (m *SchemaChange) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRaftInternal
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
			return fmt.Errorf("proto: SchemaChange: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SchemaChange: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= SchemaChangeType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Table", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
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
				return ErrInvalidLengthRaftInternal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Table = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SchemaData", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRaftInternal
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
				return ErrInvalidLengthRaftInternal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SchemaData = append(m.SchemaData[:0], dAtA[iNdEx:postIndex]...)
			if m.SchemaData == nil {
				m.SchemaData = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRaftInternal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRaftInternal
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRaftInternal
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
func skipRaftInternal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRaftInternal
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
					return 0, ErrIntOverflowRaftInternal
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRaftInternal
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
				return 0, ErrInvalidLengthRaftInternal
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthRaftInternal
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRaftInternal
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipRaftInternal(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthRaftInternal
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthRaftInternal = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRaftInternal   = fmt.Errorf("proto: integer overflow")
)
