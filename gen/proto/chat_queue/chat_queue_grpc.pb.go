// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: proto/chat_queue/chat_queue.proto

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
	ChatQueueService_GetChatQueues_FullMethodName             = "/proto.chatQueue.ChatQueueService/GetChatQueues"
	ChatQueueService_GetChatQueueById_FullMethodName          = "/proto.chatQueue.ChatQueueService/GetChatQueueById"
	ChatQueueService_InsertChatQueue_FullMethodName           = "/proto.chatQueue.ChatQueueService/InsertChatQueue"
	ChatQueueService_UpdateChatQueueById_FullMethodName       = "/proto.chatQueue.ChatQueueService/UpdateChatQueueById"
	ChatQueueService_UpdateChatQueueStatusById_FullMethodName = "/proto.chatQueue.ChatQueueService/UpdateChatQueueStatusById"
	ChatQueueService_DeleteChatQueueById_FullMethodName       = "/proto.chatQueue.ChatQueueService/DeleteChatQueueById"
)

// ChatQueueServiceClient is the client API for ChatQueueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatQueueServiceClient interface {
	GetChatQueues(ctx context.Context, in *GetChatQueuesRequest, opts ...grpc.CallOption) (*GetChatQueuesResponse, error)
	GetChatQueueById(ctx context.Context, in *GetChatQueueByIdRequest, opts ...grpc.CallOption) (*GetChatQueueByIdResponse, error)
	InsertChatQueue(ctx context.Context, in *PostChatQueueRequest, opts ...grpc.CallOption) (*PostChatQueueResponse, error)
	UpdateChatQueueById(ctx context.Context, in *PutChatQueueRequest, opts ...grpc.CallOption) (*PutChatQueueResponse, error)
	UpdateChatQueueStatusById(ctx context.Context, in *PutChatQueueStatusRequest, opts ...grpc.CallOption) (*PutChatQueueStatusResponse, error)
	DeleteChatQueueById(ctx context.Context, in *DeleteChatQueueRequest, opts ...grpc.CallOption) (*DeleteChatQueueResponse, error)
}

type chatQueueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatQueueServiceClient(cc grpc.ClientConnInterface) ChatQueueServiceClient {
	return &chatQueueServiceClient{cc}
}

