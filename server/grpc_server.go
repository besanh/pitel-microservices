package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	apiv1 "github.com/tel4vn/fins-microservices/api/v1"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	pbExample "github.com/tel4vn/fins-microservices/gen/proto/example"
	pbWebhookIncom "github.com/tel4vn/fins-microservices/gen/proto/incom"
	grpcService "github.com/tel4vn/fins-microservices/grpc"

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
			grpcMiddleware.ChainUnaryServer(grpcauth.UnaryServerInterceptor(authMiddleware.AuthMdw.GRPCAuthMiddleware)),
		),
	)
	grpcServerWebhook := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(authMiddleware.GRPCAuthInterceptor, grpcauth.UnaryServerInterceptor(authMiddleware.AuthMdw.GRPCAuthMiddleware)),
		),
	)
	pbExample.RegisterExampleServiceServer(grpcServer, grpcService.NewGRPCExample())
	pbWebhookIncom.RegisterIncomServiceServer(grpcServerWebhook, grpcService.NewGRPCIncom())
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
					EmitUnpopulated: false,
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
	// Creating a normal HTTP server
	httpServer := NewHTTPServer()
	apiv1.NewAbenlaWebhook(httpServer)
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
	md["business_unit_id"] = r.Header.Get("X-Business-Unit-Id")
	md["user_id"] = r.Header.Get("X-User-Id")
	md["username"] = r.Header.Get("X-Username")
	md["services"] = strings.Join(r.Header["X-Services"], ",")
	md["database_name"] = r.Header.Get("X-Database-Name")
	md["database_host"] = r.Header.Get("X-Database-Host")
	md["database_port"] = r.Header.Get("X-Database-Port")
	md["database_user"] = r.Header.Get("X-Database-User")
	md["database_password"] = r.Header.Get("X-Database-Password")
	md["database_es_host"] = r.Header.Get("X-Database-Es-Host")
	md["database_es_user"] = r.Header.Get("X-Database-Es-User")
	md["database_es_password"] = r.Header.Get("X-Database-Es-Password")
	md["database_es_index"] = r.Header.Get("X-Database-Es-Index")
	md["x_signature"] = r.Header.Get("X-Signature")
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
