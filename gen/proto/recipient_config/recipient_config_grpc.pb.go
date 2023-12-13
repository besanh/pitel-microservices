// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: proto/recipient_config/recipient_config.proto

package pb

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
	RecipientConfig_PostRecipientConfig_FullMethodName       = "/proto.recipient_config.RecipientConfig/PostRecipientConfig"
	RecipientConfig_GetRecipientConfigs_FullMethodName       = "/proto.recipient_config.RecipientConfig/GetRecipientConfigs"
	RecipientConfig_GetRecipientConfigById_FullMethodName    = "/proto.recipient_config.RecipientConfig/GetRecipientConfigById"
	RecipientConfig_PutRecipientConfigById_FullMethodName    = "/proto.recipient_config.RecipientConfig/PutRecipientConfigById"
	RecipientConfig_DeleteRecipientConfigById_FullMethodName = "/proto.recipient_config.RecipientConfig/DeleteRecipientConfigById"
)

// RecipientConfigClient is the client API for RecipientConfig service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecipientConfigClient interface {
	PostRecipientConfig(ctx context.Context, in *RecipientConfigBodyRequest, opts ...grpc.CallOption) (*RecipientConfigResponse, error)
	GetRecipientConfigs(ctx context.Context, in *RecipientConfigRequest, opts ...grpc.CallOption) (*RecipientConfigResponse, error)
	GetRecipientConfigById(ctx context.Context, in *RecipientConfigByIdRequest, opts ...grpc.CallOption) (*RecipientConfigByIdResponse, error)
	PutRecipientConfigById(ctx context.Context, in *RecipientConfigPutRequest, opts ...grpc.CallOption) (*RecipientConfigByIdResponse, error)
	DeleteRecipientConfigById(ctx context.Context, in *RecipientConfigByIdRequest, opts ...grpc.CallOption) (*RecipientConfigByIdResponse, error)
}

type recipientConfigClient struct {
	cc grpc.ClientConnInterface
}

func NewRecipientConfigClient(cc grpc.ClientConnInterface) RecipientConfigClient {
	return &recipientConfigClient{cc}
}

