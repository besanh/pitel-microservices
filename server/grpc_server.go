package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/tel4vn/fins-microservices/api/v1"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	pbChatAuth "github.com/tel4vn/fins-microservices/gen/proto/chat_auth"
	pbChatIntegrateSystem "github.com/tel4vn/fins-microservices/gen/proto/chat_integrate_system"
	pbChatRole "github.com/tel4vn/fins-microservices/gen/proto/chat_role"
	pbChatUser "github.com/tel4vn/fins-microservices/gen/proto/chat_user"
	pbChatVendor "github.com/tel4vn/fins-microservices/gen/proto/chat_vendor"
	pbExample "github.com/tel4vn/fins-microservices/gen/proto/example"
	grpcService "github.com/tel4vn/fins-microservices/grpc"
	"github.com/tel4vn/fins-microservices/service"

	authMiddleware "github.com/tel4vn/fins-microservices/middleware/auth"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in the allowed header, don't send the header
	return s, false
}

func NewGRPCServer(port string) {
	// Setup gRPC
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(authMiddleware.GRPCAuthInterceptor, grpcauth.UnaryServerInterceptor(authMiddleware.GRPCAuthMiddleware)),
		),
	)
	pbExample.RegisterExampleServiceServer(grpcServer, grpcService.NewGRPCExample())
	pbChatIntegrateSystem.RegisterChatIntegrateSystemServer(grpcServer, grpcService.NewGRPCChatIntegrateSystem())
	pbChatRole.RegisterChatRoleServiceServer(grpcServer, grpcService.NewChatRole())
	pbChatUser.RegisterChatUserServiceServer(grpcServer, grpcService.NewGRPCChatUser())
	pbChatAuth.RegisterChatAuthServiceServer(grpcServer, grpcService.NewGRPCChatAuth())
	pbChatAuth.RegisterChatTokenServiceServer(grpcServer, grpcService.NewGRPCChatToken())
	pbChatVendor.RegisterChatVendorServiceServer(grpcServer, grpcService.NewGRPCChatVendor())

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(handleMatchHeaders),
		runtime.WithMetadata(handleMetadata),
		runtime.WithErrorHandler(handleError),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// setting up a dial up for gRPC service by specifying endpoint/target url
	if err := pbExample.RegisterExampleServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatIntegrateSystem.RegisterChatIntegrateSystemHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatIntegrateSystem.RegisterChatIntegrateSystemHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatRole.RegisterChatRoleServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatUser.RegisterChatUserServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatAuth.RegisterChatAuthServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatAuth.RegisterChatTokenServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err := pbChatVendor.RegisterChatVendorServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}

	// Creating a normal HTTP server
	httpServer := NewHTTPServer()
	v1.APIChatVendorHandler = v1.NewChatVendor()
	httpServer.Group("bss-chat/*{grpc_gateway}").Any("", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-vendor/upload") && c.Request.Method == "POST":
			v1.APIChatVendorHandler.HandlePostChatVendorLogoUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-vendor/upload/") && c.Request.Method == "PUT":
			v1.APIChatVendorHandler.HandlePutChatVendorLogoUpload(c)
		default:
			gin.WrapH(mux)
		}
	})
	v1.NewOttMessage(httpServer, service.NewOttMessage(), service.NewChatConnectionApp(), service.NewConversation())
	v1.NewMessage(httpServer, service.NewMessage())
	v1.NewWebSocket(httpServer, service.NewSubscriberService())
	v1.NewConversation(httpServer, service.NewConversation())
	v1.NewChatApp(httpServer, service.NewChatApp())
	v1.NewChatConnectionApp(httpServer, service.NewChatConnectionApp())
	v1.NewChatRouting(httpServer, service.NewChatRouting())
	v1.NewChatQueue(httpServer, service.NewChatQueue())
	v1.NewChatQueueUser(httpServer, service.NewChatQueueUser())
	v1.NewShareInfo(httpServer, service.NewShareInfo())
	v1.NewFacebook(httpServer, service.NewFacebook())
	v1.NewManageQueue(httpServer, service.NewManageQueue())
	v1.NewAssignConversation(httpServer, service.NewAssignConversation())
	v1.NewProfile(httpServer, service.NewProfile())
	v1.NewChatEmail(httpServer, service.NewChatEmail())
	v1.NewChatConnectionQueue(httpServer, service.NewChatConnectionQueue())
	v1.NewTest(httpServer)
	v1.NewChatMsgSample(httpServer, service.NewChatMsgSample())
	v1.NewChatScript(httpServer, service.NewChatScript())
	v1.NewChatAutoScript(httpServer, service.NewChatAutoScript())
	v1.NewChatLabel(httpServer, service.NewChatLabel())
	v1.NewChatPolicySetting(httpServer, service.NewChatPolicySetting())
	// httpServer.Static("/swagger/", "swagger-ui/")
	// httpServer.Static("/swagger-doc/", "gen/openapiv2/proto/pb")
	mixedHandler := newHTTPandGRPC(httpServer, grpcServer)
	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(mixedHandler, http2Server)}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Infof("HTTP server listening on %s", lis.Addr())
	log.Infof("GRPC server listening on %s", lis.Addr())
	err = http1Server.Serve(lis)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		panic(err)
	}

}

func newHTTPandGRPC(httpHand http.Handler, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			grpcHandler.ServeHTTP(w, r)
			return
		}
		httpHand.ServeHTTP(w, r)
	})
}

func handleMatchHeaders(key string) (string, bool) {
	switch key {
	default:
		return key, false
	}
}

func handleMetadata(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)
	md["tenant_id"] = r.Header.Get("X-Tenant-Id")
	md["user_id"] = r.Header.Get("X-User-Id")
	md["username"] = r.Header.Get("X-Username")
	md["level"] = r.Header.Get("X-User-Level")
	md["role_id"] = r.Header.Get("X-Role-Id")
	md["secret_key"] = r.Header.Get("secret-key")

	return metadata.New(md)
}

func handleError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	code, response := response.HandleGRPCErrResponse(err)
	HTTPErrorHandler(ctx, marshaler, writer, request, code, response)
}

func HTTPErrorHandler(ctx context.Context, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, codes codes.Code, resp any) {
	// return Internal when Marshal failed
	const fallback = `{"code": 13, "message": "failed to marshal error message"}`

	buf, merr := marshaler.Marshal(resp)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", resp, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")
	w.Header().Set("Content-Type", "application/json")
	st := runtime.HTTPStatusFromCode(codes)
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
