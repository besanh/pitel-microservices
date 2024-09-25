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
	pbAssignConversation "github.com/tel4vn/fins-microservices/gen/proto/assign_conversation"
	pbChatApp "github.com/tel4vn/fins-microservices/gen/proto/chat_app"
	pbChatAuth "github.com/tel4vn/fins-microservices/gen/proto/chat_auth"
	pbChatAutoScript "github.com/tel4vn/fins-microservices/gen/proto/chat_auto_script"
	pbChatConnectionApp "github.com/tel4vn/fins-microservices/gen/proto/chat_connection_app"
	pbChatConnectionPipeline "github.com/tel4vn/fins-microservices/gen/proto/chat_connection_pipeline"
	pbConnectionQueue "github.com/tel4vn/fins-microservices/gen/proto/chat_connection_queue"
	pbChatEmail "github.com/tel4vn/fins-microservices/gen/proto/chat_email"
	pbChatIntegrateSystem "github.com/tel4vn/fins-microservices/gen/proto/chat_integrate_system"
	pbChatLabel "github.com/tel4vn/fins-microservices/gen/proto/chat_label"
	pbChatManageQueue "github.com/tel4vn/fins-microservices/gen/proto/chat_manage_queue"
	pbChatMessageSample "github.com/tel4vn/fins-microservices/gen/proto/chat_message_sample"
	pbChatPolicySetting "github.com/tel4vn/fins-microservices/gen/proto/chat_policy_setting"
	pbChatQueue "github.com/tel4vn/fins-microservices/gen/proto/chat_queue"
	pbChatQueueUser "github.com/tel4vn/fins-microservices/gen/proto/chat_queue_user"
	pbChatReport "github.com/tel4vn/fins-microservices/gen/proto/chat_report"
	pbChatRole "github.com/tel4vn/fins-microservices/gen/proto/chat_role"
	pbChatRouting "github.com/tel4vn/fins-microservices/gen/proto/chat_routing"
	pbChatScript "github.com/tel4vn/fins-microservices/gen/proto/chat_script"
	pbChatTenant "github.com/tel4vn/fins-microservices/gen/proto/chat_tenant"
	pbChatUser "github.com/tel4vn/fins-microservices/gen/proto/chat_user"
	pbChatVendor "github.com/tel4vn/fins-microservices/gen/proto/chat_vendor"
	pbConversation "github.com/tel4vn/fins-microservices/gen/proto/conversation"
	pbExample "github.com/tel4vn/fins-microservices/gen/proto/example"
	pbMessage "github.com/tel4vn/fins-microservices/gen/proto/message"
	pbNotesList "github.com/tel4vn/fins-microservices/gen/proto/note_list"
	pbProfile "github.com/tel4vn/fins-microservices/gen/proto/profile"
	pbShareInfo "github.com/tel4vn/fins-microservices/gen/proto/share_info"
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
	pbChatTenant.RegisterChatTenantServiceServer(grpcServer, grpcService.NewGRPCChatTenant())
	pbChatVendor.RegisterChatVendorServiceServer(grpcServer, grpcService.NewGRPCChatVendor())
	pbChatApp.RegisterChatAppServiceServer(grpcServer, grpcService.NewGRPCChatApp())
	pbAssignConversation.RegisterAssignConversationServiceServer(grpcServer, grpcService.NewGRPCAssignConversation())
	pbChatMessageSample.RegisterMessageSampleServiceServer(grpcServer, grpcService.NewGRPCChatMessageSample())
	pbChatScript.RegisterChatScriptServiceServer(grpcServer, grpcService.NewGRPCChatScript())
	pbChatAutoScript.RegisterChatAutoScriptServiceServer(grpcServer, grpcService.NewGRPCChatAutoScript())
	pbChatEmail.RegisterChatEmailServiceServer(grpcServer, grpcService.NewGRPCChatEmail())
	pbChatConnectionPipeline.RegisterChatConnectionPipelineServiceServer(grpcServer, grpcService.NewGRPCChatConnectionPipeline())
	pbChatConnectionApp.RegisterChatConnectionAppServiceServer(grpcServer, grpcService.NewGRPCChatConnectionApp())
	pbChatQueue.RegisterChatQueueServiceServer(grpcServer, grpcService.NewGRPCChatQueue())
	pbConversation.RegisterConversationServiceServer(grpcServer, grpcService.NewGRPCConversation())
	pbMessage.RegisterMessageServiceServer(grpcServer, grpcService.NewGRPCMessage())
	pbChatLabel.RegisterChatLabelServiceServer(grpcServer, grpcService.NewGRPCChatLabel())
	pbChatPolicySetting.RegisterChatPolicySettingServiceServer(grpcServer, grpcService.NewGRPCChatPolicySetting())
	pbChatRouting.RegisterChatRoutingServiceServer(grpcServer, grpcService.NewGRPCChatRouting())
	pbProfile.RegisterProfileServiceServer(grpcServer, grpcService.NewGRPCProfile())
	pbConnectionQueue.RegisterChatConnectionQueueServiceServer(grpcServer, grpcService.NewGRPCChatConnectionQueue())
	pbChatManageQueue.RegisterChatManageQueueServiceServer(grpcServer, grpcService.NewGRPCChatManageQueue())
	pbChatQueueUser.RegisterChatQueueUserServiceServer(grpcServer, grpcService.NewGRPCChatQueueUser())
	pbShareInfo.RegisterShareInfoServiceServer(grpcServer, grpcService.NewGRPCShareInfo())
	pbNotesList.RegisterNotesListServiceServer(grpcServer, grpcService.NewGRPCNotesList())
	pbChatReport.RegisterChatReportServiceServer(grpcServer, grpcService.NewGRPCChatReport())

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
	var err error
	if err = pbExample.RegisterExampleServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatIntegrateSystem.RegisterChatIntegrateSystemHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatIntegrateSystem.RegisterChatIntegrateSystemHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatRole.RegisterChatRoleServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatUser.RegisterChatUserServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatAuth.RegisterChatAuthServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatAuth.RegisterChatTokenServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatVendor.RegisterChatVendorServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatTenant.RegisterChatTenantServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatApp.RegisterChatAppServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbAssignConversation.RegisterAssignConversationServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatMessageSample.RegisterMessageSampleServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatScript.RegisterChatScriptServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatAutoScript.RegisterChatAutoScriptServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatEmail.RegisterChatEmailServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatConnectionApp.RegisterChatConnectionAppServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatConnectionPipeline.RegisterChatConnectionPipelineServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatQueue.RegisterChatQueueServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbConversation.RegisterConversationServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbMessage.RegisterMessageServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatPolicySetting.RegisterChatPolicySettingServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatLabel.RegisterChatLabelServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatRouting.RegisterChatRoutingServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbProfile.RegisterProfileServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbConnectionQueue.RegisterChatConnectionQueueServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatManageQueue.RegisterChatManageQueueServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatQueueUser.RegisterChatQueueUserServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbShareInfo.RegisterShareInfoServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbNotesList.RegisterNotesListServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}
	if err = pbChatReport.RegisterChatReportServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts); err != nil {
		log.Fatal(err)
	}

	// Creating a normal HTTP server
	httpServer := NewHTTPServer()
	httpServer.Group("bss-chat/*{grpc_gateway}").Any("", func(c *gin.Context) {
		switch {
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-vendor/upload") && c.Request.Method == "POST":
			v1.APIChatVendorHandler.HandlePostChatVendorLogoUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-vendor/upload/") && c.Request.Method == "PUT":
			v1.APIChatVendorHandler.HandlePutChatVendorLogoUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-sample/upload") && c.Request.Method == "POST":
			v1.APIChatMessageSampleHandler.HandlePostChatMessageSampleUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-sample/upload/") && c.Request.Method == "PUT":
			v1.APIChatMessageSampleHandler.HandlePutChatMessageSampleUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-script/upload") && c.Request.Method == "POST":
			v1.APIChatScript.HandlePostChatScriptUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-script/upload/") && c.Request.Method == "PUT":
			v1.APIChatScript.HandlePutChatScriptUpload(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/message/send") && c.Request.Method == "POST":
			v1.APIMessage.HandlePostSendMessage(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/share-info/config") && c.Request.Method == "POST":
			v1.APIShareInfo.HandlePostConfigForm(c) // create form
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/share-info/image/") && c.Request.Method == "GET":
			v1.APIShareInfo.HandleGetImageShareInfo(c)
		case strings.HasPrefix(c.Request.RequestURI, "/bss-chat/v1/share-info/") && c.Request.Method == "PUT":
			v1.APIShareInfo.HandlePutShareInfoById(c)
		default:
			gin.WrapH(mux)(c)
		}
	})

	v1.APIChatVendorHandler = v1.NewChatVendor()
	v1.APIChatMessageSampleHandler = v1.NewChatMessageSample()
	v1.APIChatScript = v1.NewAPIChatScript()
	v1.APIMessage = v1.NewAPIMessage()
	v1.APIShareInfo = v1.NewAPIShareInfo()
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
	md["system_id"] = r.Header.Get("system-key")

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
