// Code generated by protoc-gen-go.
// source: gdalinfo.proto
// DO NOT EDIT!

/*
Package gdalinfo is a generated protocol buffer package.

It is generated from these files:
	gdalinfo.proto

It has these top-level messages:
	Request
	GDALFile
*/
package gdalinfo

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

type Request struct {
	FilePath string `protobuf:"bytes,1,opt,name=filePath" json:"filePath,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

type GDALFile struct {
	FileName string `protobuf:"bytes,1,opt,name=fileName" json:"fileName,omitempty"`
	Output   string `protobuf:"bytes,2,opt,name=output" json:"output,omitempty"`
}

func (m *GDALFile) Reset()                    { *m = GDALFile{} }
func (m *GDALFile) String() string            { return proto.CompactTextString(m) }
func (*GDALFile) ProtoMessage()               {}
func (*GDALFile) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GDALFile) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *GDALFile) GetOutput() string {
	if m != nil {
		return m.Output
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "gdalinfo.Request")
	proto.RegisterType((*GDALFile)(nil), "gdalinfo.GDALFile")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GDALInfo service

type GDALInfoClient interface {
	Extract(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GDALFile, error)
}

type gDALInfoClient struct {
	cc *grpc.ClientConn
}

func NewGDALInfoClient(cc *grpc.ClientConn) GDALInfoClient {
	return &gDALInfoClient{cc}
}

func (c *gDALInfoClient) Extract(ctx context.Context, in *Request, opts ...grpc.CallOption) (*GDALFile, error) {
	out := new(GDALFile)
	err := grpc.Invoke(ctx, "/gdalinfo.GDALInfo/Extract", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for GDALInfo service

type GDALInfoServer interface {
	Extract(context.Context, *Request) (*GDALFile, error)
}

func RegisterGDALInfoServer(s *grpc.Server, srv GDALInfoServer) {
	s.RegisterService(&_GDALInfo_serviceDesc, srv)
}

func _GDALInfo_Extract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GDALInfoServer).Extract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdalinfo.GDALInfo/Extract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GDALInfoServer).Extract(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _GDALInfo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gdalinfo.GDALInfo",
	HandlerType: (*GDALInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Extract",
			Handler:    _GDALInfo_Extract_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gdalinfo.proto",
}

func init() { proto.RegisterFile("gdalinfo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 144 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x4f, 0x49, 0xcc,
	0xc9, 0xcc, 0x4b, 0xcb, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0xf1, 0x95, 0xa4,
	0xb9, 0xd8, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x04, 0xb8, 0x38, 0xd2, 0x32, 0x73,
	0x52, 0x03, 0x12, 0x4b, 0x32, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x95, 0x74, 0xb8, 0x38, 0xdc,
	0x5d, 0x1c, 0x7d, 0xdc, 0x32, 0x73, 0x52, 0x61, 0xb2, 0x7e, 0x89, 0xb9, 0xa9, 0x10, 0x59, 0x21,
	0x3e, 0x2e, 0xb6, 0xfc, 0xd2, 0x92, 0x82, 0xd2, 0x12, 0x09, 0x26, 0x10, 0xdf, 0xc8, 0x06, 0xa2,
	0xda, 0x33, 0x2f, 0x2d, 0x5f, 0xc8, 0x80, 0x8b, 0xdd, 0xb5, 0xa2, 0xa4, 0x28, 0x31, 0xb9, 0x44,
	0x48, 0x50, 0x0f, 0x6e, 0x39, 0xd4, 0x26, 0x29, 0x21, 0x84, 0x10, 0xcc, 0xfc, 0x24, 0x36, 0xb0,
	0xcb, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x26, 0xc5, 0x1c, 0xab, 0x00, 0x00, 0x00,
}
