// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.1
// source: gateway/v1/gateway.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Gateway_GetVersion_FullMethodName     = "/api.gateway.v1.Gateway/GetVersion"
	Gateway_ListMinerSet_FullMethodName   = "/api.gateway.v1.Gateway/ListMinerSet"
	Gateway_CreateMinerSet_FullMethodName = "/api.gateway.v1.Gateway/CreateMinerSet"
	Gateway_CreateMiner_FullMethodName    = "/api.gateway.v1.Gateway/CreateMiner"
)

// GatewayClient is the client API for Gateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayClient interface {
	// GetVersion
	GetVersion(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetVersionResponse, error)
	// ListMinerSet
	ListMinerSet(ctx context.Context, in *ListMinerSetRequest, opts ...grpc.CallOption) (*ListMinerSetResponse, error)
	// CreateMinerSet
	CreateMinerSet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// CreateMiner
	CreateMiner(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type gatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayClient(cc grpc.ClientConnInterface) GatewayClient {
	return &gatewayClient{cc}
}

func (c *gatewayClient) GetVersion(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetVersionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetVersionResponse)
	err := c.cc.Invoke(ctx, Gateway_GetVersion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayClient) ListMinerSet(ctx context.Context, in *ListMinerSetRequest, opts ...grpc.CallOption) (*ListMinerSetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListMinerSetResponse)
	err := c.cc.Invoke(ctx, Gateway_ListMinerSet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayClient) CreateMinerSet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gateway_CreateMinerSet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gatewayClient) CreateMiner(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gateway_CreateMiner_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServer is the server API for Gateway service.
// All implementations must embed UnimplementedGatewayServer
// for forward compatibility.
type GatewayServer interface {
	// GetVersion
	GetVersion(context.Context, *emptypb.Empty) (*GetVersionResponse, error)
	// ListMinerSet
	ListMinerSet(context.Context, *ListMinerSetRequest) (*ListMinerSetResponse, error)
	// CreateMinerSet
	CreateMinerSet(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// CreateMiner
	CreateMiner(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedGatewayServer()
}

// UnimplementedGatewayServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGatewayServer struct{}

func (UnimplementedGatewayServer) GetVersion(context.Context, *emptypb.Empty) (*GetVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedGatewayServer) ListMinerSet(context.Context, *ListMinerSetRequest) (*ListMinerSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMinerSet not implemented")
}
func (UnimplementedGatewayServer) CreateMinerSet(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMinerSet not implemented")
}
func (UnimplementedGatewayServer) CreateMiner(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMiner not implemented")
}
func (UnimplementedGatewayServer) mustEmbedUnimplementedGatewayServer() {}
func (UnimplementedGatewayServer) testEmbeddedByValue()                 {}

// UnsafeGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayServer will
// result in compilation errors.
type UnsafeGatewayServer interface {
	mustEmbedUnimplementedGatewayServer()
}

func RegisterGatewayServer(s grpc.ServiceRegistrar, srv GatewayServer) {
	// If the following call pancis, it indicates UnimplementedGatewayServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Gateway_ServiceDesc, srv)
}

func _Gateway_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gateway_GetVersion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).GetVersion(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gateway_ListMinerSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMinerSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).ListMinerSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gateway_ListMinerSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).ListMinerSet(ctx, req.(*ListMinerSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gateway_CreateMinerSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).CreateMinerSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gateway_CreateMinerSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).CreateMinerSet(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gateway_CreateMiner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).CreateMiner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gateway_CreateMiner_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).CreateMiner(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Gateway_ServiceDesc is the grpc.ServiceDesc for Gateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.gateway.v1.Gateway",
	HandlerType: (*GatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _Gateway_GetVersion_Handler,
		},
		{
			MethodName: "ListMinerSet",
			Handler:    _Gateway_ListMinerSet_Handler,
		},
		{
			MethodName: "CreateMinerSet",
			Handler:    _Gateway_CreateMinerSet_Handler,
		},
		{
			MethodName: "CreateMiner",
			Handler:    _Gateway_CreateMiner_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway/v1/gateway.proto",
}
