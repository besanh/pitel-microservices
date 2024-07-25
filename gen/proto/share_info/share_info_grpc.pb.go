// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/share_info/share_info.proto

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

// ShareInfoServiceClient is the client API for ShareInfoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShareInfoServiceClient interface {
	// Submit send to ott
	PostShareInfo(ctx context.Context, in *PostShareInfoRequest, opts ...grpc.CallOption) (*PostShareInfoResponse, error)
	GetShareInfos(ctx context.Context, in *GetShareInfoRequest, opts ...grpc.CallOption) (*GetShareInfoResponse, error)
	GetShareInfoById(ctx context.Context, in *GetShareInfoByIdRequest, opts ...grpc.CallOption) (*GetShareInfoByIdResponse, error)
	DeleteShareInfoById(ctx context.Context, in *DeleteShareInfoRequest, opts ...grpc.CallOption) (*DeleteShareInfoResponse, error)
}

type shareInfoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShareInfoServiceClient(cc grpc.ClientConnInterface) ShareInfoServiceClient {
	return &shareInfoServiceClient{cc}
}

func (c *shareInfoServiceClient) PostShareInfo(ctx context.Context, in *PostShareInfoRequest, opts ...grpc.CallOption) (*PostShareInfoResponse, error) {
	out := new(PostShareInfoResponse)
	err := c.cc.Invoke(ctx, "/proto.shareInfo.ShareInfoService/PostShareInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shareInfoServiceClient) GetShareInfos(ctx context.Context, in *GetShareInfoRequest, opts ...grpc.CallOption) (*GetShareInfoResponse, error) {
	out := new(GetShareInfoResponse)
	err := c.cc.Invoke(ctx, "/proto.shareInfo.ShareInfoService/GetShareInfos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shareInfoServiceClient) GetShareInfoById(ctx context.Context, in *GetShareInfoByIdRequest, opts ...grpc.CallOption) (*GetShareInfoByIdResponse, error) {
	out := new(GetShareInfoByIdResponse)
	err := c.cc.Invoke(ctx, "/proto.shareInfo.ShareInfoService/GetShareInfoById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shareInfoServiceClient) DeleteShareInfoById(ctx context.Context, in *DeleteShareInfoRequest, opts ...grpc.CallOption) (*DeleteShareInfoResponse, error) {
	out := new(DeleteShareInfoResponse)
	err := c.cc.Invoke(ctx, "/proto.shareInfo.ShareInfoService/DeleteShareInfoById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShareInfoServiceServer is the server API for ShareInfoService service.
// All implementations should embed UnimplementedShareInfoServiceServer
// for forward compatibility
type ShareInfoServiceServer interface {
	// Submit send to ott
	PostShareInfo(context.Context, *PostShareInfoRequest) (*PostShareInfoResponse, error)
	GetShareInfos(context.Context, *GetShareInfoRequest) (*GetShareInfoResponse, error)
	GetShareInfoById(context.Context, *GetShareInfoByIdRequest) (*GetShareInfoByIdResponse, error)
	DeleteShareInfoById(context.Context, *DeleteShareInfoRequest) (*DeleteShareInfoResponse, error)
}

// UnimplementedShareInfoServiceServer should be embedded to have forward compatible implementations.
type UnimplementedShareInfoServiceServer struct {
}

func (UnimplementedShareInfoServiceServer) PostShareInfo(context.Context, *PostShareInfoRequest) (*PostShareInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostShareInfo not implemented")
}
func (UnimplementedShareInfoServiceServer) GetShareInfos(context.Context, *GetShareInfoRequest) (*GetShareInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShareInfos not implemented")
}
func (UnimplementedShareInfoServiceServer) GetShareInfoById(context.Context, *GetShareInfoByIdRequest) (*GetShareInfoByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShareInfoById not implemented")
}
func (UnimplementedShareInfoServiceServer) DeleteShareInfoById(context.Context, *DeleteShareInfoRequest) (*DeleteShareInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteShareInfoById not implemented")
}

// UnsafeShareInfoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShareInfoServiceServer will
// result in compilation errors.
type UnsafeShareInfoServiceServer interface {
	mustEmbedUnimplementedShareInfoServiceServer()
}

func RegisterShareInfoServiceServer(s grpc.ServiceRegistrar, srv ShareInfoServiceServer) {
	s.RegisterService(&ShareInfoService_ServiceDesc, srv)
}

func _ShareInfoService_PostShareInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostShareInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShareInfoServiceServer).PostShareInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.shareInfo.ShareInfoService/PostShareInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShareInfoServiceServer).PostShareInfo(ctx, req.(*PostShareInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShareInfoService_GetShareInfos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShareInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShareInfoServiceServer).GetShareInfos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.shareInfo.ShareInfoService/GetShareInfos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShareInfoServiceServer).GetShareInfos(ctx, req.(*GetShareInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShareInfoService_GetShareInfoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShareInfoByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShareInfoServiceServer).GetShareInfoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.shareInfo.ShareInfoService/GetShareInfoById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShareInfoServiceServer).GetShareInfoById(ctx, req.(*GetShareInfoByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShareInfoService_DeleteShareInfoById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteShareInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShareInfoServiceServer).DeleteShareInfoById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.shareInfo.ShareInfoService/DeleteShareInfoById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShareInfoServiceServer).DeleteShareInfoById(ctx, req.(*DeleteShareInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShareInfoService_ServiceDesc is the grpc.ServiceDesc for ShareInfoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShareInfoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.shareInfo.ShareInfoService",
	HandlerType: (*ShareInfoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostShareInfo",
			Handler:    _ShareInfoService_PostShareInfo_Handler,
		},
		{
			MethodName: "GetShareInfos",
			Handler:    _ShareInfoService_GetShareInfos_Handler,
		},
		{
			MethodName: "GetShareInfoById",
			Handler:    _ShareInfoService_GetShareInfoById_Handler,
		},
		{
			MethodName: "DeleteShareInfoById",
			Handler:    _ShareInfoService_DeleteShareInfoById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/share_info/share_info.proto",
}
