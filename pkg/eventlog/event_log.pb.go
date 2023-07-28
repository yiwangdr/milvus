// Code generated by protoc-gen-go. DO NOT EDIT.
// source: event_log.proto

package eventlog

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Level int32

const (
	Level_Undefined Level = 0
	Level_Debug     Level = 1
	Level_Info      Level = 2
	Level_Warn      Level = 3
	Level_Error     Level = 4
)

var Level_name = map[int32]string{
	0: "Undefined",
	1: "Debug",
	2: "Info",
	3: "Warn",
	4: "Error",
}

var Level_value = map[string]int32{
	"Undefined": 0,
	"Debug":     1,
	"Info":      2,
	"Warn":      3,
	"Error":     4,
}

func (x Level) String() string {
	return proto.EnumName(Level_name, int32(x))
}

func (Level) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_443313318a2fd90c, []int{0}
}

type ListenRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListenRequest) Reset()         { *m = ListenRequest{} }
func (m *ListenRequest) String() string { return proto.CompactTextString(m) }
func (*ListenRequest) ProtoMessage()    {}
func (*ListenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_443313318a2fd90c, []int{0}
}

func (m *ListenRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListenRequest.Unmarshal(m, b)
}
func (m *ListenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListenRequest.Marshal(b, m, deterministic)
}
func (m *ListenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListenRequest.Merge(m, src)
}
func (m *ListenRequest) XXX_Size() int {
	return xxx_messageInfo_ListenRequest.Size(m)
}
func (m *ListenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListenRequest proto.InternalMessageInfo

type Event struct {
	Level                Level    `protobuf:"varint,1,opt,name=level,proto3,enum=milvus.proto.eventlog.Level" json:"level,omitempty"`
	Type                 int32    `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Ts                   int64    `protobuf:"varint,4,opt,name=ts,proto3" json:"ts,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_443313318a2fd90c, []int{1}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetLevel() Level {
	if m != nil {
		return m.Level
	}
	return Level_Undefined
}

func (m *Event) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Event) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event) GetTs() int64 {
	if m != nil {
		return m.Ts
	}
	return 0
}

func init() {
	proto.RegisterEnum("milvus.proto.eventlog.Level", Level_name, Level_value)
	proto.RegisterType((*ListenRequest)(nil), "milvus.proto.eventlog.ListenRequest")
	proto.RegisterType((*Event)(nil), "milvus.proto.eventlog.Event")
}

func init() { proto.RegisterFile("event_log.proto", fileDescriptor_443313318a2fd90c) }

