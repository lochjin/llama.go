// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hello.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type HelloWorldRequest struct {
	Referer              string   `protobuf:"bytes,1,opt,name=referer,proto3" json:"referer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloWorldRequest) Reset()         { *m = HelloWorldRequest{} }
func (m *HelloWorldRequest) String() string { return proto.CompactTextString(m) }
func (*HelloWorldRequest) ProtoMessage()    {}
func (*HelloWorldRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_hello_ba147fb7f43978ce, []int{0}
}
func (m *HelloWorldRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloWorldRequest.Unmarshal(m, b)
}
func (m *HelloWorldRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloWorldRequest.Marshal(b, m, deterministic)
}
func (dst *HelloWorldRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloWorldRequest.Merge(dst, src)
}
func (m *HelloWorldRequest) XXX_Size() int {
	return xxx_messageInfo_HelloWorldRequest.Size(m)
}
func (m *HelloWorldRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloWorldRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloWorldRequest proto.InternalMessageInfo

func (m *HelloWorldRequest) GetReferer() string {
	if m != nil {
		return m.Referer
	}
	return ""
}

type HelloWorldResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloWorldResponse) Reset()         { *m = HelloWorldResponse{} }
func (m *HelloWorldResponse) String() string { return proto.CompactTextString(m) }
func (*HelloWorldResponse) ProtoMessage()    {}
func (*HelloWorldResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_hello_ba147fb7f43978ce, []int{1}
}
func (m *HelloWorldResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloWorldResponse.Unmarshal(m, b)
}
func (m *HelloWorldResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloWorldResponse.Marshal(b, m, deterministic)
}
func (dst *HelloWorldResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloWorldResponse.Merge(dst, src)
}
func (m *HelloWorldResponse) XXX_Size() int {
	return xxx_messageInfo_HelloWorldResponse.Size(m)
}
func (m *HelloWorldResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloWorldResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HelloWorldResponse proto.InternalMessageInfo

func (m *HelloWorldResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloWorldRequest)(nil), "proto.HelloWorldRequest")
	proto.RegisterType((*HelloWorldResponse)(nil), "proto.HelloWorldResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// HelloWorldClient is the client API for HelloWorld service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type HelloWorldClient interface {
	SayHelloWorld(ctx context.Context, in *HelloWorldRequest, opts ...grpc.CallOption) (*HelloWorldResponse, error)
}

type helloWorldClient struct {
	cc *grpc.ClientConn
}

func NewHelloWorldClient(cc *grpc.ClientConn) HelloWorldClient {
	return &helloWorldClient{cc}
}

func (c *helloWorldClient) SayHelloWorld(ctx context.Context, in *HelloWorldRequest, opts ...grpc.CallOption) (*HelloWorldResponse, error) {
	out := new(HelloWorldResponse)
	err := c.cc.Invoke(ctx, "/proto.HelloWorld/SayHelloWorld", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloWorldServer is the server API for HelloWorld service.
type HelloWorldServer interface {
	SayHelloWorld(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error)
}

func RegisterHelloWorldServer(s *grpc.Server, srv HelloWorldServer) {
	s.RegisterService(&_HelloWorld_serviceDesc, srv)
}

func _HelloWorld_SayHelloWorld_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloWorldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloWorldServer).SayHelloWorld(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.HelloWorld/SayHelloWorld",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloWorldServer).SayHelloWorld(ctx, req.(*HelloWorldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _HelloWorld_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.HelloWorld",
	HandlerType: (*HelloWorldServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHelloWorld",
			Handler:    _HelloWorld_SayHelloWorld_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hello.proto",
}

func init() { proto.RegisterFile("hello.proto", fileDescriptor_hello_ba147fb7f43978ce) }

var fileDescriptor_hello_ba147fb7f43978ce = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0x48, 0xcd, 0xc9,
	0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0x32, 0xe9, 0xf9, 0xf9,
	0xe9, 0x39, 0xa9, 0xfa, 0x89, 0x05, 0x99, 0xfa, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99,
	0xf9, 0x79, 0xc5, 0x10, 0x45, 0x4a, 0xba, 0x5c, 0x82, 0x1e, 0x20, 0x3d, 0xe1, 0xf9, 0x45, 0x39,
	0x29, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x12, 0x5c, 0xec, 0x45, 0xa9, 0x69, 0xa9,
	0x45, 0xa9, 0x45, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x30, 0xae, 0x92, 0x1e, 0x97, 0x10,
	0xb2, 0xf2, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x90, 0xfa, 0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4,
	0x54, 0x98, 0x7a, 0x28, 0xd7, 0x28, 0x9b, 0x8b, 0x0b, 0xa1, 0x5e, 0x28, 0x96, 0x8b, 0x37, 0x38,
	0xb1, 0x12, 0x49, 0x40, 0x02, 0xe2, 0x0a, 0x3d, 0x0c, 0x27, 0x48, 0x49, 0x62, 0x91, 0x81, 0xd8,
	0xa6, 0x24, 0xde, 0x74, 0xf9, 0xc9, 0x64, 0x26, 0x41, 0x25, 0x1e, 0x7d, 0xb0, 0x6f, 0xe3, 0xcb,
	0x41, 0xb2, 0x56, 0x8c, 0x5a, 0x49, 0x6c, 0x60, 0x2d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x4c, 0xe7, 0x98, 0xc4, 0x06, 0x01, 0x00, 0x00,
}
