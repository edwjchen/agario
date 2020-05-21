// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blob.proto

/*
Package blob is a generated protocol buffer package.

It is generated from these files:
	blob.proto

It has these top-level messages:
	BlobRequest
	BlobResponse
*/
package blob

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

type BlobRequest struct {
	Id uint32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	X  float64 `protobuf:"fixed64,2,opt,name=x" json:"x,omitempty"`
	Y  float64 `protobuf:"fixed64,3,opt,name=y" json:"y,omitempty"`
}

func (m *BlobRequest) Reset()                    { *m = BlobRequest{} }
func (m *BlobRequest) String() string            { return proto.CompactTextString(m) }
func (*BlobRequest) ProtoMessage()               {}
func (*BlobRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *BlobRequest) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *BlobRequest) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *BlobRequest) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

type BlobResponse struct {
	X       float64 `protobuf:"fixed64,1,opt,name=x" json:"x,omitempty"`
	Y       float64 `protobuf:"fixed64,2,opt,name=y" json:"y,omitempty"`
	Alive   bool    `protobuf:"varint,3,opt,name=alive" json:"alive,omitempty"`
	Mass    int32   `protobuf:"varint,4,opt,name=mass" json:"mass,omitempty"`
	Players []byte  `protobuf:"bytes,5,opt,name=players,proto3" json:"players,omitempty"`
	Food    []byte  `protobuf:"bytes,6,opt,name=food,proto3" json:"food,omitempty"`
}

func (m *BlobResponse) Reset()                    { *m = BlobResponse{} }
func (m *BlobResponse) String() string            { return proto.CompactTextString(m) }
func (*BlobResponse) ProtoMessage()               {}
func (*BlobResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BlobResponse) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *BlobResponse) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *BlobResponse) GetAlive() bool {
	if m != nil {
		return m.Alive
	}
	return false
}

func (m *BlobResponse) GetMass() int32 {
	if m != nil {
		return m.Mass
	}
	return 0
}

func (m *BlobResponse) GetPlayers() []byte {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *BlobResponse) GetFood() []byte {
	if m != nil {
		return m.Food
	}
	return nil
}

func init() {
	proto.RegisterType((*BlobRequest)(nil), "blob.BlobRequest")
	proto.RegisterType((*BlobResponse)(nil), "blob.BlobResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Blob service

type BlobClient interface {
	Move(ctx context.Context, in *BlobRequest, opts ...grpc.CallOption) (*BlobResponse, error)
}

type blobClient struct {
	cc *grpc.ClientConn
}

func NewBlobClient(cc *grpc.ClientConn) BlobClient {
	return &blobClient{cc}
}

func (c *blobClient) Move(ctx context.Context, in *BlobRequest, opts ...grpc.CallOption) (*BlobResponse, error) {
	out := new(BlobResponse)
	err := grpc.Invoke(ctx, "/blob.Blob/Move", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Blob service

type BlobServer interface {
	Move(context.Context, *BlobRequest) (*BlobResponse, error)
}

func RegisterBlobServer(s *grpc.Server, srv BlobServer) {
	s.RegisterService(&_Blob_serviceDesc, srv)
}

func _Blob_Move_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobServer).Move(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blob.Blob/Move",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobServer).Move(ctx, req.(*BlobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Blob_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blob.Blob",
	HandlerType: (*BlobServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Move",
			Handler:    _Blob_Move_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blob.proto",
}

func init() { proto.RegisterFile("blob.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4b, 0xc7, 0x30,
	0x10, 0xc5, 0xbd, 0x9a, 0x56, 0x39, 0xab, 0xe0, 0xe1, 0x10, 0x9c, 0x42, 0xa7, 0x4c, 0x15, 0x74,
	0x10, 0x57, 0x77, 0x97, 0x7c, 0x83, 0x86, 0x46, 0x28, 0x44, 0xaf, 0x36, 0xb5, 0x34, 0xab, 0x9f,
	0x5c, 0x92, 0x28, 0xf4, 0xbf, 0xbd, 0xdf, 0xf1, 0x1e, 0xbc, 0x77, 0x88, 0xd6, 0xb3, 0xed, 0xe7,
	0x85, 0x57, 0x26, 0x91, 0x74, 0xf7, 0x82, 0x57, 0xaf, 0x9e, 0xad, 0x71, 0x5f, 0xdf, 0x2e, 0xac,
	0x74, 0x83, 0xd5, 0x34, 0x4a, 0x50, 0xa0, 0xaf, 0x4d, 0x35, 0x8d, 0xd4, 0x22, 0xec, 0xb2, 0x52,
	0xa0, 0xc1, 0xc0, 0x9e, 0x28, 0xca, 0xf3, 0x42, 0xb1, 0xfb, 0x01, 0x6c, 0x4b, 0x36, 0xcc, 0xfc,
	0x19, 0x5c, 0x31, 0xc3, 0x89, 0xf9, 0x2f, 0x1a, 0xe9, 0x0e, 0xeb, 0xc1, 0x4f, 0x9b, 0xcb, 0xf1,
	0x4b, 0x53, 0x80, 0x08, 0xc5, 0xc7, 0x10, 0x82, 0x14, 0x0a, 0x74, 0x6d, 0xb2, 0x26, 0x89, 0x17,
	0xb3, 0x1f, 0xa2, 0x5b, 0x82, 0xac, 0x15, 0xe8, 0xd6, 0xfc, 0x63, 0x72, 0xbf, 0x33, 0x8f, 0xb2,
	0xc9, 0xe7, 0xac, 0x1f, 0x9f, 0x51, 0xa4, 0x0e, 0xf4, 0x80, 0xe2, 0x8d, 0x37, 0x47, 0xb7, 0x7d,
	0x9e, 0x78, 0xd8, 0x74, 0x4f, 0xc7, 0x53, 0xa9, 0xda, 0x9d, 0xd9, 0x26, 0x7f, 0xe1, 0xe9, 0x37,
	0x00, 0x00, 0xff, 0xff, 0x12, 0x17, 0xbb, 0x57, 0x13, 0x01, 0x00, 0x00,
}
