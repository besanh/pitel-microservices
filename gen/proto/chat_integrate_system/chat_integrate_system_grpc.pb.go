// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: proto/chat_integrate_system/chat_integrate_system.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ChatIntegrateSystem_GetChatIntegrateSystems_FullMethodName       = "/proto.chatIntegrateSystem.ChatIntegrateSystem/GetChatIntegrateSystems"
	ChatIntegrateSystem_PostChatIntegrateSystem_FullMethodName       = "/proto.chatIntegrateSystem.ChatIntegrateSystem/PostChatIntegrateSystem"
	ChatIntegrateSystem_GetChatIntegrateSystemById_FullMethodName    = "/proto.chatIntegrateSystem.ChatIntegrateSystem/GetChatIntegrateSystemById"
	ChatIntegrateSystem_UpdateChatIntegrateSystemById_FullMethodName = "/proto.chatIntegrateSystem.ChatIntegrateSystem/UpdateChatIntegrateSystemById"
	ChatIntegrateSystem_DeleteChatIntegrateSystemById_FullMethodName = "/proto.chatIntegrateSystem.ChatIntegrateSystem/DeleteChatIntegrateSystemById"
)

// ChatIntegrateSystemClient is the client API for ChatIntegrateSystem service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatIntegrateSystemClient interface {
	GetChatIntegrateSystems(ctx context.Context, in *GetChatIntegrateSystemRequest, opts ...grpc.CallOption) (*GetChatIntegrateSystemResponse, error)
	PostChatIntegrateSystem(ctx context.Context, in *PostChatIntegrateSystemRequest, opts ...grpc.CallOption) (*PostChatIntegrateSystemResponse, error)
	GetChatIntegrateSystemById(ctx context.Context, in *GetChatIntegrateSystemByIdRequest, opts ...grpc.CallOption) (*GetChatIntegrateSystemByIdResponse, error)
	UpdateChatIntegrateSystemById(ctx context.Context, in *UpdateChatIntegrateSystemByIdRequest, opts ...grpc.CallOption) (*UpdateChatIntegrateSystemByIdResponse, error)
	DeleteChatIntegrateSystemById(ctx context.Context, in *DeleteChatIntegrateSystemByIdRequest, opts ...grpc.CallOption) (*DeleteChatIntegrateSystemByIdResponse, error)
}

type chatIntegrateSystemClient struct {
	cc grpc.ClientConnInterface
}

func NewChatIntegrateSystemClient(cc grpc.ClientConnInterface) ChatIntegrateSystemClient {
	return &chatIntegrateSystemClient{cc}
}

