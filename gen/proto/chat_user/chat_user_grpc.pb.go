// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/chat_user/chat_user.proto

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

// ChatUserServiceClient is the client API for ChatUserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatUserServiceClient interface {
	PostChatUser(ctx context.Context, in *PostChatUserRequest, opts ...grpc.CallOption) (*PostChatUserResponse, error)
	UpdateChatUserStatusById(ctx context.Context, in *PutChatUserStatusRequest, opts ...grpc.CallOption) (*PutChatUserResponse, error)
	GetChatUserStatusById(ctx context.Context, in *GetChatUserStatusRequest, opts ...grpc.CallOption) (*GetChatUserResponse, error)
}

type chatUserServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatUserServiceClient(cc grpc.ClientConnInterface) ChatUserServiceClient {
	return &chatUserServiceClient{cc}
}

func (c *chatUserServiceClient) PostChatUser(ctx context.Context, in *PostChatUserRequest, opts ...grpc.CallOption) (*PostChatUserResponse, error) {
	out := new(PostChatUserResponse)
	err := c.cc.Invoke(ctx, "/proto.chatUser.ChatUserService/PostChatUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatUserServiceClient) UpdateChatUserStatusById(ctx context.Context, in *PutChatUserStatusRequest, opts ...grpc.CallOption) (*PutChatUserResponse, error) {
	out := new(PutChatUserResponse)
	err := c.cc.Invoke(ctx, "/proto.chatUser.ChatUserService/UpdateChatUserStatusById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatUserServiceClient) GetChatUserStatusById(ctx context.Context, in *GetChatUserStatusRequest, opts ...grpc.CallOption) (*GetChatUserResponse, error) {
	out := new(GetChatUserResponse)
	err := c.cc.Invoke(ctx, "/proto.chatUser.ChatUserService/GetChatUserStatusById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatUserServiceServer is the server API for ChatUserService service.
// All implementations should embed UnimplementedChatUserServiceServer
// for forward compatibility
type ChatUserServiceServer interface {
	PostChatUser(context.Context, *PostChatUserRequest) (*PostChatUserResponse, error)
	UpdateChatUserStatusById(context.Context, *PutChatUserStatusRequest) (*PutChatUserResponse, error)
	GetChatUserStatusById(context.Context, *GetChatUserStatusRequest) (*GetChatUserResponse, error)
}

// UnimplementedChatUserServiceServer should be embedded to have forward compatible implementations.
type UnimplementedChatUserServiceServer struct {
}

func (UnimplementedChatUserServiceServer) PostChatUser(context.Context, *PostChatUserRequest) (*PostChatUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostChatUser not implemented")
}
func (UnimplementedChatUserServiceServer) UpdateChatUserStatusById(context.Context, *PutChatUserStatusRequest) (*PutChatUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChatUserStatusById not implemented")
}
func (UnimplementedChatUserServiceServer) GetChatUserStatusById(context.Context, *GetChatUserStatusRequest) (*GetChatUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatUserStatusById not implemented")
}

// UnsafeChatUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatUserServiceServer will
// result in compilation errors.
type UnsafeChatUserServiceServer interface {
	mustEmbedUnimplementedChatUserServiceServer()
}

func RegisterChatUserServiceServer(s grpc.ServiceRegistrar, srv ChatUserServiceServer) {
	s.RegisterService(&ChatUserService_ServiceDesc, srv)
}

func _ChatUserService_PostChatUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostChatUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatUserServiceServer).PostChatUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.chatUser.ChatUserService/PostChatUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatUserServiceServer).PostChatUser(ctx, req.(*PostChatUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatUserService_UpdateChatUserStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutChatUserStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatUserServiceServer).UpdateChatUserStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.chatUser.ChatUserService/UpdateChatUserStatusById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatUserServiceServer).UpdateChatUserStatusById(ctx, req.(*PutChatUserStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatUserService_GetChatUserStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatUserStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatUserServiceServer).GetChatUserStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.chatUser.ChatUserService/GetChatUserStatusById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatUserServiceServer).GetChatUserStatusById(ctx, req.(*GetChatUserStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatUserService_ServiceDesc is the grpc.ServiceDesc for ChatUserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatUserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.chatUser.ChatUserService",
	HandlerType: (*ChatUserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostChatUser",
			Handler:    _ChatUserService_PostChatUser_Handler,
		},
		{
			MethodName: "UpdateChatUserStatusById",
			Handler:    _ChatUserService_UpdateChatUserStatusById_Handler,
		},
		{
			MethodName: "GetChatUserStatusById",
			Handler:    _ChatUserService_GetChatUserStatusById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/chat_user/chat_user.proto",
}
