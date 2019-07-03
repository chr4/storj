// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: orders.proto

package pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import time "time"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// PieceAction is an enumeration of all possible executed actions on storage node
type PieceAction int32

const (
	PieceAction_INVALID    PieceAction = 0
	PieceAction_PUT        PieceAction = 1
	PieceAction_GET        PieceAction = 2
	PieceAction_GET_AUDIT  PieceAction = 3
	PieceAction_GET_REPAIR PieceAction = 4
	PieceAction_PUT_REPAIR PieceAction = 5
	PieceAction_DELETE     PieceAction = 6
)

var PieceAction_name = map[int32]string{
	0: "INVALID",
	1: "PUT",
	2: "GET",
	3: "GET_AUDIT",
	4: "GET_REPAIR",
	5: "PUT_REPAIR",
	6: "DELETE",
}
var PieceAction_value = map[string]int32{
	"INVALID":    0,
	"PUT":        1,
	"GET":        2,
	"GET_AUDIT":  3,
	"GET_REPAIR": 4,
	"PUT_REPAIR": 5,
	"DELETE":     6,
}

func (x PieceAction) String() string {
	return proto.EnumName(PieceAction_name, int32(x))
}
func (PieceAction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{0}
}

type SettlementResponse_Status int32

const (
	SettlementResponse_INVALID  SettlementResponse_Status = 0
	SettlementResponse_ACCEPTED SettlementResponse_Status = 1
	SettlementResponse_REJECTED SettlementResponse_Status = 2
)

var SettlementResponse_Status_name = map[int32]string{
	0: "INVALID",
	1: "ACCEPTED",
	2: "REJECTED",
}
var SettlementResponse_Status_value = map[string]int32{
	"INVALID":  0,
	"ACCEPTED": 1,
	"REJECTED": 2,
}

func (x SettlementResponse_Status) String() string {
	return proto.EnumName(SettlementResponse_Status_name, int32(x))
}
func (SettlementResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{4, 0}
}