func (c *chatIntegrateSystemClient) GetChatIntegrateSystems(ctx context.Context, in *GetChatIntegrateSystemRequest, opts ...grpc.CallOption) (*GetChatIntegrateSystemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChatIntegrateSystemResponse)
	err := c.cc.Invoke(ctx, ChatIntegrateSystem_GetChatIntegrateSystems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatIntegrateSystemClient) PostChatIntegrateSystem(ctx context.Context, in *PostChatIntegrateSystemRequest, opts ...grpc.CallOption) (*PostChatIntegrateSystemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostChatIntegrateSystemResponse)
	err := c.cc.Invoke(ctx, ChatIntegrateSystem_PostChatIntegrateSystem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatIntegrateSystemClient) GetChatIntegrateSystemById(ctx context.Context, in *GetChatIntegrateSystemByIdRequest, opts ...grpc.CallOption) (*GetChatIntegrateSystemByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChatIntegrateSystemByIdResponse)
	err := c.cc.Invoke(ctx, ChatIntegrateSystem_GetChatIntegrateSystemById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatIntegrateSystemClient) UpdateChatIntegrateSystemById(ctx context.Context, in *UpdateChatIntegrateSystemByIdRequest, opts ...grpc.CallOption) (*UpdateChatIntegrateSystemByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateChatIntegrateSystemByIdResponse)
	err := c.cc.Invoke(ctx, ChatIntegrateSystem_UpdateChatIntegrateSystemById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatIntegrateSystemClient) DeleteChatIntegrateSystemById(ctx context.Context, in *DeleteChatIntegrateSystemByIdRequest, opts ...grpc.CallOption) (*DeleteChatIntegrateSystemByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteChatIntegrateSystemByIdResponse)
	err := c.cc.Invoke(ctx, ChatIntegrateSystem_DeleteChatIntegrateSystemById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatIntegrateSystemServer is the server API for ChatIntegrateSystem service.
// All implementations should embed UnimplementedChatIntegrateSystemServer
// for forward compatibility
type ChatIntegrateSystemServer interface {
	GetChatIntegrateSystems(context.Context, *GetChatIntegrateSystemRequest) (*GetChatIntegrateSystemResponse, error)
	PostChatIntegrateSystem(context.Context, *PostChatIntegrateSystemRequest) (*PostChatIntegrateSystemResponse, error)
	GetChatIntegrateSystemById(context.Context, *GetChatIntegrateSystemByIdRequest) (*GetChatIntegrateSystemByIdResponse, error)
	UpdateChatIntegrateSystemById(context.Context, *UpdateChatIntegrateSystemByIdRequest) (*UpdateChatIntegrateSystemByIdResponse, error)
	DeleteChatIntegrateSystemById(context.Context, *DeleteChatIntegrateSystemByIdRequest) (*DeleteChatIntegrateSystemByIdResponse, error)
}

// UnimplementedChatIntegrateSystemServer should be embedded to have forward compatible implementations.
type UnimplementedChatIntegrateSystemServer struct {
}

func (UnimplementedChatIntegrateSystemServer) GetChatIntegrateSystems(context.Context, *GetChatIntegrateSystemRequest) (*GetChatIntegrateSystemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatIntegrateSystems not implemented")
}
func (UnimplementedChatIntegrateSystemServer) PostChatIntegrateSystem(context.Context, *PostChatIntegrateSystemRequest) (*PostChatIntegrateSystemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostChatIntegrateSystem not implemented")
}
func (UnimplementedChatIntegrateSystemServer) GetChatIntegrateSystemById(context.Context, *GetChatIntegrateSystemByIdRequest) (*GetChatIntegrateSystemByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatIntegrateSystemById not implemented")
}
func (UnimplementedChatIntegrateSystemServer) UpdateChatIntegrateSystemById(context.Context, *UpdateChatIntegrateSystemByIdRequest) (*UpdateChatIntegrateSystemByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChatIntegrateSystemById not implemented")
}
func (UnimplementedChatIntegrateSystemServer) DeleteChatIntegrateSystemById(context.Context, *DeleteChatIntegrateSystemByIdRequest) (*DeleteChatIntegrateSystemByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChatIntegrateSystemById not implemented")
}

// UnsafeChatIntegrateSystemServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatIntegrateSystemServer will
// result in compilation errors.
type UnsafeChatIntegrateSystemServer interface {
	mustEmbedUnimplementedChatIntegrateSystemServer()
}

func RegisterChatIntegrateSystemServer(s grpc.ServiceRegistrar, srv ChatIntegrateSystemServer) {
	s.RegisterService(&ChatIntegrateSystem_ServiceDesc, srv)
}

func _ChatIntegrateSystem_GetChatIntegrateSystems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatIntegrateSystemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatIntegrateSystemServer).GetChatIntegrateSystems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatIntegrateSystem_GetChatIntegrateSystems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatIntegrateSystemServer).GetChatIntegrateSystems(ctx, req.(*GetChatIntegrateSystemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatIntegrateSystem_PostChatIntegrateSystem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostChatIntegrateSystemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatIntegrateSystemServer).PostChatIntegrateSystem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatIntegrateSystem_PostChatIntegrateSystem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatIntegrateSystemServer).PostChatIntegrateSystem(ctx, req.(*PostChatIntegrateSystemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatIntegrateSystem_GetChatIntegrateSystemById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatIntegrateSystemByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatIntegrateSystemServer).GetChatIntegrateSystemById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatIntegrateSystem_GetChatIntegrateSystemById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatIntegrateSystemServer).GetChatIntegrateSystemById(ctx, req.(*GetChatIntegrateSystemByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatIntegrateSystem_UpdateChatIntegrateSystemById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateChatIntegrateSystemByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatIntegrateSystemServer).UpdateChatIntegrateSystemById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatIntegrateSystem_UpdateChatIntegrateSystemById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatIntegrateSystemServer).UpdateChatIntegrateSystemById(ctx, req.(*UpdateChatIntegrateSystemByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatIntegrateSystem_DeleteChatIntegrateSystemById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChatIntegrateSystemByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatIntegrateSystemServer).DeleteChatIntegrateSystemById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatIntegrateSystem_DeleteChatIntegrateSystemById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatIntegrateSystemServer).DeleteChatIntegrateSystemById(ctx, req.(*DeleteChatIntegrateSystemByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatIntegrateSystem_ServiceDesc is the grpc.ServiceDesc for ChatIntegrateSystem service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatIntegrateSystem_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.chatIntegrateSystem.ChatIntegrateSystem",
	HandlerType: (*ChatIntegrateSystemServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChatIntegrateSystems",
			Handler:    _ChatIntegrateSystem_GetChatIntegrateSystems_Handler,
		},
		{
			MethodName: "PostChatIntegrateSystem",
			Handler:    _ChatIntegrateSystem_PostChatIntegrateSystem_Handler,
		},
		{
			MethodName: "GetChatIntegrateSystemById",
			Handler:    _ChatIntegrateSystem_GetChatIntegrateSystemById_Handler,
		},
		{
			MethodName: "UpdateChatIntegrateSystemById",
			Handler:    _ChatIntegrateSystem_UpdateChatIntegrateSystemById_Handler,
		},
		{
			MethodName: "DeleteChatIntegrateSystemById",
			Handler:    _ChatIntegrateSystem_DeleteChatIntegrateSystemById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/chat_integrate_system/chat_integrate_system.proto",
}