func (c *recipientConfigClient) PostRecipientConfig(ctx context.Context, in *RecipientConfigBodyRequest, opts ...grpc.CallOption) (*RecipientConfigResponse, error) {
	out := new(RecipientConfigResponse)
	err := c.cc.Invoke(ctx, RecipientConfig_PostRecipientConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipientConfigClient) GetRecipientConfigs(ctx context.Context, in *RecipientConfigRequest, opts ...grpc.CallOption) (*RecipientConfigResponse, error) {
	out := new(RecipientConfigResponse)
	err := c.cc.Invoke(ctx, RecipientConfig_GetRecipientConfigs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipientConfigClient) GetRecipientConfigById(ctx context.Context, in *RecipientConfigByIdRequest, opts ...grpc.CallOption) (*RecipientConfigByIdResponse, error) {
	out := new(RecipientConfigByIdResponse)
	err := c.cc.Invoke(ctx, RecipientConfig_GetRecipientConfigById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipientConfigClient) PutRecipientConfigById(ctx context.Context, in *RecipientConfigPutRequest, opts ...grpc.CallOption) (*RecipientConfigByIdResponse, error) {
	out := new(RecipientConfigByIdResponse)
	err := c.cc.Invoke(ctx, RecipientConfig_PutRecipientConfigById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipientConfigClient) DeleteRecipientConfigById(ctx context.Context, in *RecipientConfigByIdRequest, opts ...grpc.CallOption) (*RecipientConfigByIdResponse, error) {
	out := new(RecipientConfigByIdResponse)
	err := c.cc.Invoke(ctx, RecipientConfig_DeleteRecipientConfigById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecipientConfigServer is the server API for RecipientConfig service.
// All implementations should embed UnimplementedRecipientConfigServer
// for forward compatibility
type RecipientConfigServer interface {
	PostRecipientConfig(context.Context, *RecipientConfigBodyRequest) (*RecipientConfigResponse, error)
	GetRecipientConfigs(context.Context, *RecipientConfigRequest) (*RecipientConfigResponse, error)
	GetRecipientConfigById(context.Context, *RecipientConfigByIdRequest) (*RecipientConfigByIdResponse, error)
	PutRecipientConfigById(context.Context, *RecipientConfigPutRequest) (*RecipientConfigByIdResponse, error)
	DeleteRecipientConfigById(context.Context, *RecipientConfigByIdRequest) (*RecipientConfigByIdResponse, error)
}

// UnimplementedRecipientConfigServer should be embedded to have forward compatible implementations.
type UnimplementedRecipientConfigServer struct {
}

func (UnimplementedRecipientConfigServer) PostRecipientConfig(context.Context, *RecipientConfigBodyRequest) (*RecipientConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostRecipientConfig not implemented")
}
func (UnimplementedRecipientConfigServer) GetRecipientConfigs(context.Context, *RecipientConfigRequest) (*RecipientConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecipientConfigs not implemented")
}
func (UnimplementedRecipientConfigServer) GetRecipientConfigById(context.Context, *RecipientConfigByIdRequest) (*RecipientConfigByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecipientConfigById not implemented")
}
func (UnimplementedRecipientConfigServer) PutRecipientConfigById(context.Context, *RecipientConfigPutRequest) (*RecipientConfigByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutRecipientConfigById not implemented")
}
func (UnimplementedRecipientConfigServer) DeleteRecipientConfigById(context.Context, *RecipientConfigByIdRequest) (*RecipientConfigByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecipientConfigById not implemented")
}

// UnsafeRecipientConfigServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecipientConfigServer will
// result in compilation errors.
type UnsafeRecipientConfigServer interface {
	mustEmbedUnimplementedRecipientConfigServer()
}

func RegisterRecipientConfigServer(s grpc.ServiceRegistrar, srv RecipientConfigServer) {
	s.RegisterService(&RecipientConfig_ServiceDesc, srv)
}

func _RecipientConfig_PostRecipientConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipientConfigBodyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipientConfigServer).PostRecipientConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipientConfig_PostRecipientConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipientConfigServer).PostRecipientConfig(ctx, req.(*RecipientConfigBodyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipientConfig_GetRecipientConfigs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipientConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipientConfigServer).GetRecipientConfigs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipientConfig_GetRecipientConfigs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipientConfigServer).GetRecipientConfigs(ctx, req.(*RecipientConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipientConfig_GetRecipientConfigById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipientConfigByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipientConfigServer).GetRecipientConfigById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipientConfig_GetRecipientConfigById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipientConfigServer).GetRecipientConfigById(ctx, req.(*RecipientConfigByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipientConfig_PutRecipientConfigById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipientConfigPutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipientConfigServer).PutRecipientConfigById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipientConfig_PutRecipientConfigById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipientConfigServer).PutRecipientConfigById(ctx, req.(*RecipientConfigPutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipientConfig_DeleteRecipientConfigById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipientConfigByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipientConfigServer).DeleteRecipientConfigById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipientConfig_DeleteRecipientConfigById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipientConfigServer).DeleteRecipientConfigById(ctx, req.(*RecipientConfigByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecipientConfig_ServiceDesc is the grpc.ServiceDesc for RecipientConfig service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecipientConfig_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.recipient_config.RecipientConfig",
	HandlerType: (*RecipientConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostRecipientConfig",
			Handler:    _RecipientConfig_PostRecipientConfig_Handler,
		},
		{
			MethodName: "GetRecipientConfigs",
			Handler:    _RecipientConfig_GetRecipientConfigs_Handler,
		},
		{
			MethodName: "GetRecipientConfigById",
			Handler:    _RecipientConfig_GetRecipientConfigById_Handler,
		},
		{
			MethodName: "PutRecipientConfigById",
			Handler:    _RecipientConfig_PutRecipientConfigById_Handler,
		},
		{
			MethodName: "DeleteRecipientConfigById",
			Handler:    _RecipientConfig_DeleteRecipientConfigById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/recipient_config/recipient_config.proto",
}
