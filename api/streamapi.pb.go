// Code generated by protoc-gen-go.
// source: streamapi.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	streamapi.proto

It has these top-level messages:
	OpenStreamArguments
	Message
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type OpenStreamArguments struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Offset int64  `protobuf:"varint,2,opt,name=offset" json:"offset,omitempty"`
}

func (m *OpenStreamArguments) Reset()                    { *m = OpenStreamArguments{} }
func (m *OpenStreamArguments) String() string            { return proto.CompactTextString(m) }
func (*OpenStreamArguments) ProtoMessage()               {}
func (*OpenStreamArguments) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Message struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*OpenStreamArguments)(nil), "api.OpenStreamArguments")
	proto.RegisterType((*Message)(nil), "api.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for StreamApi service

type StreamApiClient interface {
	OpenStream(ctx context.Context, in *OpenStreamArguments, opts ...grpc.CallOption) (StreamApi_OpenStreamClient, error)
}

type streamApiClient struct {
	cc *grpc.ClientConn
}

func NewStreamApiClient(cc *grpc.ClientConn) StreamApiClient {
	return &streamApiClient{cc}
}

func (c *streamApiClient) OpenStream(ctx context.Context, in *OpenStreamArguments, opts ...grpc.CallOption) (StreamApi_OpenStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_StreamApi_serviceDesc.Streams[0], c.cc, "/api.StreamApi/OpenStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamApiOpenStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamApi_OpenStreamClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type streamApiOpenStreamClient struct {
	grpc.ClientStream
}

func (x *streamApiOpenStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for StreamApi service

type StreamApiServer interface {
	OpenStream(*OpenStreamArguments, StreamApi_OpenStreamServer) error
}

func RegisterStreamApiServer(s *grpc.Server, srv StreamApiServer) {
	s.RegisterService(&_StreamApi_serviceDesc, srv)
}

func _StreamApi_OpenStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OpenStreamArguments)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamApiServer).OpenStream(m, &streamApiOpenStreamServer{stream})
}

type StreamApi_OpenStreamServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type streamApiOpenStreamServer struct {
	grpc.ServerStream
}

func (x *streamApiOpenStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

var _StreamApi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.StreamApi",
	HandlerType: (*StreamApiServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OpenStream",
			Handler:       _StreamApi_OpenStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "streamapi.proto",
}

func init() { proto.RegisterFile("streamapi.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 163 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x2e, 0x29, 0x4a,
	0x4d, 0xcc, 0x4d, 0x2c, 0xc8, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8,
	0x54, 0xb2, 0xe5, 0x12, 0xf6, 0x2f, 0x48, 0xcd, 0x0b, 0x06, 0xcb, 0x39, 0x16, 0xa5, 0x97, 0xe6,
	0xa6, 0xe6, 0x95, 0x14, 0x0b, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70,
	0x06, 0x31, 0x65, 0xa6, 0x08, 0x89, 0x71, 0xb1, 0xe5, 0xa7, 0xa5, 0x15, 0xa7, 0x96, 0x48, 0x30,
	0x29, 0x30, 0x6a, 0x30, 0x07, 0x41, 0x79, 0x4a, 0xb2, 0x5c, 0xec, 0xbe, 0xa9, 0xc5, 0xc5, 0x89,
	0xe9, 0xa9, 0x42, 0x42, 0x5c, 0x2c, 0x29, 0x89, 0x25, 0x89, 0x60, 0x4d, 0x3c, 0x41, 0x60, 0xb6,
	0x91, 0x2b, 0x17, 0x27, 0xd4, 0xe4, 0x82, 0x4c, 0x21, 0x0b, 0x2e, 0x2e, 0x84, 0x55, 0x42, 0x12,
	0x7a, 0x20, 0x97, 0x60, 0xb1, 0x5b, 0x8a, 0x07, 0x2c, 0x03, 0x35, 0x56, 0x89, 0xc1, 0x80, 0x31,
	0x89, 0x0d, 0xec, 0x60, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x65, 0x1f, 0x4a, 0xc3,
	0x00, 0x00, 0x00,
}