func (c *chatQueueServiceClient) GetChatQueues(ctx context.Context, in *GetChatQueuesRequest, opts ...grpc.CallOption) (*GetChatQueuesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChatQueuesResponse)
	err := c.cc.Invoke(ctx, ChatQueueService_GetChatQueues_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatQueueServiceClient) GetChatQueueById(ctx context.Context, in *GetChatQueueByIdRequest, opts ...grpc.CallOption) (*GetChatQueueByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetChatQueueByIdResponse)
	err := c.cc.Invoke(ctx, ChatQueueService_GetChatQueueById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatQueueServiceClient) InsertChatQueue(ctx context.Context, in *PostChatQueueRequest, opts ...grpc.CallOption) (*PostChatQueueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PostChatQueueResponse)
	err := c.cc.Invoke(ctx, ChatQueueService_InsertChatQueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatQueueServiceClient) UpdateChatQueueById(ctx context.Context, in *PutChatQueueRequest, opts ...grpc.CallOption) (*PutChatQueueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PutChatQueueResponse)
	err := c.cc.Invoke(ctx, ChatQueueService_UpdateChatQueueById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatQueueServiceClient) UpdateChatQueueStatusById(ctx context.Context, in *PutChatQueueStatusRequest, opts ...grpc.CallOption) (*PutChatQueueStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PutChatQueueStatusResponse)
	err := c.cc.Invoke(ctx, ChatQueueService_UpdateChatQueueStatusById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatQueueServiceClient) DeleteChatQueueById(ctx context.Context, in *DeleteChatQueueRequest, opts ...grpc.CallOption) (*DeleteChatQueueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteChatQueueResponse)
	err := c.cc.Invoke(ctx, ChatQueueService_DeleteChatQueueById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatQueueServiceServer is the server API for ChatQueueService service.
// All implementations should embed UnimplementedChatQueueServiceServer
// for forward compatibility
type ChatQueueServiceServer interface {
	GetChatQueues(context.Context, *GetChatQueuesRequest) (*GetChatQueuesResponse, error)
	GetChatQueueById(context.Context, *GetChatQueueByIdRequest) (*GetChatQueueByIdResponse, error)
	InsertChatQueue(context.Context, *PostChatQueueRequest) (*PostChatQueueResponse, error)
	UpdateChatQueueById(context.Context, *PutChatQueueRequest) (*PutChatQueueResponse, error)
	UpdateChatQueueStatusById(context.Context, *PutChatQueueStatusRequest) (*PutChatQueueStatusResponse, error)
	DeleteChatQueueById(context.Context, *DeleteChatQueueRequest) (*DeleteChatQueueResponse, error)
}

// UnimplementedChatQueueServiceServer should be embedded to have forward compatible implementations.
type UnimplementedChatQueueServiceServer struct {
}

func (UnimplementedChatQueueServiceServer) GetChatQueues(context.Context, *GetChatQueuesRequest) (*GetChatQueuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatQueues not implemented")
}
func (UnimplementedChatQueueServiceServer) GetChatQueueById(context.Context, *GetChatQueueByIdRequest) (*GetChatQueueByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChatQueueById not implemented")
}
func (UnimplementedChatQueueServiceServer) InsertChatQueue(context.Context, *PostChatQueueRequest) (*PostChatQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertChatQueue not implemented")
}
func (UnimplementedChatQueueServiceServer) UpdateChatQueueById(context.Context, *PutChatQueueRequest) (*PutChatQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChatQueueById not implemented")
}
func (UnimplementedChatQueueServiceServer) UpdateChatQueueStatusById(context.Context, *PutChatQueueStatusRequest) (*PutChatQueueStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChatQueueStatusById not implemented")
}
func (UnimplementedChatQueueServiceServer) DeleteChatQueueById(context.Context, *DeleteChatQueueRequest) (*DeleteChatQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChatQueueById not implemented")
}

// UnsafeChatQueueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatQueueServiceServer will
// result in compilation errors.
type UnsafeChatQueueServiceServer interface {
	mustEmbedUnimplementedChatQueueServiceServer()
}

func RegisterChatQueueServiceServer(s grpc.ServiceRegistrar, srv ChatQueueServiceServer) {
	s.RegisterService(&ChatQueueService_ServiceDesc, srv)
}

func _ChatQueueService_GetChatQueues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatQueuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatQueueServiceServer).GetChatQueues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatQueueService_GetChatQueues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatQueueServiceServer).GetChatQueues(ctx, req.(*GetChatQueuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatQueueService_GetChatQueueById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatQueueByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatQueueServiceServer).GetChatQueueById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatQueueService_GetChatQueueById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatQueueServiceServer).GetChatQueueById(ctx, req.(*GetChatQueueByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatQueueService_InsertChatQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostChatQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatQueueServiceServer).InsertChatQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatQueueService_InsertChatQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatQueueServiceServer).InsertChatQueue(ctx, req.(*PostChatQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatQueueService_UpdateChatQueueById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutChatQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatQueueServiceServer).UpdateChatQueueById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatQueueService_UpdateChatQueueById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatQueueServiceServer).UpdateChatQueueById(ctx, req.(*PutChatQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatQueueService_UpdateChatQueueStatusById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutChatQueueStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatQueueServiceServer).UpdateChatQueueStatusById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatQueueService_UpdateChatQueueStatusById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatQueueServiceServer).UpdateChatQueueStatusById(ctx, req.(*PutChatQueueStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatQueueService_DeleteChatQueueById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChatQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatQueueServiceServer).DeleteChatQueueById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatQueueService_DeleteChatQueueById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatQueueServiceServer).DeleteChatQueueById(ctx, req.(*DeleteChatQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatQueueService_ServiceDesc is the grpc.ServiceDesc for ChatQueueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatQueueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.chatQueue.ChatQueueService",
	HandlerType: (*ChatQueueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChatQueues",
			Handler:    _ChatQueueService_GetChatQueues_Handler,
		},
		{
			MethodName: "GetChatQueueById",
			Handler:    _ChatQueueService_GetChatQueueById_Handler,
		},
		{
			MethodName: "InsertChatQueue",
			Handler:    _ChatQueueService_InsertChatQueue_Handler,
		},
		{
			MethodName: "UpdateChatQueueById",
			Handler:    _ChatQueueService_UpdateChatQueueById_Handler,
		},
		{
			MethodName: "UpdateChatQueueStatusById",
			Handler:    _ChatQueueService_UpdateChatQueueStatusById_Handler,
		},
		{
			MethodName: "DeleteChatQueueById",
			Handler:    _ChatQueueService_DeleteChatQueueById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/chat_queue/chat_queue.proto",
}
