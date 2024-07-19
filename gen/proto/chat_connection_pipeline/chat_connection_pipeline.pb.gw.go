// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: proto/chat_connection_pipeline/chat_connection_pipeline.proto

/*
Package pb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package pb

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_ChatConnectionPipelineService_AttachConnectionQueueToApp_0(ctx context.Context, marshaler runtime.Marshaler, client ChatConnectionPipelineServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ChatConnectionPipelineQueueRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.AttachConnectionQueueToApp(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ChatConnectionPipelineService_AttachConnectionQueueToApp_0(ctx context.Context, marshaler runtime.Marshaler, server ChatConnectionPipelineServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ChatConnectionPipelineQueueRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.AttachConnectionQueueToApp(ctx, &protoReq)
	return msg, metadata, err

}

func request_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0(ctx context.Context, marshaler runtime.Marshaler, client ChatConnectionPipelineServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq UpsertQueueInConnectionAppByIdRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	msg, err := client.UpsertQueueInConnectionAppById(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0(ctx context.Context, marshaler runtime.Marshaler, server ChatConnectionPipelineServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq UpsertQueueInConnectionAppByIdRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	msg, err := server.UpsertQueueInConnectionAppById(ctx, &protoReq)
	return msg, metadata, err

}

func request_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0(ctx context.Context, marshaler runtime.Marshaler, client ChatConnectionPipelineServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq PutChatConnectionAppStatusRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	msg, err := client.UpdateChatConnectionAppStatusById(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0(ctx context.Context, marshaler runtime.Marshaler, server ChatConnectionPipelineServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq PutChatConnectionAppStatusRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "id")
	}

	protoReq.Id, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "id", err)
	}

	msg, err := server.UpdateChatConnectionAppStatusById(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterChatConnectionPipelineServiceHandlerServer registers the http handlers for service ChatConnectionPipelineService to "mux".
// UnaryRPC     :call ChatConnectionPipelineServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterChatConnectionPipelineServiceHandlerFromEndpoint instead.
func RegisterChatConnectionPipelineServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ChatConnectionPipelineServiceServer) error {

	mux.Handle("POST", pattern_ChatConnectionPipelineService_AttachConnectionQueueToApp_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.chatConnectionPipeline.ChatConnectionPipelineService/AttachConnectionQueueToApp", runtime.WithHTTPPathPattern("/bss-chat/v1/chat-connection-pipeline"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ChatConnectionPipelineService_AttachConnectionQueueToApp_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ChatConnectionPipelineService_AttachConnectionQueueToApp_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.chatConnectionPipeline.ChatConnectionPipelineService/UpsertQueueInConnectionAppById", runtime.WithHTTPPathPattern("/bss-chat/v1/chat-connection-pipeline/edit/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/proto.chatConnectionPipeline.ChatConnectionPipelineService/UpdateChatConnectionAppStatusById", runtime.WithHTTPPathPattern("/bss-chat/v1/chat-connection-pipeline/edit/status/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterChatConnectionPipelineServiceHandlerFromEndpoint is same as RegisterChatConnectionPipelineServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterChatConnectionPipelineServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.NewClient(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Errorf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterChatConnectionPipelineServiceHandler(ctx, mux, conn)
}

// RegisterChatConnectionPipelineServiceHandler registers the http handlers for service ChatConnectionPipelineService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterChatConnectionPipelineServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterChatConnectionPipelineServiceHandlerClient(ctx, mux, NewChatConnectionPipelineServiceClient(conn))
}

// RegisterChatConnectionPipelineServiceHandlerClient registers the http handlers for service ChatConnectionPipelineService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ChatConnectionPipelineServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ChatConnectionPipelineServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ChatConnectionPipelineServiceClient" to call the correct interceptors.
func RegisterChatConnectionPipelineServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ChatConnectionPipelineServiceClient) error {

	mux.Handle("POST", pattern_ChatConnectionPipelineService_AttachConnectionQueueToApp_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/proto.chatConnectionPipeline.ChatConnectionPipelineService/AttachConnectionQueueToApp", runtime.WithHTTPPathPattern("/bss-chat/v1/chat-connection-pipeline"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ChatConnectionPipelineService_AttachConnectionQueueToApp_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ChatConnectionPipelineService_AttachConnectionQueueToApp_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/proto.chatConnectionPipeline.ChatConnectionPipelineService/UpsertQueueInConnectionAppById", runtime.WithHTTPPathPattern("/bss-chat/v1/chat-connection-pipeline/edit/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/proto.chatConnectionPipeline.ChatConnectionPipelineService/UpdateChatConnectionAppStatusById", runtime.WithHTTPPathPattern("/bss-chat/v1/chat-connection-pipeline/edit/status/{id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ChatConnectionPipelineService_AttachConnectionQueueToApp_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"bss-chat", "v1", "chat-connection-pipeline"}, ""))

	pattern_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 1, 0, 4, 1, 5, 4}, []string{"bss-chat", "v1", "chat-connection-pipeline", "edit", "id"}, ""))

	pattern_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 2, 4, 1, 0, 4, 1, 5, 5}, []string{"bss-chat", "v1", "chat-connection-pipeline", "edit", "status", "id"}, ""))
)

var (
	forward_ChatConnectionPipelineService_AttachConnectionQueueToApp_0 = runtime.ForwardResponseMessage

	forward_ChatConnectionPipelineService_UpsertQueueInConnectionAppById_0 = runtime.ForwardResponseMessage

	forward_ChatConnectionPipelineService_UpdateChatConnectionAppStatusById_0 = runtime.ForwardResponseMessage
)
