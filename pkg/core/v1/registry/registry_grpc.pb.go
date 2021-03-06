// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package registry

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NetworkServiceEndpointRegistryClient is the client API for NetworkServiceEndpointRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkServiceEndpointRegistryClient interface {
	Register(ctx context.Context, in *NetworkServiceEndpoint, opts ...grpc.CallOption) (*NetworkServiceEndpoint, error)
	Find(ctx context.Context, in *NetworkServiceEndpointQuery, opts ...grpc.CallOption) (NetworkServiceEndpointRegistry_FindClient, error)
	Unregister(ctx context.Context, in *NetworkServiceEndpoint, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type networkServiceEndpointRegistryClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkServiceEndpointRegistryClient(cc grpc.ClientConnInterface) NetworkServiceEndpointRegistryClient {
	return &networkServiceEndpointRegistryClient{cc}
}

func (c *networkServiceEndpointRegistryClient) Register(ctx context.Context, in *NetworkServiceEndpoint, opts ...grpc.CallOption) (*NetworkServiceEndpoint, error) {
	out := new(NetworkServiceEndpoint)
	err := c.cc.Invoke(ctx, "/v1.registry.NetworkServiceEndpointRegistry/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceEndpointRegistryClient) Find(ctx context.Context, in *NetworkServiceEndpointQuery, opts ...grpc.CallOption) (NetworkServiceEndpointRegistry_FindClient, error) {
	stream, err := c.cc.NewStream(ctx, &NetworkServiceEndpointRegistry_ServiceDesc.Streams[0], "/v1.registry.NetworkServiceEndpointRegistry/Find", opts...)
	if err != nil {
		return nil, err
	}
	x := &networkServiceEndpointRegistryFindClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NetworkServiceEndpointRegistry_FindClient interface {
	Recv() (*NetworkServiceEndpointResponse, error)
	grpc.ClientStream
}

type networkServiceEndpointRegistryFindClient struct {
	grpc.ClientStream
}

func (x *networkServiceEndpointRegistryFindClient) Recv() (*NetworkServiceEndpointResponse, error) {
	m := new(NetworkServiceEndpointResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *networkServiceEndpointRegistryClient) Unregister(ctx context.Context, in *NetworkServiceEndpoint, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/v1.registry.NetworkServiceEndpointRegistry/Unregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServiceEndpointRegistryServer is the server API for NetworkServiceEndpointRegistry service.
// All implementations should embed UnimplementedNetworkServiceEndpointRegistryServer
// for forward compatibility
type NetworkServiceEndpointRegistryServer interface {
	Register(context.Context, *NetworkServiceEndpoint) (*NetworkServiceEndpoint, error)
	Find(*NetworkServiceEndpointQuery, NetworkServiceEndpointRegistry_FindServer) error
	Unregister(context.Context, *NetworkServiceEndpoint) (*emptypb.Empty, error)
}

// UnimplementedNetworkServiceEndpointRegistryServer should be embedded to have forward compatible implementations.
type UnimplementedNetworkServiceEndpointRegistryServer struct {
}

func (UnimplementedNetworkServiceEndpointRegistryServer) Register(context.Context, *NetworkServiceEndpoint) (*NetworkServiceEndpoint, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedNetworkServiceEndpointRegistryServer) Find(*NetworkServiceEndpointQuery, NetworkServiceEndpointRegistry_FindServer) error {
	return status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedNetworkServiceEndpointRegistryServer) Unregister(context.Context, *NetworkServiceEndpoint) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unregister not implemented")
}

// UnsafeNetworkServiceEndpointRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkServiceEndpointRegistryServer will
// result in compilation errors.
type UnsafeNetworkServiceEndpointRegistryServer interface {
	mustEmbedUnimplementedNetworkServiceEndpointRegistryServer()
}

func RegisterNetworkServiceEndpointRegistryServer(s grpc.ServiceRegistrar, srv NetworkServiceEndpointRegistryServer) {
	s.RegisterService(&NetworkServiceEndpointRegistry_ServiceDesc, srv)
}

func _NetworkServiceEndpointRegistry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkServiceEndpoint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceEndpointRegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.registry.NetworkServiceEndpointRegistry/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceEndpointRegistryServer).Register(ctx, req.(*NetworkServiceEndpoint))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkServiceEndpointRegistry_Find_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NetworkServiceEndpointQuery)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NetworkServiceEndpointRegistryServer).Find(m, &networkServiceEndpointRegistryFindServer{stream})
}

type NetworkServiceEndpointRegistry_FindServer interface {
	Send(*NetworkServiceEndpointResponse) error
	grpc.ServerStream
}

type networkServiceEndpointRegistryFindServer struct {
	grpc.ServerStream
}

func (x *networkServiceEndpointRegistryFindServer) Send(m *NetworkServiceEndpointResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NetworkServiceEndpointRegistry_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkServiceEndpoint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceEndpointRegistryServer).Unregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.registry.NetworkServiceEndpointRegistry/Unregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceEndpointRegistryServer).Unregister(ctx, req.(*NetworkServiceEndpoint))
	}
	return interceptor(ctx, in, info, handler)
}

