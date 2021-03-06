// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blob.proto

/*
Package blob is a generated protocol buffer package.

It is generated from these files:
	blob.proto

It has these top-level messages:
	InitRequest
	InitResponse
	MoveRequest
	MoveResponse
	RegionRequest
	RegionResponse
	Player
	Food
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

type InitRequest struct {
}

func (m *InitRequest) Reset()                    { *m = InitRequest{} }
func (m *InitRequest) String() string            { return proto.CompactTextString(m) }
func (*InitRequest) ProtoMessage()               {}
func (*InitRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type InitResponse struct {
	Id   string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	X    float64 `protobuf:"fixed64,2,opt,name=x" json:"x,omitempty"`
	Y    float64 `protobuf:"fixed64,3,opt,name=y" json:"y,omitempty"`
	Mass int32   `protobuf:"varint,4,opt,name=mass" json:"mass,omitempty"`
}

func (m *InitResponse) Reset()                    { *m = InitResponse{} }
func (m *InitResponse) String() string            { return proto.CompactTextString(m) }
func (*InitResponse) ProtoMessage()               {}
func (*InitResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *InitResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *InitResponse) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *InitResponse) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *InitResponse) GetMass() int32 {
	if m != nil {
		return m.Mass
	}
	return 0
}

type MoveRequest struct {
	Id string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	X  float64 `protobuf:"fixed64,2,opt,name=x" json:"x,omitempty"`
	Y  float64 `protobuf:"fixed64,3,opt,name=y" json:"y,omitempty"`
}

func (m *MoveRequest) Reset()                    { *m = MoveRequest{} }
func (m *MoveRequest) String() string            { return proto.CompactTextString(m) }
func (*MoveRequest) ProtoMessage()               {}
func (*MoveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MoveRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *MoveRequest) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *MoveRequest) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

type MoveResponse struct {
	X     float64 `protobuf:"fixed64,1,opt,name=x" json:"x,omitempty"`
	Y     float64 `protobuf:"fixed64,2,opt,name=y" json:"y,omitempty"`
	Alive bool    `protobuf:"varint,3,opt,name=alive" json:"alive,omitempty"`
	Mass  int32   `protobuf:"varint,4,opt,name=mass" json:"mass,omitempty"`
}

func (m *MoveResponse) Reset()                    { *m = MoveResponse{} }
func (m *MoveResponse) String() string            { return proto.CompactTextString(m) }
func (*MoveResponse) ProtoMessage()               {}
func (*MoveResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MoveResponse) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *MoveResponse) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *MoveResponse) GetAlive() bool {
	if m != nil {
		return m.Alive
	}
	return false
}

func (m *MoveResponse) GetMass() int32 {
	if m != nil {
		return m.Mass
	}
	return 0
}

type RegionRequest struct {
	Id string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	X  float64 `protobuf:"fixed64,2,opt,name=x" json:"x,omitempty"`
	Y  float64 `protobuf:"fixed64,3,opt,name=y" json:"y,omitempty"`
}

func (m *RegionRequest) Reset()                    { *m = RegionRequest{} }
func (m *RegionRequest) String() string            { return proto.CompactTextString(m) }
func (*RegionRequest) ProtoMessage()               {}
func (*RegionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RegionRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RegionRequest) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *RegionRequest) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

type RegionResponse struct {
	Players []*Player `protobuf:"bytes,1,rep,name=players" json:"players,omitempty"`
	Foods   []*Food   `protobuf:"bytes,2,rep,name=foods" json:"foods,omitempty"`
}

func (m *RegionResponse) Reset()                    { *m = RegionResponse{} }
func (m *RegionResponse) String() string            { return proto.CompactTextString(m) }
func (*RegionResponse) ProtoMessage()               {}
func (*RegionResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RegionResponse) GetPlayers() []*Player {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *RegionResponse) GetFoods() []*Food {
	if m != nil {
		return m.Foods
	}
	return nil
}

type Player struct {
	Id    string  `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	X     float64 `protobuf:"fixed64,2,opt,name=x" json:"x,omitempty"`
	Y     float64 `protobuf:"fixed64,3,opt,name=y" json:"y,omitempty"`
	Alive bool    `protobuf:"varint,4,opt,name=alive" json:"alive,omitempty"`
	Mass  int32   `protobuf:"varint,5,opt,name=mass" json:"mass,omitempty"`
}

func (m *Player) Reset()                    { *m = Player{} }
func (m *Player) String() string            { return proto.CompactTextString(m) }
func (*Player) ProtoMessage()               {}
func (*Player) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Player) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Player) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Player) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Player) GetAlive() bool {
	if m != nil {
		return m.Alive
	}
	return false
}

func (m *Player) GetMass() int32 {
	if m != nil {
		return m.Mass
	}
	return 0
}

type Food struct {
	X float64 `protobuf:"fixed64,1,opt,name=x" json:"x,omitempty"`
	Y float64 `protobuf:"fixed64,2,opt,name=y" json:"y,omitempty"`
}

func (m *Food) Reset()                    { *m = Food{} }
func (m *Food) String() string            { return proto.CompactTextString(m) }
func (*Food) ProtoMessage()               {}
func (*Food) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Food) GetX() float64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Food) GetY() float64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func init() {
	proto.RegisterType((*InitRequest)(nil), "blob.InitRequest")
	proto.RegisterType((*InitResponse)(nil), "blob.InitResponse")
	proto.RegisterType((*MoveRequest)(nil), "blob.MoveRequest")
	proto.RegisterType((*MoveResponse)(nil), "blob.MoveResponse")
	proto.RegisterType((*RegionRequest)(nil), "blob.RegionRequest")
	proto.RegisterType((*RegionResponse)(nil), "blob.RegionResponse")
	proto.RegisterType((*Player)(nil), "blob.Player")
	proto.RegisterType((*Food)(nil), "blob.Food")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Blob service

type BlobClient interface {
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error)
	Move(ctx context.Context, in *MoveRequest, opts ...grpc.CallOption) (*MoveResponse, error)
	Region(ctx context.Context, in *RegionRequest, opts ...grpc.CallOption) (*RegionResponse, error)
}

type blobClient struct {
	cc *grpc.ClientConn
}

func NewBlobClient(cc *grpc.ClientConn) BlobClient {
	return &blobClient{cc}
}

func (c *blobClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	out := new(InitResponse)
	err := grpc.Invoke(ctx, "/blob.Blob/Init", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobClient) Move(ctx context.Context, in *MoveRequest, opts ...grpc.CallOption) (*MoveResponse, error) {
	out := new(MoveResponse)
	err := grpc.Invoke(ctx, "/blob.Blob/Move", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blobClient) Region(ctx context.Context, in *RegionRequest, opts ...grpc.CallOption) (*RegionResponse, error) {
	out := new(RegionResponse)
	err := grpc.Invoke(ctx, "/blob.Blob/Region", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Blob service

type BlobServer interface {
	Init(context.Context, *InitRequest) (*InitResponse, error)
	Move(context.Context, *MoveRequest) (*MoveResponse, error)
	Region(context.Context, *RegionRequest) (*RegionResponse, error)
}

func RegisterBlobServer(s *grpc.Server, srv BlobServer) {
	s.RegisterService(&_Blob_serviceDesc, srv)
}

func _Blob_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blob.Blob/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blob_Move_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveRequest)
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
		return srv.(BlobServer).Move(ctx, req.(*MoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blob_Region_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlobServer).Region(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blob.Blob/Region",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlobServer).Region(ctx, req.(*RegionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Blob_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blob.Blob",
	HandlerType: (*BlobServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _Blob_Init_Handler,
		},
		{
			MethodName: "Move",
			Handler:    _Blob_Move_Handler,
		},
		{
			MethodName: "Region",
			Handler:    _Blob_Region_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blob.proto",
}

func init() { proto.RegisterFile("blob.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 318 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0xcf, 0x4b, 0xc3, 0x30,
	0x14, 0x5e, 0xba, 0x74, 0xea, 0x5b, 0x37, 0x30, 0xee, 0x10, 0x76, 0x2a, 0x39, 0x48, 0x4f, 0x13,
	0x26, 0x1e, 0xc4, 0x9b, 0x07, 0xc1, 0x83, 0x22, 0xc1, 0x93, 0x07, 0xa1, 0xa5, 0x51, 0x0a, 0xb5,
	0xa9, 0x4b, 0x1d, 0xeb, 0x5f, 0xe3, 0xbf, 0x2a, 0xc9, 0x6b, 0xa4, 0x85, 0x21, 0xcc, 0x5b, 0xde,
	0x8f, 0x2f, 0xef, 0x7b, 0xdf, 0xf7, 0x00, 0xb2, 0x52, 0x67, 0xab, 0x7a, 0xa3, 0x1b, 0xcd, 0xa8,
	0x7d, 0x8b, 0x19, 0x4c, 0xef, 0xab, 0xa2, 0x91, 0xea, 0xf3, 0x4b, 0x99, 0x46, 0x3c, 0x42, 0x84,
	0xa1, 0xa9, 0x75, 0x65, 0x14, 0x9b, 0x43, 0x50, 0xe4, 0x9c, 0xc4, 0x24, 0x39, 0x91, 0x41, 0x91,
	0xb3, 0x08, 0xc8, 0x8e, 0x07, 0x31, 0x49, 0x88, 0x24, 0x3b, 0x1b, 0xb5, 0x7c, 0x8c, 0x51, 0xcb,
	0x18, 0xd0, 0x8f, 0xd4, 0x18, 0x4e, 0x63, 0x92, 0x84, 0xd2, 0xbd, 0xc5, 0x35, 0x4c, 0x1f, 0xf4,
	0x56, 0x75, 0xdf, 0x1f, 0xf2, 0x9d, 0x78, 0x86, 0x08, 0xa1, 0x1d, 0x15, 0xd7, 0x4b, 0x06, 0xbd,
	0x81, 0x1f, 0xbd, 0x80, 0x30, 0x2d, 0x8b, 0xad, 0x72, 0xe8, 0x63, 0x89, 0xc1, 0x5e, 0x42, 0x37,
	0x30, 0x93, 0xea, 0xbd, 0xd0, 0xd5, 0x7f, 0x28, 0xbd, 0xc0, 0xdc, 0x83, 0x3b, 0x52, 0xe7, 0x70,
	0x54, 0x97, 0x69, 0xab, 0x36, 0x86, 0x93, 0x78, 0x9c, 0x4c, 0xd7, 0xd1, 0xca, 0x49, 0xfc, 0xe4,
	0x92, 0xd2, 0x17, 0x59, 0x0c, 0xe1, 0x9b, 0xd6, 0xb9, 0xe1, 0x81, 0xeb, 0x02, 0xec, 0xba, 0xd3,
	0x3a, 0x97, 0x58, 0x10, 0xaf, 0x30, 0x41, 0xd0, 0x41, 0x9a, 0xff, 0x2e, 0x4e, 0xf7, 0x2d, 0x1e,
	0xf6, 0x16, 0x17, 0x40, 0xed, 0xb8, 0xbf, 0x64, 0x5c, 0x7f, 0x13, 0xa0, 0xb7, 0xa5, 0xce, 0xd8,
	0x05, 0x50, 0x7b, 0x06, 0xec, 0x14, 0x79, 0xf6, 0x2e, 0x64, 0xc9, 0xfa, 0x29, 0x54, 0x41, 0x8c,
	0x2c, 0xc0, 0x9a, 0xe5, 0x01, 0x3d, 0xcf, 0x3d, 0xa0, 0xef, 0xa5, 0x18, 0xb1, 0x2b, 0x98, 0xa0,
	0x94, 0xec, 0x0c, 0xeb, 0x03, 0x57, 0x96, 0x8b, 0x61, 0xd2, 0xc3, 0xb2, 0x89, 0xbb, 0xdd, 0xcb,
	0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xda, 0x23, 0x6b, 0x37, 0xc9, 0x02, 0x00, 0x00,
}