var fileDescriptor_443313318a2fd90c = []byte{
	// 276 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0x5f, 0x4b, 0xc3, 0x30,
	0x14, 0xc5, 0x97, 0xae, 0x1d, 0xee, 0xe2, 0xb6, 0x12, 0x10, 0x8a, 0xf8, 0x50, 0x8a, 0x0f, 0x65,
	0x60, 0x2b, 0xf5, 0x0b, 0x88, 0xb8, 0x07, 0xa1, 0x0f, 0x52, 0x11, 0xc1, 0x17, 0xe9, 0x9f, 0xbb,
	0x18, 0xec, 0x92, 0x9a, 0xa4, 0x05, 0xbf, 0xbd, 0x34, 0x55, 0xf0, 0xc1, 0xbd, 0x9d, 0xe4, 0x9c,
	0xfb, 0xe3, 0x1c, 0xd8, 0xe0, 0x80, 0xc2, 0xbc, 0xb5, 0x92, 0x25, 0x9d, 0x92, 0x46, 0xd2, 0xb3,
	0x03, 0x6f, 0x87, 0x5e, 0x4f, 0xaf, 0xc4, 0xba, 0xad, 0x64, 0xd1, 0x06, 0x56, 0x39, 0xd7, 0x06,
	0x45, 0x81, 0x9f, 0x3d, 0x6a, 0x13, 0x69, 0xf0, 0x76, 0xa3, 0x49, 0x33, 0xf0, 0x5a, 0x1c, 0xb0,
	0x0d, 0x48, 0x48, 0xe2, 0x75, 0x76, 0x91, 0xfc, 0x0b, 0x48, 0xf2, 0x31, 0x53, 0x4c, 0x51, 0x4a,
	0xc1, 0x35, 0x5f, 0x1d, 0x06, 0x4e, 0x48, 0x62, 0xaf, 0xb0, 0x7a, 0xfc, 0x6b, 0x4a, 0x53, 0x06,
	0xf3, 0x90, 0xc4, 0xa7, 0x85, 0xd5, 0x74, 0x0d, 0x8e, 0xd1, 0x81, 0x1b, 0x92, 0x78, 0x5e, 0x38,
	0x46, 0x6f, 0x6f, 0xc1, 0xb3, 0x1c, 0xba, 0x82, 0xe5, 0xb3, 0x68, 0x70, 0xcf, 0x05, 0x36, 0xfe,
	0x8c, 0x2e, 0xc1, 0xbb, 0xc7, 0xaa, 0x67, 0x3e, 0xa1, 0x27, 0xe0, 0x3e, 0x88, 0xbd, 0xf4, 0x9d,
	0x51, 0xbd, 0x94, 0x4a, 0xf8, 0xf3, 0xd1, 0xde, 0x29, 0x25, 0x95, 0xef, 0x66, 0x35, 0x6c, 0x6c,
	0xed, 0x5c, 0xb2, 0x27, 0x54, 0x03, 0xaf, 0x91, 0x3e, 0xc2, 0x62, 0x9a, 0x46, 0x2f, 0x8f, 0x75,
	0xff, 0xbb, 0xfc, 0xfc, 0xd8, 0x42, 0xcb, 0x8d, 0x66, 0xd7, 0xe4, 0x6e, 0xfb, 0x1a, 0x33, 0x6e,
	0xde, 0xfb, 0x2a, 0xa9, 0xe5, 0x21, 0x9d, 0xd2, 0x57, 0x5c, 0xfe, 0xa8, 0xb4, 0xfb, 0x60, 0xe9,
	0xef, 0x55, 0xb5, 0xb0, 0x94, 0x9b, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x46, 0x60, 0x11,
	0x89, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EventLogServiceClient is the client API for EventLogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EventLogServiceClient interface {
	Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (EventLogService_ListenClient, error)
}

type eventLogServiceClient struct {
	cc *grpc.ClientConn
}

func NewEventLogServiceClient(cc *grpc.ClientConn) EventLogServiceClient {
	return &eventLogServiceClient{cc}
}

func (c *eventLogServiceClient) Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (EventLogService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EventLogService_serviceDesc.Streams[0], "/milvus.proto.eventlog.EventLogService/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventLogServiceListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EventLogService_ListenClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventLogServiceListenClient struct {
	grpc.ClientStream
}

func (x *eventLogServiceListenClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventLogServiceServer is the server API for EventLogService service.
type EventLogServiceServer interface {
	Listen(*ListenRequest, EventLogService_ListenServer) error
}

// UnimplementedEventLogServiceServer can be embedded to have forward compatible implementations.
type UnimplementedEventLogServiceServer struct {
}

func (*UnimplementedEventLogServiceServer) Listen(req *ListenRequest, srv EventLogService_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}

func RegisterEventLogServiceServer(s *grpc.Server, srv EventLogServiceServer) {
	s.RegisterService(&_EventLogService_serviceDesc, srv)
}

func _EventLogService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventLogServiceServer).Listen(m, &eventLogServiceListenServer{stream})
}

type EventLogService_ListenServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type eventLogServiceListenServer struct {
	grpc.ServerStream
}

func (x *eventLogServiceListenServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

var _EventLogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "milvus.proto.eventlog.EventLogService",
	HandlerType: (*EventLogServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _EventLogService_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "event_log.proto",
}
