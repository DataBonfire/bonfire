// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: examples/planet/api/planet/v1/influencer.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Influencer_CreateInfluencer_FullMethodName = "/examples.planet.api.influencer.v1.Influencer/CreateInfluencer"
	Influencer_UpdateInfluencer_FullMethodName = "/examples.planet.api.influencer.v1.Influencer/UpdateInfluencer"
	Influencer_DeleteInfluencer_FullMethodName = "/examples.planet.api.influencer.v1.Influencer/DeleteInfluencer"
	Influencer_GetInfluencer_FullMethodName    = "/examples.planet.api.influencer.v1.Influencer/GetInfluencer"
	Influencer_ListInfluencer_FullMethodName   = "/examples.planet.api.influencer.v1.Influencer/ListInfluencer"
)

// InfluencerClient is the client API for Influencer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InfluencerClient interface {
	CreateInfluencer(ctx context.Context, in *CreateInfluencerRequest, opts ...grpc.CallOption) (*CreateInfluencerReply, error)
	UpdateInfluencer(ctx context.Context, in *UpdateInfluencerRequest, opts ...grpc.CallOption) (*UpdateInfluencerReply, error)
	DeleteInfluencer(ctx context.Context, in *DeleteInfluencerRequest, opts ...grpc.CallOption) (*DeleteInfluencerReply, error)
	GetInfluencer(ctx context.Context, in *GetInfluencerRequest, opts ...grpc.CallOption) (*GetInfluencerReply, error)
	ListInfluencer(ctx context.Context, in *ListInfluencerRequest, opts ...grpc.CallOption) (*ListInfluencerReply, error)
}

type influencerClient struct {
	cc grpc.ClientConnInterface
}

func NewInfluencerClient(cc grpc.ClientConnInterface) InfluencerClient {
	return &influencerClient{cc}
}

func (c *influencerClient) CreateInfluencer(ctx context.Context, in *CreateInfluencerRequest, opts ...grpc.CallOption) (*CreateInfluencerReply, error) {
	out := new(CreateInfluencerReply)
	err := c.cc.Invoke(ctx, Influencer_CreateInfluencer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *influencerClient) UpdateInfluencer(ctx context.Context, in *UpdateInfluencerRequest, opts ...grpc.CallOption) (*UpdateInfluencerReply, error) {
	out := new(UpdateInfluencerReply)
	err := c.cc.Invoke(ctx, Influencer_UpdateInfluencer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *influencerClient) DeleteInfluencer(ctx context.Context, in *DeleteInfluencerRequest, opts ...grpc.CallOption) (*DeleteInfluencerReply, error) {
	out := new(DeleteInfluencerReply)
	err := c.cc.Invoke(ctx, Influencer_DeleteInfluencer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *influencerClient) GetInfluencer(ctx context.Context, in *GetInfluencerRequest, opts ...grpc.CallOption) (*GetInfluencerReply, error) {
	out := new(GetInfluencerReply)
	err := c.cc.Invoke(ctx, Influencer_GetInfluencer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *influencerClient) ListInfluencer(ctx context.Context, in *ListInfluencerRequest, opts ...grpc.CallOption) (*ListInfluencerReply, error) {
	out := new(ListInfluencerReply)
	err := c.cc.Invoke(ctx, Influencer_ListInfluencer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InfluencerServer is the server API for Influencer service.
// All implementations must embed UnimplementedInfluencerServer
// for forward compatibility
type InfluencerServer interface {
	CreateInfluencer(context.Context, *CreateInfluencerRequest) (*CreateInfluencerReply, error)
	UpdateInfluencer(context.Context, *UpdateInfluencerRequest) (*UpdateInfluencerReply, error)
	DeleteInfluencer(context.Context, *DeleteInfluencerRequest) (*DeleteInfluencerReply, error)
	GetInfluencer(context.Context, *GetInfluencerRequest) (*GetInfluencerReply, error)
	ListInfluencer(context.Context, *ListInfluencerRequest) (*ListInfluencerReply, error)
	mustEmbedUnimplementedInfluencerServer()
}

// UnimplementedInfluencerServer must be embedded to have forward compatible implementations.
type UnimplementedInfluencerServer struct {
}

func (UnimplementedInfluencerServer) CreateInfluencer(context.Context, *CreateInfluencerRequest) (*CreateInfluencerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInfluencer not implemented")
}
func (UnimplementedInfluencerServer) UpdateInfluencer(context.Context, *UpdateInfluencerRequest) (*UpdateInfluencerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateInfluencer not implemented")
}
func (UnimplementedInfluencerServer) DeleteInfluencer(context.Context, *DeleteInfluencerRequest) (*DeleteInfluencerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteInfluencer not implemented")
}
func (UnimplementedInfluencerServer) GetInfluencer(context.Context, *GetInfluencerRequest) (*GetInfluencerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfluencer not implemented")
}
func (UnimplementedInfluencerServer) ListInfluencer(context.Context, *ListInfluencerRequest) (*ListInfluencerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInfluencer not implemented")
}
func (UnimplementedInfluencerServer) mustEmbedUnimplementedInfluencerServer() {}

// UnsafeInfluencerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InfluencerServer will
// result in compilation errors.
type UnsafeInfluencerServer interface {
	mustEmbedUnimplementedInfluencerServer()
}

func RegisterInfluencerServer(s grpc.ServiceRegistrar, srv InfluencerServer) {
	s.RegisterService(&Influencer_ServiceDesc, srv)
}

func _Influencer_CreateInfluencer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInfluencerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfluencerServer).CreateInfluencer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Influencer_CreateInfluencer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfluencerServer).CreateInfluencer(ctx, req.(*CreateInfluencerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Influencer_UpdateInfluencer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateInfluencerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfluencerServer).UpdateInfluencer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Influencer_UpdateInfluencer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfluencerServer).UpdateInfluencer(ctx, req.(*UpdateInfluencerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Influencer_DeleteInfluencer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInfluencerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfluencerServer).DeleteInfluencer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Influencer_DeleteInfluencer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfluencerServer).DeleteInfluencer(ctx, req.(*DeleteInfluencerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Influencer_GetInfluencer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfluencerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfluencerServer).GetInfluencer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Influencer_GetInfluencer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfluencerServer).GetInfluencer(ctx, req.(*GetInfluencerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Influencer_ListInfluencer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInfluencerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfluencerServer).ListInfluencer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Influencer_ListInfluencer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfluencerServer).ListInfluencer(ctx, req.(*ListInfluencerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Influencer_ServiceDesc is the grpc.ServiceDesc for Influencer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Influencer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examples.planet.api.influencer.v1.Influencer",
	HandlerType: (*InfluencerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateInfluencer",
			Handler:    _Influencer_CreateInfluencer_Handler,
		},
		{
			MethodName: "UpdateInfluencer",
			Handler:    _Influencer_UpdateInfluencer_Handler,
		},
		{
			MethodName: "DeleteInfluencer",
			Handler:    _Influencer_DeleteInfluencer_Handler,
		},
		{
			MethodName: "GetInfluencer",
			Handler:    _Influencer_GetInfluencer_Handler,
		},
		{
			MethodName: "ListInfluencer",
			Handler:    _Influencer_ListInfluencer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/planet/api/planet/v1/influencer.proto",
}