// NetworkServiceEndpointRegistry_ServiceDesc is the grpc.ServiceDesc for NetworkServiceEndpointRegistry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkServiceEndpointRegistry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.registry.NetworkServiceEndpointRegistry",
	HandlerType: (*NetworkServiceEndpointRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _NetworkServiceEndpointRegistry_Register_Handler,
		},
		{
			MethodName: "Unregister",
			Handler:    _NetworkServiceEndpointRegistry_Unregister_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Find",
			Handler:       _NetworkServiceEndpointRegistry_Find_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/core/v1/registry/registry.proto",
}

// NetworkServiceRegistryClient is the client API for NetworkServiceRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetworkServiceRegistryClient interface {
	Register(ctx context.Context, in *NetworkService, opts ...grpc.CallOption) (*NetworkService, error)
	Find(ctx context.Context, in *NetworkServiceQuery, opts ...grpc.CallOption) (NetworkServiceRegistry_FindClient, error)
	Unregister(ctx context.Context, in *NetworkService, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type networkServiceRegistryClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkServiceRegistryClient(cc grpc.ClientConnInterface) NetworkServiceRegistryClient {
	return &networkServiceRegistryClient{cc}
}

func (c *networkServiceRegistryClient) Register(ctx context.Context, in *NetworkService, opts ...grpc.CallOption) (*NetworkService, error) {
	out := new(NetworkService)
	err := c.cc.Invoke(ctx, "/v1.registry.NetworkServiceRegistry/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *networkServiceRegistryClient) Find(ctx context.Context, in *NetworkServiceQuery, opts ...grpc.CallOption) (NetworkServiceRegistry_FindClient, error) {
	stream, err := c.cc.NewStream(ctx, &NetworkServiceRegistry_ServiceDesc.Streams[0], "/v1.registry.NetworkServiceRegistry/Find", opts...)
	if err != nil {
		return nil, err
	}
	x := &networkServiceRegistryFindClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NetworkServiceRegistry_FindClient interface {
	Recv() (*NetworkServiceResponse, error)
	grpc.ClientStream
}

type networkServiceRegistryFindClient struct {
	grpc.ClientStream
}

func (x *networkServiceRegistryFindClient) Recv() (*NetworkServiceResponse, error) {
	m := new(NetworkServiceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *networkServiceRegistryClient) Unregister(ctx context.Context, in *NetworkService, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/v1.registry.NetworkServiceRegistry/Unregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkServiceRegistryServer is the server API for NetworkServiceRegistry service.
// All implementations should embed UnimplementedNetworkServiceRegistryServer
// for forward compatibility
type NetworkServiceRegistryServer interface {
	Register(context.Context, *NetworkService) (*NetworkService, error)
	Find(*NetworkServiceQuery, NetworkServiceRegistry_FindServer) error
	Unregister(context.Context, *NetworkService) (*emptypb.Empty, error)
}

// UnimplementedNetworkServiceRegistryServer should be embedded to have forward compatible implementations.
type UnimplementedNetworkServiceRegistryServer struct {
}

func (UnimplementedNetworkServiceRegistryServer) Register(context.Context, *NetworkService) (*NetworkService, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedNetworkServiceRegistryServer) Find(*NetworkServiceQuery, NetworkServiceRegistry_FindServer) error {
	return status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedNetworkServiceRegistryServer) Unregister(context.Context, *NetworkService) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unregister not implemented")
}

// UnsafeNetworkServiceRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetworkServiceRegistryServer will
// result in compilation errors.
type UnsafeNetworkServiceRegistryServer interface {
	mustEmbedUnimplementedNetworkServiceRegistryServer()
}

func RegisterNetworkServiceRegistryServer(s grpc.ServiceRegistrar, srv NetworkServiceRegistryServer) {
	s.RegisterService(&NetworkServiceRegistry_ServiceDesc, srv)
}

func _NetworkServiceRegistry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceRegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.registry.NetworkServiceRegistry/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceRegistryServer).Register(ctx, req.(*NetworkService))
	}
	return interceptor(ctx, in, info, handler)
}

func _NetworkServiceRegistry_Find_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NetworkServiceQuery)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NetworkServiceRegistryServer).Find(m, &networkServiceRegistryFindServer{stream})
}

type NetworkServiceRegistry_FindServer interface {
	Send(*NetworkServiceResponse) error
	grpc.ServerStream
}

type networkServiceRegistryFindServer struct {
	grpc.ServerStream
}

func (x *networkServiceRegistryFindServer) Send(m *NetworkServiceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _NetworkServiceRegistry_Unregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NetworkService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkServiceRegistryServer).Unregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.registry.NetworkServiceRegistry/Unregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkServiceRegistryServer).Unregister(ctx, req.(*NetworkService))
	}
	return interceptor(ctx, in, info, handler)
}

// NetworkServiceRegistry_ServiceDesc is the grpc.ServiceDesc for NetworkServiceRegistry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetworkServiceRegistry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.registry.NetworkServiceRegistry",
	HandlerType: (*NetworkServiceRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _NetworkServiceRegistry_Register_Handler,
		},
		{
			MethodName: "Unregister",
			Handler:    _NetworkServiceRegistry_Unregister_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Find",
			Handler:       _NetworkServiceRegistry_Find_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/core/v1/registry/registry.proto",
}
