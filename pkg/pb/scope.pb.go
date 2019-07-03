// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: scope.proto

package pb

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type Scope struct {
	SatelliteAddr        string            `protobuf:"bytes,1,opt,name=satellite_addr,json=satelliteAddr,proto3" json:"satellite_addr,omitempty"`
	ApiKey               []byte            `protobuf:"bytes,2,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	EncryptionAccess     *EncryptionAccess `protobuf:"bytes,3,opt,name=encryption_access,json=encryptionAccess,proto3" json:"encryption_access,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Scope) Reset()         { *m = Scope{} }
func (m *Scope) String() string { return proto.CompactTextString(m) }
func (*Scope) ProtoMessage()    {}
func (*Scope) Descriptor() ([]byte, []int) {
	return fileDescriptor_c67276d5d71daf81, []int{0}
}
func (m *Scope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Scope.Unmarshal(m, b)
}
func (m *Scope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Scope.Marshal(b, m, deterministic)
}
func (m *Scope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Scope.Merge(m, src)
}
func (m *Scope) XXX_Size() int {
	return xxx_messageInfo_Scope.Size(m)
}
func (m *Scope) XXX_DiscardUnknown() {
	xxx_messageInfo_Scope.DiscardUnknown(m)
}

var xxx_messageInfo_Scope proto.InternalMessageInfo

func (m *Scope) GetSatelliteAddr() string {
	if m != nil {
		return m.SatelliteAddr
	}
	return ""
}

func (m *Scope) GetApiKey() []byte {
	if m != nil {
		return m.ApiKey
	}
	return nil
}

func (m *Scope) GetEncryptionAccess() *EncryptionAccess {
	if m != nil {
		return m.EncryptionAccess
	}
	return nil
}

func init() {
	proto.RegisterType((*Scope)(nil), "scope.Scope")
}

func init() { proto.RegisterFile("scope.proto", fileDescriptor_c67276d5d71daf81) }

var fileDescriptor_c67276d5d71daf81 = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x4e, 0xce, 0x2f,
	0x48, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0xa4, 0xb8, 0xd2, 0xf3, 0xd3,
	0xf3, 0x21, 0x42, 0x52, 0xe2, 0xa9, 0x79, 0xc9, 0x45, 0x95, 0x05, 0x25, 0x99, 0xf9, 0x79, 0xf1,
	0x89, 0xc9, 0xc9, 0xa9, 0xc5, 0xc5, 0x10, 0x09, 0xa5, 0x99, 0x8c, 0x5c, 0xac, 0xc1, 0x20, 0xe5,
	0x42, 0xaa, 0x5c, 0x7c, 0xc5, 0x89, 0x25, 0xa9, 0x39, 0x39, 0x99, 0x25, 0xa9, 0xf1, 0x89, 0x29,
	0x29, 0x45, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0xbc, 0x70, 0x51, 0xc7, 0x94, 0x94, 0x22,
	0x21, 0x71, 0x2e, 0xf6, 0xc4, 0x82, 0xcc, 0xf8, 0xec, 0xd4, 0x4a, 0x09, 0x26, 0x05, 0x46, 0x0d,
	0x9e, 0x20, 0xb6, 0xc4, 0x82, 0x4c, 0xef, 0xd4, 0x4a, 0xa1, 0x00, 0x2e, 0x41, 0x0c, 0x4b, 0x24,
	0x98, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0x94, 0xf5, 0x30, 0xad, 0x77, 0x85, 0x8b, 0x38, 0x82, 0x05,
	0x82, 0x04, 0x52, 0xd1, 0x44, 0x9c, 0x58, 0xa2, 0x98, 0x0a, 0x92, 0x92, 0xd8, 0xc0, 0x0e, 0x35,
	0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x83, 0x7e, 0x90, 0x36, 0xe3, 0x00, 0x00, 0x00,
}
