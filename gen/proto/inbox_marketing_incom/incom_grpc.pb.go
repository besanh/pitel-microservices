// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: proto/inbox_marketing_incom/incom.proto

package inbox_marketing

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
	IncomService_HandleWebhook_FullMethodName = "/proto.inbox_marketing_incom.IncomService/HandleWebhook"
)

// IncomServiceClient is the client API for IncomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IncomServiceClient interface {
	HandleWebhook(ctx context.Context, in *IncomBodyRequest, opts ...grpc.CallOption) (*IncomResponse, error)
}

type incomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIncomServiceClient(cc grpc.ClientConnInterface) IncomServiceClient {
	return &incomServiceClient{cc}
}

func (c *incomServiceClient) HandleWebhook(ctx context.Context, in *IncomBodyRequest, opts ...grpc.CallOption) (*IncomResponse, error) {
	out := new(IncomResponse)
	err := c.cc.Invoke(ctx, IncomService_HandleWebhook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IncomServiceServer is the server API for IncomService service.
// All implementations should embed UnimplementedIncomServiceServer
// for forward compatibility
type IncomServiceServer interface {
	HandleWebhook(context.Context, *IncomBodyRequest) (*IncomResponse, error)
}

// UnimplementedIncomServiceServer should be embedded to have forward compatible implementations.
type UnimplementedIncomServiceServer struct {
}

func (UnimplementedIncomServiceServer) HandleWebhook(context.Context, *IncomBodyRequest) (*IncomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleWebhook not implemented")
}

// UnsafeIncomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IncomServiceServer will
// result in compilation errors.
type UnsafeIncomServiceServer interface {
	mustEmbedUnimplementedIncomServiceServer()
}

func RegisterIncomServiceServer(s grpc.ServiceRegistrar, srv IncomServiceServer) {
	s.RegisterService(&IncomService_ServiceDesc, srv)
}

func _IncomService_HandleWebhook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncomBodyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IncomServiceServer).HandleWebhook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IncomService_HandleWebhook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IncomServiceServer).HandleWebhook(ctx, req.(*IncomBodyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IncomService_ServiceDesc is the grpc.ServiceDesc for IncomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IncomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.inbox_marketing_incom.IncomService",
	HandlerType: (*IncomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleWebhook",
			Handler:    _IncomService_HandleWebhook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/inbox_marketing_incom/incom.proto",
}