// OrderLimit is provided by satellite to execute specific action on storage node within some limits
type OrderLimit struct {
	// unique serial to avoid replay attacks
	SerialNumber SerialNumber `protobuf:"bytes,1,opt,name=serial_number,json=serialNumber,proto3,customtype=SerialNumber" json:"serial_number"`
	// satellite who issued this order limit allowing orderer to do the specified action
	SatelliteId NodeID `protobuf:"bytes,2,opt,name=satellite_id,json=satelliteId,proto3,customtype=NodeID" json:"satellite_id"`
	// uplink who requested or whom behalf the order limit to do an action
	UplinkId NodeID `protobuf:"bytes,3,opt,name=uplink_id,json=uplinkId,proto3,customtype=NodeID" json:"uplink_id"`
	// storage node who can reclaim the order limit specified by serial
	StorageNodeId NodeID `protobuf:"bytes,4,opt,name=storage_node_id,json=storageNodeId,proto3,customtype=NodeID" json:"storage_node_id"`
	// piece which is allowed to be touched
	PieceId PieceID `protobuf:"bytes,5,opt,name=piece_id,json=pieceId,proto3,customtype=PieceID" json:"piece_id"`
	// limit in bytes how much can be changed
	Limit              int64                `protobuf:"varint,6,opt,name=limit,proto3" json:"limit,omitempty"`
	Action             PieceAction          `protobuf:"varint,7,opt,name=action,proto3,enum=orders.PieceAction" json:"action,omitempty"`
	PieceExpiration    *timestamp.Timestamp `protobuf:"bytes,8,opt,name=piece_expiration,json=pieceExpiration,proto3" json:"piece_expiration,omitempty"`
	OrderExpiration    *timestamp.Timestamp `protobuf:"bytes,9,opt,name=order_expiration,json=orderExpiration,proto3" json:"order_expiration,omitempty"`
	OrderCreation      time.Time            `protobuf:"bytes,12,opt,name=order_creation,json=orderCreation,proto3,stdtime" json:"order_creation"`
	SatelliteSignature []byte               `protobuf:"bytes,10,opt,name=satellite_signature,json=satelliteSignature,proto3" json:"satellite_signature,omitempty"`
	// satellites aren't necessarily discoverable in kademlia. this allows
	// a storage node to find a satellite and handshake with it to get its key.
	SatelliteAddress     *NodeAddress `protobuf:"bytes,11,opt,name=satellite_address,json=satelliteAddress,proto3" json:"satellite_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *OrderLimit) Reset()         { *m = OrderLimit{} }
func (m *OrderLimit) String() string { return proto.CompactTextString(m) }
func (*OrderLimit) ProtoMessage()    {}
func (*OrderLimit) Descriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{0}
}
func (m *OrderLimit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderLimit.Unmarshal(m, b)
}
func (m *OrderLimit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderLimit.Marshal(b, m, deterministic)
}
func (dst *OrderLimit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderLimit.Merge(dst, src)
}
func (m *OrderLimit) XXX_Size() int {
	return xxx_messageInfo_OrderLimit.Size(m)
}
func (m *OrderLimit) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderLimit.DiscardUnknown(m)
}

var xxx_messageInfo_OrderLimit proto.InternalMessageInfo

func (m *OrderLimit) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *OrderLimit) GetAction() PieceAction {
	if m != nil {
		return m.Action
	}
	return PieceAction_INVALID
}

func (m *OrderLimit) GetPieceExpiration() *timestamp.Timestamp {
	if m != nil {
		return m.PieceExpiration
	}
	return nil
}

func (m *OrderLimit) GetOrderExpiration() *timestamp.Timestamp {
	if m != nil {
		return m.OrderExpiration
	}
	return nil
}

func (m *OrderLimit) GetOrderCreation() time.Time {
	if m != nil {
		return m.OrderCreation
	}
	return time.Time{}
}

func (m *OrderLimit) GetSatelliteSignature() []byte {
	if m != nil {
		return m.SatelliteSignature
	}
	return nil
}

func (m *OrderLimit) GetSatelliteAddress() *NodeAddress {
	if m != nil {
		return m.SatelliteAddress
	}
	return nil
}

// Order is a one step of fullfilling Amount number of bytes from an OrderLimit with SerialNumber
type Order struct {
	// serial of the order limit that was signed
	SerialNumber SerialNumber `protobuf:"bytes,1,opt,name=serial_number,json=serialNumber,proto3,customtype=SerialNumber" json:"serial_number"`
	// amount to be signed for
	Amount int64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
	// signature
	UplinkSignature      []byte   `protobuf:"bytes,3,opt,name=uplink_signature,json=uplinkSignature,proto3" json:"uplink_signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{1}
}
func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (dst *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(dst, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *Order) GetUplinkSignature() []byte {
	if m != nil {
		return m.UplinkSignature
	}
	return nil
}

type PieceHash struct {
	// piece id
	PieceId PieceID `protobuf:"bytes,1,opt,name=piece_id,json=pieceId,proto3,customtype=PieceID" json:"piece_id"`
	// hash of the piece that was/is uploaded
	Hash []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	// signature either satellite or storage node
	Signature []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	// size of uploaded piece
	PieceSize int64 `protobuf:"varint,4,opt,name=piece_size,json=pieceSize,proto3" json:"piece_size,omitempty"`
	// timestamp when upload occur
	Timestamp            time.Time `protobuf:"bytes,5,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PieceHash) Reset()         { *m = PieceHash{} }
func (m *PieceHash) String() string { return proto.CompactTextString(m) }
func (*PieceHash) ProtoMessage()    {}
func (*PieceHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{2}
}
func (m *PieceHash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PieceHash.Unmarshal(m, b)
}
func (m *PieceHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PieceHash.Marshal(b, m, deterministic)
}
func (dst *PieceHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PieceHash.Merge(dst, src)
}
func (m *PieceHash) XXX_Size() int {
	return xxx_messageInfo_PieceHash.Size(m)
}
func (m *PieceHash) XXX_DiscardUnknown() {
	xxx_messageInfo_PieceHash.DiscardUnknown(m)
}

var xxx_messageInfo_PieceHash proto.InternalMessageInfo

func (m *PieceHash) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *PieceHash) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *PieceHash) GetPieceSize() int64 {
	if m != nil {
		return m.PieceSize
	}
	return 0
}

func (m *PieceHash) GetTimestamp() time.Time {
	if m != nil {
		return m.Timestamp
	}
	return time.Time{}
}

type SettlementRequest struct {
	Limit                *OrderLimit `protobuf:"bytes,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Order                *Order      `protobuf:"bytes,2,opt,name=order,proto3" json:"order,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *SettlementRequest) Reset()         { *m = SettlementRequest{} }
func (m *SettlementRequest) String() string { return proto.CompactTextString(m) }
func (*SettlementRequest) ProtoMessage()    {}
func (*SettlementRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{3}
}
func (m *SettlementRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SettlementRequest.Unmarshal(m, b)
}
func (m *SettlementRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SettlementRequest.Marshal(b, m, deterministic)
}
func (dst *SettlementRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SettlementRequest.Merge(dst, src)
}
func (m *SettlementRequest) XXX_Size() int {
	return xxx_messageInfo_SettlementRequest.Size(m)
}
func (m *SettlementRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SettlementRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SettlementRequest proto.InternalMessageInfo

func (m *SettlementRequest) GetLimit() *OrderLimit {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *SettlementRequest) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

type SettlementResponse struct {
	SerialNumber         SerialNumber              `protobuf:"bytes,1,opt,name=serial_number,json=serialNumber,proto3,customtype=SerialNumber" json:"serial_number"`
	Status               SettlementResponse_Status `protobuf:"varint,2,opt,name=status,proto3,enum=orders.SettlementResponse_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *SettlementResponse) Reset()         { *m = SettlementResponse{} }
func (m *SettlementResponse) String() string { return proto.CompactTextString(m) }
func (*SettlementResponse) ProtoMessage()    {}
func (*SettlementResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_orders_3cfba8bd45227bd2, []int{4}
}
func (m *SettlementResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SettlementResponse.Unmarshal(m, b)
}
func (m *SettlementResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SettlementResponse.Marshal(b, m, deterministic)
}
func (dst *SettlementResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SettlementResponse.Merge(dst, src)
}
func (m *SettlementResponse) XXX_Size() int {
	return xxx_messageInfo_SettlementResponse.Size(m)
}
func (m *SettlementResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SettlementResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SettlementResponse proto.InternalMessageInfo

func (m *SettlementResponse) GetStatus() SettlementResponse_Status {
	if m != nil {
		return m.Status
	}
	return SettlementResponse_INVALID
}

func init() {
	proto.RegisterType((*OrderLimit)(nil), "orders.OrderLimit")
	proto.RegisterType((*Order)(nil), "orders.Order")
	proto.RegisterType((*PieceHash)(nil), "orders.PieceHash")
	proto.RegisterType((*SettlementRequest)(nil), "orders.SettlementRequest")
	proto.RegisterType((*SettlementResponse)(nil), "orders.SettlementResponse")
	proto.RegisterEnum("orders.PieceAction", PieceAction_name, PieceAction_value)
	proto.RegisterEnum("orders.SettlementResponse_Status", SettlementResponse_Status_name, SettlementResponse_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrdersClient is the client API for Orders service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrdersClient interface {
	Settlement(ctx context.Context, opts ...grpc.CallOption) (Orders_SettlementClient, error)
}

type ordersClient struct {
	cc *grpc.ClientConn
}

func NewOrdersClient(cc *grpc.ClientConn) OrdersClient {
	return &ordersClient{cc}
}

func (c *ordersClient) Settlement(ctx context.Context, opts ...grpc.CallOption) (Orders_SettlementClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Orders_serviceDesc.Streams[0], "/orders.Orders/Settlement", opts...)
	if err != nil {
		return nil, err
	}
	x := &ordersSettlementClient{stream}
	return x, nil
}

type Orders_SettlementClient interface {
	Send(*SettlementRequest) error
	Recv() (*SettlementResponse, error)
	grpc.ClientStream
}

type ordersSettlementClient struct {
	grpc.ClientStream
}

func (x *ordersSettlementClient) Send(m *SettlementRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ordersSettlementClient) Recv() (*SettlementResponse, error) {
	m := new(SettlementResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OrdersServer is the server API for Orders service.
type OrdersServer interface {
	Settlement(Orders_SettlementServer) error
}

func RegisterOrdersServer(s *grpc.Server, srv OrdersServer) {
	s.RegisterService(&_Orders_serviceDesc, srv)
}

func _Orders_Settlement_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OrdersServer).Settlement(&ordersSettlementServer{stream})
}

type Orders_SettlementServer interface {
	Send(*SettlementResponse) error
	Recv() (*SettlementRequest, error)
	grpc.ServerStream
}

type ordersSettlementServer struct {
	grpc.ServerStream
}

func (x *ordersSettlementServer) Send(m *SettlementResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ordersSettlementServer) Recv() (*SettlementRequest, error) {
	m := new(SettlementRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Orders_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orders.Orders",
	HandlerType: (*OrdersServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Settlement",
			Handler:       _Orders_Settlement_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "orders.proto",
}

func init() { proto.RegisterFile("orders.proto", fileDescriptor_orders_3cfba8bd45227bd2) }

var fileDescriptor_orders_3cfba8bd45227bd2 = []byte{
	// 727 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xed, 0xe6, 0xc3, 0x49, 0x26, 0x5f, 0xee, 0xb6, 0x42, 0x21, 0x02, 0x25, 0x84, 0x4b, 0x68,
	0xa5, 0x94, 0x06, 0x09, 0xa9, 0x17, 0xa4, 0x7c, 0x58, 0xc5, 0x50, 0x95, 0x68, 0x93, 0x72, 0xe0,
	0x12, 0x39, 0xf5, 0x92, 0x5a, 0x24, 0x76, 0xf0, 0x6e, 0x24, 0xd4, 0x3b, 0x77, 0xce, 0xfc, 0x17,
	0xee, 0x1c, 0xf8, 0x05, 0x1c, 0xca, 0x5f, 0x41, 0x3b, 0xeb, 0xc4, 0x29, 0x14, 0x15, 0xa9, 0xb7,
	0x7d, 0x33, 0xef, 0xcd, 0x78, 0x76, 0xde, 0x1a, 0x0a, 0x41, 0xe8, 0xf2, 0x50, 0xb4, 0x16, 0x61,
	0x20, 0x03, 0x6a, 0x68, 0x54, 0x85, 0x69, 0x30, 0x0d, 0x74, 0xac, 0x5a, 0x9b, 0x06, 0xc1, 0x74,
	0xc6, 0x0f, 0x10, 0x4d, 0x96, 0xef, 0x0f, 0xa4, 0x37, 0xe7, 0x42, 0x3a, 0xf3, 0x45, 0x44, 0x00,
	0x3f, 0x70, 0xb9, 0x3e, 0x37, 0xbe, 0xa6, 0x01, 0xde, 0xa8, 0x1a, 0x27, 0xde, 0xdc, 0x93, 0xf4,
	0x08, 0x8a, 0x82, 0x87, 0x9e, 0x33, 0x1b, 0xfb, 0xcb, 0xf9, 0x84, 0x87, 0x15, 0x52, 0x27, 0xcd,
	0x42, 0x77, 0xf7, 0xfb, 0x55, 0x6d, 0xeb, 0xe7, 0x55, 0xad, 0x30, 0xc4, 0xe4, 0x29, 0xe6, 0x58,
	0x41, 0x6c, 0x20, 0x7a, 0x08, 0x05, 0xe1, 0x48, 0x3e, 0x9b, 0x79, 0x92, 0x8f, 0x3d, 0xb7, 0x92,
	0x40, 0x65, 0x29, 0x52, 0x1a, 0xa7, 0x81, 0xcb, 0xed, 0x3e, 0xcb, 0xaf, 0x39, 0xb6, 0x4b, 0xf7,
	0x21, 0xb7, 0x5c, 0xcc, 0x3c, 0xff, 0x83, 0xe2, 0x27, 0x6f, 0xe4, 0x67, 0x35, 0xc1, 0x76, 0xe9,
	0x73, 0x28, 0x0b, 0x19, 0x84, 0xce, 0x94, 0x8f, 0xd5, 0xf7, 0x2b, 0x49, 0xea, 0x46, 0x49, 0x31,
	0xa2, 0x21, 0x74, 0xe9, 0x1e, 0x64, 0x17, 0x1e, 0x3f, 0x47, 0x41, 0x1a, 0x05, 0xe5, 0x48, 0x90,
	0x19, 0xa8, 0xb8, 0xdd, 0x67, 0x19, 0x24, 0xd8, 0x2e, 0xdd, 0x85, 0xf4, 0x4c, 0xdd, 0x43, 0xc5,
	0xa8, 0x93, 0x66, 0x92, 0x69, 0x40, 0xf7, 0xc1, 0x70, 0xce, 0xa5, 0x17, 0xf8, 0x95, 0x4c, 0x9d,
	0x34, 0x4b, 0xed, 0x9d, 0x56, 0xb4, 0x03, 0xd4, 0x77, 0x30, 0xc5, 0x22, 0x0a, 0xb5, 0xc0, 0xd4,
	0xed, 0xf8, 0xa7, 0x85, 0x17, 0x3a, 0x28, 0xcb, 0xd6, 0x49, 0x33, 0xdf, 0xae, 0xb6, 0xf4, 0x62,
	0x5a, 0xab, 0xc5, 0xb4, 0x46, 0xab, 0xc5, 0xb0, 0x32, 0x6a, 0xac, 0xb5, 0x44, 0x95, 0xc1, 0x26,
	0x9b, 0x65, 0x72, 0xb7, 0x97, 0x41, 0xcd, 0x46, 0x99, 0xd7, 0x50, 0xd2, 0x65, 0xce, 0x43, 0xae,
	0x8b, 0x14, 0x6e, 0x2b, 0xd2, 0xcd, 0xaa, 0xeb, 0xf9, 0xf2, 0xab, 0x46, 0x58, 0x11, 0xb5, 0xbd,
	0x48, 0x4a, 0x0f, 0x60, 0x27, 0xde, 0xb0, 0xf0, 0xa6, 0xbe, 0x23, 0x97, 0x21, 0xaf, 0x80, 0xba,
	0x54, 0x46, 0xd7, 0xa9, 0xe1, 0x2a, 0x43, 0x5f, 0xc0, 0x76, 0x2c, 0x70, 0x5c, 0x37, 0xe4, 0x42,
	0x54, 0xf2, 0xf8, 0x01, 0xdb, 0x2d, 0x34, 0xa1, 0xda, 0x51, 0x47, 0x27, 0x98, 0xb9, 0xe6, 0x46,
	0x91, 0xc6, 0x67, 0x02, 0x69, 0x34, 0xe7, 0x5d, 0x7c, 0x79, 0x0f, 0x0c, 0x67, 0x1e, 0x2c, 0x7d,
	0x89, 0x8e, 0x4c, 0xb2, 0x08, 0xd1, 0x27, 0x60, 0x46, 0xe6, 0x8b, 0x47, 0x41, 0x0f, 0xb2, 0xb2,
	0x8e, 0xaf, 0xe7, 0x68, 0xfc, 0x20, 0x90, 0xc3, 0x5d, 0xbf, 0x74, 0xc4, 0xc5, 0x35, 0x43, 0x91,
	0x5b, 0x0c, 0x45, 0x21, 0x75, 0xe1, 0x88, 0x0b, 0xfd, 0x18, 0x18, 0x9e, 0xe9, 0x03, 0xc8, 0xfd,
	0xd9, 0x31, 0x0e, 0xd0, 0x87, 0x00, 0xba, 0xba, 0xf0, 0x2e, 0x39, 0x3a, 0x3c, 0xc9, 0x72, 0x18,
	0x19, 0x7a, 0x97, 0x9c, 0x76, 0x21, 0xb7, 0x7e, 0xce, 0x68, 0xe7, 0xff, 0xdd, 0x65, 0x2c, 0x6b,
	0x4c, 0x60, 0x7b, 0xc8, 0xa5, 0x9c, 0xf1, 0x39, 0xf7, 0x25, 0xe3, 0x1f, 0x97, 0x5c, 0x48, 0xda,
	0x5c, 0x59, 0x9f, 0x60, 0x51, 0xba, 0xf2, 0x78, 0xfc, 0x73, 0x58, 0x3d, 0x87, 0xc7, 0x90, 0xc6,
	0x1c, 0x0e, 0x95, 0x6f, 0x17, 0xaf, 0x31, 0x99, 0xce, 0x35, 0xbe, 0x11, 0xa0, 0x9b, 0x4d, 0xc4,
	0x22, 0xf0, 0x05, 0xbf, 0xcb, 0x1e, 0x8f, 0xc0, 0x10, 0xd2, 0x91, 0x4b, 0x81, 0x7d, 0x4b, 0xed,
	0x47, 0xab, 0xbe, 0x7f, 0xb7, 0x69, 0x0d, 0x91, 0xc8, 0x22, 0x41, 0xe3, 0x10, 0x0c, 0x1d, 0xa1,
	0x79, 0xc8, 0xd8, 0xa7, 0x6f, 0x3b, 0x27, 0x76, 0xdf, 0xdc, 0xa2, 0x05, 0xc8, 0x76, 0x7a, 0x3d,
	0x6b, 0x30, 0xb2, 0xfa, 0x26, 0x51, 0x88, 0x59, 0xaf, 0xac, 0x9e, 0x42, 0x89, 0xbd, 0x29, 0xe4,
	0x37, 0x5e, 0xf7, 0x75, 0x5d, 0x06, 0x92, 0x83, 0xb3, 0x91, 0x49, 0xd4, 0xe1, 0xd8, 0x1a, 0x99,
	0x09, 0x5a, 0x84, 0xdc, 0xb1, 0x35, 0x1a, 0x77, 0xce, 0xfa, 0xf6, 0xc8, 0x4c, 0xd2, 0x12, 0x80,
	0x82, 0xcc, 0x1a, 0x74, 0x6c, 0x66, 0xa6, 0x14, 0x1e, 0x9c, 0xad, 0x71, 0x9a, 0x02, 0x18, 0x7d,
	0xeb, 0xc4, 0x1a, 0x59, 0xa6, 0xd1, 0x1e, 0x82, 0x81, 0x17, 0x27, 0xa8, 0x0d, 0x10, 0x8f, 0x42,
	0xef, 0xdf, 0x34, 0x1e, 0xae, 0xaa, 0x5a, 0xfd, 0xf7, 0xe4, 0x8d, 0xad, 0x26, 0x79, 0x4a, 0xba,
	0xa9, 0x77, 0x89, 0xc5, 0x64, 0x62, 0xa0, 0x21, 0x9e, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x71,
	0xd3, 0xb9, 0x60, 0x33, 0x06, 0x00, 0x00,
}
